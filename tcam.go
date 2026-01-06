package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/s5i/tcam/enum"
	"github.com/s5i/tcam/gamedata"
	"github.com/s5i/tcam/loader"
	"github.com/s5i/tcam/network"
	"github.com/s5i/tcam/parser"
	"golang.org/x/sync/errgroup"
)

var (
	inputDir = flag.String("dir", ".", "Directory to process.")
	loggers  = flag.String("loggers", "", "Comma-separated list of loggers to enable (main, loader, parser).")
	tibiaDat = flag.String("dat", "", "Path to Tibia.dat.")
)

var Logger = log.New(io.Discard, "[MAIN] ", 0)
var MsgLogger = log.New(io.Discard, "[MSG] ", 0)

func main() {
	ctx := context.Background()

	flag.Parse()

	for l := range strings.SplitSeq(*loggers, ",") {
		switch l {
		case "":
		case "main":
			Logger.SetOutput(os.Stderr)
		case "loader":
			loader.Logger.SetOutput(os.Stderr)
		case "parser":
			parser.Logger.SetOutput(os.Stderr)
		case "gamedata":
			gamedata.Logger.SetOutput(os.Stderr)
		case "msg":
			MsgLogger.SetOutput(os.Stderr)
		default:
			fmt.Fprintf(os.Stderr, "Unknown logger: %q\n", l)
			os.Exit(1)
		}
	}

	f, err := os.Open(*tibiaDat)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open Tibia.dat: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	if err := gamedata.ReadFile(ctx, *tibiaDat); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load %q: %v\n", *tibiaDat, err)
		os.Exit(1)
	}

	if err := processDir(ctx, *inputDir); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func processDir(ctx context.Context, dirPath string) error {
	return filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		t := time.Now()
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if strings.ToLower(filepath.Ext(path)) != ".cam" {
			return nil
		}

		defer func() {
			Logger.Printf("Processed %q in %s", path, time.Since(t))
		}()

		eg, ctx := errgroup.WithContext(ctx)

		loaderOutCh, loaderErrCh := loader.ReadFile(ctx, path)
		parserInCh := make(chan *network.Packet)
		parserOutCh, parserErrCh := parser.ParsePackets(ctx, parserInCh)

		eg.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()

				case err := <-loaderErrCh:
					return err

				case pkt, ok := <-loaderOutCh:
					if !ok {
						close(parserInCh)
						return nil
					}

					select {
					case <-ctx.Done():
						return ctx.Err()

					case parserInCh <- pkt:
					}
				}
			}
		})

		eg.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()

				case err := <-parserErrCh:
					return err

				case x, ok := <-parserOutCh:
					if !ok {
						return nil
					}

					switch x := x.(type) {
					case enum.OpCode:
						Logger.Printf("unhandled opcode: %s", x)
					case *parser.Talk:
						if x.Mode == enum.MessageModeMessageSay {
							MsgLogger.Printf("[%10v] %s %s: %s", x.TimeOffset.Truncate(time.Second), x.Offset(), x.Name, x.Msg)
						}
					}
				}
			}
		})

		if err := eg.Wait(); err != nil {
			return fmt.Errorf("error processing %q: %v", path, err)
		}

		return nil
	})
}

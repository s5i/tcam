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

	"github.com/s5i/tcam/loader"
	"github.com/s5i/tcam/network"
	"github.com/s5i/tcam/parser"
	"golang.org/x/sync/errgroup"
)

var (
	inputDir = flag.String("dir", ".", "Directory to process.")
	loggers  = flag.String("loggers", "", "Comma-separated list of loggers to enable (main, loader, parser).")
)

var Logger = log.New(io.Discard, "[MAIN] ", 0)

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
		default:
			fmt.Fprintf(os.Stderr, "Unknown logger: %q\n", l)
			os.Exit(1)
		}
	}

	err := processDir(ctx, *inputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing directory: %v\n", err)
		os.Exit(1)
	}
}

func processDir(ctx context.Context, dirPath string) error {
	return filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			Logger.Printf("Skipping %q due to an error: %v", path, err)
			return err
		}

		if d.IsDir() {
			return nil
		}

		if strings.ToLower(filepath.Ext(path)) != ".cam" {
			Logger.Printf("Skipping %q: non-cam file", path)
			return nil
		}

		Logger.Printf("Processing: %s", path)

		eg, ctx := errgroup.WithContext(ctx)

		loaderOutCh, loaderErrCh := loader.ReadFile(ctx, path)
		parserInCh := make(chan network.Packet)
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
					case *parser.Talk:
						Logger.Printf("%s: %s", x.Name, x.Msg)
					}
				}
			}
		})

		return eg.Wait()
	})
}

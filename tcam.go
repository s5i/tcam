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
	loggers  = flag.String("loggers", "", "Comma-separated list of loggers to enable (main, loader, parser, msg, nooutput).")
	// tibiaDat = flag.String("dat", "", "Path to Tibia.dat.")
	output = flag.String("output", "", "Output file.")
	player = flag.String("player", "", "Player name.")
	npcs   = flag.String("npcs", "", "Comma-separated list of NPCs to process.")
)

var Logger = log.New(io.Discard, "[MAIN] ", 0)
var MsgLogger = log.New(io.Discard, "[MSG] ", 0)
var OutputLogger = log.New(os.Stdout, "", 0)

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
		case "nooutput":
			OutputLogger.SetOutput(io.Discard)
		default:
			fmt.Fprintf(os.Stderr, "Unknown logger: %q\n", l)
			os.Exit(1)
		}
	}

	// f, err := os.Open(*tibiaDat)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Failed to open Tibia.dat: %v\n", err)
	// 	os.Exit(1)
	// }
	// defer f.Close()

	// if err := gamedata.ReadFile(ctx, *tibiaDat); err != nil {
	// 	fmt.Fprintf(os.Stderr, "Failed to load %q: %v\n", *tibiaDat, err)
	// 	os.Exit(1)
	// }

	if *output != "" {
		f, err := os.Create(*output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create %q: %v\n", *output, err)
			os.Exit(1)
		}
		defer f.Close()

		OutputLogger.SetOutput(f)
		var args []string
		for _, arg := range os.Args {
			args = append(args, fmt.Sprintf("%q", arg))
		}
		OutputLogger.Printf("%s", strings.Join(args, " "))
		OutputLogger.Printf("Timestamp: %s", time.Now().Format("2006-01-02 15:04:05 MST"))
	}

	if err := processDir(ctx, *inputDir, *player, *npcs); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func processDir(ctx context.Context, dirPath string, player string, npcs string) error {
	npc := map[string]bool{}
	for n := range strings.SplitSeq(npcs, ",") {
		if n != "" {
			npc[strings.ToLower(n)] = true
		}
	}
	player = strings.ToLower(player)

	dialogueCamSep := true
	var lastDialogueOffset time.Duration
	var lastDialogueNPC string
	return filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		defer func() { dialogueCamSep = true; lastDialogueOffset = 0 }()

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
			var playerMsg string
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
					case *parser.UnhandledPacket:
						Logger.Printf("unhandled packet: %s", x.Packet)
					case *parser.Talk:
						if x.Mode == enum.MessageModeMessageSay {
							MsgLogger.Printf("[%10v] %s %s: %s", x.TimeOffset.Truncate(time.Second), x.Offset(), x.Name, x.Msg)

							n := strings.ToLower(x.Name)

							dialogueOffsetSep := x.TimeOffset-lastDialogueOffset > 5*time.Minute
							dialogueNPCSep := lastDialogueNPC != n

							if npc[n] {
								// OutputLogger.Printf("Last dialogue NPC: %q, n: %q", lastDialogueNPC, n)
								if dialogueCamSep || dialogueOffsetSep || dialogueNPCSep {
									OutputLogger.Printf("--------------------------------------------------------------------------------")
									dialogueCamSep = false
								}

								lastDialogueOffset = x.TimeOffset
								lastDialogueNPC = n

								if playerMsg != "" {
									OutputLogger.Printf("%s", playerMsg)
									playerMsg = ""
								}

								OutputLogger.Printf("%s: %s", x.Name, x.Msg)
							}

							if n == player {
								playerMsg = fmt.Sprintf("%s: %s", x.Name, x.Msg)
							}
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

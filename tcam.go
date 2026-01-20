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
	"github.com/s5i/tcam/msgcontext"
	"github.com/s5i/tcam/network"
	"github.com/s5i/tcam/npc"
	"github.com/s5i/tcam/parser"
	"golang.org/x/sync/errgroup"
)

var (
	inputDir    = flag.String("dir", ".", "Directory to process.")
	loggers     = flag.String("loggers", "", "Comma-separated list of loggers to enable (debug, loader, parser, msg, nooutput).")
	output      = flag.String("output", "", "Output file.")
	player      = flag.String("player", "", "[manual-interaction] Player name.")
	npcs        = flag.String("npcs", "", "[manual-interaction] Comma-separated list of NPCs to process.")
	contextSize = flag.Int("context_size", 1, "[non-retail] Conversation context length.")
	mode        = flag.String("mode", "non-retail", "Mode to use (non-retail, manual-interaction).")
)

func main() {
	ctx := context.Background()

	flag.Parse()

	for l := range strings.SplitSeq(*loggers, ",") {
		switch l {
		case "":
		case "debug":
			Logger.SetOutput(os.Stderr)
		case "loader":
			loader.Logger.SetOutput(os.Stdout)
		case "parser":
			parser.Logger.SetOutput(os.Stdout)
		case "gamedata":
			gamedata.Logger.SetOutput(os.Stdout)
		case "msg":
			MsgLogger.SetOutput(os.Stdout)
		case "nooutput":
			OutputLogger.SetOutput(io.Discard)
		default:
			fmt.Fprintf(os.Stderr, "Unknown logger: %q\n", l)
			os.Exit(1)
		}
	}

	// TODO(s5i): restore this once fixed.
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

	switch *mode {
	case "non-retail":
		if err := nonRetail(ctx, *inputDir, *contextSize); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	case "manual-interaction":
		if err := manualInteraction(ctx, *inputDir, *player, *npcs); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown mode: %q\n", *mode)
		os.Exit(1)
	}
}

func manualInteraction(ctx context.Context, dirPath string, player string, npcs string) error {
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
		Logger.Printf("Processing %q\n", path)
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
			Logger.Printf("Processed %q in %s\n", path, time.Since(t))
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
					case *parser.Talk:
						if x.Mode != enum.MessageModeMessageSay {
							continue
						}

						MsgLogger.Printf("[%10v] %s %s: %s", x.TimeOffset.Truncate(time.Second), x.Offset(), x.Name, x.Msg)

						n := strings.ToLower(x.Name)

						dialogueOffsetSep := x.TimeOffset-lastDialogueOffset > 5*time.Minute
						dialogueNPCSep := lastDialogueNPC != n

						if npc[n] {
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
		})

		if err := eg.Wait(); err != nil {
			return fmt.Errorf("error processing %q: %v", path, err)
		}

		return nil
	})
}

func nonRetail(ctx context.Context, dirPath string, contextSize int) error {
	seen := map[string]bool{}

	return filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		Logger.Printf("Processing %q\n", path)
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
			Logger.Printf("Processed %q in %s\n", path, time.Since(t))
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
			msgCtx := msgcontext.NewContext(contextSize)
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
						if x.Mode != enum.MessageModeMessageSay {
							continue
						}

						if npc.IsNPC(x.Name) && !npc.IsRetailResponse(x.Name, x.Msg) {
							resp := fmt.Sprintf("%s: %s", x.Name, x.Msg)
							if seen[resp] {
								continue
							}
							seen[resp] = true

							OutputLogger.Printf("--------------------------------------------------------------------------------")
							for _, msg := range msgCtx.Pop() {
								OutputLogger.Printf("%s: %s", msg.Name, msg.Message)
							}
							OutputLogger.Printf("%s: %s", x.Name, x.Msg)
						}

						msgCtx.Put(x.Name, x.Msg)
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

var Logger = log.New(io.Discard, "[MAIN] ", 0)
var MsgLogger = log.New(io.Discard, "[MSG] ", 0)
var OutputLogger = log.New(os.Stdout, "", 0)

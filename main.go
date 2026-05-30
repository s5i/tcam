package main

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

func main() {
	var camDir string
	var outputFile string
	var statsCmdNoFilter bool
	var locationRadius int

	output := func() (io.Writer, func(), error) {
		if outputFile == "" {
			return os.Stdout, func() {}, nil
		}

		if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
			return nil, nil, err
		}

		f, err := os.Create(outputFile)
		if err != nil {
			return nil, nil, err
		}

		return f, func() { f.Close() }, nil
	}

	rootCmd := &cobra.Command{
		Use:   "tcam",
		Short: "Tibiantis CAM processor",
	}

	dialoguesCmd := &cobra.Command{
		Use:  "dialogues --camdir=... [--out=f] npc-name",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			npcName := args[0]
			out, outClose, err := output()
			if err != nil {
				return err
			}
			defer outClose()

			return Dialogues(ctx, camDir, out, npcName, time.Minute)
		},
	}

	dialogueTreeCmd := &cobra.Command{
		Use:  "dialogue-tree --camdir=... [--out=f]",
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			out, outClose, err := output()
			if err != nil {
				return err
			}
			defer outClose()

			return DialogueTree(ctx, camDir, out)
		},
	}

	locationCmd := &cobra.Command{
		Use:  "location --camdir=... [--out=f] [--radius=n] x y z",
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			x, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			y, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			z, err := strconv.Atoi(args[2])
			if err != nil {
				return err
			}

			out, outClose, err := output()
			if err != nil {
				return err
			}
			defer outClose()

			return Location(ctx, camDir, out, x, y, z, locationRadius, 10*time.Minute)
		},
	}

	creatureCmd := &cobra.Command{
		Use:  "creature --camdir=... [--out=f] name",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			name := args[0]

			out, outClose, err := output()
			if err != nil {
				return err
			}
			defer outClose()

			return Creature(ctx, camDir, out, name)
		},
	}

	statsCmd := &cobra.Command{
		Use:  "parse-stats --camdir=... [--out=f] [--nofilter]",
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			return ParseStats(ctx, camDir, os.Stderr, statsCmdNoFilter)
		},
	}

	rootCmd.PersistentFlags().StringVar(&camDir, "camdir", "", "Directory to search for .cam files.")
	rootCmd.PersistentFlags().StringVar(&outputFile, "out", "", "Output file path; stdout when empty.")
	statsCmd.PersistentFlags().BoolVar(&statsCmdNoFilter, "nofilter", false, "When true, skip the optype filter optimization.")
	locationCmd.PersistentFlags().IntVar(&locationRadius, "radius", 7, "Max difference to be considered 'in location' for X and Y parameters.")
	rootCmd.AddCommand(dialoguesCmd)
	rootCmd.AddCommand(dialogueTreeCmd)
	rootCmd.AddCommand(locationCmd)
	rootCmd.AddCommand(creatureCmd)
	rootCmd.AddCommand(statsCmd)

	if rootCmd.Execute() != nil {
		os.Exit(1)
	}
}

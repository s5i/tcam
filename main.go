package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "tcam",
		Short: "Tibiantis CAM processor",
	}

	dialoguesCmd := &cobra.Command{
		Use:  "dialogues npc-name output-file cam-directory",
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			npcName := args[0]
			outFile := args[1]
			camDir := args[2]

			if err := os.MkdirAll(filepath.Dir(outFile), 0755); err != nil {
				return err
			}

			f, err := os.Create(outFile)
			if err != nil {
				return err
			}

			return Dialogues(ctx, camDir, f, npcName, time.Minute)
		},
	}

	rootCmd.AddCommand(dialoguesCmd)

	if rootCmd.Execute() != nil {
		os.Exit(1)
	}
}

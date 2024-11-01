package cli

import (
	"github.com/ondrovic/video-organizer/internal/types"
	"github.com/spf13/cobra"
)

var (
	options = types.CliFlags{}
	RootCmd = &cobra.Command{
		Use:   "video-organizer",
		Short: "A CLI tool to organize videos by duration",
	}
)

func InitializeCommands() {
	RootCmd.AddCommand(checkCmd)
	RootCmd.AddCommand(organizeCmd)
}

func Execute() error {
	if err := RootCmd.Execute(); err != nil {
		return err
	}

	return nil
}

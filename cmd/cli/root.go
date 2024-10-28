package cli

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "video-organizer",
	Short: "A CLI tool to organize videos by duration",
}

func Execute() error {
	if err := RootCmd.Execute(); err != nil {
		return err
	}

	return nil
}

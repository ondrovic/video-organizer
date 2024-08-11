// app/cmd/organizer/video-organizer.go
package main

import (
	"os"
	"runtime"
	"video-organizer/internal/video"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	commonUtils "github.com/ondrovic/common/utils"
)

type Options struct {
	RootDirectoryPath string
}

func main() {
	commonUtils.ClearTerminalScreen(runtime.GOOS)

	var opts Options

	rootCmd := &cobra.Command{
		Use:   "video-sorter [root-directory]",
		Short: "A CLI tool that organizes video files in a directory based on their duration",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			opts.RootDirectoryPath = args[0]
			if err := video.OrganizeVideos(opts.RootDirectoryPath); err != nil {
				pterm.Error.Printf("Error organizing videos: %v\n", err)
			}
		},
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

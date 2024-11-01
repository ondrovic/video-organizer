package cli

import (
	"github.com/ondrovic/video-organizer/internal/video"
	"github.com/spf13/cobra"
)

var (
	organizeCmd = &cobra.Command{
		Use:   "organize <directory>",
		Short: "Organize the directory",
		Args:  cobra.ExactArgs(1),
		RunE:  runOrganize,
	}
)

func runOrganize(cmd *cobra.Command, args []string) error {
	options.Directory = video.FormatDirectory(args[0])

	if err := video.OrganizeVideos(options.Directory); err != nil {
		return err
	}

	return nil
}

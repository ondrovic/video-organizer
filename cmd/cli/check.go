package cli

import (
	"github.com/ondrovic/video-organizer/internal/video"
	"github.com/spf13/cobra"
)

var (
	checkCmd = &cobra.Command{
		Use:   "check <directory>",
		Short: "Check a directory to see if it's been organized",
		Args:  cobra.ExactArgs(1),
		RunE:  runCheck,
	}
)

func runCheck(cmd *cobra.Command, args []string) error {
	options.Directory = video.FormatDirectory(args[0])

	if err := video.CheckDirectories(options.Directory); err != nil {
		return err
	}

	return nil
}

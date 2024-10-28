package cli

import (
	"github.com/ondrovic/video-organizer/internal/types"
	"github.com/ondrovic/video-organizer/internal/video"
	"github.com/spf13/cobra"
)

var (
	options     = types.CliFlags{}
	organizeCmd = &cobra.Command{}
)

func init() {
	organizeCmd = &cobra.Command{
		Use:   "organize <directory>",
		Short: "Organize the directory",
		Args:  cobra.ExactArgs(1),
		RunE:  run,
	}

	RootCmd.AddCommand(organizeCmd)
}

func run(cmd *cobra.Command, args []string) error {
	options.Directory = video.FormatDirectory(args[0])

	if err := video.OrganizeVideos(options.Directory); err != nil {
		return err
	}

	return nil
}

// cmd/organizer/main.go
package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"video-organizer/internal/video"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type Options struct {
	RootDirectoryPath string
}

func clearConsole() {
	var clearCmd *exec.Cmd

	switch runtime.GOOS {
	case "linux", "darwin":
		clearCmd = exec.Command("clear")
	case "windows":
		clearCmd = exec.Command("cmd", "/c", "cls")
	default:
		fmt.Println("Unsupported platform")
		return
	}

	clearCmd.Stdout = os.Stdout
	clearCmd.Run()
}

func main() {
	clearConsole()
	
	var opts Options

	rootCmd := &cobra.Command{
		Use:   "video-organizer",
		Short: "A CLI tool to organize videos based on their duration",
		Run: func(cmd *cobra.Command, args []string) {
			if err := video.OrganizeVideos(opts.RootDirectoryPath); err != nil {
				// log.Fatalf("Error organizing videos: %v", err)
				pterm.Error.Printf("Error organizing videos: %v\n", err)
			}
		},
	}

	flags := rootCmd.Flags()
	flags.StringVarP(&opts.RootDirectoryPath, "root-directory", "r", "", "Root directory you want to organize")
	rootCmd.MarkFlagRequired("root-directory")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
 
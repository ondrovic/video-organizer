// internal/video/organizer.go
package video

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/pterm/pterm"
	"github.com/xfrr/goffmpeg/transcoder"
)

// Define a struct to hold folder information
type FolderInfo struct {
	Name        string
	MaxDuration float64
}

var folderStructure = []FolderInfo{
	{Name: "Micro", MaxDuration: 15},                       // 0 <= duration <= 15 seconds
	{Name: "Mini", MaxDuration: 60},                        // 15 < duration <= 60 seconds
	{Name: "Short", MaxDuration: 5 * 60},                   // 60 < duration <= 5*60 seconds (5 minutes)
	{Name: "Medium", MaxDuration: 15 * 60},                 // 560 < duration <= 1560 seconds (15 minutes)
	{Name: "Long", MaxDuration: 30 * 60},                   // 1560 < duration <= 3060 seconds (30 minutes)
	{Name: "Extended", MaxDuration: 60 * 60},               // 3060 < duration <= 6060 seconds (60 minutes)
	{Name: "Feature", MaxDuration: 120 * 60},               // 6060 < duration <= 12060 seconds (120 minutes)
	{Name: "Epic", MaxDuration: float64(^uint(0) >> 1)},    // > 120*60 seconds
}

func getDurationInSeconds(filePath string) (float64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	trans := new(transcoder.Transcoder)

	err = trans.Initialize(filePath, ".")
	if err != nil {
		panic(err)
	}

	durationStr := trans.MediaFile().Metadata().Format.Duration

	if durationStr == "" {
		return 0, fmt.Errorf("no duration found")
	}

	durationFlt, err := strconv.ParseFloat(durationStr, 64)

	if err != nil {
		return 0, fmt.Errorf("no duration found %v", err)
	}

	return durationFlt, nil
}

func isTargetFolder(folder string) bool {
	for _, info := range folderStructure {
		if folder == info.Name {
			return true
		}
	}
	return false
}

func getTargetFolder(duration float64) string {
	for _, info := range folderStructure {
		if duration < info.MaxDuration {
			return info.Name
		}
	}
	return folderStructure[len(folderStructure)-1].Name // Default to the last folder
}

func OrganizeVideos(rootDirectory string) error {
	videoFiles, err := getVideoFiles(rootDirectory)
	if err != nil {
		return err
	}

	var totalFiles int
	for _, files := range videoFiles {
		totalFiles += len(files)
	}

	progressbar, _ := pterm.DefaultProgressbar.WithTotal(totalFiles).WithTitle("Organizing videos").Start()

	for subDir, files := range videoFiles {
		for _, filePath := range files {
			currentFolder := filepath.Base(subDir)
			if isTargetFolder(currentFolder) {
				progressbar.Increment()
				continue
			}

			duration, err := getDurationInSeconds(filePath)
			if err != nil {
				return err
			}

			targetFolder := getTargetFolder(duration)

			// Create the target path within the same subdirectory
			targetPath := filepath.Join(filepath.Dir(filePath), targetFolder, filepath.Base(filePath))

			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return err
			}

			if err := os.Rename(filePath, targetPath); err != nil {
				return err
			}

			progressbar.Increment()
		}
	}

	progressbar.Stop()

	return nil
}

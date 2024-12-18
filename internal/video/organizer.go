// internal/video/organizer.go
package video

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pterm/pterm"
	"github.com/xfrr/goffmpeg/transcoder"
)

// Define a struct to hold folder information
type FolderInfo struct {
	Name        string
	MaxDuration float64
}

var (
	folderStructure = []FolderInfo{
		{Name: "Micro", MaxDuration: 15},                    // 0 <= duration <= 15 seconds
		{Name: "Mini", MaxDuration: 60},                     // 15 < duration <= 60 seconds
		{Name: "Short", MaxDuration: 5 * 60},                // 60 < duration <= 5*60 seconds (5 minutes)
		{Name: "Medium", MaxDuration: 15 * 60},              // 560 < duration <= 1560 seconds (15 minutes)
		{Name: "Long", MaxDuration: 30 * 60},                // 1560 < duration <= 3060 seconds (30 minutes)
		{Name: "Extended", MaxDuration: 60 * 60},            // 3060 < duration <= 6060 seconds (60 minutes)
		{Name: "Feature", MaxDuration: 120 * 60},            // 6060 < duration <= 12060 seconds (120 minutes)
		{Name: "Epic", MaxDuration: float64(^uint(0) >> 1)}, // > 120*60 seconds
	}
)

// var

func getDurationInSeconds(filePath string) (float64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	trans := new(transcoder.Transcoder)

	// Attempt to initialize the transcoder with the file path
	err = trans.Initialize(filePath, ".")
	if err != nil {
		// Check if the error contains the specific invalid data message
		if strings.Contains(err.Error(), "Invalid data found when processing input") {
			// skip as this is likely a bad video file
			return 0, nil
		}
		return 0, err // Return other errors as-is
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

	progressbar, _ := pterm.DefaultProgressbar.WithTotal(totalFiles).WithTitle("Organizing videos").WithRemoveWhenDone(true).Start()

	for subDir, files := range videoFiles {
		for _, filePath := range files {
			currentFolder := filepath.Base(subDir)
			if isTargetFolder(currentFolder) {
				progressbar.Increment()
				continue
			}

			duration, err := getDurationInSeconds(filePath)
			if err != nil {
				// Log the error and continue with the next file
				fmt.Printf("Error getting duration for %s: %v\n", filePath, err)
				progressbar.Increment()
				continue
			}

			targetFolder := getTargetFolder(duration)

			// Create the target path within the same subdirectory
			targetPath := filepath.Join(filepath.Dir(filePath), targetFolder, filepath.Base(filePath))

			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				// Log the error and continue with the next file
				fmt.Printf("Error creating directory for %s: %v\n", targetPath, err)
				progressbar.Increment()
				continue
			}

			if err := os.Rename(filePath, targetPath); err != nil {
				// Check if the error is due to the file being in use
				if os.IsPermission(err) || strings.Contains(err.Error(), "being used by another process") {
					fmt.Printf("Skipping file in use: %s\n", filePath)
				} else {
					fmt.Printf("Error renaming %s to %s: %v\n", filePath, targetPath, err)
				}
				progressbar.Increment()
				continue
			}

			progressbar.Increment()
		}
	}

	progressbar.Stop()

	return nil
}

func CheckDirectories(rootDirectory string) error {
	// Prepare table data for pterm
	rows := pterm.TableData{
		{"Directory", "Organized"},
	}

	// List top-level directories in rootDirectory
	entries, err := os.ReadDir(rootDirectory)
	if err != nil {
		return err
	}

	// Check each child directory for a "videos" subdirectory
	for _, entry := range entries {
		if entry.IsDir() {
			childDir := filepath.Join(rootDirectory, entry.Name())
			videoDir := filepath.Join(childDir, "videos")

			// Check if the "videos" directory exists within the child directory (case-insensitive)
			videoDirExists := false
			entries, err := os.ReadDir(childDir)
			if err != nil {
				return err
			}
			for _, e := range entries {
				if e.IsDir() && strings.EqualFold(e.Name(), "videos") {
					videoDirExists = true
					videoDir = filepath.Join(childDir, e.Name())
					break
				}
			}

			// If "videos" folder doesn't exist, add a "False" row and continue
			if !videoDirExists {
				rows = append(rows, []string{childDir, pterm.FgRed.Sprintf("False")})
				continue
			}

			matchFound := false

			// Traverse the "videos" directory to look for target folders
			videoEntries, err := os.ReadDir(videoDir)
			if err != nil {
				return err
			}
			for _, videoEntry := range videoEntries {
				if videoEntry.IsDir() && isTargetFolder(videoEntry.Name()) {
					matchFound = true
					break
				}
			}

			// only display which directories aren't organized
			if !matchFound {
				rows = append(rows, []string{childDir, pterm.FgRed.Sprintf("False")})
			}
		}
	}

	if len(rows) > 1 {
		// Render table with pterm
		pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(rows).Render()
	} else {
		pterm.Info.Printf("All sub-directories in '%s' are organized", rootDirectory)
	}
	return nil
}

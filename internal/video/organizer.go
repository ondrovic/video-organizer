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
    Name     string
    MaxDuration float64
}

// Define the folder structure once
var folderStructure = []FolderInfo{
    {Name: "Very Short", MaxDuration: 60},
    {Name: "Short", MaxDuration: 5 * 60},
    {Name: "Medium", MaxDuration: 15 * 60},
    {Name: "Long", MaxDuration: 30 * 60},
    {Name: "Very Long", MaxDuration: 60 * 60},
    {Name: "Super Long", MaxDuration: float64(^uint(0) >> 1)}, // Max float64 value
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

// func getTargetFolder(duration float64) string {
// 	switch {
// 	case duration < 60:
// 		return "Very Short"
// 	case duration < 5*60:
// 		return "Short"
// 	case duration < 15*60:
// 		return "Medium"
// 	case duration < 30*60:
// 		return "Long"
// 	case duration < 60*60:
// 		return "Very Long"
// 	default:
// 		return "Super Long"
// 	}
// }

// func OrganizeVideos(rootDirectory string) error {
// 	videoFiles, err := getVideoFiles(rootDirectory)
// 	if err != nil {
// 		return err
// 	}

// 	var totalFiles int
// 	for _, files := range videoFiles {
// 		totalFiles += len(files)
// 	}

// 	bar := progressbar.Default(int64(totalFiles))
// 	defer bar.Close()

// 	for subDir, files := range videoFiles {
// 		for _, filePath := range files {
// 			// Check if the file is already in a target folder
// 			currentFolder := filepath.Base(subDir)
// 			if isTargetFolder(currentFolder) {
// 				bar.Add(1)
// 				continue // Skip this file as it's already organized
// 			}

// 			duration, err := getDurationInSeconds(filePath)
// 			if err != nil {
// 				return err
// 			}

// 			targetFolder := getTargetFolder(duration)
// 			targetPath := filepath.Join(rootDirectory, targetFolder, filepath.Base(filePath))

// 			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
// 				return err
// 			}

// 			if err := os.Rename(filePath, targetPath); err != nil {
// 				return err
// 			}

// 			bar.Add(1)
// 		}
// 	}

// 	return nil
// }

// func isTargetFolder(folder string) bool {
// 	targetFolders := []string{"Very Short", "Short", "Medium", "Long", "Very Long", "Super Long"}
// 	for _, target := range targetFolders {
// 		if folder == target {
// 			return true
// 		}
// 	}
// 	return false
// }

// func OrganizeVideos(rootDirectory string) error {
// 	videoFiles, err := getVideoFiles(rootDirectory)
// 	if err != nil {
// 		return err
// 	}

// 	var totalFiles int
// 	for _, files := range videoFiles {
// 		totalFiles += len(files)
// 	}

// 	progressbar, _ := pterm.DefaultProgressbar.WithTotal(totalFiles).WithTitle("Organizing videos").Start()

// 	for subDir, files := range videoFiles {
// 		for _, filePath := range files {
// 			currentFolder := filepath.Base(subDir)
// 			if isTargetFolder(currentFolder) {
// 				progressbar.Increment()
// 				continue
// 			}

// 			duration, err := getDurationInSeconds(filePath)
// 			if err != nil {
// 				return err
// 			}

// 			targetFolder := getTargetFolder(duration)
// 			targetPath := filepath.Join(rootDirectory, targetFolder, filepath.Base(filePath))

// 			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
// 				return err
// 			}

// 			if err := os.Rename(filePath, targetPath); err != nil {
// 				return err
// 			}

// 			progressbar.Increment()
// 		}
// 	}

// 	progressbar.Stop()

// 	return nil
// }
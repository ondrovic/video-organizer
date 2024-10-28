// internal/video/file_utils.go
package video

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/pterm/pterm"

	commonTypes "github.com/ondrovic/common/types"
	cUtils "github.com/ondrovic/common/utils"
	cFormatter "github.com/ondrovic/common/utils/formatters"
)

func formatSkippedFilesMessage(skippedFiles int64) (string, error) {
	fileMsg, err := cFormatter.Pluralize(skippedFiles, "file", "files")
	if err != nil {
		return "", err
	}

	actionMsg, err := cFormatter.Pluralize(skippedFiles, "was", "were")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("\t\tskipped %v %s that %s already in target directories",
		skippedFiles,
		fileMsg,
		actionMsg), nil
}

func getVideoFiles(rootDirectory string) (map[string][]string, error) {
	videoFiles := make(map[string][]string)
	var mu sync.Mutex
	var wg sync.WaitGroup

	semaphore := make(chan struct{}, runtime.NumCPU())

	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Scanning %v for video files... ", rootDirectory))

	var totalFiles int64
	var skippedFiles int64

	err := walkDir(rootDirectory, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && cUtils.IsExtensionValid(commonTypes.FileTypes.Video, info.Name()) {
			relPath, err := filepath.Rel(rootDirectory, filepath.Dir(path))
			if err != nil {
				pterm.Error.Printf("Error getting relPath: %v\n", err)
				return nil
			}

			// Skip files that are already in target directories
			if isTargetFolder(filepath.Base(relPath)) {
				atomic.AddInt64(&skippedFiles, 1)
				return nil
			}

			atomic.AddInt64(&totalFiles, 1)
			wg.Add(1)
			go func(p string) {
				defer wg.Done()
				semaphore <- struct{}{}
				defer func() { <-semaphore }()

				mu.Lock()
				videoFiles[relPath] = append(videoFiles[relPath], p)
				mu.Unlock()
			}(path)
		}
		return nil
	})

	wg.Wait()

	spinner.Success(fmt.Sprintf("%v scan completed!", rootDirectory))

	if skippedFiles > 0 {
		skippedMsg, err := formatSkippedFilesMessage(skippedFiles)
		if err != nil {
			pterm.Error.Println(err)
			return nil, err
		}
		pterm.Info.Println(skippedMsg)
	}
	pterm.Success.Printf("Found %v video files in %v\n", totalFiles, rootDirectory)
	return videoFiles, err
}

func walkDir(dir string, fn func(path string, info os.DirEntry, err error) error) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			if err := walkDir(path, fn); err != nil {
				return err
			}
		} else {
			if err := fn(path, entry, nil); err != nil {
				return err
			}
		}
	}

	return nil
}

func FormatDirectory(directory string) string {
	// Expand '~' to home directory for Linux/macOS
	if runtime.GOOS != "windows" && directory[:1] == "~" {
		if home, err := os.UserHomeDir(); err == nil {
			directory = filepath.Join(home, directory[1:])
		}
	}
	// Ensure path is clean and uses the appropriate separators
	return filepath.Clean(directory)
}

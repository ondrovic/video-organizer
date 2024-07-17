// internal/video/file_utils.go
package video

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/pterm/pterm"
)

var validExtensions = []string{".mp4", ".avi", ".mov", ".mkv", ".ts", "mov", ".m4v", ".wmv", ".flv", ".webm", ".mpg", ".mpeg"}

func isValidVideoFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, validExt := range validExtensions {
		if ext == validExt {
			return true
		}
	}
	return false
}

func getVideoFiles(rootDirectory string) (map[string][]string, error) {
	videoFiles := make(map[string][]string)
	var mu sync.Mutex
	var wg sync.WaitGroup

	semaphore := make(chan struct{}, runtime.NumCPU())

	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Scanning %v for video files... ", rootDirectory))

	var totalFiles int64

	err := filepath.Walk(rootDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && isValidVideoFile(info.Name()) {
			totalFiles++
			wg.Add(1)
			go func(p string) {
				defer wg.Done()
				semaphore <- struct{}{}
				defer func() { <-semaphore }()

				relPath, err := filepath.Rel(rootDirectory, filepath.Dir(p))
				if err != nil {
					pterm.Error.Printf("Error getting relPath: %v\n", err)
					return
				}

				mu.Lock()
				videoFiles[relPath] = append(videoFiles[relPath], p)
				mu.Unlock()
			}(path)
		}
		return nil
	})

	wg.Wait()

	spinner.Success(fmt.Sprintf("%v scan completed!", rootDirectory))

	pterm.Success.Printf("Found %v video files in %v!\n", totalFiles, rootDirectory)

	return videoFiles, err
}

// internal/video/file_utils.go
package video

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"

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

func pluralize(count int64, singular string, plural string) string {
	if count == 1 {
		return singular
	}
	return plural
}

func formatSkippedFilesMessage(skippedFiles int64) string {
	return fmt.Sprintf("\t\tskipped %v %s that %s already in target directories",
		skippedFiles,
		pluralize(skippedFiles, "file", "files"),
		pluralize(skippedFiles, "was", "were"))
}

func getVideoFiles(rootDirectory string) (map[string][]string, error) {
	videoFiles := make(map[string][]string)
	var mu sync.Mutex
	var wg sync.WaitGroup

	semaphore := make(chan struct{}, runtime.NumCPU())

	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Scanning %v for video files... ", rootDirectory))

	var totalFiles int64
	var skippedFiles int64

	err := filepath.Walk(rootDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && isValidVideoFile(info.Name()) {
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
		pterm.Info.Println(formatSkippedFilesMessage(skippedFiles))
	}

	pterm.Success.Printf("Found %v video files in %v\n", totalFiles, rootDirectory)
	return videoFiles, err
}

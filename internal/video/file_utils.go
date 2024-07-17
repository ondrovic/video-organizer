// internal/video/file_utils.go
package video

import (
	// "log"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	// "sync/atomic"
	// "time"

	"github.com/pterm/pterm"
)

var validExtensions = []string{".mp4", ".avi", ".mov", ".mkv", ".ts", "mov", ".m4v"}

func isValidVideoFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, validExt := range validExtensions {
		if ext == validExt {
			return true
		}
	}
	return false
}

// func getVideoFiles(rootDirectory string) (map[string][]string, error) {
//     videoFiles := make(map[string][]string)
//     err := filepath.Walk(rootDirectory, func(path string, info os.FileInfo, err error) error {
//         if err != nil {
//             return err
//         }
//         if !info.IsDir() && isValidVideoFile(info.Name()) {
//             relPath, err := filepath.Rel(rootDirectory, filepath.Dir(path))
//             if err != nil {
//                 return err
//             }
//             videoFiles[relPath] = append(videoFiles[relPath], path)
//         }
//         return nil
//     })
//     return videoFiles, err
// }

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

// func getVideoFiles(rootDirectory string) (map[string][]string, error) {
// 	videoFiles := make(map[string][]string)
// 	var mu sync.Mutex
// 	var wg sync.WaitGroup

// 	semaphore := make(chan struct{}, runtime.NumCPU())

// 	var totalFiles int64
// 	var processedFiles int64

// 	// Create a new progress bar
// 	progressbar, _ := pterm.DefaultProgressbar.
// 		WithTotal(0).  // We don't know the total yet, so set it to 0
// 		WithTitle("Scanning for video files").
// 		WithRemoveWhenDone(true).
// 		Start()

// 	// Use a ticker to update the progress bar every 100ms
// 	ticker := time.NewTicker(100 * time.Millisecond)
// 	defer ticker.Stop()

// 	go func() {
// 		for range ticker.C {
// 			current := atomic.LoadInt64(&processedFiles)
// 			total := atomic.LoadInt64(&totalFiles)
// 			progressbar.Total = int(current) // Update total processed files
// 			progressbar.Current = int(total) // Update current video files found
// 			_ = progressbar.Add(0) // This triggers an update of the display
// 		}
// 	}()

// 	err := filepath.Walk(rootDirectory, func(path string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			return err
// 		}

// 		atomic.AddInt64(&processedFiles, 1)

// 		if !info.IsDir() && isValidVideoFile(info.Name()) {
// 			atomic.AddInt64(&totalFiles, 1)
// 			wg.Add(1)
// 			go func(p string) {
// 				defer wg.Done()
// 				semaphore <- struct{}{}
// 				defer func() { <-semaphore }()

// 				relPath, err := filepath.Rel(rootDirectory, filepath.Dir(p))
// 				if err != nil {
// 					pterm.Error.Printf("Error getting relPath: %v\n", err)
// 					return
// 				}

// 				mu.Lock()
// 				videoFiles[relPath] = append(videoFiles[relPath], p)
// 				mu.Unlock()
// 			}(path)
// 		}
// 		return nil
// 	})

// 	wg.Wait()
// 	ticker.Stop()

// 	progressbar.Stop()

// 	// pterm.Success.Printf("Video file scanning completed! Found %d video files\n", totalFiles)

// 	return videoFiles, err
// }
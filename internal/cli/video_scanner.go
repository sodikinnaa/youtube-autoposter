package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type FoundVideo struct {
	Path     string
	RelPath  string
	SizeBytes int64
}

// ScanVideoFiles recursively scans rootDir and subdirectories for video files
func ScanVideoFiles(rootDir string) ([]FoundVideo, error) {
	videoExts := map[string]bool{
		".mp4":  true,
		".mkv":  true,
		".mov":  true,
		".avi":  true,
		".webm": true,
		".flv":  true,
		".wmv":  true,
		".m4v":  true,
	}

	var found []FoundVideo
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip unreadable directories
		}

		// Skip hidden folders like .git
		if info.IsDir() && strings.HasPrefix(info.Name(), ".") && info.Name() != "." {
			return filepath.SkipDir
		}

		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			if videoExts[ext] {
				relPath, rErr := filepath.Rel(rootDir, path)
				if rErr != nil {
					relPath = path
				}
				if !strings.HasPrefix(relPath, "./") {
					relPath = "./" + relPath
				}
				found = append(found, FoundVideo{
					Path:      path,
					RelPath:   relPath,
					SizeBytes: info.Size(),
				})
			}
		}
		return nil
	})

	return found, err
}

func FormatFileSize(bytes int64) string {
	mb := float64(bytes) / (1024 * 1024)
	if mb >= 1000 {
		gb := mb / 1024
		return fmt.Sprintf("%.2f GB", gb)
	}
	return fmt.Sprintf("%.2f MB", mb)
}

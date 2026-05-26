package quantitative

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// ScanResult holds quantitative metrics for a directory.
type ScanResult struct {
	TotalLOC      int
	TotalFiles    int
	ByExtension   map[string]int
}

var codeExtensions = map[string]bool{
	".go": true, ".py": true, ".java": true, ".ts": true, ".js": true,
	".cpp": true, ".rs": true, ".kt": true, ".swift": true, ".rb": true,
}

// Scan walks a directory and counts lines of non-blank, non-comment code.
func Scan(root string) ScanResult {
	result := ScanResult{ByExtension: make(map[string]int)}
	filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if !codeExtensions[ext] {
			return nil
		}
		loc := countLOC(path)
		result.TotalLOC += loc
		result.TotalFiles++
		result.ByExtension[ext] += loc
		return nil
	})
	return result
}

func countLOC(path string) int {
	f, err := os.Open(path)
	if err != nil {
		return 0
	}
	defer f.Close()
	n := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "#") {
			continue
		}
		n++
	}
	return n
}

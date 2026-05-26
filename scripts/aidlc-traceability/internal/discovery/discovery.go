package discovery

import (
	"os"
	"path/filepath"
	"strings"
)

var sourceExtensions = map[string]bool{
	".go": true, ".py": true, ".java": true, ".ts": true,
	".cpp": true, ".rs": true, ".kt": true, ".swift": true,
}

// FindAidlcDocs looks for an aidlc-docs/ directory directly under projectRoot.
// Returns the full path or empty string if absent.
func FindAidlcDocs(projectRoot string) string {
	candidate := filepath.Join(projectRoot, "aidlc-docs")
	if info, err := os.Stat(candidate); err == nil && info.IsDir() {
		return candidate
	}
	return ""
}

// DiscoverSourceCode returns all non-test source files under root.
func DiscoverSourceCode(root string) []string {
	var files []string
	filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if !sourceExtensions[ext] {
			return nil
		}
		base := filepath.Base(path)
		if strings.HasSuffix(base, "_test.go") || strings.HasSuffix(base, "_test.py") {
			return nil
		}
		files = append(files, path)
		return nil
	})
	return files
}

package parsers

import (
	"bufio"
	"os"
	"path/filepath"

	"github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/models"
)

// ParseCodeFile returns a CODE artifact for the given source file.
// The artifact ID is the file path relative to projectRoot.
func ParseCodeFile(filePath, projectRoot string) *models.Artifact {
	rel, err := filepath.Rel(projectRoot, filePath)
	if err != nil {
		rel = filePath
	}
	loc := countLines(filePath)
	return &models.Artifact{
		ID:         rel,
		Title:      filepath.Base(filePath),
		Type:       models.ArtifactTypeCode,
		SourceFile: filePath,
		Metadata:   map[string]any{"loc": loc},
	}
}

func countLines(path string) int {
	f, err := os.Open(path)
	if err != nil {
		return 0
	}
	defer f.Close()
	n := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		n++
	}
	return n
}

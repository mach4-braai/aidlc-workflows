package parsers

import (
	"regexp"

	"github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/models"
)

var unitHeading = regexp.MustCompile(`^##\s+(UNIT-\d+|[A-Za-z][A-Za-z0-9]*(?:---[A-Za-z][A-Za-z0-9]*)*):\s+(.+)$`)

// ParseUnits extracts UNIT artifacts from a markdown file.
func ParseUnits(filePath string) []models.Artifact {
	return parseHeadings(filePath, unitHeading, models.ArtifactTypeUnit)
}

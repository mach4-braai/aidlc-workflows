package parsers

import (
	"regexp"

	"github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/models"
)

var componentHeading = regexp.MustCompile(`^##\s+([A-Z][a-zA-Z]+(?:\s+[A-Z][a-zA-Z]+)*)$`)

// ParseComponents extracts COMPONENT artifacts from an application-components.md file.
func ParseComponents(filePath string) []models.Artifact {
	return parseHeadings(filePath, componentHeading, models.ArtifactTypeComponent)
}

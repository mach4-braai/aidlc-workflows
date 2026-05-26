package parsers

import (
	"regexp"

	"github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/models"
)

var storyHeading = regexp.MustCompile(`^##\s+(US-\d+(?:\.\d+)?):\s+(.+)$`)

// ParseStories extracts STORY artifacts from a markdown file.
func ParseStories(filePath string) []models.Artifact {
	return parseHeadings(filePath, storyHeading, models.ArtifactTypeStory)
}

package parsers

import (
	"regexp"

	"github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/models"
)

var codePlanHeading = regexp.MustCompile(`^##\s+(CP-\d+|[A-Za-z][A-Za-z0-9-]+):\s+(.+)$`)

// ParseCodePlans extracts CODE_PLAN artifacts from a markdown file.
func ParseCodePlans(filePath string) []models.Artifact {
	return parseHeadings(filePath, codePlanHeading, models.ArtifactTypeCodePlan)
}

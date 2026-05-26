package parsers

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/models"
)

var reqHeading = regexp.MustCompile(`^##\s+((FR|NFR|AR)-\d+):\s+(.+)$`)

// ParseRequirements extracts REQUIREMENT artifacts from a markdown file.
func ParseRequirements(filePath string) []models.Artifact {
	return parseHeadings(filePath, reqHeading, models.ArtifactTypeRequirement)
}

func parseHeadings(filePath string, re *regexp.Regexp, artifactType models.ArtifactType) []models.Artifact {
	f, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer f.Close()

	var artifacts []models.Artifact
	scanner := bufio.NewScanner(f)
	lineNum := 0
	var current *models.Artifact
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if m := re.FindStringSubmatch(line); m != nil {
			if current != nil {
				artifacts = append(artifacts, *current)
			}
			current = &models.Artifact{
				ID:         m[1],
				Title:      strings.TrimSpace(m[len(m)-1]),
				Type:       artifactType,
				SourceFile: filePath,
				SourceLine: lineNum,
			}
		} else if current != nil && strings.TrimSpace(line) != "" && !strings.HasPrefix(line, "#") {
			if current.Desc == "" {
				current.Desc = strings.TrimSpace(line)
			}
		}
	}
	if current != nil {
		artifacts = append(artifacts, *current)
	}
	return artifacts
}

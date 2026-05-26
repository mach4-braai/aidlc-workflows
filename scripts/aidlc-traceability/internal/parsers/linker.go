package parsers

import (
	"regexp"
	"strings"

	"github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/models"
)

// InferRequirementStoryLinks finds explicit requirement ID mentions in story
// Title or Desc fields and creates traces_to relationships.
func InferRequirementStoryLinks(reqs, stories []models.Artifact) []models.Relationship {
	return inferLinks(reqs, stories, "traces_to")
}

// InferStoryUnitLinks finds explicit story ID mentions in unit artifacts.
func InferStoryUnitLinks(stories, units []models.Artifact) []models.Relationship {
	return inferLinks(stories, units, "implemented_by")
}

// InferUnitComponentLinks finds explicit unit ID mentions in component artifacts.
func InferUnitComponentLinks(units, components []models.Artifact) []models.Relationship {
	return inferLinks(units, components, "realized_by")
}

// InferComponentCodeLinks finds explicit component ID mentions in code artifacts.
func InferComponentCodeLinks(components, codeFiles []models.Artifact) []models.Relationship {
	return inferLinks(components, codeFiles, "implemented_in")
}

func inferLinks(sources, targets []models.Artifact, relType string) []models.Relationship {
	var rels []models.Relationship
	for _, src := range sources {
		pat := regexp.MustCompile(`\b` + regexp.QuoteMeta(src.ID) + `\b`)
		for _, tgt := range targets {
			haystack := tgt.Title + " " + tgt.Desc
			if pat.MatchString(haystack) || strings.Contains(haystack, src.ID) {
				rels = append(rels, models.Relationship{
					SourceID:         src.ID,
					TargetID:         tgt.ID,
					RelationshipType: relType,
				})
			}
		}
	}
	return rels
}

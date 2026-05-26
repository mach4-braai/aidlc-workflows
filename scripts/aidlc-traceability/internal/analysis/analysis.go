package analysis

import (
	"github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/graph"
	"github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/models"
)

// DetectGaps returns coverage gaps: artifacts that are missing their expected
// downstream link (e.g. a requirement with no story, a story with no unit).
func DetectGaps(artifacts []models.Artifact, g *graph.TraceGraph) []models.CoverageGap {
	var gaps []models.CoverageGap
	for _, a := range artifacts {
		switch a.Type {
		case models.ArtifactTypeRequirement:
			if !graph.HasSuccessor(g, a.ID) {
				gaps = append(gaps, models.CoverageGap{
					ArtifactID:    a.ID,
					ArtifactTitle: a.Title,
					ArtifactType:  a.Type,
					GapType:       "missing_story",
					Description:   "Requirement has no linked user story",
				})
			}
		case models.ArtifactTypeStory:
			if !graph.HasSuccessor(g, a.ID) {
				gaps = append(gaps, models.CoverageGap{
					ArtifactID:    a.ID,
					ArtifactTitle: a.Title,
					ArtifactType:  a.Type,
					GapType:       "missing_unit",
					Description:   "Story has no linked implementation unit",
				})
			}
		case models.ArtifactTypeUnit:
			if !graph.HasSuccessor(g, a.ID) {
				gaps = append(gaps, models.CoverageGap{
					ArtifactID:    a.ID,
					ArtifactTitle: a.Title,
					ArtifactType:  a.Type,
					GapType:       "missing_code",
					Description:   "Unit has no linked code artifact",
				})
			}
		case models.ArtifactTypeCode:
			if !graph.HasSuccessor(g, a.ID) {
				gaps = append(gaps, models.CoverageGap{
					ArtifactID:    a.ID,
					ArtifactTitle: a.Title,
					ArtifactType:  a.Type,
					GapType:       "missing_test",
					Description:   "Code file has no linked test",
				})
			}
		}
	}
	return gaps
}

// CalculateMetrics counts artifacts by type and coverage relationships.
func CalculateMetrics(artifacts []models.Artifact, rels []models.Relationship) models.CoverageMetrics {
	m := models.CoverageMetrics{}

	// Count by type.
	for _, a := range artifacts {
		switch a.Type {
		case models.ArtifactTypeRequirement:
			m.TotalRequirements++
		case models.ArtifactTypeStory:
			m.TotalStories++
		case models.ArtifactTypeUnit:
			m.TotalUnits++
		case models.ArtifactTypeCode:
			m.TotalCodeFiles++
		case models.ArtifactTypeTest:
			m.TotalTests++
		}
	}

	// Build successor sets per source ID.
	storyTargets := targetSet(rels, models.ArtifactTypeRequirement, models.ArtifactTypeStory, artifacts)
	unitTargets := targetSet(rels, models.ArtifactTypeStory, models.ArtifactTypeUnit, artifacts)
	codeTargets := targetSet(rels, models.ArtifactTypeUnit, models.ArtifactTypeCode, artifacts)
	testTargets := targetSet(rels, models.ArtifactTypeCode, models.ArtifactTypeTest, artifacts)

	m.RequirementsWithStories = len(storyTargets)
	m.StoriesWithUnits = len(unitTargets)
	m.UnitsWithCode = len(codeTargets)
	m.CodeWithTests = len(testTargets)

	return m
}

// targetSet returns the set of source IDs that have at least one relationship
// to a target of the given artifact type.
func targetSet(rels []models.Relationship, srcType, tgtType models.ArtifactType, artifacts []models.Artifact) map[string]struct{} {
	typeOf := make(map[string]models.ArtifactType, len(artifacts))
	for _, a := range artifacts {
		typeOf[a.ID] = a.Type
	}
	set := make(map[string]struct{})
	for _, r := range rels {
		if typeOf[r.SourceID] == srcType && typeOf[r.TargetID] == tgtType {
			set[r.SourceID] = struct{}{}
		}
	}
	return set
}

package analysis_test

import (
	"testing"

	"github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/analysis"
	"github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/graph"
	"github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/models"
)

func TestDetectGapsFindsOrphanedRequirement(t *testing.T) {
	artifacts := []models.Artifact{
		{ID: "FR-001", Type: models.ArtifactTypeRequirement},
	}
	g := graph.Build(artifacts, nil)
	gaps := analysis.DetectGaps(artifacts, g)
	if len(gaps) == 0 {
		t.Fatal("expected gap for orphaned requirement")
	}
}

func TestCalculateMetricsCountsCorrectly(t *testing.T) {
	artifacts := []models.Artifact{
		{ID: "FR-001", Type: models.ArtifactTypeRequirement},
		{ID: "US-1.1", Type: models.ArtifactTypeStory},
	}
	rels := []models.Relationship{
		{SourceID: "FR-001", TargetID: "US-1.1"},
	}
	m := analysis.CalculateMetrics(artifacts, rels)
	if m.TotalRequirements != 1 {
		t.Fatalf("expected 1 requirement, got %d", m.TotalRequirements)
	}
	if m.RequirementsWithStories != 1 {
		t.Fatalf("expected 1 requirement with story, got %d", m.RequirementsWithStories)
	}
}

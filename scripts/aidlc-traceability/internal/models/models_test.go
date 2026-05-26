package models_test

import (
	"testing"

	"github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/models"
)

func TestArtifactIDIsNonEmpty(t *testing.T) {
	a := models.Artifact{ID: "FR-001", Title: "Login", Type: models.ArtifactTypeRequirement}
	if a.ID == "" {
		t.Fatal("artifact ID must not be empty")
	}
}

func TestCoverageMetricsZeroValue(t *testing.T) {
	m := models.CoverageMetrics{}
	if m.TotalRequirements != 0 {
		t.Fatal("zero value should be 0")
	}
}

func TestTraceabilityReportHoldsArtifacts(t *testing.T) {
	r := models.TraceabilityReport{
		Artifacts: []models.Artifact{
			{ID: "FR-001", Title: "Login", Type: models.ArtifactTypeRequirement},
		},
	}
	if len(r.Artifacts) != 1 {
		t.Fatalf("expected 1 artifact, got %d", len(r.Artifacts))
	}
}

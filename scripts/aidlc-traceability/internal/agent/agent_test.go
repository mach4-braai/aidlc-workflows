package agent_test

import (
	"context"
	"testing"

	"github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/agent"
	"github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/models"
)

func TestParseAgentJSONExtractsRelationships(t *testing.T) {
	validIDs := map[string]bool{"FR-001": true, "US-1.1": true}
	responseText := `{"relationships": [{"source_id": "FR-001", "target_id": "US-1.1", "relationship_type": "traces_to"}], "insights": "good match"}`
	rels, insights := agent.ParseAgentJSON(responseText, validIDs)
	if len(rels) != 1 {
		t.Fatalf("expected 1 relationship, got %d", len(rels))
	}
	if len(insights) == 0 {
		t.Fatal("expected insights")
	}
}

func TestParseAgentJSONSkipsInvalidIDs(t *testing.T) {
	validIDs := map[string]bool{"FR-001": true}
	responseText := `{"relationships": [{"source_id": "INVALID", "target_id": "US-1.1", "relationship_type": "traces_to"}]}`
	rels, _ := agent.ParseAgentJSON(responseText, validIDs)
	if len(rels) != 0 {
		t.Fatal("expected 0 relationships for invalid IDs")
	}
}

func TestParseAgentJSONHandlesMalformedJSON(t *testing.T) {
	validIDs := map[string]bool{"FR-001": true}
	rels, _ := agent.ParseAgentJSON("not json", validIDs)
	if len(rels) != 0 {
		t.Fatal("expected 0 relationships for malformed JSON")
	}
}

func TestRunReqStoryAnalysis_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	ctx := context.Background()
	reqs := []models.Artifact{{ID: "FR-001", Title: "User Login", Type: models.ArtifactTypeRequirement}}
	stories := []models.Artifact{{ID: "US-1.1", Title: "Implements FR-001 login form", Type: models.ArtifactTypeStory}}
	rels, _ := agent.RunReqStoryAnalysis(ctx, reqs, stories, "", "us-east-1")
	_ = rels
}

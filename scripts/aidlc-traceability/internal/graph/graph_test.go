package graph_test

import (
	"testing"

	"github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/graph"
	"github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/models"
)

func TestBuildGraphHasNodes(t *testing.T) {
	artifacts := []models.Artifact{
		{ID: "FR-001", Type: models.ArtifactTypeRequirement},
		{ID: "US-1.1", Type: models.ArtifactTypeStory},
	}
	rels := []models.Relationship{
		{SourceID: "FR-001", TargetID: "US-1.1", RelationshipType: "traces_to"},
	}
	g := graph.Build(artifacts, rels)
	if graph.NodeCount(g) != 2 {
		t.Fatalf("expected 2 nodes, got %d", graph.NodeCount(g))
	}
	if graph.EdgeCount(g) != 1 {
		t.Fatalf("expected 1 edge, got %d", graph.EdgeCount(g))
	}
}

func TestBuildGraphHandlesOrphanedArtifacts(t *testing.T) {
	artifacts := []models.Artifact{
		{ID: "FR-001", Type: models.ArtifactTypeRequirement},
	}
	g := graph.Build(artifacts, nil)
	if graph.NodeCount(g) != 1 {
		t.Fatalf("expected 1 node, got %d", graph.NodeCount(g))
	}
}

func TestHasSuccessorReturnsTrueWhenEdgeExists(t *testing.T) {
	artifacts := []models.Artifact{
		{ID: "FR-001", Type: models.ArtifactTypeRequirement},
		{ID: "US-1.1", Type: models.ArtifactTypeStory},
	}
	rels := []models.Relationship{
		{SourceID: "FR-001", TargetID: "US-1.1"},
	}
	g := graph.Build(artifacts, rels)
	if !graph.HasSuccessor(g, "FR-001") {
		t.Fatal("FR-001 should have a successor")
	}
	if graph.HasSuccessor(g, "US-1.1") {
		t.Fatal("US-1.1 should have no successor")
	}
}

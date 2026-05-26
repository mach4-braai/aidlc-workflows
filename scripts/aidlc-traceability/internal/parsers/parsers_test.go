package parsers_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/models"
	"github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/parsers"
)

func TestParseRequirementsExtractsIDs(t *testing.T) {
	content := "# Requirements\n\n## FR-001: User Login\n\nUsers must be able to log in.\n\n## FR-002: Dashboard\n\nShow dashboard.\n"
	tmp := t.TempDir()
	f := filepath.Join(tmp, "requirements.md")
	os.WriteFile(f, []byte(content), 0644)
	artifacts := parsers.ParseRequirements(f)
	if len(artifacts) != 2 {
		t.Fatalf("expected 2 requirements, got %d", len(artifacts))
	}
	if artifacts[0].ID != "FR-001" {
		t.Fatalf("expected FR-001, got %s", artifacts[0].ID)
	}
}

func TestParseStoriesExtractsIDs(t *testing.T) {
	content := "# Stories\n\n## US-1.1: Login Form\n\nAs a user I want to log in.\n"
	tmp := t.TempDir()
	f := filepath.Join(tmp, "stories.md")
	os.WriteFile(f, []byte(content), 0644)
	artifacts := parsers.ParseStories(f)
	if len(artifacts) == 0 {
		t.Fatal("expected at least one story")
	}
	if artifacts[0].Type != models.ArtifactTypeStory {
		t.Fatalf("expected STORY type, got %s", artifacts[0].Type)
	}
}

func TestParseCodeFileReturnsArtifact(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "main.go")
	os.WriteFile(f, []byte("package main\n\nfunc main() {}\n"), 0644)
	a := parsers.ParseCodeFile(f, tmp)
	if a == nil {
		t.Fatal("expected artifact, got nil")
	}
	if a.Type != models.ArtifactTypeCode {
		t.Fatalf("expected CODE, got %s", a.Type)
	}
}

func TestInferLinksFindsExplicitReferences(t *testing.T) {
	reqs := []models.Artifact{{ID: "FR-001", Type: models.ArtifactTypeRequirement}}
	stories := []models.Artifact{{ID: "US-1.1", Title: "Implements FR-001", Type: models.ArtifactTypeStory}}
	links := parsers.InferRequirementStoryLinks(reqs, stories)
	if len(links) == 0 {
		t.Fatal("expected at least one inferred link")
	}
}

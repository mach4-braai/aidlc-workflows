package pipeline_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/pipeline"
)

func TestRunPipelineOnEmptyProject(t *testing.T) {
	tmp := t.TempDir()
	report, err := pipeline.Run(pipeline.Config{
		ProjectRoot: tmp,
		UseAI:       false,
		OutputDir:   tmp,
		Format:      "markdown",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.ProjectName == "" {
		t.Fatal("report must have a project name")
	}
}

func TestRunPipelineWritesOutputFile(t *testing.T) {
	tmp := t.TempDir()
	docsDir := filepath.Join(tmp, "aidlc-docs")
	os.MkdirAll(docsDir, 0755)
	os.WriteFile(filepath.Join(docsDir, "requirements.md"), []byte("# Requirements\n\n## FR-001: Login\n\nUsers must log in.\n"), 0644)

	_, err := pipeline.Run(pipeline.Config{
		ProjectRoot: tmp,
		UseAI:       false,
		OutputDir:   tmp,
		Format:      "markdown",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries, _ := os.ReadDir(tmp)
	hasMarkdown := false
	for _, e := range entries {
		if filepath.Ext(e.Name()) == ".md" {
			hasMarkdown = true
		}
	}
	if !hasMarkdown {
		t.Fatal("expected markdown output file")
	}
}

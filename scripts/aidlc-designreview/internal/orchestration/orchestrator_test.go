package orchestration_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mach4-braai/aidlc-workflows/aidlc-designreview/internal/orchestration"
)

func TestReviewOrchestratorFailsOnMissingDir(t *testing.T) {
	orch := orchestration.NewReviewOrchestrator(orchestration.Config{})
	_, err := orch.Run("/nonexistent/path")
	if err == nil {
		t.Fatal("expected error for missing directory")
	}
}

func TestReviewOrchestratorRunsOnMinimalProject(t *testing.T) {
	tmp := t.TempDir()
	docsDir := filepath.Join(tmp, "aidlc-docs", "construction")
	os.MkdirAll(docsDir, 0755)
	os.WriteFile(filepath.Join(docsDir, "application-design.md"),
		[]byte("# Application Design\n\n## Architecture\n\nMonolith.\n"), 0644)

	orch := orchestration.NewReviewOrchestrator(orchestration.Config{
		MockAI: true,
	})
	result, err := orch.Run(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

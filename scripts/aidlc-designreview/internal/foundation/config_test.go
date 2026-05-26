package foundation_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mach4-braai/aidlc-workflows/aidlc-designreview/internal/foundation"
)

func TestLoadConfigFromYAML(t *testing.T) {
	yaml := `
models:
  critique: us.anthropic.claude-sonnet-4-20250514-v1:0
  alternatives: us.anthropic.claude-sonnet-4-20250514-v1:0
  gap: us.anthropic.claude-sonnet-4-20250514-v1:0
aws:
  region: us-east-1
`
	tmp := t.TempDir()
	f := filepath.Join(tmp, "config.yaml")
	os.WriteFile(f, []byte(yaml), 0644)
	cfg, err := foundation.LoadConfig(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.AWS.Region != "us-east-1" {
		t.Fatalf("expected us-east-1, got %s", cfg.AWS.Region)
	}
}

func TestLoadConfigUsesDefaults(t *testing.T) {
	cfg, err := foundation.LoadDefaultConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.AWS.Region == "" {
		t.Fatal("default region must be set")
	}
}

func TestPatternLibraryLoadsPatterns(t *testing.T) {
	lib, err := foundation.LoadPatternLibrary()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lib.Patterns) == 0 {
		t.Fatal("expected at least one design pattern")
	}
}

func TestPromptManagerBuildsPrompt(t *testing.T) {
	pm, err := foundation.LoadPromptManager()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	prompt, err := pm.BuildAgentPrompt("critique", map[string]string{"design_content": "test content"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if prompt == "" {
		t.Fatal("prompt must not be empty")
	}
}

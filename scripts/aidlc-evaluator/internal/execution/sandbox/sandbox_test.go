package sandbox_test

import (
	"testing"

	"github.com/mach4-braai/aidlc-workflows/aidlc-evaluator/internal/execution/sandbox"
)

func TestSandboxConfigDefaultImage(t *testing.T) {
	cfg := sandbox.DefaultConfig()
	if cfg.Image == "" {
		t.Fatal("default config must specify an image")
	}
}

// TestBuildAndRunSandbox_Integration requires Docker daemon.
func TestBuildAndRunSandbox_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping docker integration test")
	}
	cfg := sandbox.DefaultConfig()
	err := sandbox.Build(cfg)
	if err != nil {
		t.Fatalf("build failed: %v", err)
	}
	result, err := sandbox.Run(cfg, "echo hello")
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if result.ExitCode != 0 {
		t.Fatalf("expected exit 0, got %d: %s", result.ExitCode, result.Stderr)
	}
}

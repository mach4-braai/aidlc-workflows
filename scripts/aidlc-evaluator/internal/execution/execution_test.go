package execution_test

import (
	"os/exec"
	"testing"

	"github.com/mach4-braai/aidlc-workflows/aidlc-evaluator/internal/execution"
)

func TestRunCommandExecutesBashEcho(t *testing.T) {
	if _, err := exec.LookPath("echo"); err != nil {
		t.Skip("echo not available")
	}
	result := execution.RunCommand("echo hello")
	if result.ExitCode != 0 {
		t.Fatalf("expected exit 0, got %d", result.ExitCode)
	}
	if result.Stdout != "hello\n" {
		t.Fatalf("expected 'hello\\n', got %q", result.Stdout)
	}
}

func TestRunnerConfigValidation(t *testing.T) {
	cfg := execution.RunnerConfig{}
	if err := cfg.Validate(); err == nil {
		t.Fatal("empty config should fail validation")
	}
}

func TestMetricsCollectorTracksTokens(t *testing.T) {
	m := execution.NewMetricsCollector()
	m.AddTokens("executor", 100, 50)
	stats := m.Summary()
	if stats.TotalInputTokens != 100 {
		t.Fatalf("expected 100 input tokens, got %d", stats.TotalInputTokens)
	}
}

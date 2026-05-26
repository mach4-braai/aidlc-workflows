package execution

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// CommandResult holds the output of a shell command execution.
type CommandResult struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Duration time.Duration
}

// RunCommand executes a shell command string and returns its result.
func RunCommand(cmd string) CommandResult {
	start := time.Now()
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return CommandResult{ExitCode: 1, Stderr: "empty command"}
	}
	c := exec.Command(parts[0], parts[1:]...)
	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr
	err := c.Run()
	exitCode := 0
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}
	return CommandResult{
		ExitCode: exitCode,
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		Duration: time.Since(start),
	}
}

// RunnerConfig holds the configuration for a single evaluation run.
type RunnerConfig struct {
	VisionPath     string
	TechEnvPath    string
	ExecutorModel  string
	SimulatorModel string
	OutputDir      string
	AWSProfile     string
	AWSRegion      string
}

// Validate returns an error if the config is incomplete.
func (c RunnerConfig) Validate() error {
	if c.VisionPath == "" && c.OutputDir == "" {
		return errors.New("at least VisionPath or OutputDir must be set")
	}
	return nil
}

// MetricsCollector accumulates token usage across agents.
type MetricsCollector struct {
	mu      sync.Mutex
	entries map[string]tokenEntry
}

type tokenEntry struct {
	InputTokens  int
	OutputTokens int
}

// MetricsSummary is the aggregated token usage across all agents.
type MetricsSummary struct {
	TotalInputTokens  int
	TotalOutputTokens int
	ByAgent           map[string]tokenEntry
}

// NewMetricsCollector creates an empty MetricsCollector.
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{entries: make(map[string]tokenEntry)}
}

// AddTokens records token usage for the named agent.
func (m *MetricsCollector) AddTokens(agent string, input, output int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	e := m.entries[agent]
	e.InputTokens += input
	e.OutputTokens += output
	m.entries[agent] = e
}

// Summary returns the aggregated token usage.
func (m *MetricsCollector) Summary() MetricsSummary {
	m.mu.Lock()
	defer m.mu.Unlock()
	s := MetricsSummary{ByAgent: make(map[string]tokenEntry)}
	for k, v := range m.entries {
		s.TotalInputTokens += v.InputTokens
		s.TotalOutputTokens += v.OutputTokens
		s.ByAgent[k] = v
	}
	return s
}

package cliharness

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// Adapter runs a CLI tool with a prompt and returns its output.
type Adapter interface {
	Name() string
	Run(scenario, prompt string) (string, error)
}

var ansiRe = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

// Normalize strips ANSI escape codes from text.
func Normalize(text string) string {
	return ansiRe.ReplaceAllString(text, "")
}

// BaseAdapter provides shared Run logic for CLI-based adapters.
type BaseAdapter struct {
	name    string
	binary  string
	buildArgs func(scenario, prompt string) []string
}

func (a *BaseAdapter) Name() string { return a.name }

func (a *BaseAdapter) Run(scenario, prompt string) (string, error) {
	args := a.buildArgs(scenario, prompt)
	cmd := exec.Command(a.binary, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%s: %w", a.name, err)
	}
	return Normalize(out.String()), nil
}

// ClaudeCodeAdapter wraps the claude-code CLI.
type ClaudeCodeAdapter struct{ BaseAdapter }

func newClaudeCodeAdapter() Adapter {
	return &ClaudeCodeAdapter{BaseAdapter{
		name:   "claude-code",
		binary: "claude",
		buildArgs: func(scenario, prompt string) []string {
			return []string{"--print", prompt}
		},
	}}
}

// KiroCLIAdapter wraps the kiro CLI.
type KiroCLIAdapter struct{ BaseAdapter }

func newKiroCLIAdapter() Adapter {
	return &KiroCLIAdapter{BaseAdapter{
		name:   "kiro-cli",
		binary: "kiro",
		buildArgs: func(scenario, prompt string) []string {
			return []string{"run", "--prompt", prompt}
		},
	}}
}

// Registry maps adapter names to factory functions.
type Registry struct {
	factories map[string]func() Adapter
}

// NewRegistry creates a registry with all known CLI adapters.
func NewRegistry() *Registry {
	r := &Registry{factories: make(map[string]func() Adapter)}
	r.factories["claude-code"] = newClaudeCodeAdapter
	r.factories["kiro-cli"] = newKiroCLIAdapter
	return r
}

// List returns all registered adapter names.
func (r *Registry) List() []string {
	names := make([]string, 0, len(r.factories))
	for k := range r.factories {
		names = append(names, k)
	}
	return names
}

// Get returns the named adapter or an error if unknown.
func (r *Registry) Get(name string) (Adapter, error) {
	f, ok := r.factories[name]
	if !ok {
		return nil, fmt.Errorf("unknown CLI adapter: %s (known: %s)", name, strings.Join(r.List(), ", "))
	}
	return f(), nil
}

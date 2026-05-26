package ideharness

import (
	"fmt"
	"strings"
)

// Adapter represents an IDE-based harness adapter.
type Adapter interface {
	Name() string
	Run(scenario, prompt string) (string, error)
}

// stubAdapter is a placeholder for IDE adapters that require GUI interaction.
type stubAdapter struct{ name string }

func (a *stubAdapter) Name() string { return a.name }
func (a *stubAdapter) Run(scenario, prompt string) (string, error) {
	return "", fmt.Errorf("%s: IDE adapters require manual or scripted interaction", a.name)
}

// Registry maps IDE adapter names to factory functions.
type Registry struct {
	factories map[string]func() Adapter
}

// NewRegistry creates a registry with all known IDE adapters.
func NewRegistry() *Registry {
	r := &Registry{factories: make(map[string]func() Adapter)}
	for _, name := range []string{"cursor", "cline", "kiro", "copilot", "windsurf"} {
		n := name
		r.factories[n] = func() Adapter { return &stubAdapter{name: n} }
	}
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
		return nil, fmt.Errorf("unknown IDE adapter: %s (known: %s)", name, strings.Join(r.List(), ", "))
	}
	return f(), nil
}

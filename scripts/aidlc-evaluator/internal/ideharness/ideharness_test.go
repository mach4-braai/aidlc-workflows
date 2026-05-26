package ideharness_test

import (
	"testing"

	"github.com/mach4-braai/aidlc-workflows/aidlc-evaluator/internal/ideharness"
)

func TestIDERegistryListsKnownAdapters(t *testing.T) {
	r := ideharness.NewRegistry()
	adapters := r.List()
	if len(adapters) < 2 {
		t.Fatalf("expected at least 2 IDE adapters, got %d", len(adapters))
	}
}

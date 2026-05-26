package cliharness_test

import (
	"testing"

	"github.com/mach4-braai/aidlc-workflows/aidlc-evaluator/internal/cliharness"
)

func TestRegistryListsKnownAdapters(t *testing.T) {
	r := cliharness.NewRegistry()
	adapters := r.List()
	if len(adapters) < 2 {
		t.Fatalf("expected at least 2 adapters, got %d", len(adapters))
	}
}

func TestNormalizerStripsANSIEscapeCodes(t *testing.T) {
	input := "\x1b[32mGreen text\x1b[0m"
	output := cliharness.Normalize(input)
	if output != "Green text" {
		t.Fatalf("expected stripped text, got %q", output)
	}
}

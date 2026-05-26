package shared_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mach4-braai/aidlc-workflows/aidlc-evaluator/internal/shared"
)

func TestScrubCredentialsRemovesAWSKey(t *testing.T) {
	input := "token: AKIAIOSFODNN7EXAMPLE and more text"
	output := shared.ScrubCredentials(input)
	if output == input {
		t.Fatal("should scrub AWS access key")
	}
	if strings.Contains(output, "AKIAIOSFODNN7EXAMPLE") {
		t.Fatal("scrubbed output must not contain raw key")
	}
}

func TestLoadScenarioFromYAML(t *testing.T) {
	yaml := "name: sci-calc\ndescription: Scientific calculator\n"
	tmp := t.TempDir()
	f := filepath.Join(tmp, "scenario.yaml")
	os.WriteFile(f, []byte(yaml), 0644)
	s, err := shared.LoadScenario(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Name != "sci-calc" {
		t.Fatalf("expected sci-calc, got %s", s.Name)
	}
}

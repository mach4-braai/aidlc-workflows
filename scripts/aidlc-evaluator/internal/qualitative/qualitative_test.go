package qualitative_test

import (
	"testing"

	"github.com/mach4-braai/aidlc-workflows/aidlc-evaluator/internal/qualitative"
)

func TestScoreReturnsOneForIdenticalDocs(t *testing.T) {
	content := "# Requirements\n\n## FR-001: Login\n\nUsers must log in.\n"
	score := qualitative.Score(content, content)
	if score.Percent < 95.0 {
		t.Fatalf("identical docs should score ≥ 95%%, got %.1f", score.Percent)
	}
}

func TestScoreReturnsZeroForEmptyActual(t *testing.T) {
	golden := "# Requirements\n\n## FR-001: Login\n\nUsers must log in.\n"
	score := qualitative.Score(golden, "")
	if score.Percent > 10.0 {
		t.Fatalf("empty actual should score < 10%%, got %.1f", score.Percent)
	}
}

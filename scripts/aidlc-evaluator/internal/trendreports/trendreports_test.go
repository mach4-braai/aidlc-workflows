package trendreports_test

import (
	"testing"

	"github.com/mach4-braai/aidlc-workflows/aidlc-evaluator/internal/trendreports"
)

func TestSparklineRendersCorrectWidth(t *testing.T) {
	values := []float64{0.5, 0.6, 0.7, 0.8, 0.75}
	line := trendreports.Sparkline(values)
	if len([]rune(line)) != len(values) {
		t.Fatalf("sparkline width %d != data points %d", len([]rune(line)), len(values))
	}
}

func TestGatePassesWhenAboveThreshold(t *testing.T) {
	passed := trendreports.CheckGate(0.85, 0.80)
	if !passed {
		t.Fatal("0.85 should pass 0.80 threshold")
	}
}

func TestGateFailsWhenBelowThreshold(t *testing.T) {
	passed := trendreports.CheckGate(0.75, 0.80)
	if passed {
		t.Fatal("0.75 should fail 0.80 threshold")
	}
}

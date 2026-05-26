package aireview

import (
	"context"
	"encoding/json"

	"github.com/mach4-braai/aidlc-workflows/aidlc-designreview/internal/parsing"
)

// GapAgent identifies missing design concerns.
type GapAgent struct {
	base *BaseAgent
}

// NewGapAgent creates a GapAgent.
func NewGapAgent(base *BaseAgent) *GapAgent {
	return &GapAgent{base: base}
}

// Execute runs the gap analysis.
func (a *GapAgent) Execute(ctx context.Context, data parsing.DesignData, sysPrompt, userPrompt string) (GapAnalysisResult, error) {
	if a.base == nil {
		return GapAnalysisResult{}, nil
	}
	text, usage, err := a.base.InvokeModel(ctx, sysPrompt, userPrompt)
	if err != nil {
		return GapAnalysisResult{Usage: usage}, err
	}
	var resp struct {
		Findings []GapFinding `json:"findings"`
		Summary  string       `json:"summary"`
	}
	if err := json.Unmarshal([]byte(ExtractJSONFromMarkdown(text)), &resp); err != nil {
		return GapAnalysisResult{Usage: usage}, nil
	}
	return GapAnalysisResult{Findings: resp.Findings, Summary: resp.Summary, Usage: usage}, nil
}

package aireview

import (
	"context"
	"encoding/json"

	"github.com/mach4-braai/aidlc-workflows/aidlc-designreview/internal/parsing"
)

// CritiqueAgent reviews design documents for issues.
type CritiqueAgent struct {
	base *BaseAgent
}

// NewCritiqueAgent creates a CritiqueAgent.
func NewCritiqueAgent(base *BaseAgent) *CritiqueAgent {
	return &CritiqueAgent{base: base}
}

// Execute runs the critique analysis and returns findings.
func (a *CritiqueAgent) Execute(ctx context.Context, data parsing.DesignData, sysPrompt, userPrompt string) (CritiqueResult, error) {
	if a.base == nil {
		return CritiqueResult{}, nil
	}
	text, usage, err := a.base.InvokeModel(ctx, sysPrompt, userPrompt)
	if err != nil {
		return CritiqueResult{Usage: usage}, err
	}
	var resp struct {
		Findings []CritiqueFinding `json:"findings"`
		Summary  string            `json:"summary"`
	}
	if err := json.Unmarshal([]byte(ExtractJSONFromMarkdown(text)), &resp); err != nil {
		return CritiqueResult{Usage: usage}, nil
	}
	return CritiqueResult{Findings: resp.Findings, Summary: resp.Summary, Usage: usage}, nil
}

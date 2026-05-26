package aireview

import (
	"context"
	"encoding/json"

	"github.com/mach4-braai/aidlc-workflows/aidlc-designreview/internal/parsing"
)

// AlternativesAgent suggests alternative design approaches.
type AlternativesAgent struct {
	base *BaseAgent
}

// NewAlternativesAgent creates an AlternativesAgent.
func NewAlternativesAgent(base *BaseAgent) *AlternativesAgent {
	return &AlternativesAgent{base: base}
}

// Execute runs the alternatives analysis.
func (a *AlternativesAgent) Execute(ctx context.Context, data parsing.DesignData, sysPrompt, userPrompt string) (AlternativesResult, error) {
	if a.base == nil {
		return AlternativesResult{}, nil
	}
	text, usage, err := a.base.InvokeModel(ctx, sysPrompt, userPrompt)
	if err != nil {
		return AlternativesResult{Usage: usage}, err
	}
	var resp struct {
		Suggestions []AlternativeSuggestion `json:"suggestions"`
		Summary     string                  `json:"summary"`
	}
	if err := json.Unmarshal([]byte(ExtractJSONFromMarkdown(text)), &resp); err != nil {
		return AlternativesResult{Usage: usage}, nil
	}
	return AlternativesResult{Suggestions: resp.Suggestions, Summary: resp.Summary, Usage: usage}, nil
}

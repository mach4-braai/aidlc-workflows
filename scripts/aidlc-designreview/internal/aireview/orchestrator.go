package aireview

import (
	"context"
	"fmt"

	"github.com/mach4-braai/aidlc-workflows/aidlc-designreview/internal/foundation"
	"github.com/mach4-braai/aidlc-workflows/aidlc-designreview/internal/parsing"
)

// AIOrchestrator runs all three review agents in sequence.
type AIOrchestrator struct {
	base    *BaseAgent
	pm      *foundation.PromptManager
	cfg     foundation.Config
}

// NewAIOrchestrator creates an orchestrator. Pass nil base to run in mock mode.
func NewAIOrchestrator(base *BaseAgent, pm *foundation.PromptManager, cfg foundation.Config) *AIOrchestrator {
	return &AIOrchestrator{base: base, pm: pm, cfg: cfg}
}

// Run executes all three agents and aggregates their results.
func (o *AIOrchestrator) Run(ctx context.Context, data parsing.DesignData) (ReviewResult, error) {
	var result ReviewResult
	var designContent string
	if data.AppDesign != nil {
		designContent = data.AppDesign.RawContent
	}

	vars := map[string]string{"design_content": designContent}

	sysPrompt := "You are an expert software architect reviewing an AI-DLC design document. Respond with valid JSON only."

	// Critique
	critiquePrompt, _ := o.pm.BuildAgentPrompt("critique", vars)
	critiqueAgent := NewCritiqueAgent(o.base)
	cr, err := critiqueAgent.Execute(ctx, data, sysPrompt, critiquePrompt)
	if err != nil {
		return result, fmt.Errorf("critique agent: %w", err)
	}
	result.Critique = cr
	result.TotalUsage.InputTokens += cr.Usage.InputTokens
	result.TotalUsage.OutputTokens += cr.Usage.OutputTokens

	// Alternatives
	if o.cfg.Review.EnableAlternatives {
		altPrompt, _ := o.pm.BuildAgentPrompt("alternatives", vars)
		altAgent := NewAlternativesAgent(o.base)
		ar, err := altAgent.Execute(ctx, data, sysPrompt, altPrompt)
		if err != nil {
			return result, fmt.Errorf("alternatives agent: %w", err)
		}
		result.Alternatives = ar
		result.TotalUsage.InputTokens += ar.Usage.InputTokens
		result.TotalUsage.OutputTokens += ar.Usage.OutputTokens
	}

	// Gap analysis
	if o.cfg.Review.EnableGapAnalysis {
		gapPrompt, _ := o.pm.BuildAgentPrompt("gap", vars)
		gapAgent := NewGapAgent(o.base)
		gr, err := gapAgent.Execute(ctx, data, sysPrompt, gapPrompt)
		if err != nil {
			return result, fmt.Errorf("gap agent: %w", err)
		}
		result.GapAnalysis = gr
		result.TotalUsage.InputTokens += gr.Usage.InputTokens
		result.TotalUsage.OutputTokens += gr.Usage.OutputTokens
	}

	return result, nil
}

package orchestration

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mach4-braai/aidlc-workflows/aidlc-designreview/internal/aireview"
	"github.com/mach4-braai/aidlc-workflows/aidlc-designreview/internal/foundation"
	"github.com/mach4-braai/aidlc-workflows/aidlc-designreview/internal/parsing"
	"github.com/mach4-braai/aidlc-workflows/aidlc-designreview/internal/reporting"
	"github.com/mach4-braai/aidlc-workflows/aidlc-designreview/internal/validation"
)

// Config holds orchestrator configuration.
type Config struct {
	AIDLCDocsPath string
	OutputPath    string
	ConfigPath    string
	Format        string // "markdown", "html", "both"
	AWSProfile    string
	AWSRegion     string
	MockAI        bool // skip Bedrock calls for testing
}

// ReviewOrchestrator runs the full design review pipeline.
type ReviewOrchestrator struct {
	cfg foundation.Config
	raw Config
}

// NewReviewOrchestrator creates an orchestrator with the given raw config.
func NewReviewOrchestrator(raw Config) *ReviewOrchestrator {
	var cfg foundation.Config
	if raw.ConfigPath != "" {
		cfg, _ = foundation.LoadConfig(raw.ConfigPath)
	} else {
		cfg, _ = foundation.LoadDefaultConfig()
	}
	return &ReviewOrchestrator{cfg: cfg, raw: raw}
}

// RunResult holds all outputs of a review run.
type RunResult struct {
	ReportData  reporting.ReportData
	OutputFiles []string
}

// Run executes the 6-stage review pipeline for the given project root.
func (o *ReviewOrchestrator) Run(projectRoot string) (*RunResult, error) {
	// Stage 1: Validate
	docsPath := o.raw.AIDLCDocsPath
	if docsPath == "" {
		docsPath = filepath.Join(projectRoot, "aidlc-docs")
	}
	if _, err := os.Stat(docsPath); err != nil {
		return nil, fmt.Errorf("aidlc-docs not found at %s: %w", docsPath, err)
	}
	vr := validation.ValidateStructure(projectRoot)
	if !vr.IsValid {
		return nil, fmt.Errorf("project structure invalid: %v", vr.Errors)
	}

	// Stage 2: Parse
	data := parsing.LoadDesignData(docsPath)

	// Stage 3: AI review (or mock)
	var reviewResult aireview.ReviewResult
	if !o.raw.MockAI {
		pm, err := foundation.LoadPromptManager()
		if err != nil {
			return nil, fmt.Errorf("load prompt manager: %w", err)
		}
		ctx := context.Background()
		base, err := aireview.NewBaseAgent(ctx, o.cfg.Models.DefaultModel, o.raw.AWSProfile, o.raw.AWSRegion)
		if err != nil {
			return nil, fmt.Errorf("init bedrock agent: %w", err)
		}
		orch := aireview.NewAIOrchestrator(base, pm, o.cfg)
		reviewResult, err = orch.Run(ctx, data)
		if err != nil {
			return nil, fmt.Errorf("AI review: %w", err)
		}
	}

	// Stage 4: Build report
	reportData := reporting.BuildReport(reviewResult)

	// Stage 5: Write outputs
	outputDir := o.raw.OutputPath
	if outputDir == "" {
		outputDir = projectRoot
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, err
	}

	var files []string
	format := o.raw.Format
	if format == "" {
		format = "markdown"
	}

	if format == "markdown" || format == "both" {
		out, err := reporting.RenderMarkdown(reportData)
		if err != nil {
			return nil, err
		}
		p := filepath.Join(outputDir, "design-review.md")
		if err := os.WriteFile(p, []byte(out), 0644); err != nil {
			return nil, err
		}
		files = append(files, p)
	}
	if format == "html" || format == "both" {
		out, err := reporting.RenderHTML(reportData)
		if err != nil {
			return nil, err
		}
		p := filepath.Join(outputDir, "design-review.html")
		if err := os.WriteFile(p, []byte(out), 0644); err != nil {
			return nil, err
		}
		files = append(files, p)
	}

	return &RunResult{ReportData: reportData, OutputFiles: files}, nil
}

package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/pipeline"
)

func runGenerate(cmd *cobra.Command, args []string) error {
	cfg := pipeline.Config{
		ProjectRoot: flagInput,
		OutputDir:   flagOutput,
		Format:      flagFormat,
		UseAI:       !flagNoAI,
		AWSProfile:  flagProfile,
		AWSRegion:   flagRegion,
		Verbose:     flagVerbose,
	}
	report, err := pipeline.Run(cfg)
	if err != nil {
		return err
	}
	fmt.Printf("Generated traceability report for %q\n", report.ProjectName)
	fmt.Printf("Artifacts: %d | Relationships: %d | Gaps: %d\n",
		len(report.Artifacts), len(report.Relationships), len(report.Gaps))
	return nil
}

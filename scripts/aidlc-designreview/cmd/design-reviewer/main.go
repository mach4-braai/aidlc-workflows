package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mach4-braai/aidlc-workflows/aidlc-designreview/internal/orchestration"
)

var version = "dev"

var (
	flagAIDLCDocs string
	flagOutput    string
	flagConfig    string
	flagFormat    string
	flagProfile   string
	flagRegion    string
)

var rootCmd = &cobra.Command{
	Use:     "design-reviewer",
	Short:   "AI-DLC design reviewer — AI-powered design document analysis",
	Version: version,
	RunE:    runReview,
}

func init() {
	rootCmd.Flags().StringVar(&flagAIDLCDocs, "aidlc-docs", "", "Path to aidlc-docs directory (default: <project>/aidlc-docs)")
	rootCmd.Flags().StringVar(&flagOutput, "output", "", "Output directory for generated reports (default: project root)")
	rootCmd.Flags().StringVar(&flagConfig, "config", "", "Path to config.yaml")
	rootCmd.Flags().StringVar(&flagFormat, "format", "markdown", "Output format: markdown, html, or both")
	rootCmd.Flags().StringVar(&flagProfile, "profile", "", "AWS profile name")
	rootCmd.Flags().StringVar(&flagRegion, "region", "us-east-1", "AWS region")
}

func runReview(cmd *cobra.Command, args []string) error {
	projectRoot := "."
	if len(args) > 0 {
		projectRoot = args[0]
	}

	cfg := orchestration.Config{
		AIDLCDocsPath: flagAIDLCDocs,
		OutputPath:    flagOutput,
		ConfigPath:    flagConfig,
		Format:        flagFormat,
		AWSProfile:    flagProfile,
		AWSRegion:     flagRegion,
	}

	orch := orchestration.NewReviewOrchestrator(cfg)
	result, err := orch.Run(projectRoot)
	if err != nil {
		return err
	}

	fmt.Printf("Design review complete.\n")
	fmt.Printf("Findings: %d total (%d critical, %d high, %d medium, %d low)\n",
		result.ReportData.Summary.TotalFindings,
		result.ReportData.Summary.CriticalCount,
		result.ReportData.Summary.HighCount,
		result.ReportData.Summary.MediumCount,
		result.ReportData.Summary.LowCount,
	)
	for _, f := range result.OutputFiles {
		fmt.Printf("Report written: %s\n", f)
	}
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

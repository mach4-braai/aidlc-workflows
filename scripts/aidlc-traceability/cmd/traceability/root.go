package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:     "traceability",
	Short:   "AI-DLC traceability matrix generator",
	Version: version,
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a traceability matrix for a project",
	RunE:  runGenerate,
}

var (
	flagInput   string
	flagOutput  string
	flagFormat  string
	flagNoAI    bool
	flagProfile string
	flagRegion  string
	flagVerbose bool
)

func init() {
	generateCmd.Flags().StringVar(&flagInput, "input", ".", "Project root directory")
	generateCmd.Flags().StringVar(&flagOutput, "output", ".", "Output directory for generated reports")
	generateCmd.Flags().StringVar(&flagFormat, "format", "markdown", "Output format: markdown, html, or both")
	generateCmd.Flags().BoolVar(&flagNoAI, "no-ai", false, "Skip AI-based relationship inference")
	generateCmd.Flags().StringVar(&flagProfile, "profile", "", "AWS profile name")
	generateCmd.Flags().StringVar(&flagRegion, "region", "us-east-1", "AWS region")
	generateCmd.Flags().BoolVarP(&flagVerbose, "verbose", "v", false, "Verbose output")
	rootCmd.AddCommand(generateCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

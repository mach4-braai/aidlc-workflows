package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	fullVisionPath  string
	fullTechEnvPath string
	fullGoldenPath  string
	fullOpenAPIPath string
	fullConfig      string
	fullProfile     string
	fullRegion      string
)

var fullCmd = &cobra.Command{
	Use:   "full",
	Short: "Run a full evaluation (execution + qualitative + quantitative + contract)",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Running full evaluation\n")
		fmt.Printf("  vision: %s\n", fullVisionPath)
		fmt.Printf("  tech-env: %s\n", fullTechEnvPath)
		return nil
	},
}

func init() {
	fullCmd.Flags().StringVar(&fullVisionPath, "vision", "", "Path to vision document")
	fullCmd.Flags().StringVar(&fullTechEnvPath, "tech-env", "", "Path to technical environment document")
	fullCmd.Flags().StringVar(&fullGoldenPath, "golden", "", "Path to golden output directory")
	fullCmd.Flags().StringVar(&fullOpenAPIPath, "openapi", "", "Path to OpenAPI spec")
	fullCmd.Flags().StringVar(&fullConfig, "config", "", "Path to evaluator config")
	fullCmd.Flags().StringVar(&fullProfile, "aws-profile", "", "AWS profile name")
	fullCmd.Flags().StringVar(&fullRegion, "aws-region", "us-east-1", "AWS region")
}

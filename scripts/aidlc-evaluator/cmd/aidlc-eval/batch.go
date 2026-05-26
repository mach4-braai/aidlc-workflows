package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	batchModels   []string
	batchScenario string
	batchConfig   string
)

var batchCmd = &cobra.Command{
	Use:   "batch",
	Short: "Run evaluation across multiple models",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Running batch evaluation across %d models, scenario=%s\n", len(batchModels), batchScenario)
		return nil
	},
}

func init() {
	batchCmd.Flags().StringSliceVar(&batchModels, "models", nil, "Comma-separated list of model IDs to evaluate")
	batchCmd.Flags().StringVar(&batchScenario, "scenario", "", "Path to scenario directory")
	batchCmd.Flags().StringVar(&batchConfig, "config", "", "Path to evaluator config")
}

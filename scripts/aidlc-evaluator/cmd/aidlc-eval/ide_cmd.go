package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	ideAdapter  string
	ideScenario string
	ideConfig   string
)

var ideCmd = &cobra.Command{
	Use:   "ide",
	Short: "Evaluate an IDE-based AI coding tool",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Running IDE harness evaluation with adapter=%s scenario=%s\n", ideAdapter, ideScenario)
		return nil
	},
}

func init() {
	ideCmd.Flags().StringVar(&ideAdapter, "ide", "cursor", "IDE adapter name (cursor, cline, kiro, windsurf, copilot)")
	ideCmd.Flags().StringVar(&ideScenario, "scenario", "", "Path to scenario directory")
	ideCmd.Flags().StringVar(&ideConfig, "config", "", "Path to evaluator config")
}

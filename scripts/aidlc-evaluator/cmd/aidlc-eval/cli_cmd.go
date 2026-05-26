package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	cliAdapter  string
	cliScenario string
	cliConfig   string
)

var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "Evaluate a CLI-based AI coding tool",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Running CLI harness evaluation with adapter=%s scenario=%s\n", cliAdapter, cliScenario)
		return nil
	},
}

func init() {
	cliCmd.Flags().StringVar(&cliAdapter, "cli", "claude-code", "CLI adapter name (claude-code, kiro-cli)")
	cliCmd.Flags().StringVar(&cliScenario, "scenario", "", "Path to scenario directory")
	cliCmd.Flags().StringVar(&cliConfig, "config", "", "Path to evaluator config")
}

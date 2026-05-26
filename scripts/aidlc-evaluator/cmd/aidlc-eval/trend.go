package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	trendBaseline string
	trendRunsDir  string
	trendOutput   string
)

var trendCmd = &cobra.Command{
	Use:   "trend",
	Short: "Generate trend reports from historical evaluation runs",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Generating trend report: baseline=%s runs=%s output=%s\n", trendBaseline, trendRunsDir, trendOutput)
		return nil
	},
}

func init() {
	trendCmd.Flags().StringVar(&trendBaseline, "baseline", "", "Path to baseline run YAML")
	trendCmd.Flags().StringVar(&trendRunsDir, "runs-dir", ".", "Directory containing historical run results")
	trendCmd.Flags().StringVar(&trendOutput, "output", ".", "Output directory for trend reports")
}

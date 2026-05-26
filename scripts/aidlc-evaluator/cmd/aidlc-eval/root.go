package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:     "aidlc-eval",
	Short:   "AI-DLC evaluator — multi-mode evaluation framework",
	Version: version,
}

func init() {
	rootCmd.AddCommand(fullCmd)
	rootCmd.AddCommand(cliCmd)
	rootCmd.AddCommand(ideCmd)
	rootCmd.AddCommand(batchCmd)
	rootCmd.AddCommand(trendCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

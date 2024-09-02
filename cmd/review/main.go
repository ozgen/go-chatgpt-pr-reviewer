package main

import (
	"fmt"
	"go-chatgpt-pr-reviewer/review"
	"os"

	"github.com/spf13/cobra"
)

// Global variables to store flag values
var (
	localDir string
	prNumber int
)

func main() {
	// Create a new Cobra command
	var rootCmd = &cobra.Command{
		Use:   "review",
		Short: "review is a CLI tool to review GitHub PRs using ChatGPT",
		Run: func(cmd *cobra.Command, args []string) {
			review.RunReview(localDir, prNumber)
		},
	}

	// Define flags
	rootCmd.Flags().StringVar(&localDir, "local", "", "Local git repository directory")
	rootCmd.Flags().IntVar(&prNumber, "pr", 0, "Pull Request number to review")
	rootCmd.MarkFlagRequired("local")
	rootCmd.MarkFlagRequired("pr")

	// Execute the command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

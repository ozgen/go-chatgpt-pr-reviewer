package main

import (
	"fmt"
	"github.com/ozgen/go-chatgpt-pr-reviewer/review"
	"os"

	"github.com/spf13/cobra"
)

// Global variables to store flag values
var (
	localDir     string
	prNumber     int
	postComments bool // Default is false
)

func main() {
	// Create a new Cobra command
	var rootCmd = &cobra.Command{
		Use:   "review",
		Short: "review is a CLI tool to review GitHub PRs using ChatGPT",
		Run: func(cmd *cobra.Command, args []string) {
			review.RunReview(localDir, prNumber, postComments)
		},
	}

	// Define flags
	rootCmd.Flags().StringVar(&localDir, "local", "", "Local git repository directory")
	rootCmd.Flags().IntVar(&prNumber, "pr", 0, "Pull Request number to review")
	rootCmd.Flags().BoolVar(&postComments, "post-comments", false, "Post review comments to GitHub (default: false)")
	rootCmd.MarkFlagRequired("local")
	rootCmd.MarkFlagRequired("pr")

	// Execute the command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

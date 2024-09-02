package review

import (
	"context"
	"fmt"
	"github.com/ozgen/go-chatgpt-pr-reviewer/chatgpt"
	"github.com/ozgen/go-chatgpt-pr-reviewer/config"
	"github.com/ozgen/go-chatgpt-pr-reviewer/github"
	"github.com/ozgen/go-chatgpt-pr-reviewer/types"
	"log"
)

func RunReview(localDir string, prNumber int) {
	// Load configuration
	apiKey := config.Envs.OpenAIApiKey
	organizationId := config.Envs.OrganizationId
	projectId := config.Envs.ProjectId
	githubToken := config.Envs.GithubToken
	ctx := context.Background()

	// Get GitHub repository information
	owner, repo, err := github.GetGitRemoteInfo(localDir)
	if err != nil {
		fmt.Printf("Error getting git remote info: %v\n", err)
		return
	}

	fmt.Printf("Owner: %s, Repo: %s\n", owner, repo)

	// Set up GitHub client
	githubClient := github.SetupGitHubClient(ctx, githubToken)

	// Get PR changes
	files, err := github.GetPRChanges(ctx, githubClient, owner, repo, prNumber)
	if err != nil {
		log.Fatalf("Failed to get PR files: %v", err)
	}

	// Display the files with changes
	for _, file := range files {
		fmt.Printf("File: %s, Changes: +%d -%d\n", file.GetFilename(), file.GetAdditions(), file.GetDeletions())
	}

	// Set up ChatGPT client
	client := chatgpt.NewChatGPTClient(apiKey, organizationId, projectId)

	// Process each file and send the modified blocks to ChatGPT for review
	for _, file := range files {
		modifiedLines := github.ExtractModifiedLinesWithNumbers(file.GetPatch())

		for _, modifiedLine := range modifiedLines {
			// Send each modified block to ChatGPT
			prompt := fmt.Sprintf("Code Review Request: Review the following block in file %s starting at line %d. Suggest any improvements:\n\n%s", file.GetFilename(), modifiedLine.LineNumber, modifiedLine.Content)
			feedback, err := client.SendRequest(types.Payload{
				Prompt:    prompt,
				MaxTokens: 500,
			})

			if err != nil {
				log.Printf("Error during ChatGPT review: %v", err)
				continue
			}

			fmt.Printf("Feedback for file %s at line %d:\n%s\n", file.GetFilename(), modifiedLine.LineNumber, feedback)

			if feedback != "" {
				commentBody := fmt.Sprintf("ChatGPT suggests:\n%s", feedback)
				path := file.GetFilename()
				err = github.PostReviewComment(ctx, githubClient, owner, repo, prNumber, commentBody, path, modifiedLine.LineNumber)
				if err != nil {
					log.Printf("Failed to post comment: %v", err)
				}
			}
		}
	}
}

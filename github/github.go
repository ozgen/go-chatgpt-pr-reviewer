package github

import (
	"context"
	"fmt"
	"github.com/google/go-github/v42/github"
	"go-chatgpt-pr-reviewer/types"
	"golang.org/x/oauth2"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// SetupGitHubClient creates a GitHub client with the provided token.
func SetupGitHubClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

// GetPRChanges fetches file changes for a given pull request.
func GetPRChanges(ctx context.Context, client *github.Client, owner, repo string, prNumber int) ([]*github.CommitFile, error) {
	opts := &github.ListOptions{PerPage: 10}
	files, _, err := client.PullRequests.ListFiles(ctx, owner, repo, prNumber, opts)
	if err != nil {
		return nil, err
	}
	return files, nil
}

// GetGitRemoteInfo executes git command to get remote URL and extracts owner and repo.
// The directory parameter specifies the path to the Git repository.
func GetGitRemoteInfo(directory string) (string, string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = directory // Set the working directory to the specified directory
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", "", fmt.Errorf("failed to execute git command: %v, output: %s", err, string(output))
	}
	return parseGitURL(strings.TrimSpace(string(output)))
}

// PostReviewComment posts a comment on a pull request at a specified position within a file.
func PostReviewComment(ctx context.Context, client *github.Client, owner, repo string, prNumber int, body, path string, line int) error {
	// Retrieve the pull request to get the latest commit ID
	pr, _, err := client.PullRequests.Get(ctx, owner, repo, prNumber)
	if err != nil {
		return fmt.Errorf("failed to retrieve PR information: %v", err)
	}
	commitID := pr.GetHead().GetSHA()

	// Create a new review comment
	comment := &github.PullRequestComment{
		Body:     &body,
		Path:     &path,
		CommitID: &commitID,
		Line:     &line,
		Side:     github.String("RIGHT"), // Usually you want to comment on the "RIGHT" side of the diff
	}

	_, _, err = client.PullRequests.CreateComment(ctx, owner, repo, prNumber, comment)
	if err != nil {
		return fmt.Errorf("failed to post review comment: %v", err)
	}
	return nil
}

// ExtractModifiedLinesWithNumbers captures entire blocks of changes instead of line by line.
func ExtractModifiedLinesWithNumbers(patch string) []types.ModifiedLine {
	lines := strings.Split(patch, "\n")
	var modifiedLines []types.ModifiedLine
	var currentLineOriginal, currentLineModified int
	var currentBlock []string
	var blockStartLine int

	// Regex to match @@ -a,b +c,d @@ lines
	hunkHeaderRegex := regexp.MustCompile(`@@ -\d+,\d+ \+(\d+),\d+ @@`)
	for _, line := range lines {
		if strings.HasPrefix(line, "@@") {
			// Process the previous block if any
			if len(currentBlock) > 0 {
				// Append the block with the starting line number
				modifiedLines = append(modifiedLines, types.ModifiedLine{LineNumber: blockStartLine, Content: strings.Join(currentBlock, "\n")})
				currentBlock = nil // Reset the block
			}

			// Parse hunk header to get starting line number in modified file
			matches := hunkHeaderRegex.FindStringSubmatch(line)
			if len(matches) == 2 {
				currentLineModified, _ = strconv.Atoi(matches[1])
				currentLineOriginal = currentLineModified // Start the original line counter from here
			}
		} else if strings.HasPrefix(line, "+") || strings.HasPrefix(line, "-") {
			if len(currentBlock) == 0 {
				blockStartLine = currentLineModified // Record the starting line for the block
			}
			// Append lines to the current block
			currentBlock = append(currentBlock, line)
			if strings.HasPrefix(line, "+") {
				currentLineModified++
			} else if strings.HasPrefix(line, "-") {
				currentLineOriginal++
			}
		} else {
			if len(currentBlock) > 0 {
				// End of a block detected, store it
				modifiedLines = append(modifiedLines, types.ModifiedLine{LineNumber: blockStartLine, Content: strings.Join(currentBlock, "\n")})
				currentBlock = nil // Reset the block
			}
			currentLineOriginal++
			currentLineModified++
		}
	}

	// Process the last block if any
	if len(currentBlock) > 0 {
		modifiedLines = append(modifiedLines, types.ModifiedLine{LineNumber: blockStartLine, Content: strings.Join(currentBlock, "\n")})
	}

	return modifiedLines
}

// parseGitURL parses the git URL to extract the owner and repository name.
func parseGitURL(url string) (string, string, error) {
	var sshPrefix = "git@github.com:"
	var httpsPrefix = "https://github.com/"
	var suffix = ".git"

	if strings.HasPrefix(url, sshPrefix) {
		trimmed := strings.TrimPrefix(url, sshPrefix)
		trimmed = strings.TrimSuffix(trimmed, suffix)
		parts := strings.Split(trimmed, "/")
		if len(parts) == 2 {
			return parts[0], parts[1], nil
		}
	} else if strings.HasPrefix(url, httpsPrefix) {
		trimmed := strings.TrimPrefix(url, httpsPrefix)
		trimmed = strings.TrimSuffix(trimmed, suffix)
		parts := strings.Split(trimmed, "/")
		if len(parts) == 2 {
			return parts[0], parts[1], nil
		}
	}
	return "", "", fmt.Errorf("unknown git url format")
}

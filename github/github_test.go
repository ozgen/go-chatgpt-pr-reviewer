package github

import (
	"context"
	"github.com/ozgen/go-chatgpt-pr-reviewer/types"
	"golang.org/x/oauth2"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/google/go-github/v42/github"
)

// TestSetupGitHubClient tests the SetupGitHubClient function.
func TestSetupGitHubClient(t *testing.T) {
	ctx := context.Background()
	token := "fake-token"

	client := SetupGitHubClient(ctx, token)

	if client == nil {
		t.Fatal("Expected GitHub client, got nil")
	}

	// Verify that the client has a valid OAuth2 transport
	if _, ok := client.Client().Transport.(*oauth2.Transport); !ok {
		t.Errorf("Expected OAuth2 transport, got %T", client.Client().Transport)
	}
}

// TestGetPRChanges tests the GetPRChanges function.
func TestGetPRChanges(t *testing.T) {
	// Mock the GitHub API server
	mockResponse := `[
		{
			"filename": "example.go",
			"additions": 10,
			"deletions": 2,
			"changes": 12,
			"status": "modified",
			"patch": "@@ -1,2 +1,2 @@\n-func Old() {}\n+func New() {}"
		}
	]`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	ctx := context.Background()
	client := github.NewClient(nil)

	// Set the BaseURL of the client to the mock server URL
	baseURL, _ := url.Parse(server.URL + "/")
	client.BaseURL = baseURL

	owner := "owner"
	repo := "repo"
	prNumber := 1

	files, err := GetPRChanges(ctx, client, owner, repo, prNumber)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("Expected 1 file, got %d", len(files))
	}

	if files[0].GetFilename() != "example.go" {
		t.Errorf("Expected filename 'example.go', got '%s'", files[0].GetFilename())
	}
}

// TestGetGitRemoteInfo tests the GetGitRemoteInfo function.
func TestGetGitRemoteInfo(t *testing.T) {
	directory, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	owner, repo, err := GetGitRemoteInfo(directory)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if owner == "" || repo == "" {
		t.Fatalf("Expected owner and repo, got owner: '%s', repo: '%s'", owner, repo)
	}
}

// TestPostReviewComment tests the PostReviewComment function.
func TestPostReviewComment(t *testing.T) {
	// Mock the GitHub API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	ctx := context.Background()
	client := github.NewClient(nil)

	// Set the BaseURL of the client to the mock server URL
	baseURL, _ := url.Parse(server.URL + "/")
	client.BaseURL = baseURL

	owner := "owner"
	repo := "repo"
	prNumber := 1
	body := "This is a review comment"
	path := "example.go"
	line := 10

	err := PostReviewComment(ctx, client, owner, repo, prNumber, body, path, line)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

// TestExtractModifiedLinesWithNumbers tests the ExtractModifiedLinesWithNumbers function.
func TestExtractModifiedLinesWithNumbers(t *testing.T) {
	patch := `@@ -1,2 +1,2 @@
-func Old() {}
+func New() {}
@@ -4,6 +4,7 @@
 var x = 10
+var y = 20
`

	expected := []types.ModifiedLine{
		{
			LineNumber: 1,
			Content:    "-func Old() {}\n+func New() {}",
		},
		{
			LineNumber: 5,
			Content:    "+var y = 20",
		},
	}

	result := ExtractModifiedLinesWithNumbers(patch)
	if len(result) != len(expected) {
		t.Fatalf("Expected %d modified lines, got %d", len(expected), len(result))
	}

	for i, modifiedLine := range result {
		if modifiedLine.LineNumber != expected[i].LineNumber || modifiedLine.Content != expected[i].Content {
			t.Errorf("At index %d, expected %v, got %v", i, expected[i], modifiedLine)
		}
	}
}

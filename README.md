
# ChatGPT PR Reviewer

ChatGPT PR Reviewer is a Go-based command-line tool that integrates with GitHub and OpenAI's ChatGPT to perform automated code reviews on pull requests (PRs). The tool fetches changes from a specified PR, sends the modified code blocks to ChatGPT for review, and posts the feedback as comments on the PR.

## Features

- Fetches pull request changes from a GitHub repository.
- Sends modified code blocks to ChatGPT for review.
- Posts ChatGPT's feedback as comments on the pull request.

## Prerequisites

- Go 1.17 or later
- A GitHub account with a personal access token.
- An OpenAI account with API access.

## Installation

1. **Install the CLI Tool**

   Use the following command to install the `review` CLI tool:

   ```bash
   go install github.com/ozgen/go-chatgpt-pr-reviewer/cmd/review@v1.0.3
   ```

This will download and build the CLI tool, placing the `review` binary in your `GOPATH/bin` or `GOBIN` directory.

2. **Set Environment Variables**

   Before running the tool, you need to set the following environment variables:

   ```bash
   export OPENAI_API_KEY=<your_openai_api_key>
   export ORGANIZATION_ID=<your_organization_id>
   export PROJECT_ID=<your_project_id>
   export GITHUB_TOKEN=<your_github_token>
   ```

    - **`OPENAI_API_KEY`**: Your OpenAI API key.
    - **`ORGANIZATION_ID`**: Your OpenAI organization ID.
    - **`PROJECT_ID`**: Your OpenAI project ID.
    - **`GITHUB_TOKEN`**: Your GitHub personal access token with `repo` scope.

## Usage

Once installed, you can use the `review` command to perform a code review on a pull request:

```bash
review --local "/path/to/local/repo" --pr 1
```

- `--local` specifies the path to your local Git repository.
- `--pr` specifies the pull request number you want to review.

### Example

1. **Set Environment Variables**:

   ```bash
   export OPENAI_API_KEY=sk-...
   export ORGANIZATION_ID=org-...
   export PROJECT_ID=proj-...
   export GITHUB_TOKEN=ghp_...
   ```

2. **Run the Review Command**:

   ```bash
   review --local "/Users/ozgen/Desktop/checkstyle-runner" --pr 1
   ```

   This command will fetch the changes from PR #1 in the specified local repository, send the changes to ChatGPT for review, and post the feedback as comments on the PR.

## Development

### Building from Source

If you want to build the project from source:

1. Clone the repository:

   ```bash
   git clone https://github.com/ozgen/go-chatgpt-pr-reviewer.git
   cd go-chatgpt-pr-reviewer
   ```

2. Build the CLI tool:

   ```bash
   go build -o pr-reviewer ./cmd/review
   ```

3. Run the tool:

   ```bash
   ./pr-reviewer --local "/path/to/local/repo" --pr 1
   ```

### Running Tests

To run tests:

```bash
go test -v ./...
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

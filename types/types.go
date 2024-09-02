package types

// Payload represents the data sent to OpenAI.
type Payload struct {
	Prompt    string `json:"prompt"`
	MaxTokens int    `json:"max_tokens"`
}

// Response represents the structure of data received from OpenAI.
type Response struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

// ModifiedLine represents a line in the diff with its line number and content.
type ModifiedLine struct {
	LineNumber int
	Content    string
}

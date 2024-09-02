package chatgpt

import (
	"encoding/json"
	"go-chatgpt-pr-reviewer/types"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestSendRequestSuccess tests the successful case of sending a request
func TestSendRequestSuccess(t *testing.T) {
	// Mock the server
	mockResponse := types.Response{
		Choices: []struct {
			Text string `json:"text"`
		}{
			{Text: "This is a test response from ChatGPT."},
		},
	}
	responseBody, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(responseBody)
	}))
	defer server.Close()

	// Create a client with the mock server URL
	client := NewChatGPTClient("fake-api-key", "fake-org-id", "fake-project-id", server.URL)

	// Define a sample payload
	payload := types.Payload{
		Prompt:    "Test Prompt",
		MaxTokens: 50,
	}

	// Call the method
	response, err := client.SendRequest(payload)

	// Validate the response
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedResponse := "This is a test response from ChatGPT."
	if response != expectedResponse {
		t.Errorf("Expected response '%s', got '%s'", expectedResponse, response)
	}
}

// TestSendRequestError tests the case where the API returns an error
func TestSendRequestError(t *testing.T) {
	// Mock the server to return an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer server.Close()

	// Create a client with the mock server URL
	client := NewChatGPTClient("fake-api-key", "fake-org-id", "fake-project-id", server.URL)

	// Define a sample payload
	payload := types.Payload{
		Prompt:    "Test Prompt",
		MaxTokens: 50,
	}

	// Call the method
	_, err := client.SendRequest(payload)

	// Validate the error
	if err == nil {
		t.Fatal("Expected an error, but got none")
	}
}

// TestSendRequestInvalidJSON tests the case where the JSON response is invalid
func TestSendRequestInvalidJSON(t *testing.T) {
	// Mock the server to return invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{invalid json}`))
	}))
	defer server.Close()

	// Create a client with the mock server URL
	client := NewChatGPTClient("fake-api-key", "fake-org-id", "fake-project-id", server.URL)

	// Define a sample payload
	payload := types.Payload{
		Prompt:    "Test Prompt",
		MaxTokens: 50,
	}

	// Call the method
	_, err := client.SendRequest(payload)

	// Validate the error
	if err == nil {
		t.Fatal("Expected a JSON unmarshalling error, but got none")
	}
}

// TestSendRequestNoChoices tests the case where the API response has no choices
func TestSendRequestNoChoices(t *testing.T) {
	// Mock the server to return a response with no choices
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"choices":[]}`))
	}))
	defer server.Close()

	// Create a client with the mock server URL
	client := NewChatGPTClient("fake-api-key", "fake-org-id", "fake-project-id", server.URL)

	// Define a sample payload
	payload := types.Payload{
		Prompt:    "Test Prompt",
		MaxTokens: 50,
	}

	// Call the method
	response, err := client.SendRequest(payload)

	// Validate the response
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if response != "" {
		t.Errorf("Expected empty response, got '%s'", response)
	}
}

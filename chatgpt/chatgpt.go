package chatgpt

import (
	"bytes"
	"encoding/json"
	"go-chatgpt-pr-reviewer/types"
	"io/ioutil"
	"log"
	"net/http"
)

const defaultAPIURL = "https://api.openai.com/v1/engines/gpt-3.5-turbo-instruct/completions"

// ChatGPTClient holds the configuration for the API client
type ChatGPTClient struct {
	APIKey         string
	OrganizationID string
	ProjectID      string
	APIURL         string
}

// NewChatGPTClient creates a new client with the given API key
func NewChatGPTClient(apiKey, organizationID, projectID string, apiURL ...string) *ChatGPTClient {
	url := defaultAPIURL
	if len(apiURL) > 0 {
		url = apiURL[0]
	}
	return &ChatGPTClient{
		APIKey:         apiKey,
		OrganizationID: organizationID,
		ProjectID:      projectID,
		APIURL:         url,
	}
}

// SendRequest sends a code review request to ChatGPT and returns the response
func (c *ChatGPTClient) SendRequest(payload types.Payload) (string, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling payload: %v", err)
		return "", err
	}

	req, err := http.NewRequest("POST", c.APIURL, bytes.NewReader(payloadBytes))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("OpenAI-Organization", c.OrganizationID)
	req.Header.Set("OpenAI-Project", c.ProjectID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return "", err
	}

	var response types.Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("Error unmarshaling response: %v. Body: %s", err, string(body))
		return "", err
	}
	if resp.StatusCode == 200 {
		var apiWarning struct {
			Warnings []string `json:"warnings"`
		}
		if err := json.Unmarshal(body, &apiWarning); err == nil {
			for _, warning := range apiWarning.Warnings {
				log.Printf("API Warning: %s", warning)
			}
		}
	}

	if len(response.Choices) > 0 {
		log.Printf("Response received: %s", response.Choices[0].Text)
		return response.Choices[0].Text, nil
	}
	log.Println("No choices in response")
	return "", nil
}

package config

import (
	"github.com/joho/godotenv"
	"go-chatgpt-pr-reviewer/utils"
)

type Config struct {
	OpenAIApiKey   string
	OrganizationId string
	ProjectId      string
	GithubToken    string
}

var Envs = initConfig()

func initConfig() Config {

	// Load env variables
	godotenv.Load()

	return Config{

		OpenAIApiKey:   utils.GetEnv("OPENAI_API_KEY", ""),
		OrganizationId: utils.GetEnv("ORGANIZATION_ID", ""),
		ProjectId:      utils.GetEnv("PROJECT_ID", ""),
		GithubToken:    utils.GetEnv("GITHUB_TOKEN", ""),
	}
}

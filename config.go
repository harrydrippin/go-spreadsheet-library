package main

import (
	"os"

	"github.com/joho/godotenv"
)

// Config is a struct to hold all the configuration values from dotenv
type Config struct {
	GoogleOAuthClientID     string
	GoogleOAuthClientSecret string
	GoogleSpreadsheetID     string
	SlackToken              string
}

// NewConfig creates a new Config object
func NewConfig() *Config {
	// Load the .env file, but don't fail if it doesn't exist
	godotenv.Load()

	return &Config{
		GoogleOAuthClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		GoogleOAuthClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		GoogleSpreadsheetID:     os.Getenv("GOOGLE_SPREADSHEET_ID"),
		SlackToken:              os.Getenv("SLACK_TOKEN"),
	}
}

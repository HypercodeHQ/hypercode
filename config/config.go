package config

import (
	"os"
	"path/filepath"

	"github.com/hypercommithq/hypercommit/env"
)

type Config struct {
	HTTPAddr           string
	DatabasePath       string
	SigningSecret      string
	ReposBasePath      string
	GitHubClientID     string
	GitHubClientSecret string
	GitHubCallbackURL  string
}

func New() Config {
	return Config{
		HTTPAddr:           env.GetVar("HTTP_ADDR", ":3000"),
		DatabasePath:       env.GetVar("DATABASE_PATH", "hypercommit.db"),
		SigningSecret:      getSigningSecret(),
		ReposBasePath:      env.GetVar("REPOS_BASE_PATH", "repos"),
		GitHubClientID:     env.GetVar("GITHUB_OAUTH_CLIENT_ID", ""),
		GitHubClientSecret: env.GetVar("GITHUB_OAUTH_CLIENT_SECRET", ""),
		GitHubCallbackURL:  env.GetVar("GITHUB_CALLBACK_URL", "http://localhost:3000/auth/github/callback"),
	}
}

func getSigningSecret() string {
	if credsDir := os.Getenv("CREDENTIALS_DIRECTORY"); credsDir != "" {
		secretPath := filepath.Join(credsDir, "signing_secret")
		if data, err := os.ReadFile(secretPath); err == nil {
			return string(data)
		}
	}

	return env.GetVar("SIGNING_SECRET", "insecure-dev-secret")
}

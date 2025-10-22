package config

import (
	"os"
	"path/filepath"

	"github.com/hyperstitieux/hypercode/env"
)

type Config struct {
	HTTPAddr      string
	DatabasePath  string
	SigningSecret string
	ReposBasePath string
}

func New() Config {
	return Config{
		HTTPAddr:      env.GetVar("HTTP_ADDR", ":3000"),
		DatabasePath:  env.GetVar("DATABASE_PATH", "hypercode.db"),
		SigningSecret: getSigningSecret(),
		ReposBasePath: env.GetVar("REPOS_BASE_PATH", "repos"),
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

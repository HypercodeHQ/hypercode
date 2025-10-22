package env

import (
	"log/slog"
	"os"
)

func GetVar(key, fallback string) string {
	if value, found := os.LookupEnv(key); found {
		return value
	}
	slog.Warn("env var undefined, using fallback", "key", key, "fallback", fallback)
	return fallback
}

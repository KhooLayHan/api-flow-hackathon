package main

import (
	"log/slog"
	"os"
)

func NewLogger() *slog.Logger {
	// Defaults to log level INFO. Can be configured via other means, e.g. environment variables.
	logLevel := slog.LevelInfo
	if os.Getenv("LOG_LEVEL") == "DEBUG" {
		logLevel = slog.LevelDebug
	}

	// Creates a JSON handler that writes to standard output to allow Railway and Docker to capture logs.
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	return slog.New(handler)
}

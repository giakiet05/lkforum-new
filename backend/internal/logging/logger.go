package logging

import (
	"log/slog"
	"os"
	"strings"
)

func NewLogger() *slog.Logger {
	logLevel := slog.LevelInfo

	if strings.ToLower(os.Getenv("LOG_LEVEL")) == "debug" {
		logLevel = slog.LevelDebug
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)

	if logLevel == slog.LevelDebug {
		logger.Info("Debug logging enabled")
	} else {
		logger.Info("Info logging enabled")
	}

	return logger
}

func ConfigureDefault() *slog.Logger {
	logger := NewLogger()
	slog.SetDefault(logger)
	return logger
}

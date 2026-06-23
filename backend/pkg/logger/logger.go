package logger

import (
	"log/slog"
	"os"
	"strings"
)

func Init(serviceName string) *slog.Logger {
	level := getLogLevel()

	var handler slog.Handler

	if os.Getenv("ENV") == "production" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	}

	logger := slog.New(handler).With(slog.String("service", serviceName))

	slog.SetDefault(logger)

	return logger
}

func getLogLevel() slog.Level {
	envLevel := os.Getenv("LOG_LEVEL")

	switch strings.ToUpper(envLevel) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

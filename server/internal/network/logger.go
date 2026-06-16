package network

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

type contextKey string

const requestIDKey contextKey = "request_id"

func SetupLogger(level string) *slog.Logger {
	var logLevel slog.Level
	switch strings.ToLower(level) {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}

func LoggerFromContext(ctx context.Context) *slog.Logger {
	logger := slog.Default()
	if reqID, ok := ctx.Value(requestIDKey).(string); ok && reqID != "" {
		logger = logger.With("request_id", reqID)
	}
	return logger
}

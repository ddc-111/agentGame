package network

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"
)

func TestSetupLogger(t *testing.T) {
	tests := []struct {
		name  string
		level string
		want  slog.Level
	}{
		{"debug", "debug", slog.LevelDebug},
		{"info", "info", slog.LevelInfo},
		{"warn", "warn", slog.LevelWarn},
		{"error", "error", slog.LevelError},
		{"default", "invalid", slog.LevelInfo},
		{"empty", "", slog.LevelInfo},
		{"uppercase", "DEBUG", slog.LevelDebug},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := SetupLogger(tt.level)
			if logger == nil {
				t.Fatal("SetupLogger returned nil")
			}
		})
	}
}

func TestSetupLoggerJSONFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := &slog.HandlerOptions{Level: slog.LevelDebug}
	handler := slog.NewJSONHandler(&buf, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	logger.Info("test message", "key", "value")

	var entry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("log output is not valid JSON: %v", err)
	}
	if entry["msg"] != "test message" {
		t.Errorf("expected msg 'test message', got %v", entry["msg"])
	}
	if entry["key"] != "value" {
		t.Errorf("expected key 'value', got %v", entry["key"])
	}
	if _, ok := entry["time"]; !ok {
		t.Error("expected 'time' field in log output")
	}
	if _, ok := entry["level"]; !ok {
		t.Error("expected 'level' field in log output")
	}
}

func TestLoggerFromContextWithRequestID(t *testing.T) {
	var buf bytes.Buffer
	opts := &slog.HandlerOptions{Level: slog.LevelDebug}
	handler := slog.NewJSONHandler(&buf, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	ctx := context.WithValue(context.Background(), requestIDKey, "test-req-123")
	reqLogger := LoggerFromContext(ctx)
	reqLogger.Info("test message")

	var entry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("log output is not valid JSON: %v", err)
	}
	if entry["request_id"] != "test-req-123" {
		t.Errorf("expected request_id 'test-req-123', got %v", entry["request_id"])
	}
	if entry["msg"] != "test message" {
		t.Errorf("expected msg 'test message', got %v", entry["msg"])
	}
}

func TestLoggerFromContextWithoutRequestID(t *testing.T) {
	ctx := context.Background()
	logger := LoggerFromContext(ctx)
	if logger == nil {
		t.Fatal("LoggerFromContext returned nil")
	}
}

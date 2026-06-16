package network

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequestLoggingMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var buf bytes.Buffer
	opts := &slog.HandlerOptions{Level: slog.LevelDebug}
	handler := slog.NewJSONHandler(&buf, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.Use(RequestLoggingMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test?foo=bar", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var entry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("log output is not valid JSON: %v", err)
	}

	if entry["method"] != "GET" {
		t.Errorf("expected method 'GET', got %v", entry["method"])
	}
	if entry["path"] != "/test" {
		t.Errorf("expected path '/test', got %v", entry["path"])
	}
	if entry["query"] != "foo=bar" {
		t.Errorf("expected query 'foo=bar', got %v", entry["query"])
	}
	status, ok := entry["status"].(float64)
	if !ok || int(status) != 200 {
		t.Errorf("expected status 200, got %v", entry["status"])
	}
	if _, ok := entry["latency"]; !ok {
		t.Error("expected 'latency' field in log output")
	}
	if _, ok := entry["client_ip"]; !ok {
		t.Error("expected 'client_ip' field in log output")
	}
	if _, ok := entry["request_id"]; !ok {
		t.Error("expected 'request_id' field in log output")
	}
}

func TestRequestLoggingMiddlewareServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var buf bytes.Buffer
	opts := &slog.HandlerOptions{Level: slog.LevelDebug}
	handler := slog.NewJSONHandler(&buf, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.Use(RequestLoggingMiddleware())
	router.GET("/error", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "fail"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/error", nil)
	router.ServeHTTP(w, req)

	var entry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("log output is not valid JSON: %v", err)
	}

	if entry["level"] != "ERROR" {
		t.Errorf("expected level 'ERROR' for 500 status, got %v", entry["level"])
	}
}

func TestRequestLoggingMiddlewareClientError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var buf bytes.Buffer
	opts := &slog.HandlerOptions{Level: slog.LevelDebug}
	handler := slog.NewJSONHandler(&buf, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.Use(RequestLoggingMiddleware())
	router.GET("/bad", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/bad", nil)
	router.ServeHTTP(w, req)

	var entry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("log output is not valid JSON: %v", err)
	}

	if entry["level"] != "WARN" {
		t.Errorf("expected level 'WARN' for 400 status, got %v", entry["level"])
	}
}

func TestRequestIDMiddlewarePropagatesHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Request-ID", "custom-id-42")
	router.ServeHTTP(w, req)

	if w.Header().Get("X-Request-ID") != "custom-id-42" {
		t.Errorf("expected X-Request-ID 'custom-id-42', got %q", w.Header().Get("X-Request-ID"))
	}
}

func TestRequestIDMiddlewareGeneratesID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	id := w.Header().Get("X-Request-ID")
	if id == "" {
		t.Error("expected X-Request-ID header to be set")
	}
	if len(id) != 16 {
		t.Errorf("expected 16-char hex request ID, got %q (len=%d)", id, len(id))
	}
}

package network

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const jwtIssuer = "agentgame"

type rateLimiter struct {
	visitors    map[string]*visitor
	mu          sync.Mutex
	rate        float64
	burst       float64
	now         func() time.Time
	lastCleanup time.Time
}

type visitor struct {
	tokens     float64
	lastRefill time.Time
	lastSeen   time.Time
}

func newRateLimiter(rate, burst int) *rateLimiter {
	return newRateLimiterWithClock(rate, burst, time.Now)
}

func newRateLimiterWithClock(rate, burst int, now func() time.Time) *rateLimiter {
	if rate <= 0 {
		rate = 1
	}
	if burst <= 0 {
		burst = 1
	}
	if now == nil {
		now = time.Now
	}
	current := now()
	return &rateLimiter{
		visitors:    make(map[string]*visitor),
		rate:        float64(rate),
		burst:       float64(burst),
		now:         now,
		lastCleanup: current,
	}
}

func (rl *rateLimiter) allow(ip string) bool {
	return rl.allowAt(ip, rl.now())
}

func (rl *rateLimiter) allowAt(ip string, now time.Time) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if now.Sub(rl.lastCleanup) >= 3*time.Minute {
		for key, entry := range rl.visitors {
			if now.Sub(entry.lastSeen) >= 3*time.Minute {
				delete(rl.visitors, key)
			}
		}
		rl.lastCleanup = now
	}

	entry, exists := rl.visitors[ip]
	if !exists {
		rl.visitors[ip] = &visitor{
			tokens:     rl.burst - 1,
			lastRefill: now,
			lastSeen:   now,
		}
		return true
	}

	elapsed := now.Sub(entry.lastRefill).Seconds()
	if elapsed > 0 {
		entry.tokens = math.Min(rl.burst, entry.tokens+elapsed*rl.rate)
		entry.lastRefill = now
	}
	entry.lastSeen = now

	if entry.tokens < 1 {
		return false
	}
	entry.tokens--
	return true
}

func RateLimitMiddleware(rate, burst int) gin.HandlerFunc {
	rl := newRateLimiter(rate, burst)
	return func(c *gin.Context) {
		if !rl.allow(c.ClientIP()) {
			c.Header("Retry-After", "1")
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": gin.H{
					"code":    "RATE_LIMITED",
					"message": "too many requests, please try again later",
				},
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func CORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	originSet := make(map[string]bool, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		originSet[origin] = true
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if originSet[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

type JWTClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(secret, username, role string, expiryHours int) (string, error) {
	claims := JWTClaims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiryHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    jwtIssuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			respondUnauthorized(c, "missing authorization header")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
			respondUnauthorized(c, "invalid authorization format")
			return
		}

		claims := &JWTClaims{}
		token, err := jwt.ParseWithClaims(
			parts[1],
			claims,
			func(token *jwt.Token) (interface{}, error) {
				if token.Method != jwt.SigningMethodHS256 {
					return nil, fmt.Errorf("unexpected signing method: %s", token.Method.Alg())
				}
				return []byte(jwtSecret), nil
			},
			jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
			jwt.WithIssuer(jwtIssuer),
		)
		if err != nil || !token.Valid || claims.Username == "" {
			respondUnauthorized(c, "invalid or expired token")
			return
		}

		c.Set("gm_username", claims.Username)
		c.Set("gm_role", claims.Role)
		c.Next()
	}
}

func ManagementAuthMiddleware(jwtSecret string) gin.HandlerFunc {
	authenticate := AuthMiddleware(jwtSecret)
	return func(c *gin.Context) {
		if !isManagementPath(c.Request.URL.Path) {
			c.Next()
			return
		}
		authenticate(c)
	}
}

func isManagementPath(path string) bool {
	if path == "/api/gm/login" {
		return false
	}
	if path == "/mcp" || strings.HasPrefix(path, "/mcp/") {
		return true
	}

	managementPrefixes := []string{
		"/api/gm/",
		"/api/generator/",
		"/api/mcp/",
		"/api/scenes",
		"/api/npcs",
		"/api/agents",
		"/api/llm/",
		"/api/prompts",
		"/api/shops",
		"/api/items",
		"/api/tasks",
		"/api/flows",
		"/api/players",
		"/api/conversations",
		"/api/config",
		"/api/export",
		"/api/import",
	}
	for _, prefix := range managementPrefixes {
		if strings.HasSuffix(prefix, "/") {
			if strings.HasPrefix(path, prefix) {
				return true
			}
			continue
		}
		if path == prefix || strings.HasPrefix(path, prefix+"/") {
			return true
		}
	}
	return false
}

func respondUnauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"error": gin.H{
			"code":    "UNAUTHORIZED",
			"message": message,
		},
	})
	c.Abort()
}

func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			reqID = generateRequestID()
		}
		ctx := context.WithValue(c.Request.Context(), requestIDKey, reqID)
		c.Request = c.Request.WithContext(ctx)
		c.Set("request_id", reqID)
		c.Header("X-Request-ID", reqID)
		c.Next()
	}
}

func RequestLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		logger := LoggerFromContext(c.Request.Context())
		attrs := []slog.Attr{
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.Int("status", status),
			slog.Duration("latency", latency),
			slog.String("client_ip", c.ClientIP()),
		}
		if query != "" {
			attrs = append(attrs, slog.String("query", query))
		}
		if len(c.Errors) > 0 {
			attrs = append(attrs, slog.String("errors", c.Errors.String()))
		}

		level := slog.LevelInfo
		if status >= 500 {
			level = slog.LevelError
		} else if status >= 400 {
			level = slog.LevelWarn
		}
		logger.LogAttrs(c.Request.Context(), level, "HTTP request", attrs...)
	}
}

func generateRequestID() string {
	var buf [8]byte
	_, _ = rand.Read(buf[:])
	return hex.EncodeToString(buf[:])
}

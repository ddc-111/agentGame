package network

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type rateLimiter struct {
	visitors map[string]*visitor
	mu       sync.Mutex
	rate     int
	burst    int
}

type visitor struct {
	tokens   int
	lastSeen time.Time
}

func newRateLimiter(rate, burst int) *rateLimiter {
	rl := &rateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		burst:    burst,
	}
	go rl.cleanup()
	return rl
}

func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists || time.Since(v.lastSeen) > time.Minute {
		rl.visitors[ip] = &visitor{tokens: rl.burst - 1, lastSeen: time.Now()}
		return true
	}

	if v.tokens <= 0 {
		return false
	}

	v.tokens--
	v.lastSeen = time.Now()
	return true
}

func (rl *rateLimiter) cleanup() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func RateLimitMiddleware(rate, burst int) gin.HandlerFunc {
	rl := newRateLimiter(rate, burst)
	return func(c *gin.Context) {
		if !rl.allow(c.ClientIP()) {
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
	for _, o := range allowedOrigins {
		originSet[o] = true
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
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
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
			Issuer:    "agentgame",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "missing authorization header",
				},
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "invalid authorization format",
				},
			})
			c.Abort()
			return
		}

		claims := &JWTClaims{}
		token, err := jwt.ParseWithClaims(parts[1], claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "invalid or expired token",
				},
			})
			c.Abort()
			return
		}

		c.Set("gm_username", claims.Username)
		c.Set("gm_role", claims.Role)
		c.Next()
	}
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
		c.Header("X-Request-ID", reqID)
		c.Next()
	}
}

func generateRequestID() string {
	var buf [8]byte
	_, _ = rand.Read(buf[:])
	return hex.EncodeToString(buf[:])
}

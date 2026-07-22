package network

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func TestRateLimiterRefillsContinuously(t *testing.T) {
	now := time.Date(2026, time.July, 22, 0, 0, 0, 0, time.UTC)
	limiter := newRateLimiterWithClock(2, 2, func() time.Time { return now })

	if !limiter.allow("client") || !limiter.allow("client") {
		t.Fatal("expected initial burst to be allowed")
	}
	if limiter.allow("client") {
		t.Fatal("expected burst to be exhausted")
	}

	now = now.Add(500 * time.Millisecond)
	if !limiter.allow("client") {
		t.Fatal("expected one token to refill after 500ms at two tokens per second")
	}
	if limiter.allow("client") {
		t.Fatal("expected refilled token to be consumed")
	}
}

func TestManagementAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	const secret = "01234567890123456789012345678901"

	router := gin.New()
	router.Use(ManagementAuthMiddleware(secret))
	router.GET("/api/scenes", func(c *gin.Context) { c.Status(http.StatusOK) })
	router.GET("/api/game/init", func(c *gin.Context) { c.Status(http.StatusOK) })
	router.POST("/mcp", func(c *gin.Context) { c.Status(http.StatusOK) })
	router.POST("/api/gm/login", func(c *gin.Context) { c.Status(http.StatusOK) })

	assertStatus(t, router, http.MethodGet, "/api/game/init", "", http.StatusOK)
	assertStatus(t, router, http.MethodPost, "/api/gm/login", "", http.StatusOK)
	assertStatus(t, router, http.MethodGet, "/api/scenes", "", http.StatusUnauthorized)
	assertStatus(t, router, http.MethodPost, "/mcp", "", http.StatusUnauthorized)

	token, err := GenerateJWT(secret, "admin", "gm", 1)
	if err != nil {
		t.Fatalf("GenerateJWT() error = %v", err)
	}
	assertStatus(t, router, http.MethodGet, "/api/scenes", token, http.StatusOK)
	assertStatus(t, router, http.MethodPost, "/mcp", token, http.StatusOK)
}

func TestAuthMiddlewareRejectsWrongAlgorithmAndIssuer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	const secret = "01234567890123456789012345678901"

	router := gin.New()
	router.Use(AuthMiddleware(secret))
	router.GET("/protected", func(c *gin.Context) { c.Status(http.StatusOK) })

	wrongAlgorithm := jwt.NewWithClaims(jwt.SigningMethodHS384, JWTClaims{
		Username: "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			Issuer:    jwtIssuer,
		},
	})
	wrongAlgorithmToken, err := wrongAlgorithm.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign wrong-algorithm token: %v", err)
	}
	assertStatus(t, router, http.MethodGet, "/protected", wrongAlgorithmToken, http.StatusUnauthorized)

	wrongIssuer := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		Username: "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			Issuer:    "other-service",
		},
	})
	wrongIssuerToken, err := wrongIssuer.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign wrong-issuer token: %v", err)
	}
	assertStatus(t, router, http.MethodGet, "/protected", wrongIssuerToken, http.StatusUnauthorized)
}

func assertStatus(t *testing.T, handler http.Handler, method, path, token string, expected int) {
	t.Helper()
	req := httptest.NewRequest(method, path, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	if resp.Code != expected {
		t.Fatalf("%s %s status = %d, want %d; body=%s", method, path, resp.Code, expected, resp.Body.String())
	}
}

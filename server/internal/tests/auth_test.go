package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/ddc-111/agentGame/server/internal/network"
	"github.com/golang-jwt/jwt/v5"
)

func TestGMLoginSuccess(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]string{
		"username": "admin",
		"password": "admin123",
	}
	resp, err := makeRequest("POST", ts.URL+"/api/gm/login", body)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if result["code"] != float64(0) {
		t.Errorf("期望 code=0, 得到 %v", result["code"])
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		t.Fatal("响应中缺少 data 字段")
	}

	token, ok := data["token"].(string)
	if !ok || token == "" {
		t.Fatal("响应中缺少 token")
	}

	if token == "gm-token-placeholder" {
		t.Error("token 不应再是占位符")
	}
}

func TestGMLoginWrongPassword(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]string{
		"username": "admin",
		"password": "wrong",
	}
	resp, err := makeRequest("POST", ts.URL+"/api/gm/login", body)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusUnauthorized)
}

func TestGMLoginMissingFields(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]string{
		"username": "admin",
	}
	resp, err := makeRequest("POST", ts.URL+"/api/gm/login", body)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusBadRequest)
}

func TestGMProtectedEndpointWithoutToken(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	req, err := http.NewRequest(http.MethodGet, ts.URL+"/api/gm/me", nil)
	if err != nil {
		t.Fatalf("创建请求失败: %v", err)
	}
	req.Header.Set(skipTestAuthHeader, "true")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusUnauthorized)
}

func TestGMProtectedEndpointWithToken(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	loginBody := map[string]string{
		"username": "admin",
		"password": "admin123",
	}
	loginResp, err := makeRequest("POST", ts.URL+"/api/gm/login", loginBody)
	if err != nil {
		t.Fatalf("登录请求失败: %v", err)
	}
	defer loginResp.Body.Close()

	var loginResult map[string]interface{}
	json.NewDecoder(loginResp.Body).Decode(&loginResult)
	token := loginResult["data"].(map[string]interface{})["token"].(string)

	req, _ := http.NewRequest("GET", ts.URL+"/api/gm/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	data := result["data"].(map[string]interface{})
	if data["username"] != "admin" {
		t.Errorf("期望 username=admin, 得到 %v", data["username"])
	}
	if data["role"] != "gm" {
		t.Errorf("期望 role=gm, 得到 %v", data["role"])
	}
}

func TestGMTokenClaims(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]string{
		"username": "admin",
		"password": "admin123",
	}
	resp, err := makeRequest("POST", ts.URL+"/api/gm/login", body)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	tokenString := result["data"].(map[string]interface{})["token"].(string)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("change-me-in-production"), nil
	})
	if err != nil || !token.Valid {
		t.Fatalf("token 验证失败: %v", err)
	}

	claims := token.Claims.(jwt.MapClaims)
	if claims["username"] != "admin" {
		t.Errorf("期望 username claim=admin, 得到 %v", claims["username"])
	}
	if claims["role"] != "gm" {
		t.Errorf("期望 role claim=gm, 得到 %v", claims["role"])
	}
	if claims["iss"] != "agentgame" {
		t.Errorf("期望 issuer=agentgame, 得到 %v", claims["iss"])
	}
	if _, ok := claims["exp"]; !ok {
		t.Error("token 缺少 exp claim")
	}
}

func TestGenerateJWT(t *testing.T) {
	token, err := network.GenerateJWT("test-secret", "testuser", "gm", 1)
	if err != nil {
		t.Fatalf("生成 JWT 失败: %v", err)
	}
	if token == "" {
		t.Fatal("生成的 token 为空")
	}

	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})
	if err != nil || !parsed.Valid {
		t.Fatalf("解析生成的 JWT 失败: %v", err)
	}
}

func TestAuthMiddlewareExpiredToken(t *testing.T) {
	claims := network.JWTClaims{
		Username: "admin",
		Role:     "gm",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			Issuer:    "agentgame",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("change-me-in-production"))
	if err != nil {
		t.Fatalf("签名失败: %v", err)
	}

	ts := setupTestServer()
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL+"/api/gm/me", bytes.NewReader(nil))
	req.Header.Set("Authorization", "Bearer "+tokenString)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusUnauthorized)
}

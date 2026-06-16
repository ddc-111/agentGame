package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ddc-111/agentGame/server/internal/network"
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

	resp, err := http.Get(ts.URL + "/api/gm/me")
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

func TestGMProtectedEndpointWithInvalidToken(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL+"/api/gm/me", nil)
	req.Header.Set("Authorization", "Bearer invalid-token-value")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusUnauthorized)
}

func TestGenerateJWT(t *testing.T) {
	secret := "test-secret"
	token, err := network.GenerateJWT(secret, "admin", "gm", 24)
	if err != nil {
		t.Fatalf("生成 JWT 失败: %v", err)
	}

	if token == "" {
		t.Fatal("token 不应为空")
	}

	claims := &network.JWTClaims{}
	parsed, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		t.Fatalf("解析 JWT 失败: %v", err)
	}
	if !parsed.Valid {
		t.Fatal("token 应该有效")
	}
	if claims.Username != "admin" {
		t.Errorf("期望 Username=admin, 得到 %s", claims.Username)
	}
	if claims.Role != "gm" {
		t.Errorf("期望 Role=gm, 得到 %s", claims.Role)
	}
	if claims.Issuer != "agentgame" {
		t.Errorf("期望 Issuer=agentgame, 得到 %s", claims.Issuer)
	}
}

func TestGenerateJWTExpiry(t *testing.T) {
	secret := "test-secret"
	token, err := network.GenerateJWT(secret, "admin", "gm", 1)
	if err != nil {
		t.Fatalf("生成 JWT 失败: %v", err)
	}

	claims := &network.JWTClaims{}
	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		t.Fatalf("解析 JWT 失败: %v", err)
	}

	expectedExpiry := time.Now().Add(1 * time.Hour)
	if claims.ExpiresAt.Time.Sub(expectedExpiry) > time.Minute {
		t.Errorf("过期时间偏差过大: 期望 ~%v, 得到 %v", expectedExpiry, claims.ExpiresAt.Time)
	}
}

func TestCORSMiddlewareAllowedOrigin(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	req, _ := http.NewRequest("OPTIONS", ts.URL+"/health", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, 204)
	allowOrigin := resp.Header.Get("Access-Control-Allow-Origin")
	if allowOrigin != "http://localhost:5173" {
		t.Errorf("期望 Access-Control-Allow-Origin=http://localhost:5173, 得到 %s", allowOrigin)
	}
}

func TestCORSMiddlewareDisallowedOrigin(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL+"/health", nil)
	req.Header.Set("Origin", "http://evil.com")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	allowOrigin := resp.Header.Get("Access-Control-Allow-Origin")
	if allowOrigin == "http://evil.com" {
		t.Error("不应允许未配置的 origin")
	}
}

func TestCORSHeadersPresent(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	req, _ := http.NewRequest("OPTIONS", ts.URL+"/health", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	methods := resp.Header.Get("Access-Control-Allow-Methods")
	if methods == "" {
		t.Error("缺少 Access-Control-Allow-Methods header")
	}
	headers := resp.Header.Get("Access-Control-Allow-Headers")
	if headers == "" {
		t.Error("缺少 Access-Control-Allow-Headers header")
	}
	credentials := resp.Header.Get("Access-Control-Allow-Credentials")
	if credentials != "true" {
		t.Errorf("期望 Access-Control-Allow-Credentials=true, 得到 %s", credentials)
	}
}

func TestGMLoginReturnsValidJWT(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]string{
		"username": "admin",
		"password": "admin123",
	}
	loginResp, err := makeRequest("POST", ts.URL+"/api/gm/login", body)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer loginResp.Body.Close()

	var loginResult map[string]interface{}
	json.NewDecoder(loginResp.Body).Decode(&loginResult)
	token := loginResult["data"].(map[string]interface{})["token"].(string)

	req, _ := http.NewRequest("GET", ts.URL+"/api/gm/me", bytes.NewBuffer(nil))
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
		t.Errorf("JWT 中的 username 不匹配: 期望 admin, 得到 %v", data["username"])
	}
}

package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"
)

func TestPaginationDefaults(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	endpoints := []string{
		"/api/scenes",
		"/api/npcs",
		"/api/agents",
		"/api/items",
		"/api/tasks",
		"/api/flows",
		"/api/shops",
		"/api/skills",
		"/api/players",
		"/api/llm/providers",
		"/api/prompts",
	}

	for _, endpoint := range endpoints {
		t.Run(endpoint, func(t *testing.T) {
			resp, err := http.Get(ts.URL + endpoint)
			if err != nil {
				t.Fatalf("请求失败: %v", err)
			}
			defer resp.Body.Close()

			assertStatusCode(t, resp.StatusCode, http.StatusOK)

			totalHeader := resp.Header.Get("X-Total-Count")
			if totalHeader == "" {
				t.Errorf("%s 缺少 X-Total-Count 响应头", endpoint)
			}

			if totalHeader != "" {
				_, err := strconv.ParseInt(totalHeader, 10, 64)
				if err != nil {
					t.Errorf("%s X-Total-Count 不是有效数字: %v", endpoint, err)
				}
			}
		})
	}
}

func TestPaginationCustomParams(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/scenes?page=1&page_size=2")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	totalHeader := resp.Header.Get("X-Total-Count")
	if totalHeader == "" {
		t.Error("缺少 X-Total-Count 响应头")
	}

	var result map[string]interface{}
	if err := parseResponse(resp, &result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if _, ok := result["total"]; !ok {
		t.Error("响应缺少 total 字段")
	}

	data, ok := result["data"].([]interface{})
	if !ok {
		t.Fatal("响应 data 字段格式错误")
	}

	total, _ := strconv.ParseInt(totalHeader, 10, 64)
	if total > 2 && len(data) > 2 {
		t.Errorf("page_size=2 但返回了 %d 条数据", len(data))
	}
}

func TestPaginationSecondPage(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp1, err := http.Get(ts.URL + "/api/scenes?page=1&page_size=1")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp1.Body.Close()

	var result1 map[string]interface{}
	parseResponse(resp1, &result1)
	data1 := result1["data"].([]interface{})

	totalHeader := resp1.Header.Get("X-Total-Count")
	total, _ := strconv.ParseInt(totalHeader, 10, 64)
	if total <= 1 {
		t.Skip("种子数据不足2条，跳过分页测试")
	}

	resp2, err := http.Get(ts.URL + "/api/scenes?page=2&page_size=1")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp2.Body.Close()

	var result2 map[string]interface{}
	parseResponse(resp2, &result2)
	data2 := result2["data"].([]interface{})

	if len(data1) > 0 && len(data2) > 0 {
		id1 := data1[0].(map[string]interface{})["id"]
		id2 := data2[0].(map[string]interface{})["id"]
		if id1 == id2 {
			t.Error("第1页和第2页返回了相同的数据")
		}
	}
}

func TestPaginationPageSizeCap(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/scenes?page=1&page_size=200")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := parseResponse(resp, &result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	data := result["data"].([]interface{})
	if len(data) > 100 {
		t.Errorf("page_size 被限制为100，但返回了 %d 条数据", len(data))
	}
}

func TestPaginationInvalidParams(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/scenes?page=-1&page_size=0")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	totalHeader := resp.Header.Get("X-Total-Count")
	if totalHeader == "" {
		t.Error("无效参数时缺少 X-Total-Count 响应头")
	}
}

func TestPaginationConversations(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/conversations?page=1&page_size=10")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	totalHeader := resp.Header.Get("X-Total-Count")
	if totalHeader == "" {
		t.Error("conversations 缺少 X-Total-Count 响应头")
	}

	var result map[string]interface{}
	if err := parseResponse(resp, &result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if _, ok := result["total"]; !ok {
		t.Error("conversations 响应缺少 total 字段")
	}
}

func TestPaginationAchievements(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	createResp, err := makeRequest("POST", ts.URL+"/api/player/create", TestPlayer)
	if err != nil {
		t.Fatalf("创建玩家失败: %v", err)
	}
	defer createResp.Body.Close()

	var createResult map[string]interface{}
	json.NewDecoder(createResp.Body).Decode(&createResult)
	playerData := createResult["data"].(map[string]interface{})
	playerID := uint(playerData["id"].(float64))

	resp, err := http.Get(fmt.Sprintf("%s/api/achievements/%d?page=1&page_size=5", ts.URL, playerID))
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	totalHeader := resp.Header.Get("X-Total-Count")
	if totalHeader == "" {
		t.Error("achievements 缺少 X-Total-Count 响应头")
	}
}

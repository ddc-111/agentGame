package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/config"
	"github.com/ddc-111/agentGame/server/internal/network"
)

// TestPlayer 测试玩家数据
var TestPlayer = map[string]interface{}{
	"name":    "测试玩家",
	"account": "test_player_001",
}

// TestPlayer2 测试玩家2数据
var TestPlayer2 = map[string]interface{}{
	"name":    "测试玩家二号",
	"account": "test_player_002",
}

// TestNPC 测试NPC数据
var TestNPC = map[string]interface{}{
	"name":        "测试NPC",
	"code":        "npc_test",
	"title":       "测试称号",
	"description": "这是一个测试NPC",
}

// TestScene 测试场景数据
var TestScene = map[string]interface{}{
	"name":        "测试场景",
	"code":        "scene_test",
	"description": "这是一个测试场景",
	"width":       1920,
	"height":      1080,
}

// TestAgent 测试智能体数据
var TestAgent = map[string]interface{}{
	"name":           "测试智能体",
	"code":           "agent_test",
	"description":    "这是一个测试智能体",
	"system_prompt":  "你是一个测试智能体",
	"llm_model":      "gpt-4",
}

// TestShop 测试商店数据
var TestShop = map[string]interface{}{
	"name":        "测试商店",
	"code":        "shop_test",
	"type":        "general",
	"description": "这是一个测试商店",
	"owner_npc":   "npc_test",
}

// TestItem 测试道具数据
var TestItem = map[string]interface{}{
	"name":        "测试道具",
	"code":        "item_test",
	"category":    "consumable",
	"description": "这是一个测试道具",
}

// TestTask 测试任务数据
var TestTask = map[string]interface{}{
	"name":        "测试任务",
	"code":        "task_test",
	"type":        "main",
	"description": "这是一个测试任务",
}

// TestFlow 测试流程数据
var TestFlow = map[string]interface{}{
	"name":        "测试流程",
	"code":        "flow_test",
	"description": "这是一个测试流程",
}

// TestTemplate 测试提示词模板数据
var TestTemplate = map[string]interface{}{
	"name":    "测试模板",
	"code":    "template_test",
	"content": "这是一个测试提示词模板",
}

// MCPRequest JSON-RPC请求结构
type MCPRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// setupTestRouter 创建测试用的Gin路由
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	cfg := config.Default()
	cfg.Database.DSN = "file::memory:?cache=shared"
	cfg.Database.Driver = "sqlite"

	server := network.NewServer(cfg)
	return server.GetRouter()
}

// setupTestServer 创建测试用的HTTP服务器
func setupTestServer() *httptest.Server {
	router := setupTestRouter()
	return httptest.NewServer(router)
}

// makeRequest 发送HTTP请求并返回响应
func makeRequest(method, url string, body interface{}) (*http.Response, error) {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return client.Do(req)
}

// parseResponse 解析JSON响应
func parseResponse(resp *http.Response, result interface{}) error {
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(result)
}

// assertStatusCode 断言状态码
func assertStatusCode(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("期望状态码 %d, 得到 %d", want, got)
	}
}

// createMCPRequestBody 创建MCP请求体
func createMCPRequestBody(method string, params interface{}) MCPRequest {
	return MCPRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  method,
		Params:  params,
	}
}

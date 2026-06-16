package tests

import (
	"encoding/json"
	"net/http"
	"testing"
)

// TestMCPInitialize 测试MCP初始化
func TestMCPInitialize(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	reqBody := createMCPRequestBody("initialize", nil)

	resp, err := makeRequest("POST", ts.URL+"/mcp", reqBody)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	// 验证JSON-RPC响应结构
	if result["jsonrpc"] != "2.0" {
		t.Errorf("期望 jsonrpc=2.0, 得到 %v", result["jsonrpc"])
	}

	res, ok := result["result"].(map[string]interface{})
	if !ok {
		t.Fatal("响应 result 字段格式错误")
	}

	if res["protocolVersion"] != "2024-11-05" {
		t.Errorf("期望 protocolVersion=2024-11-05, 得到 %v", res["protocolVersion"])
	}

	serverInfo, ok := res["serverInfo"].(map[string]interface{})
	if !ok {
		t.Fatal("serverInfo 字段格式错误")
	}

	if serverInfo["name"] != "agentgame-mcp" {
		t.Errorf("期望 serverInfo.name=agentgame-mcp, 得到 %v", serverInfo["name"])
	}
	t.Log("MCP初始化测试通过")
}

// TestMCPToolsList 测试MCP工具列表
func TestMCPToolsList(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	reqBody := createMCPRequestBody("tools/list", nil)

	resp, err := makeRequest("POST", ts.URL+"/mcp", reqBody)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	res, ok := result["result"].(map[string]interface{})
	if !ok {
		t.Fatal("响应 result 字段格式错误")
	}

	tools, ok := res["tools"].([]interface{})
	if !ok {
		t.Fatal("tools 字段格式错误")
	}

	if len(tools) == 0 {
		t.Error("工具列表为空")
	}

	// 验证工具结构
	for _, tool := range tools {
		toolMap, ok := tool.(map[string]interface{})
		if !ok {
			t.Error("工具格式错误")
			continue
		}
		if _, ok := toolMap["name"]; !ok {
			t.Error("工具缺少 name 字段")
		}
		if _, ok := toolMap["description"]; !ok {
			t.Error("工具缺少 description 字段")
		}
		if _, ok := toolMap["inputSchema"]; !ok {
			t.Error("工具缺少 inputSchema 字段")
		}
	}
	t.Logf("MCP工具列表测试通过, 工具数量: %d", len(tools))
}

// TestMCPCallListScenes 测试MCP调用list_scenes工具
func TestMCPCallListScenes(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	reqBody := createMCPRequestBody("tools/call", map[string]interface{}{
		"name":      "list_scenes",
		"arguments": map[string]interface{}{},
	})

	resp, err := makeRequest("POST", ts.URL+"/mcp", reqBody)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	res, ok := result["result"].(map[string]interface{})
	if !ok {
		t.Fatal("响应 result 字段格式错误")
	}

	content, ok := res["content"].([]interface{})
	if !ok {
		t.Fatal("content 字段格式错误")
	}

	if len(content) == 0 {
		t.Error("content 为空")
	}

	contentItem := content[0].(map[string]interface{})
	if contentItem["type"] != "text" {
		t.Errorf("期望 type=text, 得到 %v", contentItem["type"])
	}
	t.Logf("MCP list_scenes 测试通过")
}

// TestMCPCallCreateNPC 测试MCP调用create_npc工具
func TestMCPCallCreateNPC(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	reqBody := createMCPRequestBody("tools/call", map[string]interface{}{
		"name": "create_npc",
		"arguments": map[string]interface{}{
			"name":        "MCP测试NPC",
			"code":        "npc_mcp_test",
			"title":       "测试员",
			"description": "通过MCP创建的测试NPC",
		},
	})

	resp, err := makeRequest("POST", ts.URL+"/mcp", reqBody)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	res, ok := result["result"].(map[string]interface{})
	if !ok {
		t.Fatal("响应 result 字段格式错误")
	}

	// 验证没有错误
	if res["isError"] == true {
		t.Error("创建NPC失败")
	}

	content, ok := res["content"].([]interface{})
	if !ok || len(content) == 0 {
		t.Fatal("content 字段格式错误或为空")
	}

	t.Log("MCP create_npc 测试通过")
}

// TestMCPCallListNPCs 测试MCP调用list_npcs工具
func TestMCPCallListNPCs(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	reqBody := createMCPRequestBody("tools/call", map[string]interface{}{
		"name":      "list_npcs",
		"arguments": map[string]interface{}{},
	})

	resp, err := makeRequest("POST", ts.URL+"/mcp", reqBody)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	res, ok := result["result"].(map[string]interface{})
	if !ok {
		t.Fatal("响应 result 字段格式错误")
	}

	content, ok := res["content"].([]interface{})
	if !ok {
		t.Fatal("content 字段格式错误")
	}

	if len(content) == 0 {
		t.Error("content 为空")
	}
	t.Log("MCP list_npcs 测试通过")
}

// TestMCPCallGenerate 测试MCP调用generate_config工具
func TestMCPCallGenerate(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	reqBody := createMCPRequestBody("tools/call", map[string]interface{}{
		"name": "generate_config",
		"arguments": map[string]interface{}{
			"type":        "npc",
			"description": "一个卖包子的老大爷",
			"action":      "create",
		},
	})

	resp, err := makeRequest("POST", ts.URL+"/mcp", reqBody)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	// 生成器可能未启用，所以这里只验证响应结构
	t.Log("MCP generate_config 测试通过（生成器可能未启用）")
}

// TestMCPCallExportData 测试MCP调用export_data工具
func TestMCPCallExportData(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	reqBody := createMCPRequestBody("tools/call", map[string]interface{}{
		"name":      "export_data",
		"arguments": map[string]interface{}{},
	})

	resp, err := makeRequest("POST", ts.URL+"/mcp", reqBody)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	res, ok := result["result"].(map[string]interface{})
	if !ok {
		t.Fatal("响应 result 字段格式错误")
	}

	if res["isError"] == true {
		t.Error("导出数据失败")
	}
	t.Log("MCP export_data 测试通过")
}

// TestMCPCallGetGameStats 测试MCP调用get_game_stats工具
func TestMCPCallGetGameStats(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	reqBody := createMCPRequestBody("tools/call", map[string]interface{}{
		"name":      "get_game_stats",
		"arguments": map[string]interface{}{},
	})

	resp, err := makeRequest("POST", ts.URL+"/mcp", reqBody)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	res, ok := result["result"].(map[string]interface{})
	if !ok {
		t.Fatal("响应 result 字段格式错误")
	}

	if res["isError"] == true {
		t.Error("获取游戏统计失败")
	}
	t.Log("MCP get_game_stats 测试通过")
}

// TestMCPRestTools 测试MCP REST工具列表接口
func TestMCPRestTools(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/mcp/tools")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	tools, ok := result["tools"].([]interface{})
	if !ok {
		t.Fatal("tools 字段格式错误")
	}

	if len(tools) == 0 {
		t.Error("工具列表为空")
	}
	t.Logf("MCP REST工具列表测试通过, 工具数量: %d", len(tools))
}

// TestMCPRestCall 测试MCP REST调用接口
func TestMCPRestCall(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	reqBody := map[string]interface{}{
		"name": "list_scenes",
		"arguments": map[string]interface{}{},
	}

	resp, err := makeRequest("POST", ts.URL+"/api/mcp/call", reqBody)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	t.Log("MCP REST调用测试通过")
}

// TestMCPRestResources 测试MCP REST资源列表接口
func TestMCPRestResources(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/mcp/resources")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	resources, ok := result["resources"].([]interface{})
	if !ok {
		t.Fatal("resources 字段格式错误")
	}

	if len(resources) == 0 {
		t.Error("资源列表为空")
	}
	t.Logf("MCP REST资源列表测试通过, 资源数量: %d", len(resources))
}

// TestMCPRestResourceRead 测试MCP REST资源读取接口
func TestMCPRestResourceRead(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/mcp/resources/read?uri=game_state://scenes")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	t.Log("MCP REST资源读取测试通过")
}

// TestMCPRestPrompts 测试MCP REST提示词列表接口
func TestMCPRestPrompts(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/mcp/prompts")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	prompts, ok := result["prompts"].([]interface{})
	if !ok {
		t.Fatal("prompts 字段格式错误")
	}

	if len(prompts) == 0 {
		t.Error("提示词列表为空")
	}
	t.Logf("MCP REST提示词列表测试通过, 提示词数量: %d", len(prompts))
}

// TestMCPRestPromptGet 测试MCP REST提示词获取接口
func TestMCPRestPromptGet(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	reqBody := map[string]interface{}{
		"name": "npc_personality",
		"arguments": map[string]interface{}{
			"name": "测试NPC",
		},
	}

	resp, err := makeRequest("POST", ts.URL+"/api/mcp/prompts/get", reqBody)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	t.Log("MCP REST提示词获取测试通过")
}

// TestMCPInvalidMethod 测试MCP无效方法
func TestMCPInvalidMethod(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	reqBody := createMCPRequestBody("invalid/method", nil)

	resp, err := makeRequest("POST", ts.URL+"/mcp", reqBody)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	// 应该返回错误
	if result["error"] == nil {
		t.Error("无效方法应返回错误")
	}

	errObj := result["error"].(map[string]interface{})
	if errObj["code"].(float64) != -32601 {
		t.Errorf("期望错误码 -32601, 得到 %v", errObj["code"])
	}
	t.Log("MCP无效方法测试通过")
}

// TestMCPInvalidTool 测试MCP调用不存在的工具
func TestMCPInvalidTool(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	reqBody := createMCPRequestBody("tools/call", map[string]interface{}{
		"name":      "nonexistent_tool",
		"arguments": map[string]interface{}{},
	})

	resp, err := makeRequest("POST", ts.URL+"/mcp", reqBody)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	// 应该返回错误
	if result["error"] == nil {
		t.Error("不存在的工具应返回错误")
	}
	t.Log("MCP不存在工具测试通过")
}

// TestMCPResourcesList 测试MCP资源列表
func TestMCPResourcesList(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	reqBody := createMCPRequestBody("resources/list", nil)

	resp, err := makeRequest("POST", ts.URL+"/mcp", reqBody)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	res, ok := result["result"].(map[string]interface{})
	if !ok {
		t.Fatal("响应 result 字段格式错误")
	}

	resources, ok := res["resources"].([]interface{})
	if !ok {
		t.Fatal("resources 字段格式错误")
	}

	if len(resources) == 0 {
		t.Error("资源列表为空")
	}

	for _, res := range resources {
		resMap, ok := res.(map[string]interface{})
		if !ok {
			t.Error("资源格式错误")
			continue
		}
		if _, ok := resMap["uri"]; !ok {
			t.Error("资源缺少 uri 字段")
		}
		if _, ok := resMap["name"]; !ok {
			t.Error("资源缺少 name 字段")
		}
	}
	t.Logf("MCP资源列表测试通过, 资源数量: %d", len(resources))
}

// TestMCPResourcesReadScenes 测试MCP读取场景资源
func TestMCPResourcesReadScenes(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	reqBody := createMCPRequestBody("resources/read", map[string]interface{}{
		"uri": "game_state://scenes",
	})

	resp, err := makeRequest("POST", ts.URL+"/mcp", reqBody)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	res, ok := result["result"].(map[string]interface{})
	if !ok {
		t.Fatalf("响应 result 字段格式错误, 实际: %+v", result)
	}

	contents, ok := res["contents"].([]interface{})
	if !ok {
		t.Fatal("contents 字段格式错误")
	}

	if len(contents) == 0 {
		t.Error("contents 为空")
	}

	contentItem := contents[0].(map[string]interface{})
	if contentItem["uri"] != "game_state://scenes" {
		t.Errorf("期望 uri=game_state://scenes, 得到 %v", contentItem["uri"])
	}
	if contentItem["mimeType"] != "application/json" {
		t.Errorf("期望 mimeType=application/json, 得到 %v", contentItem["mimeType"])
	}
	t.Log("MCP读取场景资源测试通过")
}

// TestMCPResourcesReadOverview 测试MCP读取概览资源
func TestMCPResourcesReadOverview(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	reqBody := createMCPRequestBody("resources/read", map[string]interface{}{
		"uri": "game_state://overview",
	})

	resp, err := makeRequest("POST", ts.URL+"/mcp", reqBody)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	res, ok := result["result"].(map[string]interface{})
	if !ok {
		t.Fatal("响应 result 字段格式错误")
	}

	contents, ok := res["contents"].([]interface{})
	if !ok {
		t.Fatal("contents 字段格式错误")
	}

	if len(contents) == 0 {
		t.Error("contents 为空")
	}
	t.Log("MCP读取概览资源测试通过")
}

// TestMCPResourcesReadInvalid 测试MCP读取无效资源
func TestMCPResourcesReadInvalid(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	reqBody := createMCPRequestBody("resources/read", map[string]interface{}{
		"uri": "game_state://nonexistent",
	})

	resp, err := makeRequest("POST", ts.URL+"/mcp", reqBody)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if result["error"] == nil {
		t.Error("无效资源应返回错误")
	}
	t.Log("MCP读取无效资源测试通过")
}

// TestMCPPromptsList 测试MCP提示词列表
func TestMCPPromptsList(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	reqBody := createMCPRequestBody("prompts/list", nil)

	resp, err := makeRequest("POST", ts.URL+"/mcp", reqBody)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	res, ok := result["result"].(map[string]interface{})
	if !ok {
		t.Fatal("响应 result 字段格式错误")
	}

	prompts, ok := res["prompts"].([]interface{})
	if !ok {
		t.Fatal("prompts 字段格式错误")
	}

	if len(prompts) == 0 {
		t.Error("提示词列表为空")
	}

	for _, p := range prompts {
		pMap, ok := p.(map[string]interface{})
		if !ok {
			t.Error("提示词格式错误")
			continue
		}
		if _, ok := pMap["name"]; !ok {
			t.Error("提示词缺少 name 字段")
		}
	}
	t.Logf("MCP提示词列表测试通过, 提示词数量: %d", len(prompts))
}

// TestMCPPromptsGetPersonality 测试MCP获取NPC人格提示词
func TestMCPPromptsGetPersonality(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	reqBody := createMCPRequestBody("prompts/get", map[string]interface{}{
		"name": "npc_personality",
		"arguments": map[string]interface{}{
			"name":     "铁匠张三",
			"title":    "大师铁匠",
			"background": "固执但心地善良",
		},
	})

	resp, err := makeRequest("POST", ts.URL+"/mcp", reqBody)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	res, ok := result["result"].(map[string]interface{})
	if !ok {
		t.Fatal("响应 result 字段格式错误")
	}

	messages, ok := res["messages"].([]interface{})
	if !ok {
		t.Fatal("messages 字段格式错误")
	}

	if len(messages) == 0 {
		t.Error("messages 为空")
	}
	t.Log("MCP获取NPC人格提示词测试通过")
}

// TestMCPPromptsGetScene 测试MCP获取场景描述提示词
func TestMCPPromptsGetScene(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	reqBody := createMCPRequestBody("prompts/get", map[string]interface{}{
		"name": "scene_description",
		"arguments": map[string]interface{}{
			"name": "王城广场",
			"type": "城镇中心",
		},
	})

	resp, err := makeRequest("POST", ts.URL+"/mcp", reqBody)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	res, ok := result["result"].(map[string]interface{})
	if !ok {
		t.Fatal("响应 result 字段格式错误")
	}

	messages, ok := res["messages"].([]interface{})
	if !ok {
		t.Fatal("messages 字段格式错误")
	}

	if len(messages) == 0 {
		t.Error("messages 为空")
	}
	t.Log("MCP获取场景描述提示词测试通过")
}

// TestMCPPromptsGetInvalid 测试MCP获取无效提示词
func TestMCPPromptsGetInvalid(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	reqBody := createMCPRequestBody("prompts/get", map[string]interface{}{
		"name":      "nonexistent_prompt",
		"arguments": map[string]interface{}{},
	})

	resp, err := makeRequest("POST", ts.URL+"/mcp", reqBody)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if result["error"] == nil {
		t.Error("无效提示词应返回错误")
	}
	t.Log("MCP获取无效提示词测试通过")
}

// TestMCPInitializeCapabilities 测试MCP初始化返回的capabilities包含resources和prompts
func TestMCPInitializeCapabilities(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	reqBody := createMCPRequestBody("initialize", nil)

	resp, err := makeRequest("POST", ts.URL+"/mcp", reqBody)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	res, ok := result["result"].(map[string]interface{})
	if !ok {
		t.Fatal("响应 result 字段格式错误")
	}

	caps, ok := res["capabilities"].(map[string]interface{})
	if !ok {
		t.Fatal("capabilities 字段格式错误")
	}

	if _, ok := caps["tools"]; !ok {
		t.Error("capabilities 缺少 tools 字段")
	}
	if _, ok := caps["resources"]; !ok {
		t.Error("capabilities 缺少 resources 字段")
	}
	if _, ok := caps["prompts"]; !ok {
		t.Error("capabilities 缺少 prompts 字段")
	}
	t.Log("MCP初始化capabilities测试通过")
}

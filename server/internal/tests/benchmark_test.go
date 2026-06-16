package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"
)

// BenchmarkHealthEndpoint 测试健康检查端点性能
func BenchmarkHealthEndpoint(b *testing.B) {
	ts := setupTestServer()
	defer ts.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := http.Get(ts.URL + "/health")
		if err != nil {
			b.Fatalf("请求失败: %v", err)
		}
		resp.Body.Close()
	}
}

// BenchmarkGameInit 测试游戏初始化数据加载性能
func BenchmarkGameInit(b *testing.B) {
	ts := setupTestServer()
	defer ts.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := http.Get(ts.URL + "/api/game/init")
		if err != nil {
			b.Fatalf("请求失败: %v", err)
		}
		resp.Body.Close()
	}
}

// BenchmarkGetScenes 测试获取场景列表性能
func BenchmarkGetScenes(b *testing.B) {
	ts := setupTestServer()
	defer ts.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := http.Get(ts.URL + "/api/scenes")
		if err != nil {
			b.Fatalf("请求失败: %v", err)
		}
		resp.Body.Close()
	}
}

// BenchmarkGetNPCs 测试获取NPC列表性能
func BenchmarkGetNPCs(b *testing.B) {
	ts := setupTestServer()
	defer ts.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := http.Get(ts.URL + "/api/npcs")
		if err != nil {
			b.Fatalf("请求失败: %v", err)
		}
		resp.Body.Close()
	}
}

// BenchmarkCreatePlayer 测试创建玩家性能
func BenchmarkCreatePlayer(b *testing.B) {
	ts := setupTestServer()
	defer ts.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		player := map[string]interface{}{
			"name":    fmt.Sprintf("性能测试玩家_%d", i),
			"account": fmt.Sprintf("bench_player_%d", i),
		}
		resp, err := makeRequest("POST", ts.URL+"/api/player/create", player)
		if err != nil {
			b.Fatalf("请求失败: %v", err)
		}
		resp.Body.Close()
	}
}

// BenchmarkNPCChat 测试NPC对话性能
func BenchmarkNPCChat(b *testing.B) {
	ts := setupTestServer()
	defer ts.Close()

	// 创建玩家
	createResp, err := makeRequest("POST", ts.URL+"/api/player/create", TestPlayer)
	if err != nil {
		b.Fatalf("创建玩家失败: %v", err)
	}
	defer createResp.Body.Close()

	var createResult map[string]interface{}
	json.NewDecoder(createResp.Body).Decode(&createResult)
	playerData := createResult["data"].(map[string]interface{})
	playerID := uint(playerData["id"].(float64))

	// 获取NPC
	npcsResp, err := http.Get(ts.URL + "/api/npcs")
	if err != nil {
		b.Fatalf("获取NPC失败: %v", err)
	}
	defer npcsResp.Body.Close()

	var npcsResult map[string]interface{}
	json.NewDecoder(npcsResp.Body).Decode(&npcsResult)
	npcs := npcsResult["data"].([]interface{})
	if len(npcs) == 0 {
		b.Skip("没有NPC数据")
	}
	npcID := uint(npcs[0].(map[string]interface{})["id"].(float64))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		chatData := map[string]interface{}{
			"player_id": playerID,
			"npc_id":    npcID,
			"message":   fmt.Sprintf("测试消息_%d", i),
		}
		resp, err := makeRequest("POST", ts.URL+"/api/npc/chat", chatData)
		if err != nil {
			b.Fatalf("请求失败: %v", err)
		}
		resp.Body.Close()
	}
}

// BenchmarkConcurrentPlayers 测试并发创建玩家性能
func BenchmarkConcurrentPlayers(b *testing.B) {
	ts := setupTestServer()
	defer ts.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.RunParallel(func(pb *testing.PB) {
			j := 0
			for pb.Next() {
				player := map[string]interface{}{
					"name":    fmt.Sprintf("并发测试玩家_%d_%d", i, j),
					"account": fmt.Sprintf("concurrent_player_%d_%d", i, j),
				}
				resp, err := makeRequest("POST", ts.URL+"/api/player/create", player)
				if err != nil {
					b.Fatalf("请求失败: %v", err)
				}
				resp.Body.Close()
				j++
			}
		})
	}
}

// BenchmarkConcurrentRequests 测试并发请求性能
func BenchmarkConcurrentRequests(b *testing.B) {
	ts := setupTestServer()
	defer ts.Close()

	endpoints := []string{
		"/health",
		"/api/scenes",
		"/api/npcs",
		"/api/tasks",
		"/api/items",
		"/api/shops",
		"/api/agents",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		for _, endpoint := range endpoints {
			wg.Add(1)
			go func(ep string) {
				defer wg.Done()
				resp, err := http.Get(ts.URL + ep)
				if err != nil {
					b.Errorf("请求 %s 失败: %v", ep, err)
					return
				}
				resp.Body.Close()
			}(endpoint)
		}
		wg.Wait()
	}
}

// BenchmarkMCPInitialize 测试MCP初始化性能
func BenchmarkMCPInitialize(b *testing.B) {
	ts := setupTestServer()
	defer ts.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reqBody := createMCPRequestBody("initialize", nil)
		resp, err := makeRequest("POST", ts.URL+"/mcp", reqBody)
		if err != nil {
			b.Fatalf("请求失败: %v", err)
		}
		resp.Body.Close()
	}
}

// BenchmarkMCPToolsList 测试MCP工具列表性能
func BenchmarkMCPToolsList(b *testing.B) {
	ts := setupTestServer()
	defer ts.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reqBody := createMCPRequestBody("tools/list", nil)
		resp, err := makeRequest("POST", ts.URL+"/mcp", reqBody)
		if err != nil {
			b.Fatalf("请求失败: %v", err)
		}
		resp.Body.Close()
	}
}

// BenchmarkMCPToolCall 测试MCP工具调用性能
func BenchmarkMCPToolCall(b *testing.B) {
	ts := setupTestServer()
	defer ts.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reqBody := createMCPRequestBody("tools/call", map[string]interface{}{
			"name":      "list_scenes",
			"arguments": map[string]interface{}{},
		})
		resp, err := makeRequest("POST", ts.URL+"/mcp", reqBody)
		if err != nil {
			b.Fatalf("请求失败: %v", err)
		}
		resp.Body.Close()
	}
}

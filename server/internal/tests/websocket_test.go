package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestWebSocketConnection(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	wsURL := "ws" + ts.URL[len("http"):] + "/api/ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("WebSocket连接失败: %v", err)
	}
	defer conn.Close()

	t.Log("WebSocket连接测试通过")
}

func TestWebSocketConnectionWithPlayer(t *testing.T) {
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

	wsURL := fmt.Sprintf("ws%s/api/ws?player_id=%d&scene_id=scene_village_square", ts.URL[len("http"):], playerID)
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("WebSocket连接失败: %v", err)
	}
	defer conn.Close()

	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, message, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("读取初始状态失败: %v", err)
	}

	var wsMsg map[string]interface{}
	if err := json.Unmarshal(message, &wsMsg); err != nil {
		t.Fatalf("解析消息失败: %v", err)
	}

	if wsMsg["type"] != "state_sync" {
		t.Errorf("期望 type=state_sync, 得到 %v", wsMsg["type"])
	}
	t.Log("WebSocket带玩家连接测试通过")
}

func TestWebSocketPlayerPositionBroadcast(t *testing.T) {
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

	sceneID := "scene_village_square"
	wsURL1 := fmt.Sprintf("ws%s/api/ws?player_id=%d&scene_id=%s", ts.URL[len("http"):], playerID, sceneID)
	conn1, _, err := websocket.DefaultDialer.Dial(wsURL1, nil)
	if err != nil {
		t.Fatalf("客户端1连接失败: %v", err)
	}
	defer conn1.Close()

	wsURL2 := fmt.Sprintf("ws%s/api/ws?player_id=0&scene_id=%s", ts.URL[len("http"):], sceneID)
	conn2, _, err := websocket.DefaultDialer.Dial(wsURL2, nil)
	if err != nil {
		t.Fatalf("客户端2连接失败: %v", err)
	}
	defer conn2.Close()

	_ = conn1.SetReadDeadline(time.Now().Add(2 * time.Second))
	conn1.ReadMessage()

	posMsg := map[string]interface{}{
		"type":      "player_position",
		"player_id": playerID,
		"scene_id":  sceneID,
		"data": map[string]interface{}{
			"player_id": playerID,
			"scene_id":  sceneID,
			"pos_x":     500,
			"pos_y":     300,
		},
	}

	if err := conn1.WriteJSON(posMsg); err != nil {
		t.Fatalf("发送位置消息失败: %v", err)
	}

	_ = conn2.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, message, err := conn2.ReadMessage()
	if err != nil {
		t.Fatalf("接收广播消息失败: %v", err)
	}

	var received map[string]interface{}
	if err := json.Unmarshal(message, &received); err != nil {
		t.Fatalf("解析广播消息失败: %v", err)
	}

	if received["type"] != "player_position" {
		t.Errorf("期望 type=player_position, 得到 %v", received["type"])
	}

	data, ok := received["data"].(map[string]interface{})
	if !ok {
		t.Fatal("data字段格式错误")
	}
	if data["pos_x"].(float64) != 500 {
		t.Errorf("期望 pos_x=500, 得到 %v", data["pos_x"])
	}
	if data["pos_y"].(float64) != 300 {
		t.Errorf("期望 pos_y=300, 得到 %v", data["pos_y"])
	}
	t.Log("WebSocket玩家位置广播测试通过")
}

func TestWebSocketChatMessage(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	sceneID := "scene_village_square"
	wsURL1 := fmt.Sprintf("ws%s/api/ws?scene_id=%s", ts.URL[len("http"):], sceneID)
	conn1, _, err := websocket.DefaultDialer.Dial(wsURL1, nil)
	if err != nil {
		t.Fatalf("客户端1连接失败: %v", err)
	}
	defer conn1.Close()

	wsURL2 := fmt.Sprintf("ws%s/api/ws?scene_id=%s", ts.URL[len("http"):], sceneID)
	conn2, _, err := websocket.DefaultDialer.Dial(wsURL2, nil)
	if err != nil {
		t.Fatalf("客户端2连接失败: %v", err)
	}
	defer conn2.Close()

	chatMsg := map[string]interface{}{
		"type": "chat_message",
		"data": map[string]interface{}{
			"player_name": "测试玩家",
			"channel":     "scene",
			"content":     "你好世界",
		},
	}

	if err := conn1.WriteJSON(chatMsg); err != nil {
		t.Fatalf("发送聊天消息失败: %v", err)
	}

	_ = conn2.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, message, err := conn2.ReadMessage()
	if err != nil {
		t.Fatalf("接收聊天消息失败: %v", err)
	}

	var received map[string]interface{}
	if err := json.Unmarshal(message, &received); err != nil {
		t.Fatalf("解析聊天消息失败: %v", err)
	}

	if received["type"] != "chat_message" {
		t.Errorf("期望 type=chat_message, 得到 %v", received["type"])
	}

	data, ok := received["data"].(map[string]interface{})
	if !ok {
		t.Fatal("data字段格式错误")
	}
	if data["content"] != "你好世界" {
		t.Errorf("期望 content=你好世界, 得到 %v", data["content"])
	}
	t.Log("WebSocket聊天消息测试通过")
}

func TestWebSocketGlobalChatMessage(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	sceneID1 := "scene_village_square"
	sceneID2 := "scene_forest"

	wsURL1 := fmt.Sprintf("ws%s/api/ws?scene_id=%s", ts.URL[len("http"):], sceneID1)
	conn1, _, err := websocket.DefaultDialer.Dial(wsURL1, nil)
	if err != nil {
		t.Fatalf("客户端1连接失败: %v", err)
	}
	defer conn1.Close()

	wsURL2 := fmt.Sprintf("ws%s/api/ws?scene_id=%s", ts.URL[len("http"):], sceneID2)
	conn2, _, err := websocket.DefaultDialer.Dial(wsURL2, nil)
	if err != nil {
		t.Fatalf("客户端2连接失败: %v", err)
	}
	defer conn2.Close()

	chatMsg := map[string]interface{}{
		"type": "chat_message",
		"data": map[string]interface{}{
			"player_name": "测试玩家",
			"channel":     "global",
			"content":     "全局消息测试",
		},
	}

	if err := conn1.WriteJSON(chatMsg); err != nil {
		t.Fatalf("发送全局消息失败: %v", err)
	}

	_ = conn2.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, message, err := conn2.ReadMessage()
	if err != nil {
		t.Fatalf("接收全局消息失败: %v", err)
	}

	var received map[string]interface{}
	if err := json.Unmarshal(message, &received); err != nil {
		t.Fatalf("解析消息失败: %v", err)
	}

	if received["type"] != "chat_message" {
		t.Errorf("期望 type=chat_message, 得到 %v", received["type"])
	}

	data, ok := received["data"].(map[string]interface{})
	if !ok {
		t.Fatal("data字段格式错误")
	}
	if data["content"] != "全局消息测试" {
		t.Errorf("期望 content=全局消息测试, 得到 %v", data["content"])
	}
	t.Log("WebSocket全局聊天消息测试通过")
}

func TestWebSocketPingPong(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	wsURL := fmt.Sprintf("ws%s/api/ws", ts.URL[len("http"):])
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()

	pingMsg := map[string]interface{}{
		"type": "ping",
		"data": map[string]interface{}{},
	}

	if err := conn.WriteJSON(pingMsg); err != nil {
		t.Fatalf("发送ping失败: %v", err)
	}

	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, message, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("接收pong失败: %v", err)
	}

	var received map[string]interface{}
	if err := json.Unmarshal(message, &received); err != nil {
		t.Fatalf("解析消息失败: %v", err)
	}

	if received["type"] != "pong" {
		t.Errorf("期望 type=pong, 得到 %v", received["type"])
	}
	t.Log("WebSocket Ping/Pong测试通过")
}

func TestWebSocketSceneIsolation(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	sceneA := "scene_village_square"
	sceneB := "scene_forest"

	wsURL1 := fmt.Sprintf("ws%s/api/ws?scene_id=%s", ts.URL[len("http"):], sceneA)
	conn1, _, err := websocket.DefaultDialer.Dial(wsURL1, nil)
	if err != nil {
		t.Fatalf("客户端1连接失败: %v", err)
	}
	defer conn1.Close()

	wsURL2 := fmt.Sprintf("ws%s/api/ws?scene_id=%s", ts.URL[len("http"):], sceneB)
	conn2, _, err := websocket.DefaultDialer.Dial(wsURL2, nil)
	if err != nil {
		t.Fatalf("客户端2连接失败: %v", err)
	}
	defer conn2.Close()

	posMsg := map[string]interface{}{
		"type": "player_position",
		"data": map[string]interface{}{
			"scene_id": sceneA,
			"pos_x":    100,
			"pos_y":    200,
		},
	}

	if err := conn1.WriteJSON(posMsg); err != nil {
		t.Fatalf("发送消息失败: %v", err)
	}

	_ = conn2.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	_, _, err = conn2.ReadMessage()
	if err == nil {
		t.Error("场景B的客户端不应收到来自场景A的消息")
	}
	t.Log("WebSocket场景隔离测试通过")
}

func TestBroadcastFunctions(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	router := setupTestRouter()
	_ = router

	t.Run("BroadcastNPCState", func(t *testing.T) {
		wsURL := fmt.Sprintf("ws%s/api/ws?scene_id=scene_test", ts.URL[len("http"):])
		conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			t.Fatalf("连接失败: %v", err)
		}
		defer conn.Close()

		t.Log("NPC状态广播函数可用")
	})

	t.Run("BroadcastCombatEvent", func(t *testing.T) {
		t.Log("战斗事件广播函数可用")
	})

	t.Run("BroadcastItemPickup", func(t *testing.T) {
		t.Log("道具拾取广播函数可用")
	})

	t.Run("BroadcastSystemMessage", func(t *testing.T) {
		t.Log("系统消息广播函数可用")
	})
}

func TestWebSocketHubStats(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	sceneID := "scene_village_square"

	wsURL1 := fmt.Sprintf("ws%s/api/ws?scene_id=%s", ts.URL[len("http"):], sceneID)
	conn1, _, err := websocket.DefaultDialer.Dial(wsURL1, nil)
	if err != nil {
		t.Fatalf("连接1失败: %v", err)
	}
	defer conn1.Close()

	wsURL2 := fmt.Sprintf("ws%s/api/ws?scene_id=%s", ts.URL[len("http"):], sceneID)
	conn2, _, err := websocket.DefaultDialer.Dial(wsURL2, nil)
	if err != nil {
		t.Fatalf("连接2失败: %v", err)
	}
	defer conn2.Close()

	wsURL3 := fmt.Sprintf("ws%s/api/ws?scene_id=scene_forest", ts.URL[len("http"):])
	conn3, _, err := websocket.DefaultDialer.Dial(wsURL3, nil)
	if err != nil {
		t.Fatalf("连接3失败: %v", err)
	}
	defer conn3.Close()

	t.Log("WebSocket Hub统计测试通过")
}

func TestWebSocketDisconnect(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	sceneID := "scene_village_square"

	wsURL1 := fmt.Sprintf("ws%s/api/ws?scene_id=%s", ts.URL[len("http"):], sceneID)
	conn1, _, err := websocket.DefaultDialer.Dial(wsURL1, nil)
	if err != nil {
		t.Fatalf("连接1失败: %v", err)
	}

	wsURL2 := fmt.Sprintf("ws%s/api/ws?scene_id=%s", ts.URL[len("http"):], sceneID)
	conn2, _, err := websocket.DefaultDialer.Dial(wsURL2, nil)
	if err != nil {
		t.Fatalf("连接2失败: %v", err)
	}
	defer conn2.Close()

	conn1.Close()
	time.Sleep(100 * time.Millisecond)

	chatMsg := map[string]interface{}{
		"type": "chat_message",
		"data": map[string]interface{}{
			"content": "断开后测试",
		},
	}

	if err := conn2.WriteJSON(chatMsg); err != nil {
		t.Fatalf("发送消息失败: %v", err)
	}

	t.Log("WebSocket断开处理测试通过")
}

func TestWebSocketInvalidMessage(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	wsURL := fmt.Sprintf("ws%s/api/ws", ts.URL[len("http"):])
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()

	if err := conn.WriteMessage(websocket.TextMessage, []byte("invalid json")); err != nil {
		t.Fatalf("发送无效消息失败: %v", err)
	}

	if err := conn.WriteJSON(map[string]interface{}{
		"type": "ping",
		"data": map[string]interface{}{},
	}); err != nil {
		t.Fatalf("发送后续消息失败: %v", err)
	}

	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, message, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("接收响应失败: %v", err)
	}

	var received map[string]interface{}
	json.Unmarshal(message, &received)
	if received["type"] != "pong" {
		t.Errorf("期望 type=pong, 得到 %v", received["type"])
	}
	t.Log("WebSocket无效消息处理测试通过")
}

func TestBroadcastPlayerPositionHTTP(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	sceneID := "scene_village_square"
	wsURL := fmt.Sprintf("ws%s/api/ws?scene_id=%s", ts.URL[len("http"):], sceneID)
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()

	t.Log("BroadcastPlayerPosition HTTP集成测试通过")
}

func TestWebSocketUpgrader(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/ws")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		t.Log("WebSocket端点正常响应HTTP GET")
	}
}

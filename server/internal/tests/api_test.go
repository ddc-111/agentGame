package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

// TestHealthEndpoint 测试健康检查端点
func TestHealthEndpoint(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/health")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if result["status"] != "ok" {
		t.Errorf("期望 status=ok, 得到 %v", result["status"])
	}
	t.Log("健康检查端点测试通过")
}

// TestGetGameInit 测试获取游戏初始化数据
func TestGetGameInit(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/game/init")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	// 验证返回的数据结构
	if _, ok := result["config"]; !ok {
		t.Error("响应缺少 config 字段")
	}
	if _, ok := result["scenes"]; !ok {
		t.Error("响应缺少 scenes 字段")
	}
	if _, ok := result["npcs"]; !ok {
		t.Error("响应缺少 npcs 字段")
	}
	if _, ok := result["tasks"]; !ok {
		t.Error("响应缺少 tasks 字段")
	}
	if _, ok := result["items"]; !ok {
		t.Error("响应缺少 items 字段")
	}
	t.Log("游戏初始化数据测试通过")
}

// TestCreatePlayer 测试创建玩家
func TestCreatePlayer(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	// 测试成功创建
	resp, err := makeRequest("POST", ts.URL+"/api/player/create", TestPlayer)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		t.Fatal("响应 data 字段格式错误")
	}

	if data["name"] != TestPlayer["name"] {
		t.Errorf("期望 name=%v, 得到 %v", TestPlayer["name"], data["name"])
	}
	if data["account"] != TestPlayer["account"] {
		t.Errorf("期望 account=%v, 得到 %v", TestPlayer["account"], data["account"])
	}
	t.Logf("创建玩家测试通过, ID=%v", data["id"])

	// 测试重复创建（应返回已有玩家）
	resp2, err := makeRequest("POST", ts.URL+"/api/player/create", TestPlayer)
	if err != nil {
		t.Fatalf("重复创建请求失败: %v", err)
	}
	defer resp2.Body.Close()

	assertStatusCode(t, resp2.StatusCode, http.StatusOK)

	var result2 map[string]interface{}
	if err := json.NewDecoder(resp2.Body).Decode(&result2); err != nil {
		t.Fatalf("解析重复创建响应失败: %v", err)
	}

	data2, ok := result2["data"].(map[string]interface{})
	if !ok {
		t.Fatal("重复创建响应 data 字段格式错误")
	}

	if data2["id"] != data["id"] {
		t.Errorf("重复创建应返回相同ID, 期望 %v, 得到 %v", data["id"], data2["id"])
	}
	t.Log("重复创建玩家测试通过")
}

// TestGetPlayer 测试获取玩家信息
func TestGetPlayer(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	// 先创建玩家
	createResp, err := makeRequest("POST", ts.URL+"/api/player/create", TestPlayer)
	if err != nil {
		t.Fatalf("创建玩家失败: %v", err)
	}
	defer createResp.Body.Close()

	var createResult map[string]interface{}
	json.NewDecoder(createResp.Body).Decode(&createResult)
	data := createResult["data"].(map[string]interface{})
	playerID := uint(data["id"].(float64))

	// 获取玩家信息
	resp, err := http.Get(fmt.Sprintf("%s/api/player/%d", ts.URL, playerID))
	if err != nil {
		t.Fatalf("获取玩家失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	playerData, ok := result["data"].(map[string]interface{})
	if !ok {
		t.Fatal("响应 data 字段格式错误")
	}

	if playerData["name"] != TestPlayer["name"] {
		t.Errorf("期望 name=%v, 得到 %v", TestPlayer["name"], playerData["name"])
	}
	t.Log("获取玩家信息测试通过")

	// 测试不存在的玩家
	resp404, err := http.Get(fmt.Sprintf("%s/api/player/99999", ts.URL))
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp404.Body.Close()

	assertStatusCode(t, resp404.StatusCode, http.StatusNotFound)
	t.Log("不存在玩家测试通过")
}

// TestUpdatePlayerPos 测试更新玩家位置
func TestUpdatePlayerPos(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	// 先创建玩家
	createResp, err := makeRequest("POST", ts.URL+"/api/player/create", TestPlayer)
	if err != nil {
		t.Fatalf("创建玩家失败: %v", err)
	}
	defer createResp.Body.Close()

	var createResult map[string]interface{}
	json.NewDecoder(createResp.Body).Decode(&createResult)
	data := createResult["data"].(map[string]interface{})
	playerID := uint(data["id"].(float64))

	// 更新位置
	updateData := map[string]interface{}{
		"scene_id": "scene_village_square",
		"pos_x":    500,
		"pos_y":    300,
	}

	resp, err := makeRequest("PUT", fmt.Sprintf("%s/api/player/%d/pos", ts.URL, playerID), updateData)
	if err != nil {
		t.Fatalf("更新位置失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	playerData := result["data"].(map[string]interface{})
	if playerData["scene_id"] != updateData["scene_id"] {
		t.Errorf("期望 scene_id=%v, 得到 %v", updateData["scene_id"], playerData["scene_id"])
	}
	if playerData["pos_x"].(float64) != float64(updateData["pos_x"].(int)) {
		t.Errorf("期望 pos_x=%v, 得到 %v", updateData["pos_x"], playerData["pos_x"])
	}
	t.Log("更新玩家位置测试通过")
}

// TestNPCChat 测试NPC对话
func TestNPCChat(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	// 先创建玩家
	createResp, err := makeRequest("POST", ts.URL+"/api/player/create", TestPlayer)
	if err != nil {
		t.Fatalf("创建玩家失败: %v", err)
	}
	defer createResp.Body.Close()

	var createResult map[string]interface{}
	json.NewDecoder(createResp.Body).Decode(&createResult)
	playerData := createResult["data"].(map[string]interface{})
	playerID := uint(playerData["id"].(float64))

	// 获取NPC列表（使用种子数据中的NPC）
	npcsResp, err := http.Get(ts.URL + "/api/npcs")
	if err != nil {
		t.Fatalf("获取NPC列表失败: %v", err)
	}
	defer npcsResp.Body.Close()

	var npcsResult map[string]interface{}
	json.NewDecoder(npcsResp.Body).Decode(&npcsResult)
	npcs := npcsResult["data"].([]interface{})
	if len(npcs) == 0 {
		t.Skip("没有NPC数据，跳过对话测试")
	}

	npc := npcs[0].(map[string]interface{})
	npcID := uint(npc["id"].(float64))

	// 发送对话请求
	chatData := map[string]interface{}{
		"player_id": playerID,
		"npc_id":    npcID,
		"message":   "你好",
	}

	resp, err := makeRequest("POST", ts.URL+"/api/npc/chat", chatData)
	if err != nil {
		t.Fatalf("对话请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if _, ok := result["reply"]; !ok {
		t.Error("响应缺少 reply 字段")
	}
	if _, ok := result["npc_name"]; !ok {
		t.Error("响应缺少 npc_name 字段")
	}
	t.Logf("NPC对话测试通过, 回复: %v", result["reply"])
}

// TestGetShopItems 测试获取商店商品
func TestGetShopItems(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	// 获取商店列表
	shopsResp, err := http.Get(ts.URL + "/api/shops")
	if err != nil {
		t.Fatalf("获取商店列表失败: %v", err)
	}
	defer shopsResp.Body.Close()

	var shopsResult map[string]interface{}
	json.NewDecoder(shopsResp.Body).Decode(&shopsResult)
	shops := shopsResult["data"].([]interface{})
	if len(shops) == 0 {
		t.Skip("没有商店数据，跳过商店商品测试")
	}

	shop := shops[0].(map[string]interface{})
	shopCode := shop["code"].(string)

	// 获取商店商品
	resp, err := http.Get(fmt.Sprintf("%s/api/game/shop/%s/items", ts.URL, shopCode))
	if err != nil {
		t.Fatalf("获取商店商品失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if _, ok := result["shop"]; !ok {
		t.Error("响应缺少 shop 字段")
	}
	if _, ok := result["items"]; !ok {
		t.Error("响应缺少 items 字段")
	}
	t.Logf("获取商店商品测试通过, 商店: %v", shopCode)
}

// TestBuyItem 测试购买道具
func TestBuyItem(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	// 先创建玩家
	createResp, err := makeRequest("POST", ts.URL+"/api/player/create", TestPlayer2)
	if err != nil {
		t.Fatalf("创建玩家失败: %v", err)
	}
	defer createResp.Body.Close()

	var createResult map[string]interface{}
	json.NewDecoder(createResp.Body).Decode(&createResult)
	playerData := createResult["data"].(map[string]interface{})
	playerID := uint(playerData["id"].(float64))

	// 获取商店商品
	shopsResp, err := http.Get(ts.URL + "/api/shops")
	if err != nil {
		t.Fatalf("获取商店列表失败: %v", err)
	}
	defer shopsResp.Body.Close()

	var shopsResult map[string]interface{}
	json.NewDecoder(shopsResp.Body).Decode(&shopsResult)
	shops := shopsResult["data"].([]interface{})
	if len(shops) == 0 {
		t.Skip("没有商店数据，跳过购买测试")
	}

	shop := shops[0].(map[string]interface{})
	shopCode := shop["code"].(string)

	// 获取商品列表
	itemsResp, err := http.Get(fmt.Sprintf("%s/api/game/shop/%s/items", ts.URL, shopCode))
	if err != nil {
		t.Fatalf("获取商品列表失败: %v", err)
	}
	defer itemsResp.Body.Close()

	var itemsResult map[string]interface{}
	json.NewDecoder(itemsResp.Body).Decode(&itemsResult)
	items := itemsResult["items"].([]interface{})
	if len(items) == 0 {
		t.Skip("商店没有商品，跳过购买测试")
	}

	item := items[0].(map[string]interface{})
	itemID := uint(item["item_id"].(float64))

	// 购买道具
	buyData := map[string]interface{}{
		"player_id": playerID,
		"shop_code": shopCode,
		"item_id":   itemID,
		"count":     1,
	}

	resp, err := makeRequest("POST", ts.URL+"/api/shop/buy", buyData)
	if err != nil {
		t.Fatalf("购买请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if result["message"] != "购买成功" {
		t.Errorf("期望 message=购买成功, 得到 %v", result["message"])
	}
	t.Logf("购买道具测试通过, 消费: %v", result["total_price"])
}

// TestGetScenes 测试获取场景列表
func TestGetScenes(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/scenes")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	data, ok := result["data"].([]interface{})
	if !ok {
		t.Fatal("响应 data 字段格式错误")
	}

	if len(data) == 0 {
		t.Error("场景列表为空")
	}
	t.Logf("获取场景列表测试通过, 数量: %d", len(data))
}

// TestGetSceneByCode 测试通过代码获取场景
func TestGetSceneByCode(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	// 先获取场景列表
	listResp, err := http.Get(ts.URL + "/api/scenes")
	if err != nil {
		t.Fatalf("获取场景列表失败: %v", err)
	}
	defer listResp.Body.Close()

	var listResult map[string]interface{}
	json.NewDecoder(listResp.Body).Decode(&listResult)
	scenes := listResult["data"].([]interface{})
	if len(scenes) == 0 {
		t.Skip("没有场景数据，跳过测试")
	}

	scene := scenes[0].(map[string]interface{})
	sceneCode := scene["code"].(string)

	// 通过代码获取场景
	resp, err := http.Get(fmt.Sprintf("%s/api/game/scene/%s", ts.URL, sceneCode))
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	sceneData := result["data"].(map[string]interface{})
	if sceneData["code"] != sceneCode {
		t.Errorf("期望 code=%v, 得到 %v", sceneCode, sceneData["code"])
	}
	t.Logf("通过代码获取场景测试通过, 场景: %v", sceneCode)

	// 测试不存在的场景
	resp404, err := http.Get(ts.URL + "/api/game/scene/nonexistent_scene")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp404.Body.Close()

	assertStatusCode(t, resp404.StatusCode, http.StatusNotFound)
	t.Log("不存在场景测试通过")
}

// TestGetNPCs 测试获取NPC列表
func TestGetNPCs(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/npcs")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	data, ok := result["data"].([]interface{})
	if !ok {
		t.Fatal("响应 data 字段格式错误")
	}

	if len(data) == 0 {
		t.Error("NPC列表为空")
	}
	t.Logf("获取NPC列表测试通过, 数量: %d", len(data))
}

// TestGetTasks 测试获取任务列表
func TestGetTasks(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/tasks")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	data, ok := result["data"].([]interface{})
	if !ok {
		t.Fatal("响应 data 字段格式错误")
	}

	if len(data) == 0 {
		t.Error("任务列表为空")
	}
	t.Logf("获取任务列表测试通过, 数量: %d", len(data))
}

// TestGetItems 测试获取道具列表
func TestGetItems(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/items")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	data, ok := result["data"].([]interface{})
	if !ok {
		t.Fatal("响应 data 字段格式错误")
	}

	if len(data) == 0 {
		t.Error("道具列表为空")
	}
	t.Logf("获取道具列表测试通过, 数量: %d", len(data))
}

// TestGetShops 测试获取商店列表
func TestGetShops(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/shops")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	data, ok := result["data"].([]interface{})
	if !ok {
		t.Fatal("响应 data 字段格式错误")
	}

	if len(data) == 0 {
		t.Error("商店列表为空")
	}
	t.Logf("获取商店列表测试通过, 数量: %d", len(data))
}

// TestGetAgents 测试获取智能体列表
func TestGetAgents(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/agents")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	data, ok := result["data"].([]interface{})
	if !ok {
		t.Fatal("响应 data 字段格式错误")
	}

	if len(data) == 0 {
		t.Error("智能体列表为空")
	}
	t.Logf("获取智能体列表测试通过, 数量: %d", len(data))
}

// TestGetFlows 测试获取流程列表
func TestGetFlows(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/flows")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	data, ok := result["data"].([]interface{})
	if !ok {
		t.Fatal("响应 data 字段格式错误")
	}

	if len(data) == 0 {
		t.Error("流程列表为空")
	}
	t.Logf("获取流程列表测试通过, 数量: %d", len(data))
}

// TestExportData 测试导出数据
func TestExportData(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/export")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		t.Fatal("响应 data 字段格式错误")
	}

	// 验证导出的数据包含各种类型
	expectedKeys := []string{"scenes", "npcs", "agents", "shops", "items", "tasks", "flows"}
	for _, key := range expectedKeys {
		if _, ok := data[key]; !ok {
			t.Errorf("导出数据缺少 %s 字段", key)
		}
	}
	t.Log("导出数据测试通过")
}

// TestImportDataInTransaction 测试导入数据在事务中执行
func TestImportDataInTransaction(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	exportResp, err := http.Get(ts.URL + "/api/export")
	if err != nil {
		t.Fatalf("导出请求失败: %v", err)
	}
	defer exportResp.Body.Close()

	var exportResult map[string]interface{}
	if err := json.NewDecoder(exportResp.Body).Decode(&exportResult); err != nil {
		t.Fatalf("解析导出响应失败: %v", err)
	}

	importData := exportResult["data"].(map[string]interface{})

	resp, err := makeRequest("POST", ts.URL+"/api/import", importData)
	if err != nil {
		t.Fatalf("导入请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析导入响应失败: %v", err)
	}

	if result["message"] != "Import successful" {
		t.Errorf("期望 message=Import successful, 得到 %v", result["message"])
	}

	imported, ok := result["imported"].(map[string]interface{})
	if !ok {
		t.Fatal("响应 imported 字段格式错误")
	}

	if _, ok := imported["scenes"]; !ok {
		t.Error("导入结果缺少 scenes 字段")
	}
	t.Log("事务导入数据测试通过")
}

// TestImportDataWrapped 测试导入数据带data包装
func TestImportDataWrapped(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	wrappedData := map[string]interface{}{
		"data": map[string]interface{}{
			"scenes": []interface{}{
				map[string]interface{}{
					"code":        "scene_import_test",
					"name":        "导入测试场景",
					"description": "事务导入测试",
					"width":       1920,
					"height":      1080,
				},
			},
		},
	}

	resp, err := makeRequest("POST", ts.URL+"/api/import", wrappedData)
	if err != nil {
		t.Fatalf("导入请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if result["message"] != "Import successful" {
		t.Errorf("期望 message=Import successful, 得到 %v", result["message"])
	}

	imported := result["imported"].(map[string]interface{})
	if imported["scenes"].(float64) != 1 {
		t.Errorf("期望导入1个场景, 得到 %v", imported["scenes"])
	}
	t.Log("带data包装的导入测试通过")
}

func TestStartCombat(t *testing.T) {
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

	combatData := map[string]interface{}{
		"player_id":  playerID,
		"enemy_type": "wolf",
	}

	resp, err := makeRequest("POST", ts.URL+"/api/combat/start", combatData)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if _, ok := result["data"]; !ok {
		t.Error("响应缺少 data 字段")
	}
	if result["message"] != "战斗开始" {
		t.Errorf("期望 message=战斗开始, 得到 %v", result["message"])
	}
	t.Log("开始战斗测试通过")
}

func TestGetInventory(t *testing.T) {
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

	resp, err := http.Get(fmt.Sprintf("%s/api/inventory/%d", ts.URL, playerID))
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if _, ok := result["items"]; !ok {
		t.Error("响应缺少 items 字段")
	}
	if _, ok := result["equipment"]; !ok {
		t.Error("响应缺少 equipment 字段")
	}
	if _, ok := result["gold"]; !ok {
		t.Error("响应缺少 gold 字段")
	}
	t.Log("获取背包测试通过")
}

func TestSaveAndLoadGame(t *testing.T) {
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

	saveData := map[string]interface{}{
		"player_id": playerID,
		"slot":      0,
		"name":      "测试存档",
	}
	saveResp, err := makeRequest("POST", ts.URL+"/api/save", saveData)
	if err != nil {
		t.Fatalf("保存游戏失败: %v", err)
	}
	defer saveResp.Body.Close()

	assertStatusCode(t, saveResp.StatusCode, http.StatusOK)

	var saveResult map[string]interface{}
	json.NewDecoder(saveResp.Body).Decode(&saveResult)
	if saveResult["message"] != "存档保存成功" {
		t.Errorf("期望 message=存档保存成功, 得到 %v", saveResult["message"])
	}

	savesResp, err := http.Get(fmt.Sprintf("%s/api/saves/%d", ts.URL, playerID))
	if err != nil {
		t.Fatalf("获取存档列表失败: %v", err)
	}
	defer savesResp.Body.Close()

	assertStatusCode(t, savesResp.StatusCode, http.StatusOK)

	var savesResult map[string]interface{}
	json.NewDecoder(savesResp.Body).Decode(&savesResult)
	if _, ok := savesResult["saves"]; !ok {
		t.Error("响应缺少 saves 字段")
	}
	t.Log("存档系统测试通过")
}

func TestGetSkills(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/skills")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if _, ok := result["data"]; !ok {
		t.Error("响应缺少 data 字段")
	}
	t.Log("获取技能列表测试通过")
}

func TestGetPlayerAchievements(t *testing.T) {
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

	resp, err := http.Get(fmt.Sprintf("%s/api/achievements/%d", ts.URL, playerID))
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if _, ok := result["achievements"]; !ok {
		t.Error("响应缺少 achievements 字段")
	}
	if _, ok := result["total"]; !ok {
		t.Error("响应缺少 total 字段")
	}
	if _, ok := result["unlocked"]; !ok {
		t.Error("响应缺少 unlocked 字段")
	}
	t.Log("获取玩家成就测试通过")
}

func TestCheckAchievements(t *testing.T) {
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

	checkData := map[string]interface{}{
		"player_id": playerID,
	}
	resp, err := makeRequest("POST", ts.URL+"/api/achievements/check", checkData)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if _, ok := result["new_achievements"]; !ok {
		t.Error("响应缺少 new_achievements 字段")
	}
	if _, ok := result["count"]; !ok {
		t.Error("响应缺少 count 字段")
	}
	t.Log("检查成就测试通过")
}

func TestSkillUseWithEquipment(t *testing.T) {
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

	initResp, err := http.Get(ts.URL + "/api/game/init")
	if err != nil {
		t.Fatalf("获取初始化数据失败: %v", err)
	}
	defer initResp.Body.Close()

	var initResult map[string]interface{}
	json.NewDecoder(initResp.Body).Decode(&initResult)

	var weaponID uint
	if items, ok := initResult["items"].([]interface{}); ok {
		for _, item := range items {
			itemMap := item.(map[string]interface{})
			if itemMap["code"] == "item_iron_sword" {
				weaponID = uint(itemMap["id"].(float64))
				break
			}
		}
	}
	if weaponID == 0 {
		t.Fatal("未找到铁剑道具")
	}

	var skillID uint
	if skills, ok := initResult["skills"].([]interface{}); ok && len(skills) > 0 {
		skillMap := skills[0].(map[string]interface{})
		skillID = uint(skillMap["id"].(float64))
	}
	if skillID == 0 {
		skillsResp, err := http.Get(ts.URL + "/api/skills")
		if err != nil {
			t.Fatalf("获取技能列表失败: %v", err)
		}
		defer skillsResp.Body.Close()
		var skillsResult map[string]interface{}
		json.NewDecoder(skillsResp.Body).Decode(&skillsResult)
		if skills, ok := skillsResult["data"].([]interface{}); ok && len(skills) > 0 {
			skillMap := skills[0].(map[string]interface{})
			skillID = uint(skillMap["id"].(float64))
		}
	}
	if skillID == 0 {
		t.Skip("没有技能数据，跳过技能使用测试")
	}

	buyData := map[string]interface{}{
		"player_id": playerID,
		"shop_code": "shop_blacksmith",
		"item_id":   weaponID,
		"count":     1,
	}
	buyResp, err := makeRequest("POST", ts.URL+"/api/shop/buy", buyData)
	if err != nil {
		t.Fatalf("购买道具失败: %v", err)
	}
	defer buyResp.Body.Close()
	assertStatusCode(t, buyResp.StatusCode, http.StatusOK)

	equipData := map[string]interface{}{
		"player_id": playerID,
		"item_id":   weaponID,
	}
	equipResp, err := makeRequest("POST", ts.URL+"/api/inventory/equip", equipData)
	if err != nil {
		t.Fatalf("装备失败: %v", err)
	}
	defer equipResp.Body.Close()
	assertStatusCode(t, equipResp.StatusCode, http.StatusOK)

	combatData := map[string]interface{}{
		"player_id":  playerID,
		"enemy_type": "wolf",
	}
	combatResp, err := makeRequest("POST", ts.URL+"/api/combat/start", combatData)
	if err != nil {
		t.Fatalf("开始战斗失败: %v", err)
	}
	defer combatResp.Body.Close()
	assertStatusCode(t, combatResp.StatusCode, http.StatusOK)

	var combatResult map[string]interface{}
	json.NewDecoder(combatResp.Body).Decode(&combatResult)
	combatState := combatResult["data"]

	skillData := map[string]interface{}{
		"player_id": playerID,
		"skill_id":  skillID,
		"state":     combatState,
	}
	skillResp, err := makeRequest("POST", ts.URL+"/api/skills/use", skillData)
	if err != nil {
		t.Fatalf("使用技能失败: %v", err)
	}
	defer skillResp.Body.Close()
	assertStatusCode(t, skillResp.StatusCode, http.StatusOK)

	var skillResult map[string]interface{}
	if err := json.NewDecoder(skillResp.Body).Decode(&skillResult); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}
	if _, ok := skillResult["data"]; !ok {
		t.Error("响应缺少 data 字段")
	}
	if _, ok := skillResult["skill"]; !ok {
		t.Error("响应缺少 skill 字段")
	}
	t.Log("装备后使用技能测试通过")
}

func TestCombatWithEquipment(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	combatEquipPlayer := map[string]interface{}{
		"name":    "战斗装备测试玩家",
		"account": "test_combat_equip",
	}
	createResp, err := makeRequest("POST", ts.URL+"/api/player/create", combatEquipPlayer)
	if err != nil {
		t.Fatalf("创建玩家失败: %v", err)
	}
	defer createResp.Body.Close()

	var createResult map[string]interface{}
	json.NewDecoder(createResp.Body).Decode(&createResult)
	playerData := createResult["data"].(map[string]interface{})
	playerID := uint(playerData["id"].(float64))

	initResp, err := http.Get(ts.URL + "/api/game/init")
	if err != nil {
		t.Fatalf("获取初始化数据失败: %v", err)
	}
	defer initResp.Body.Close()

	var initResult map[string]interface{}
	json.NewDecoder(initResp.Body).Decode(&initResult)

	var weaponID uint
	if items, ok := initResult["items"].([]interface{}); ok {
		for _, item := range items {
			itemMap := item.(map[string]interface{})
			if itemMap["code"] == "item_iron_sword" {
				weaponID = uint(itemMap["id"].(float64))
				break
			}
		}
	}
	if weaponID == 0 {
		t.Fatal("未找到铁剑道具")
	}

	buyData := map[string]interface{}{
		"player_id": playerID,
		"shop_code": "shop_blacksmith",
		"item_id":   weaponID,
		"count":     1,
	}
	buyResp, err := makeRequest("POST", ts.URL+"/api/shop/buy", buyData)
	if err != nil {
		t.Fatalf("购买道具失败: %v", err)
	}
	defer buyResp.Body.Close()
	assertStatusCode(t, buyResp.StatusCode, http.StatusOK)

	equipData := map[string]interface{}{
		"player_id": playerID,
		"item_id":   weaponID,
	}
	equipResp, err := makeRequest("POST", ts.URL+"/api/inventory/equip", equipData)
	if err != nil {
		t.Fatalf("装备失败: %v", err)
	}
	defer equipResp.Body.Close()
	assertStatusCode(t, equipResp.StatusCode, http.StatusOK)

	combatData := map[string]interface{}{
		"player_id":  playerID,
		"enemy_type": "wolf",
	}
	combatResp, err := makeRequest("POST", ts.URL+"/api/combat/start", combatData)
	if err != nil {
		t.Fatalf("开始战斗失败: %v", err)
	}
	defer combatResp.Body.Close()
	assertStatusCode(t, combatResp.StatusCode, http.StatusOK)

	var combatResult map[string]interface{}
	if err := json.NewDecoder(combatResp.Body).Decode(&combatResult); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}
	if _, ok := combatResult["data"]; !ok {
		t.Error("响应缺少 data 字段")
	}
	if combatResult["message"] != "战斗开始" {
		t.Errorf("期望 message=战斗开始, 得到 %v", combatResult["message"])
	}
	t.Log("装备后战斗测试通过")
}

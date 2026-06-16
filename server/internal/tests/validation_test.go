package tests

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestValidation_CreateScene_MissingRequired(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"description": "no name or code",
	}
	resp, err := makeRequest("POST", ts.URL+"/api/scenes", body)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusBadRequest)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if _, ok := result["validation_errors"]; !ok {
		t.Error("响应缺少 validation_errors 字段")
	}
	t.Log("Scene 缺少必填字段验证通过")
}

func TestValidation_CreateScene_WidthOutOfRange(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"name":        "测试场景",
		"code":        "scene_test_range",
		"description": "width too small",
		"width":       10,
		"height":      1080,
	}
	resp, err := makeRequest("POST", ts.URL+"/api/scenes", body)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusBadRequest)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	errs, ok := result["validation_errors"].([]interface{})
	if !ok {
		t.Fatal("validation_errors 字段格式错误")
	}

	found := false
	for _, e := range errs {
		errMap := e.(map[string]interface{})
		if errMap["field"] == "width" {
			found = true
			break
		}
	}
	if !found {
		t.Error("期望 width 字段验证错误")
	}
	t.Log("Scene width 范围验证通过")
}

func TestValidation_CreateScene_Success(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"name":        "测试场景",
		"code":        "scene_test_valid",
		"description": "valid scene",
		"width":       1920,
		"height":      1080,
	}
	resp, err := makeRequest("POST", ts.URL+"/api/scenes", body)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusCreated)
	t.Log("Scene 创建成功验证通过")
}

func TestValidation_CreateNPC_MissingRequired(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"description": "no name or code",
	}
	resp, err := makeRequest("POST", ts.URL+"/api/npcs", body)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusBadRequest)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if _, ok := result["validation_errors"]; !ok {
		t.Error("响应缺少 validation_errors 字段")
	}
	t.Log("NPC 缺少必填字段验证通过")
}

func TestValidation_CreateAgent_MissingRequired(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"description": "no name or code",
	}
	resp, err := makeRequest("POST", ts.URL+"/api/agents", body)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusBadRequest)
	t.Log("Agent 缺少必填字段验证通过")
}

func TestValidation_CreateAgent_MaxTokensOutOfRange(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"name":         "测试智能体",
		"code":         "agent_test_range",
		"description":  "max_tokens too large",
		"max_tokens":   999999,
		"max_messages": 20,
	}
	resp, err := makeRequest("POST", ts.URL+"/api/agents", body)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusBadRequest)
	t.Log("Agent max_tokens 范围验证通过")
}

func TestValidation_CreateTemplate_MissingContent(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"name": "测试模板",
		"code": "template_test_no_content",
	}
	resp, err := makeRequest("POST", ts.URL+"/api/prompts", body)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusBadRequest)
	t.Log("Template 缺少 content 验证通过")
}

func TestValidation_CreateShop_MissingRequired(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"description": "no name or code",
	}
	resp, err := makeRequest("POST", ts.URL+"/api/shops", body)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusBadRequest)
	t.Log("Shop 缺少必填字段验证通过")
}

func TestValidation_CreateItem_MissingRequired(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"description": "no name or code",
	}
	resp, err := makeRequest("POST", ts.URL+"/api/items", body)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusBadRequest)
	t.Log("Item 缺少必填字段验证通过")
}

func TestValidation_CreateTask_MissingRequired(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"description": "no name or code",
	}
	resp, err := makeRequest("POST", ts.URL+"/api/tasks", body)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusBadRequest)
	t.Log("Task 缺少必填字段验证通过")
}

func TestValidation_CreateFlow_MissingRequired(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"description": "no name or code",
	}
	resp, err := makeRequest("POST", ts.URL+"/api/flows", body)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusBadRequest)
	t.Log("Flow 缺少必填字段验证通过")
}

func TestValidation_CreatePlayer_MissingRequired(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"name": "test",
	}
	resp, err := makeRequest("POST", ts.URL+"/api/player/create", body)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusBadRequest)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if _, ok := result["validation_errors"]; !ok {
		t.Error("响应缺少 validation_errors 字段")
	}
	t.Log("Player 缺少 account 验证通过")
}

func TestValidation_NPCChat_MissingMessage(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"player_id": 1,
		"npc_id":    1,
	}
	resp, err := makeRequest("POST", ts.URL+"/api/npc/chat", body)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusBadRequest)
	t.Log("NPC Chat 缺少 message 验证通过")
}

func TestValidation_NPCChat_ZeroPlayerID(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"player_id": 0,
		"npc_id":    1,
		"message":   "hello",
	}
	resp, err := makeRequest("POST", ts.URL+"/api/npc/chat", body)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusBadRequest)
	t.Log("NPC Chat player_id=0 验证通过")
}

func TestValidation_BuyItem_MissingShopCode(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"player_id": 1,
		"item_id":   1,
	}
	resp, err := makeRequest("POST", ts.URL+"/api/shop/buy", body)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusBadRequest)
	t.Log("BuyItem 缺少 shop_code 验证通过")
}

func TestValidation_UpdatePlayerPos_OutOfRange(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	createResp, err := makeRequest("POST", ts.URL+"/api/player/create", map[string]interface{}{
		"name":    "pos_test_player",
		"account": "pos_test_account",
	})
	if err != nil {
		t.Fatalf("创建玩家失败: %v", err)
	}
	defer createResp.Body.Close()

	var createResult map[string]interface{}
	json.NewDecoder(createResp.Body).Decode(&createResult)
	data := createResult["data"].(map[string]interface{})
	playerID := uint(data["id"].(float64))

	body := map[string]interface{}{
		"scene_id": "scene_village",
		"pos_x":    99999,
		"pos_y":    0,
	}
	resp, err := makeRequest("PUT", ts.URL+"/api/player/"+jsonNumber(playerID)+"/pos", body)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusBadRequest)
	t.Log("UpdatePlayerPos pos_x 范围验证通过")
}

func TestValidation_StartCombat_MissingEnemyType(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"player_id": 1,
	}
	resp, err := makeRequest("POST", ts.URL+"/api/combat/start", body)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusBadRequest)
	t.Log("StartCombat 缺少 enemy_type 验证通过")
}

func TestValidation_CombatAction_InvalidAction(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"player_id": 1,
		"action":    "invalid_action",
		"state":     map[string]interface{}{"is_active": true},
	}
	resp, err := makeRequest("POST", ts.URL+"/api/combat/action", body)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusBadRequest)
	t.Log("CombatAction 无效 action 验证通过")
}

func TestValidation_EquipItem_MissingItemID(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"player_id": 1,
	}
	resp, err := makeRequest("POST", ts.URL+"/api/inventory/equip", body)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusBadRequest)
	t.Log("EquipItem 缺少 item_id 验证通过")
}

func TestValidation_UnequipItem_InvalidSlot(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"player_id": 1,
		"slot":      "invalid_slot",
	}
	resp, err := makeRequest("POST", ts.URL+"/api/inventory/unequip", body)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusBadRequest)
	t.Log("UnequipItem 无效 slot 验证通过")
}

func TestValidation_SaveGame_SlotOutOfRange(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"player_id": 1,
		"slot":      99,
	}
	resp, err := makeRequest("POST", ts.URL+"/api/save", body)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusBadRequest)
	t.Log("SaveGame slot 范围验证通过")
}

func TestValidation_CreateConversation_InvalidRole(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"player_id": 1,
		"npc_id":    1,
		"role":      "invalid_role",
		"content":   "test",
	}
	resp, err := makeRequest("POST", ts.URL+"/api/conversations", body)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusBadRequest)
	t.Log("CreateConversation 无效 role 验证通过")
}

func TestValidation_MCPCall_MissingName(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"arguments": map[string]interface{}{},
	}
	resp, err := makeRequest("POST", ts.URL+"/api/mcp/call", body)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusBadRequest)
	t.Log("MCPCall 缺少 name 验证通过")
}

func TestValidation_Generate_InvalidType(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"type":   "invalid_type",
		"action": "create",
	}
	resp, err := makeRequest("POST", ts.URL+"/api/generator/generate", body)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusBadRequest)
	t.Log("Generate 无效 type 验证通过")
}

func TestParseID_InvalidString(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := makeRequest("GET", ts.URL+"/api/scenes/abc", nil)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusBadRequest)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	errObj, ok := result["error"].(map[string]interface{})
	if !ok {
		t.Fatal("响应缺少 error 字段")
	}
	if errObj["code"] != "BAD_REQUEST" {
		t.Errorf("期望 BAD_REQUEST, 得到 %v", errObj["code"])
	}
	t.Log("parseID 非数字字符串返回 400 验证通过")
}

func TestParseID_NegativeNumber(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := makeRequest("GET", ts.URL+"/api/scenes/-1", nil)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusBadRequest)
	t.Log("parseID 负数返回 400 验证通过")
}

func TestParseID_Zero(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := makeRequest("GET", ts.URL+"/api/scenes/0", nil)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusNotFound)
	t.Log("parseID 零值通过解析但未找到资源验证通过")
}

func TestParseID_ValidID(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := makeRequest("GET", ts.URL+"/api/scenes/1", nil)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusOK)
	t.Log("parseID 合法 ID 通过解析验证通过")
}

func TestParseID_FloatString(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := makeRequest("GET", ts.URL+"/api/scenes/1.5", nil)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusBadRequest)
	t.Log("parseID 浮点字符串返回 400 验证通过")
}

func TestParseID_SpecialCharacters(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := makeRequest("GET", ts.URL+"/api/scenes/%3Cscript%3E", nil)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusBadRequest)
	t.Log("parseID 特殊字符返回 400 验证通过")
}

func jsonNumber(n uint) string {
	b, _ := json.Marshal(n)
	return string(b)
}

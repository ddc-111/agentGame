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

func TestValidation_CreateNPC_InvalidBehaviorsJSON(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"name":       "test_npc",
		"code":       "npc_test_behaviors",
		"behaviors":  "not-json",
	}
	resp, err := makeRequest("POST", ts.URL+"/api/npcs", body)
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
		if errMap["field"] == "behaviors" {
			found = true
			break
		}
	}
	if !found {
		t.Error("期望 behaviors 字段验证错误")
	}
	t.Log("NPC invalid behaviors JSON 验证通过")
}

func TestValidation_CreateNPC_InvalidScheduleJSON(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"name":     "test_npc",
		"code":     "npc_test_schedule",
		"schedule": "[invalid",
	}
	resp, err := makeRequest("POST", ts.URL+"/api/npcs", body)
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
		if errMap["field"] == "schedule" {
			found = true
			break
		}
	}
	if !found {
		t.Error("期望 schedule 字段验证错误")
	}
	t.Log("NPC invalid schedule JSON 验证通过")
}

func TestValidation_CreateNPC_ValidBehaviorsAndSchedule(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"name":       "test_npc",
		"code":       "npc_test_valid_json",
		"behaviors":  `["idle","greet"]`,
		"schedule":   `[{"time":"08:00","action":"work"}]`,
	}
	resp, err := makeRequest("POST", ts.URL+"/api/npcs", body)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusCreated)
	t.Log("NPC valid behaviors/schedule 创建成功")
}

func TestValidation_CreateItem_InvalidCategory(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"name":     "test_item",
		"code":     "item_test_cat",
		"category": "invalid_category",
	}
	resp, err := makeRequest("POST", ts.URL+"/api/items", body)
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
		if errMap["field"] == "category" {
			found = true
			break
		}
	}
	if !found {
		t.Error("期望 category 字段验证错误")
	}
	t.Log("Item invalid category 验证通过")
}

func TestValidation_CreateItem_ValidCategory(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"name":     "test_item",
		"code":     "item_test_valid_cat",
		"category": "weapon",
	}
	resp, err := makeRequest("POST", ts.URL+"/api/items", body)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusCreated)
	t.Log("Item valid category 创建成功")
}

func TestValidation_CreateItem_InvalidEffectJSON(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"name":   "test_item",
		"code":   "item_test_effect",
		"effect": "not-json",
	}
	resp, err := makeRequest("POST", ts.URL+"/api/items", body)
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
		if errMap["field"] == "effect" {
			found = true
			break
		}
	}
	if !found {
		t.Error("期望 effect 字段验证错误")
	}
	t.Log("Item invalid effect JSON 验证通过")
}

func TestValidation_CreateItem_ValidEffectJSON(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"name":   "test_item",
		"code":   "item_test_valid_effect",
		"effect": `{"hp":20}`,
	}
	resp, err := makeRequest("POST", ts.URL+"/api/items", body)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusCreated)
	t.Log("Item valid effect JSON 创建成功")
}

func TestValidation_CreateShop_InvalidType(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"name": "test_shop",
		"code": "shop_test_type",
		"type": "invalid_type",
	}
	resp, err := makeRequest("POST", ts.URL+"/api/shops", body)
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
		if errMap["field"] == "type" {
			found = true
			break
		}
	}
	if !found {
		t.Error("期望 type 字段验证错误")
	}
	t.Log("Shop invalid type 验证通过")
}

func TestValidation_CreateShop_InvalidDiscountJSON(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"name":     "test_shop",
		"code":     "shop_test_discount",
		"discount": "not-json",
	}
	resp, err := makeRequest("POST", ts.URL+"/api/shops", body)
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
		if errMap["field"] == "discount" {
			found = true
			break
		}
	}
	if !found {
		t.Error("期望 discount 字段验证错误")
	}
	t.Log("Shop invalid discount JSON 验证通过")
}

func TestValidation_CreateShop_InvalidTimeFormat(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"name":      "test_shop",
		"code":      "shop_test_time",
		"open_time": "99:99",
	}
	resp, err := makeRequest("POST", ts.URL+"/api/shops", body)
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
		if errMap["field"] == "open_time" {
			found = true
			break
		}
	}
	if !found {
		t.Error("期望 open_time 字段验证错误")
	}
	t.Log("Shop invalid time format 验证通过")
}

func TestValidation_CreateShop_InvalidTimeFormatLetters(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"name":       "test_shop",
		"code":       "shop_test_time2",
		"close_time": "abc",
	}
	resp, err := makeRequest("POST", ts.URL+"/api/shops", body)
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
		if errMap["field"] == "close_time" {
			found = true
			break
		}
	}
	if !found {
		t.Error("期望 close_time 字段验证错误")
	}
	t.Log("Shop invalid time format letters 验证通过")
}

func TestValidation_CreateTask_InvalidType(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"name": "test_task",
		"code": "task_test_type",
		"type": "invalid_type",
	}
	resp, err := makeRequest("POST", ts.URL+"/api/tasks", body)
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
		if errMap["field"] == "type" {
			found = true
			break
		}
	}
	if !found {
		t.Error("期望 type 字段验证错误")
	}
	t.Log("Task invalid type 验证通过")
}

func TestValidation_CreateTask_InvalidStatus(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"name":   "test_task",
		"code":   "task_test_status",
		"status": "invalid_status",
	}
	resp, err := makeRequest("POST", ts.URL+"/api/tasks", body)
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
		if errMap["field"] == "status" {
			found = true
			break
		}
	}
	if !found {
		t.Error("期望 status 字段验证错误")
	}
	t.Log("Task invalid status 验证通过")
}

func TestValidation_CreateTask_InvalidTriggerJSON(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"name":    "test_task",
		"code":    "task_test_trigger",
		"trigger": "not-json",
	}
	resp, err := makeRequest("POST", ts.URL+"/api/tasks", body)
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
		if errMap["field"] == "trigger" {
			found = true
			break
		}
	}
	if !found {
		t.Error("期望 trigger 字段验证错误")
	}
	t.Log("Task invalid trigger JSON 验证通过")
}

func TestValidation_CreateTask_InvalidObjectivesJSON(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"name":       "test_task",
		"code":       "task_test_obj",
		"objectives": "not-json",
	}
	resp, err := makeRequest("POST", ts.URL+"/api/tasks", body)
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
		if errMap["field"] == "objectives" {
			found = true
			break
		}
	}
	if !found {
		t.Error("期望 objectives 字段验证错误")
	}
	t.Log("Task invalid objectives JSON 验证通过")
}

func TestValidation_CreateTask_InvalidRewardsJSON(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"name":    "test_task",
		"code":    "task_test_rewards",
		"rewards": "not-json",
	}
	resp, err := makeRequest("POST", ts.URL+"/api/tasks", body)
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
		if errMap["field"] == "rewards" {
			found = true
			break
		}
	}
	if !found {
		t.Error("期望 rewards 字段验证错误")
	}
	t.Log("Task invalid rewards JSON 验证通过")
}

func TestValidation_CreateTask_ValidJSONFields(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]interface{}{
		"name":       "test_task",
		"code":       "task_test_valid_json",
		"type":       "main",
		"status":     "active",
		"trigger":    `{"type":"auto","conditions":[]}`,
		"objectives": `[{"id":"obj_1","type":"dialogue","target":"npc_1"}]`,
		"rewards":    `{"exp":50,"gold":100}`,
	}
	resp, err := makeRequest("POST", ts.URL+"/api/tasks", body)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp.StatusCode, http.StatusCreated)
	t.Log("Task valid JSON fields 创建成功")
}

func jsonNumber(n uint) string {
	b, _ := json.Marshal(n)
	return string(b)
}

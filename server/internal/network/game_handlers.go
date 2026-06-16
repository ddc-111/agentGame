package network

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/agent"
	"github.com/ddc-111/agentGame/server/internal/database/models"
)

// GetGameInit 获取游戏初始化数据
func (s *Server) handleGetGameInit(c *gin.Context) {
	// 获取配置
	configs := make(map[string]string)
	keys := []string{"game_name", "start_scene", "start_x", "start_y", "player_speed", "default_hp", "default_mp", "default_gold"}
	for _, key := range keys {
		cfg, err := s.repo.GetConfig(key)
		if err == nil {
			configs[key] = cfg.Value
		}
	}

	// 获取场景列表（带NPC和传送点）
	scenes, _ := s.repo.GetScenes()

	// 获取所有NPC
	npcs, _ := s.repo.GetNPCs()

	// 获取所有任务
	tasks, _ := s.repo.GetTasks()

	// 获取所有道具
	items, _ := s.repo.GetItems()

	// 获取所有技能
	skills, _ := s.repo.GetSkills()

	c.JSON(http.StatusOK, gin.H{
		"config": configs,
		"scenes": scenes,
		"npcs":   npcs,
		"tasks":  tasks,
		"items":  items,
		"skills": skills,
	})
}

// CreatePlayer 创建玩家
func (s *Server) handleCreatePlayer(c *gin.Context) {
	var req struct {
		Name    string `json:"name"`
		Account string `json:"account"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}

	errs := validateRequired(map[string]interface{}{"name": req.Name, "account": req.Account})
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}

	// 检查账号是否已存在
	existing, _ := s.repo.GetPlayerByAccount(req.Account)
	if existing != nil && existing.ID > 0 {
		c.JSON(http.StatusOK, gin.H{"data": existing})
		return
	}

	// 获取起始配置
	startX, _ := strconv.Atoi(getConfigValue(s, "start_x", "200"))
	startY, _ := strconv.Atoi(getConfigValue(s, "start_y", "450"))
	defaultHP, _ := strconv.Atoi(getConfigValue(s, "default_hp", "100"))
	defaultMP, _ := strconv.Atoi(getConfigValue(s, "default_mp", "50"))
	defaultGold, _ := strconv.Atoi(getConfigValue(s, "default_gold", "500"))

	player := &models.Player{
		Name:    req.Name,
		Account: req.Account,
		Level:   1,
		HP:      defaultHP,
		MP:      defaultMP,
		Gold:    defaultGold,
		SceneID: getConfigValue(s, "start_scene", "scene_village_entrance"),
		PosX:    startX,
		PosY:    startY,
		Items:   "{}",
	}

	if err := s.repo.CreatePlayer(player); err != nil {
		respondInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": player})
}

// GetPlayer 获取玩家信息
func (s *Server) handleGetPlayer(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	player, err := s.repo.GetPlayerByID(uint(id))
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Player"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": player})
}

// UpdatePlayer 更新玩家信息
func (s *Server) handleUpdatePlayerPos(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	player, err := s.repo.GetPlayerByID(uint(id))
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Player"))
		return
	}

	var req struct {
		SceneID string `json:"scene_id"`
		PosX    int    `json:"pos_x"`
		PosY    int    `json:"pos_y"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}

	errs := mergeErrors(
		validateIntRange("pos_x", req.PosX, -10000, 10000),
		validateIntRange("pos_y", req.PosY, -10000, 10000),
	)
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}

	if req.SceneID != "" {
		player.SceneID = req.SceneID
	}
	player.PosX = req.PosX
	player.PosY = req.PosY

	if err := s.repo.UpdatePlayer(player); err != nil {
		respondInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": player})
}

// GetSceneByCode 通过code获取场景
func (s *Server) handleGetSceneByCode(c *gin.Context) {
	code := c.Param("code")
	scene, err := s.repo.GetSceneByCode(code)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Scene"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": scene})
}

// GetNPCByCode 通过code获取NPC
func (s *Server) handleGetNPCByCode(c *gin.Context) {
	code := c.Param("code")
	npc, err := s.repo.GetNPCByCode(code)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("NPC"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": npc})
}

// GetTasksByPlayer 获取玩家的任务列表
func (s *Server) handleGetPlayerTasks(c *gin.Context) {
	tasks, err := s.repo.GetTasks()
	if err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tasks})
}

// NPC Chat - 与NPC对话（AI驱动）
func (s *Server) handleNPCChat(c *gin.Context) {
	var req struct {
		PlayerID uint   `json:"player_id"`
		NPCID    uint   `json:"npc_id"`
		Message  string `json:"message"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}

	errs := mergeErrors(
		validatePositiveInt("player_id", req.PlayerID),
		validatePositiveInt("npc_id", req.NPCID),
		validateRequired(map[string]interface{}{"message": req.Message}),
	)
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}

	// 获取NPC和Agent信息
	npc, err := s.repo.GetNPCByID(req.NPCID)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("NPC"))
		return
	}

	// 获取玩家信息
	player, err := s.repo.GetPlayerByID(req.PlayerID)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Player"))
		return
	}

	// 获取对话历史
	history, _ := s.repo.GetConversations(req.PlayerID, req.NPCID, 10)

	// 构建对话上下文（用于fallback）
	messages := buildChatMessages(npc, player, history, req.Message)

	// 调用AI生成回复
	var reply string
	if s.chatMgr != nil && s.chatMgr.IsEnabled() && npc.Agent != nil && npc.Agent.SystemPrompt != "" {
		// 使用真实AI对话
		persona := &agent.NPCPersona{
			Name:         npc.Name,
			Title:        npc.Title,
			Description:  npc.Description,
			SystemPrompt: npc.Agent.SystemPrompt,
		}
		playerContext := buildPlayerContext(player)
		chatHistory := convertToChatMessages(history)

		aiReply, err := s.chatMgr.ChatWithNPC(
			c.Request.Context(),
			persona,
			playerContext,
			chatHistory,
			req.Message,
			npc.Agent.MaxMessages,
		)
		if err != nil {
			log.Printf("AI chat failed, falling back to simple reply: %v", err)
			reply = generateSimpleReply(npc.Agent, messages)
		} else {
			reply = aiReply
		}
	} else if npc.Agent != nil && npc.Agent.SystemPrompt != "" {
		// AI服务未启用，使用简单回复
		reply = generateSimpleReply(npc.Agent, messages)
	} else {
		reply = "你好，客官！"
	}

	// 保存对话记录
	conv := &models.Conversation{
		PlayerID: req.PlayerID,
		NPCID:    req.NPCID,
		AgentID:  0,
		Role:     "user",
		Content:  req.Message,
	}
	if npc.Agent != nil {
		conv.AgentID = npc.Agent.ID
	}
	s.repo.CreateConversation(conv)

	convReply := &models.Conversation{
		PlayerID: req.PlayerID,
		NPCID:    req.NPCID,
		AgentID:  0,
		Role:     "assistant",
		Content:  reply,
	}
	if npc.Agent != nil {
		convReply.AgentID = npc.Agent.ID
	}
	s.repo.CreateConversation(convReply)

	// 更新内存记忆
	agent.MemoryStore.AddMessage(req.PlayerID, req.NPCID, "user", req.Message)
	agent.MemoryStore.AddMessage(req.PlayerID, req.NPCID, "assistant", reply)
	agent.MemoryStore.UpdatePlayerInfo(req.PlayerID, req.NPCID, player.Name, player.Level)

	c.JSON(http.StatusOK, gin.H{
		"reply":      reply,
		"npc_name":   npc.Name,
		"npc_title":  npc.Title,
		"npc_avatar": npc.Avatar,
		"ai_powered": s.chatMgr != nil && s.chatMgr.IsEnabled(),
	})
}

// GetShopItems 获取商店商品
func (s *Server) handleGetShopItems(c *gin.Context) {
	code := c.Param("code")
	shop, err := s.repo.GetShopByCode(code)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Shop"))
		return
	}

	// 获取商品详情
	type ShopItemDetail struct {
		models.ShopItem
		ItemName        string `json:"item_name"`
		ItemDescription string `json:"item_description"`
		ItemCategory    string `json:"item_category"`
		ItemEffect      string `json:"item_effect"`
	}

	var details []ShopItemDetail
	for _, si := range shop.Items {
		item, err := s.repo.GetItemByID(si.ItemID)
		if err != nil {
			continue
		}
		details = append(details, ShopItemDetail{
			ShopItem:        si,
			ItemName:        item.Name,
			ItemDescription: item.Description,
			ItemCategory:    item.Category,
			ItemEffect:      item.Effect,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"shop":  shop,
		"items": details,
	})
}

// BuyItem 购买道具
func (s *Server) handleBuyItem(c *gin.Context) {
	var req struct {
		PlayerID uint   `json:"player_id"`
		ShopCode string `json:"shop_code"`
		ItemID   uint   `json:"item_id"`
		Count    int    `json:"count"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}

	errs := mergeErrors(
		validatePositiveInt("player_id", req.PlayerID),
		validateRequired(map[string]interface{}{"shop_code": req.ShopCode}),
		validatePositiveInt("item_id", req.ItemID),
	)
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}

	if req.Count <= 0 {
		req.Count = 1
	}

	// 获取玩家
	player, err := s.repo.GetPlayerByID(req.PlayerID)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Player"))
		return
	}

	// 获取商店商品
	shop, err := s.repo.GetShopByCode(req.ShopCode)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Shop"))
		return
	}

	// 查找商品
	var shopItem *models.ShopItem
	for _, si := range shop.Items {
		if si.ItemID == req.ItemID {
			shopItem = &si
			break
		}
	}

	if shopItem == nil {
		respondError(c, http.StatusBadRequest, BadRequest("Item not found in shop"))
		return
	}

	// 检查库存
	if shopItem.Stock < req.Count {
		respondError(c, http.StatusBadRequest, BadRequest("Not enough stock"))
		return
	}

	// 计算总价
	totalPrice := shopItem.Price * req.Count

	// 检查金币
	if player.Gold < totalPrice {
		respondError(c, http.StatusBadRequest, BadRequest("Not enough gold"))
		return
	}

	// 扣除金币
	player.Gold -= totalPrice

	// 添加道具到背包
	var items map[string]int
	json.Unmarshal([]byte(player.Items), &items)
	if items == nil {
		items = make(map[string]int)
	}
	itemKey := strconv.Itoa(int(req.ItemID))
	items[itemKey] += req.Count
	itemsJSON, _ := json.Marshal(items)
	player.Items = string(itemsJSON)

	// 更新库存
	shopItem.Stock -= req.Count

	// 保存
	s.repo.UpdatePlayer(player)
	s.repo.SaveShopItem(shopItem)

	// 解析装备信息返回
	equipment := map[string]interface{}{
		"weapon_id": nil,
		"armor_id":  nil,
	}
	if player.Equipment != "" {
		var equip struct {
			WeaponID uint `json:"weapon_id"`
			ArmorID  uint `json:"armor_id"`
		}
		if err := json.Unmarshal([]byte(player.Equipment), &equip); err == nil {
			if equip.WeaponID > 0 {
				equipment["weapon_id"] = equip.WeaponID
			}
			if equip.ArmorID > 0 {
				equipment["armor_id"] = equip.ArmorID
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "购买成功",
		"gold":        player.Gold,
		"items":       player.Items,
		"equipment":   equipment,
		"item_name":   shopItem.Item.Name,
		"total_price": totalPrice,
	})
}

// 辅助函数
func getConfigValue(s *Server, key, defaultVal string) string {
	cfg, err := s.repo.GetConfig(key)
	if err != nil {
		return defaultVal
	}
	return cfg.Value
}

func buildChatMessages(npc *models.NPC, player *models.Player, history []models.Conversation, userMsg string) []map[string]string {
	var messages []map[string]string

	// System prompt
	if npc.Agent != nil {
		messages = append(messages, map[string]string{
			"role":    "system",
			"content": npc.Agent.SystemPrompt,
		})
	}

	// History
	for i := len(history) - 1; i >= 0; i-- {
		messages = append(messages, map[string]string{
			"role":    history[i].Role,
			"content": history[i].Content,
		})
	}

	// Current user message with context
	contextMsg := "【玩家信息】姓名：" + player.Name + "，等级：" + strconv.Itoa(player.Level) + "\n"
	contextMsg += "【玩家消息】" + userMsg

	messages = append(messages, map[string]string{
		"role":    "user",
		"content": contextMsg,
	})

	return messages
}

// convertToChatMessages converts conversation history to chat messages
func convertToChatMessages(history []models.Conversation) []agent.ChatMessage {
	var messages []agent.ChatMessage
	for i := len(history) - 1; i >= 0; i-- {
		messages = append(messages, agent.ChatMessage{
			Role:    history[i].Role,
			Content: history[i].Content,
		})
	}
	return messages
}

// buildPlayerContext builds player context string
func buildPlayerContext(player *models.Player) string {
	contextMsg := "【玩家信息】姓名：" + player.Name + "，等级：" + strconv.Itoa(player.Level)
	if player.Gold > 0 {
		contextMsg += "，金币：" + strconv.Itoa(player.Gold)
	}
	return contextMsg
}

func generateSimpleReply(agent *models.Agent, messages []map[string]string) string {
	// 获取用户最后一条消息
	userMsg := ""
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i]["role"] == "user" {
			userMsg = messages[i]["content"]
			break
		}
	}

	// 简单的关键词回复
	switch agent.Code {
	case "agent_chief_chen":
		if contains(userMsg, "你好") || contains(userMsg, "村长") {
			return "呵呵，年轻人，欢迎来到青石村！老朽是这里的村长，姓陈。你初来乍到，有什么不懂的尽管问老朽。村子虽然不大，但五脏俱全，杂货铺、铁匠铺、茶摊都有。你先四处转转，熟悉熟悉环境吧。"
		}
		if contains(userMsg, "任务") || contains(userMsg, "帮忙") {
			return "呵呵，年轻人愿意帮忙，老朽很高兴。最近村子外面出现了不少野兽，村民们有些担忧。你先去杂货铺买些补给，再去铁匠铺买把武器，然后去找猎户老周，他会告诉你详情的。"
		}
		return "呵呵，年轻人，有什么想问的尽管说。老朽在青石村生活了一辈子，对这里的一切都很熟悉。"
	case "agent_merchant_li":
		if contains(userMsg, "买") || contains(userMsg, "商品") || contains(userMsg, "价格") {
			return "客官好眼力！小店经营各种日用品和药材。草药50文一份，可恢复20点生命；灵芝200文，恢复100点生命；馒头20文，恢复体力；烧酒80文，驱寒保暖；麻绳30文，攀爬用；回城符100文，瞬间回村。客官需要些什么？"
		}
		return "客官好！欢迎光临李记杂货铺！在下经营各种日用品、药材和食材，物美价廉，童叟无欺！客官随便看看，有什么需要尽管吩咐！"
	case "agent_tea_wang":
		if contains(userMsg, "故事") || contains(userMsg, "八卦") || contains(userMsg, "传说") {
			return "哎呀，客官想听故事啊！我跟你说，这青石村可有年头了。听说几百年前，有个仙人在这里修炼，留下了一块青石，村子就是以这块石头命名的。还有啊，村后面的山上据说有个山洞，里面藏着宝贝呢！不过村长说那地方危险，不让大家去。"
		}
		return "哎呀，客官来了！快坐下喝杯茶歇歇脚。我这茶可是用山泉水泡的，香着呢！客官是新来的吧？我跟你说，咱们青石村虽然不大，但可是个好地方！"
	case "agent_blacksmith_zhang":
		if contains(userMsg, "买") || contains(userMsg, "武器") || contains(userMsg, "兵器") {
			return "俺这儿的兵器都是好货！铁剑500文，攻击+10，新手用正好；精钢刀1000文，攻击+25，好刀！猎弓400文，远程攻击；铁甲800文，防御+15；皮甲300文，轻便；铁盾600文，防御+20。你看上哪个了？"
		}
		return "嗯，客官好。俺是这儿的铁匠，姓张。俺打的兵器在方圆百里都小有名气。你要是需要兵器，尽管说。"
	case "agent_hunter_zhou":
		if contains(userMsg, "狼") || contains(userMsg, "野兽") || contains(userMsg, "危险") {
			return "嗯，最近村外的狼群越来越多了。有几只头狼特别大，很凶猛。俺一个人应付不来。你要是有本事，帮俺去看看。记住，遇到狼群不要跑，慢慢后退，用武器自卫。"
		}
		return "嗯，你好。俺是猎户老周。村外最近不太平，你出去的时候小心点。"
	default:
		return "你好，客官！"
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// convertToAgentMessages 转换对话历史为Agent消息格式
func convertToAgentMessages(history []models.Conversation, player *models.Player) []agent.Message {
	var messages []agent.Message
	for i := len(history) - 1; i >= 0; i-- {
		messages = append(messages, agent.Message{
			Role:    history[i].Role,
			Content: history[i].Content,
		})
	}
	return messages
}

// buildUserMessage 构建用户消息（带玩家上下文）
func buildUserMessage(player *models.Player, userMsg string) string {
	contextMsg := "【玩家信息】姓名：" + player.Name + "，等级：" + strconv.Itoa(player.Level)
	if player.Gold > 0 {
		contextMsg += "，金币：" + strconv.Itoa(player.Gold)
	}
	contextMsg += "\n【玩家消息】" + userMsg
	return contextMsg
}

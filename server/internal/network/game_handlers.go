package network

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/agent"
	"github.com/ddc-111/agentGame/server/internal/database/models"
)

func (s *Server) handleGetGameInit(c *gin.Context) {
	ctx := c.Request.Context()

	configs := make(map[string]string)
	keys := []string{"game_name", "start_scene", "start_x", "start_y", "player_speed", "default_hp", "default_mp", "default_gold"}
	for _, key := range keys {
		cfg, err := s.repo.GetConfig(ctx, key)
		if err == nil {
			configs[key] = cfg.Value
		}
	}

	scenes, _ := s.repo.GetScenes(ctx)
	tasks, _ := s.repo.GetTasks(ctx)
	items, _ := s.repo.GetItems(ctx)
	skills, _ := s.repo.GetSkills(ctx)

	type npcWithBehavior struct {
		models.NPC
		BehaviorState string `json:"behavior_state"`
		BehaviorMood  string `json:"behavior_mood"`
	}

	npcList, _ := s.repo.GetNPCs(ctx)
	npcsWithBehavior := make([]npcWithBehavior, 0, len(npcList))
	for _, npc := range npcList {
		behavior := s.behaviorStore.GetOrCreate(npc.Code, npc.Schedule)
		npcsWithBehavior = append(npcsWithBehavior, npcWithBehavior{
			NPC:           npc,
			BehaviorState: behavior.State,
			BehaviorMood:  behavior.Mood,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"config": configs,
		"scenes": scenes,
		"npcs":   npcsWithBehavior,
		"tasks":  tasks,
		"items":  items,
		"skills": skills,
	})
}

func (s *Server) handleCreatePlayer(c *gin.Context) {
	ctx := c.Request.Context()
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

	existing, _ := s.repo.GetPlayerByAccount(ctx, req.Account)
	if existing != nil && existing.ID > 0 {
		c.JSON(http.StatusOK, gin.H{"data": existing})
		return
	}

	startX, _ := strconv.Atoi(getConfigValue(ctx, s, "start_x", "200"))
	startY, _ := strconv.Atoi(getConfigValue(ctx, s, "start_y", "450"))
	defaultHP, _ := strconv.Atoi(getConfigValue(ctx, s, "default_hp", "100"))
	defaultMP, _ := strconv.Atoi(getConfigValue(ctx, s, "default_mp", "50"))
	defaultGold, _ := strconv.Atoi(getConfigValue(ctx, s, "default_gold", "500"))

	startScene := getConfigValue(ctx, s, "start_scene", "scene_village_entrance")
	startVisited, _ := json.Marshal([]string{startScene})

	player := &models.Player{
		Name:          req.Name,
		Account:       req.Account,
		Level:         1,
		HP:            defaultHP,
		MP:            defaultMP,
		Gold:          defaultGold,
		SceneID:       startScene,
		PosX:          startX,
		PosY:          startY,
		Items:         "{}",
		VisitedScenes: string(startVisited),
	}

	if err := s.repo.CreatePlayer(ctx, player); err != nil {
		respondInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": player})
}

func (s *Server) handleGetPlayer(c *gin.Context) {
	ctx := c.Request.Context()
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	player, err := s.repo.GetPlayerByID(ctx, id)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Player"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": player})
}

func (s *Server) handleUpdatePlayerPos(c *gin.Context) {
	ctx := c.Request.Context()
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	player, err := s.repo.GetPlayerByID(ctx, id)
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
		if player.SceneID != req.SceneID {
			var sceneCodes []string
			if player.VisitedScenes != "" {
				json.Unmarshal([]byte(player.VisitedScenes), &sceneCodes)
			}
			visited := false
			for _, code := range sceneCodes {
				if code == req.SceneID {
					visited = true
					break
				}
			}
			if !visited {
				sceneCodes = append(sceneCodes, req.SceneID)
				data, _ := json.Marshal(sceneCodes)
				player.VisitedScenes = string(data)
			}
		}
		player.SceneID = req.SceneID
	}
	player.PosX = req.PosX
	player.PosY = req.PosY

	if err := s.repo.UpdatePlayer(ctx, player); err != nil {
		respondInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": player})
}

func (s *Server) handleGetSceneByCode(c *gin.Context) {
	ctx := c.Request.Context()
	code := c.Param("code")
	scene, err := s.repo.GetSceneByCode(ctx, code)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Scene"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": scene})
}

func (s *Server) handleGetNPCByCode(c *gin.Context) {
	ctx := c.Request.Context()
	code := c.Param("code")
	npc, err := s.repo.GetNPCByCode(ctx, code)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("NPC"))
		return
	}

	behavior := s.behaviorStore.GetOrCreate(npc.Code, npc.Schedule)

	c.JSON(http.StatusOK, gin.H{
		"data":     npc,
		"behavior": behavior,
	})
}

func (s *Server) handleGetPlayerTasks(c *gin.Context) {
	ctx := c.Request.Context()
	tasks, err := s.repo.GetTasks(ctx)
	if err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tasks})
}

func (s *Server) handleNPCChat(c *gin.Context) {
	ctx := c.Request.Context()
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

	npc, err := s.repo.GetNPCByID(ctx, req.NPCID)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("NPC"))
		return
	}

	behavior := s.behaviorStore.GetOrCreate(npc.Code, npc.Schedule)
	s.behaviorMgr.ReactToPlayer(behavior, req.PlayerID, "talk")

	scenes, _ := s.repo.GetScenesByNPCID(ctx, npc.ID)
	if len(scenes) > 0 {
		s.BroadcastNPCState(npc.ID, npc.Code, npc.Name, scenes[0].Code, behavior.State, 0, 0)
	}

	player, err := s.repo.GetPlayerByID(ctx, req.PlayerID)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Player"))
		return
	}

	history, _ := s.repo.GetConversations(ctx, req.PlayerID, req.NPCID, 10)

	messages := buildChatMessages(npc, player, history, req.Message)

	dialogMood := s.behaviorMgr.GetDialogMood(behavior)
	behaviorContext := s.behaviorMgr.GetDialogContext(behavior)

	var reply string
	if s.chatMgr != nil && s.chatMgr.IsEnabled() && npc.Agent != nil && npc.Agent.SystemPrompt != "" {
		persona := &agent.NPCPersona{
			Name:         npc.Name,
			Title:        npc.Title,
			Description:  npc.Description,
			SystemPrompt: npc.Agent.SystemPrompt + "\n" + behaviorContext,
		}
		playerContext := buildPlayerContext(player)
		chatHistory := convertToChatMessages(history)

		aiReply, err := s.chatMgr.ChatWithNPC(
			ctx,
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
		reply = generateSimpleReply(npc.Agent, messages)
	} else {
		reply = "你好，客官！"
	}

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
	s.repo.CreateConversation(ctx, conv)

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
	s.repo.CreateConversation(ctx, convReply)

	agent.DefaultMemoryStore.AddMessage(req.PlayerID, req.NPCID, "user", req.Message)
	agent.DefaultMemoryStore.AddMessage(req.PlayerID, req.NPCID, "assistant", reply)
	agent.DefaultMemoryStore.UpdatePlayerInfo(req.PlayerID, req.NPCID, player.Name, player.Level)

	c.JSON(http.StatusOK, gin.H{
		"reply":        reply,
		"npc_name":     npc.Name,
		"npc_title":    npc.Title,
		"npc_avatar":   npc.Avatar,
		"ai_powered":   s.chatMgr != nil && s.chatMgr.IsEnabled(),
		"npc_state":    behavior.State,
		"npc_mood":     behavior.Mood,
		"dialog_mood":  dialogMood,
	})
}

func (s *Server) handleGetShopItems(c *gin.Context) {
	ctx := c.Request.Context()
	code := c.Param("code")
	shop, err := s.repo.GetShopByCode(ctx, code)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Shop"))
		return
	}

	type ShopItemDetail struct {
		models.ShopItem
		ItemName        string `json:"item_name"`
		ItemDescription string `json:"item_description"`
		ItemCategory    string `json:"item_category"`
		ItemEffect      string `json:"item_effect"`
	}

	var details []ShopItemDetail
	for _, si := range shop.Items {
		item, err := s.repo.GetItemByID(ctx, si.ItemID)
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

func (s *Server) handleBuyItem(c *gin.Context) {
	ctx := c.Request.Context()
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

	player, err := s.repo.GetPlayerByID(ctx, req.PlayerID)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Player"))
		return
	}

	shop, err := s.repo.GetShopByCode(ctx, req.ShopCode)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Shop"))
		return
	}

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

	if shopItem.Stock < req.Count {
		respondError(c, http.StatusBadRequest, BadRequest("Not enough stock"))
		return
	}

	totalPrice := shopItem.Price * req.Count

	if player.Gold < totalPrice {
		respondError(c, http.StatusBadRequest, BadRequest("Not enough gold"))
		return
	}

	player.Gold -= totalPrice

	var items map[string]int
	json.Unmarshal([]byte(player.Items), &items)
	if items == nil {
		items = make(map[string]int)
	}
	itemKey := strconv.Itoa(int(req.ItemID))
	items[itemKey] += req.Count
	itemsJSON, _ := json.Marshal(items)
	player.Items = string(itemsJSON)

	shopItem.Stock -= req.Count

	s.repo.UpdatePlayer(ctx, player)
	s.repo.SaveShopItem(ctx, shopItem)

	if shop.OwnerNPC != "" {
		npc, err := s.repo.GetNPCByCode(ctx, shop.OwnerNPC)
		if err == nil {
			behavior := s.behaviorStore.GetOrCreate(npc.Code, npc.Schedule)
			s.behaviorMgr.ReactToPlayer(behavior, req.PlayerID, "gift")
			scenes, _ := s.repo.GetScenesByNPCID(ctx, npc.ID)
			if len(scenes) > 0 {
				s.BroadcastNPCState(npc.ID, npc.Code, npc.Name, scenes[0].Code, behavior.State, 0, 0)
			}
		}
	}

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

func getConfigValue(ctx context.Context, s *Server, key, defaultVal string) string {
	cfg, err := s.repo.GetConfig(ctx, key)
	if err != nil {
		return defaultVal
	}
	return cfg.Value
}

func buildChatMessages(npc *models.NPC, player *models.Player, history []models.Conversation, userMsg string) []map[string]string {
	var messages []map[string]string

	if npc.Agent != nil {
		messages = append(messages, map[string]string{
			"role":    "system",
			"content": npc.Agent.SystemPrompt,
		})
	}

	for i := len(history) - 1; i >= 0; i-- {
		messages = append(messages, map[string]string{
			"role":    history[i].Role,
			"content": history[i].Content,
		})
	}

	contextMsg := "【玩家信息】姓名：" + player.Name + "，等级：" + strconv.Itoa(player.Level) + "\n"
	contextMsg += "【玩家消息】" + userMsg

	messages = append(messages, map[string]string{
		"role":    "user",
		"content": contextMsg,
	})

	return messages
}

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

func buildPlayerContext(player *models.Player) string {
	contextMsg := "【玩家信息】姓名：" + player.Name + "，等级：" + strconv.Itoa(player.Level)
	if player.Gold > 0 {
		contextMsg += "，金币：" + strconv.Itoa(player.Gold)
	}
	return contextMsg
}

func generateSimpleReply(agent *models.Agent, messages []map[string]string) string {
	userMsg := ""
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i]["role"] == "user" {
			userMsg = messages[i]["content"]
			break
		}
	}

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

func (s *Server) handleGetNPCBehavior(c *gin.Context) {
	ctx := c.Request.Context()
	code := c.Param("code")
	npc, err := s.repo.GetNPCByCode(ctx, code)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("NPC"))
		return
	}

	behavior := s.behaviorStore.GetOrCreate(npc.Code, npc.Schedule)

	c.JSON(http.StatusOK, gin.H{
		"data": behavior,
	})
}

func (s *Server) handleNPCBehaviorEvent(c *gin.Context) {
	ctx := c.Request.Context()
	code := c.Param("code")
	npc, err := s.repo.GetNPCByCode(ctx, code)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("NPC"))
		return
	}

	var req struct {
		PlayerID uint   `json:"player_id"`
		Action   string `json:"action"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}

	errs := mergeErrors(
		validatePositiveInt("player_id", req.PlayerID),
		validateRequired(map[string]interface{}{"action": req.Action}),
		validateStringIn("action", req.Action, []string{"talk", "gift", "attack"}),
	)
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}

	behavior := s.behaviorStore.GetOrCreate(npc.Code, npc.Schedule)
	s.behaviorMgr.ReactToPlayer(behavior, req.PlayerID, req.Action)

	scenes, _ := s.repo.GetScenesByNPCID(ctx, npc.ID)
	if len(scenes) > 0 {
		s.BroadcastNPCState(npc.ID, npc.Code, npc.Name, scenes[0].Code, behavior.State, 0, 0)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": behavior,
	})
}

func (s *Server) handleGameTick(c *gin.Context) {
	ctx := c.Request.Context()
	var req struct {
		Hour int `json:"hour"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}

	if req.Hour < 0 || req.Hour > 23 {
		respondError(c, http.StatusBadRequest, BadRequest("hour must be 0-23"))
		return
	}

	s.behaviorMgr.UpdateAllBehaviors(s.behaviorStore, req.Hour)

	npcs, _ := s.repo.GetNPCs(ctx)
	for _, npc := range npcs {
		behavior := s.behaviorStore.GetOrCreate(npc.Code, npc.Schedule)
		scenes, _ := s.repo.GetScenesByNPCID(ctx, npc.ID)
		if len(scenes) > 0 {
			s.BroadcastNPCState(npc.ID, npc.Code, npc.Name, scenes[0].Code, behavior.State, 0, 0)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "tick processed",
		"hour":    req.Hour,
	})
}

func (s *Server) handleGetPlayers(c *gin.Context) {
	ctx := c.Request.Context()
	p := parsePagination(c)
	players, total, err := s.repo.GetPlayersPaginated(ctx, p.Offset, p.PageSize)
	if err != nil {
		respondInternalError(c, err)
		return
	}
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
	c.JSON(http.StatusOK, gin.H{"data": players, "total": total})
}

func (s *Server) handleUpdatePlayer(c *gin.Context) {
	ctx := c.Request.Context()
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var player models.Player
	if err := c.ShouldBindJSON(&player); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	errs := validateRequired(map[string]interface{}{"name": player.Name})
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}
	player.ID = id
	if err := s.repo.UpdatePlayer(ctx, &player); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": player})
}

func (s *Server) handleGetConversations(c *gin.Context) {
	ctx := c.Request.Context()
	playerID, ok1 := parseQueryID(c, "player_id")
	if !ok1 {
		return
	}
	npcID, ok2 := parseQueryID(c, "npc_id")
	if !ok2 {
		return
	}
	p := parsePagination(c)

	conversations, total, err := s.repo.GetConversationsPaginated(ctx, playerID, npcID, p.Offset, p.PageSize)
	if err != nil {
		respondInternalError(c, err)
		return
	}
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
	c.JSON(http.StatusOK, gin.H{"data": conversations, "total": total})
}

func (s *Server) handleCreateConversation(c *gin.Context) {
	ctx := c.Request.Context()
	var conv models.Conversation
	if err := c.ShouldBindJSON(&conv); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	errs := mergeErrors(
		validatePositiveInt("player_id", conv.PlayerID),
		validatePositiveInt("npc_id", conv.NPCID),
		validateRequired(map[string]interface{}{"role": conv.Role, "content": conv.Content}),
		validateStringIn("role", conv.Role, []string{"user", "assistant", "system"}),
	)
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}
	if err := s.repo.CreateConversation(ctx, &conv); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": conv})
}

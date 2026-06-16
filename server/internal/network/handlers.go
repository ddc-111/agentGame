package network

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/database/models"
	"github.com/ddc-111/agentGame/server/internal/generator"
	"github.com/ddc-111/agentGame/server/internal/mcp"
)

// ==================== 场景API ====================

func (s *Server) handleGetScenes(c *gin.Context) {
	scenes, err := s.repo.GetScenes()
	if err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": scenes})
}

func (s *Server) handleGetScene(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	scene, err := s.repo.GetSceneByID(uint(id))
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Scene"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": scene})
}

func (s *Server) handleCreateScene(c *gin.Context) {
	var scene models.Scene
	if err := c.ShouldBindJSON(&scene); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	if err := s.repo.CreateScene(&scene); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": scene})
}

func (s *Server) handleUpdateScene(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var scene models.Scene
	if err := c.ShouldBindJSON(&scene); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	scene.ID = uint(id)
	if err := s.repo.UpdateScene(&scene); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": scene})
}

func (s *Server) handleDeleteScene(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := s.repo.DeleteScene(uint(id)); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

// ==================== NPC API ====================

func (s *Server) handleGetNPCs(c *gin.Context) {
	npcs, err := s.repo.GetNPCs()
	if err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": npcs})
}

func (s *Server) handleGetNPC(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	npc, err := s.repo.GetNPCByID(uint(id))
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("NPC"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": npc})
}

func (s *Server) handleCreateNPC(c *gin.Context) {
	var npc models.NPC
	if err := c.ShouldBindJSON(&npc); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	if err := s.repo.CreateNPC(&npc); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": npc})
}

func (s *Server) handleUpdateNPC(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var npc models.NPC
	if err := c.ShouldBindJSON(&npc); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	npc.ID = uint(id)
	if err := s.repo.UpdateNPC(&npc); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": npc})
}

func (s *Server) handleDeleteNPC(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := s.repo.DeleteNPC(uint(id)); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

// ==================== 智能体API ====================

func (s *Server) handleGetAgents(c *gin.Context) {
	agents, err := s.repo.GetAgents()
	if err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": agents})
}

func (s *Server) handleGetAgent(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	agent, err := s.repo.GetAgentByID(uint(id))
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Agent"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": agent})
}

func (s *Server) handleCreateAgent(c *gin.Context) {
	var agent models.Agent
	if err := c.ShouldBindJSON(&agent); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	if err := s.repo.CreateAgent(&agent); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": agent})
}

func (s *Server) handleUpdateAgent(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var agent models.Agent
	if err := c.ShouldBindJSON(&agent); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	agent.ID = uint(id)
	if err := s.repo.UpdateAgent(&agent); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": agent})
}

func (s *Server) handleDeleteAgent(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := s.repo.DeleteAgent(uint(id)); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

// ==================== 大模型提供商API ====================

func (s *Server) handleGetProviders(c *gin.Context) {
	providers, err := s.repo.GetProviders()
	if err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": providers})
}

func (s *Server) handleCreateProvider(c *gin.Context) {
	var provider models.LLMProvider
	if err := c.ShouldBindJSON(&provider); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	if err := s.repo.CreateProvider(&provider); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": provider})
}

func (s *Server) handleUpdateProvider(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var provider models.LLMProvider
	if err := c.ShouldBindJSON(&provider); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	provider.ID = uint(id)
	if err := s.repo.UpdateProvider(&provider); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": provider})
}

func (s *Server) handleDeleteProvider(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := s.repo.DeleteProvider(uint(id)); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

// ==================== 提示词模板API ====================

func (s *Server) handleGetTemplates(c *gin.Context) {
	templates, err := s.repo.GetTemplates()
	if err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": templates})
}

func (s *Server) handleCreateTemplate(c *gin.Context) {
	var template models.PromptTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	if err := s.repo.CreateTemplate(&template); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": template})
}

func (s *Server) handleUpdateTemplate(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var template models.PromptTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	template.ID = uint(id)
	if err := s.repo.UpdateTemplate(&template); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": template})
}

func (s *Server) handleDeleteTemplate(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := s.repo.DeleteTemplate(uint(id)); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

// ==================== 商店API ====================

func (s *Server) handleGetShops(c *gin.Context) {
	shops, err := s.repo.GetShops()
	if err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": shops})
}

func (s *Server) handleGetShop(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	shop, err := s.repo.GetShopByID(uint(id))
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Shop"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": shop})
}

func (s *Server) handleCreateShop(c *gin.Context) {
	var shop models.Shop
	if err := c.ShouldBindJSON(&shop); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	if err := s.repo.CreateShop(&shop); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": shop})
}

func (s *Server) handleUpdateShop(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var shop models.Shop
	if err := c.ShouldBindJSON(&shop); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	shop.ID = uint(id)
	if err := s.repo.UpdateShop(&shop); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": shop})
}

func (s *Server) handleDeleteShop(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := s.repo.DeleteShop(uint(id)); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

// ==================== 道具API ====================

func (s *Server) handleGetItems(c *gin.Context) {
	items, err := s.repo.GetItems()
	if err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (s *Server) handleCreateItem(c *gin.Context) {
	var item models.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	if err := s.repo.CreateItem(&item); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": item})
}

func (s *Server) handleUpdateItem(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var item models.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	item.ID = uint(id)
	if err := s.repo.UpdateItem(&item); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (s *Server) handleDeleteItem(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := s.repo.DeleteItem(uint(id)); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

// ==================== 任务API ====================

func (s *Server) handleGetTasks(c *gin.Context) {
	tasks, err := s.repo.GetTasks()
	if err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tasks})
}

func (s *Server) handleGetTask(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	task, err := s.repo.GetTaskByID(uint(id))
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Task"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": task})
}

func (s *Server) handleCreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	if err := s.repo.CreateTask(&task); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": task})
}

func (s *Server) handleUpdateTask(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	task.ID = uint(id)
	if err := s.repo.UpdateTask(&task); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": task})
}

func (s *Server) handleDeleteTask(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := s.repo.DeleteTask(uint(id)); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

// ==================== 流程API ====================

func (s *Server) handleGetFlows(c *gin.Context) {
	flows, err := s.repo.GetFlows()
	if err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": flows})
}

func (s *Server) handleCreateFlow(c *gin.Context) {
	var flow models.Flow
	if err := c.ShouldBindJSON(&flow); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	if err := s.repo.CreateFlow(&flow); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": flow})
}

func (s *Server) handleUpdateFlow(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var flow models.Flow
	if err := c.ShouldBindJSON(&flow); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	flow.ID = uint(id)
	if err := s.repo.UpdateFlow(&flow); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": flow})
}

func (s *Server) handleDeleteFlow(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := s.repo.DeleteFlow(uint(id)); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

// ==================== 玩家API ====================

func (s *Server) handleGetPlayers(c *gin.Context) {
	players, err := s.repo.GetPlayers()
	if err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": players})
}

func (s *Server) handleUpdatePlayer(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var player models.Player
	if err := c.ShouldBindJSON(&player); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	player.ID = uint(id)
	if err := s.repo.UpdatePlayer(&player); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": player})
}

// ==================== 对话记录API ====================

func (s *Server) handleGetConversations(c *gin.Context) {
	playerID, _ := strconv.ParseUint(c.Query("player_id"), 10, 32)
	npcID, _ := strconv.ParseUint(c.Query("npc_id"), 10, 32)
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	conversations, err := s.repo.GetConversations(uint(playerID), uint(npcID), limit)
	if err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": conversations})
}

func (s *Server) handleCreateConversation(c *gin.Context) {
	var conv models.Conversation
	if err := c.ShouldBindJSON(&conv); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	if err := s.repo.CreateConversation(&conv); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": conv})
}

// ==================== 数据导出导入 ====================

func (s *Server) handleExport(c *gin.Context) {
	data := make(map[string]interface{})

	scenes, _ := s.repo.GetScenes()
	data["scenes"] = scenes

	npcs, _ := s.repo.GetNPCs()
	data["npcs"] = npcs

	agents, _ := s.repo.GetAgents()
	data["agents"] = agents

	shops, _ := s.repo.GetShops()
	data["shops"] = shops

	items, _ := s.repo.GetItems()
	data["items"] = items

	tasks, _ := s.repo.GetTasks()
	data["tasks"] = tasks

	flows, _ := s.repo.GetFlows()
	data["flows"] = flows

	templates, _ := s.repo.GetTemplates()
	data["prompts"] = templates

	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (s *Server) handleImport(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}

	// Unwrap "data" key if present (matches export format {"data": {...}})
	if inner, ok := data["data"].(map[string]interface{}); ok {
		data = inner
	}

	imported := make(map[string]int)

	if scenesRaw, ok := data["scenes"].([]interface{}); ok {
		for _, sr := range scenesRaw {
			sceneData, _ := jsonMarshal(sr)
			var scene models.Scene
			if jsonUnmarshal(sceneData, &scene) == nil {
				if existing, _ := s.repo.GetSceneByCode(scene.Code); existing != nil && existing.ID > 0 {
					scene.ID = existing.ID
					s.repo.UpdateScene(&scene)
				} else {
					scene.ID = 0
					s.repo.CreateScene(&scene)
				}
				imported["scenes"]++
			}
		}
	}

	if npcsRaw, ok := data["npcs"].([]interface{}); ok {
		for _, nr := range npcsRaw {
			npcData, _ := jsonMarshal(nr)
			var npc models.NPC
			if jsonUnmarshal(npcData, &npc) == nil {
				if existing, _ := s.repo.GetNPCByCode(npc.Code); existing != nil && existing.ID > 0 {
					npc.ID = existing.ID
					s.repo.UpdateNPC(&npc)
				} else {
					npc.ID = 0
					s.repo.CreateNPC(&npc)
				}
				imported["npcs"]++
			}
		}
	}

	if agentsRaw, ok := data["agents"].([]interface{}); ok {
		for _, ar := range agentsRaw {
			agentData, _ := jsonMarshal(ar)
			var agent models.Agent
			if jsonUnmarshal(agentData, &agent) == nil {
				if existing, _ := s.repo.GetAgentByCode(agent.Code); existing != nil && existing.ID > 0 {
					agent.ID = existing.ID
					s.repo.UpdateAgent(&agent)
				} else {
					agent.ID = 0
					s.repo.CreateAgent(&agent)
				}
				imported["agents"]++
			}
		}
	}

	if shopsRaw, ok := data["shops"].([]interface{}); ok {
		for _, sr := range shopsRaw {
			shopData, _ := jsonMarshal(sr)
			var shop models.Shop
			if jsonUnmarshal(shopData, &shop) == nil {
				if existing, _ := s.repo.GetShopByCode(shop.Code); existing != nil && existing.ID > 0 {
					shop.ID = existing.ID
					s.repo.UpdateShop(&shop)
				} else {
					shop.ID = 0
					s.repo.CreateShop(&shop)
				}
				imported["shops"]++
			}
		}
	}

	if itemsRaw, ok := data["items"].([]interface{}); ok {
		for _, ir := range itemsRaw {
			itemData, _ := jsonMarshal(ir)
			var item models.Item
			if jsonUnmarshal(itemData, &item) == nil {
				if existing, _ := s.repo.GetItemByCode(item.Code); existing != nil && existing.ID > 0 {
					item.ID = existing.ID
					s.repo.UpdateItem(&item)
				} else {
					item.ID = 0
					s.repo.CreateItem(&item)
				}
				imported["items"]++
			}
		}
	}

	if tasksRaw, ok := data["tasks"].([]interface{}); ok {
		for _, tr := range tasksRaw {
			taskData, _ := jsonMarshal(tr)
			var task models.Task
			if jsonUnmarshal(taskData, &task) == nil {
				if existing, _ := s.repo.GetTaskByCode(task.Code); existing != nil && existing.ID > 0 {
					task.ID = existing.ID
					s.repo.UpdateTask(&task)
				} else {
					task.ID = 0
					s.repo.CreateTask(&task)
				}
				imported["tasks"]++
			}
		}
	}

	if flowsRaw, ok := data["flows"].([]interface{}); ok {
		for _, fr := range flowsRaw {
			flowData, _ := jsonMarshal(fr)
			var flow models.Flow
			if jsonUnmarshal(flowData, &flow) == nil {
				if existing, _ := s.repo.GetFlowByCode(flow.Code); existing != nil && existing.ID > 0 {
					flow.ID = existing.ID
					s.repo.UpdateFlow(&flow)
				} else {
					flow.ID = 0
					s.repo.CreateFlow(&flow)
				}
				imported["flows"]++
			}
		}
	}

	if templatesRaw, ok := data["prompts"].([]interface{}); ok {
		for _, tr := range templatesRaw {
			tmplData, _ := jsonMarshal(tr)
			var tmpl models.PromptTemplate
			if jsonUnmarshal(tmplData, &tmpl) == nil {
				if existing, _ := s.repo.GetTemplateByCode(tmpl.Code); existing != nil && existing.ID > 0 {
					tmpl.ID = existing.ID
					s.repo.UpdateTemplate(&tmpl)
				} else {
					tmpl.ID = 0
					s.repo.CreateTemplate(&tmpl)
				}
				imported["prompts"]++
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Import successful", "imported": imported})
}

// ==================== 生成智能体API ====================

func (s *Server) handleGenerate(c *gin.Context) {
	var req generator.GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}

	if s.generator == nil || !s.generator.IsEnabled() {
		respondError(c, http.StatusServiceUnavailable, BadRequest("生成智能体未启用，请检查配置"))
		return
	}

	resp, err := s.generator.Generate(c.Request.Context(), req)
	if err != nil {
		respondInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (s *Server) handleGeneratorStatus(c *gin.Context) {
	enabled := s.generator != nil && s.generator.IsEnabled()
	cfg := s.generator.GetConfig()

	c.JSON(http.StatusOK, gin.H{
		"enabled":  enabled,
		"provider": cfg.Provider,
		"model":    cfg.Model,
		"base_url": cfg.BaseURL,
	})
}

func (s *Server) handleGeneratorTest(c *gin.Context) {
	if s.generator == nil || !s.generator.IsEnabled() {
		respondError(c, http.StatusServiceUnavailable, BadRequest("生成智能体未启用"))
		return
	}

	// 测试生成一个简单的NPC
	req := generator.GenerateRequest{
		Type:   "npc",
		Action: "create",
		Params: map[string]interface{}{
			"description": "一个卖包子的老大爷",
			"theme":       "古风小镇",
		},
	}

	resp, err := s.generator.Generate(c.Request.Context(), req)
	if err != nil {
		respondInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ==================== MCP REST API ====================

func (s *Server) handleMCPTools(c *gin.Context) {
	tools := s.mcp.GetTools()
	c.JSON(http.StatusOK, gin.H{"tools": tools})
}

func (s *Server) handleMCPCall(c *gin.Context) {
	var req struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}

	result := s.mcp.HandleRequest(c.Request.Context(), mcp.MCPRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "tools/call",
		Params: map[string]interface{}{
			"name":      req.Name,
			"arguments": req.Arguments,
		},
	})

	c.JSON(http.StatusOK, result)
}

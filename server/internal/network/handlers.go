package network

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/database/models"
)

// ==================== 场景API ====================

func (s *Server) handleGetScenes(c *gin.Context) {
	scenes, err := s.repo.GetScenes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": scenes})
}

func (s *Server) handleGetScene(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	scene, err := s.repo.GetSceneByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Scene not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": scene})
}

func (s *Server) handleCreateScene(c *gin.Context) {
	var scene models.Scene
	if err := c.ShouldBindJSON(&scene); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.repo.CreateScene(&scene); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": scene})
}

func (s *Server) handleUpdateScene(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var scene models.Scene
	if err := c.ShouldBindJSON(&scene); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	scene.ID = uint(id)
	if err := s.repo.UpdateScene(&scene); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": scene})
}

func (s *Server) handleDeleteScene(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := s.repo.DeleteScene(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

// ==================== NPC API ====================

func (s *Server) handleGetNPCs(c *gin.Context) {
	npcs, err := s.repo.GetNPCs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": npcs})
}

func (s *Server) handleGetNPC(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	npc, err := s.repo.GetNPCByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "NPC not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": npc})
}

func (s *Server) handleCreateNPC(c *gin.Context) {
	var npc models.NPC
	if err := c.ShouldBindJSON(&npc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.repo.CreateNPC(&npc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": npc})
}

func (s *Server) handleUpdateNPC(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var npc models.NPC
	if err := c.ShouldBindJSON(&npc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	npc.ID = uint(id)
	if err := s.repo.UpdateNPC(&npc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": npc})
}

func (s *Server) handleDeleteNPC(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := s.repo.DeleteNPC(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

// ==================== 智能体API ====================

func (s *Server) handleGetAgents(c *gin.Context) {
	agents, err := s.repo.GetAgents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": agents})
}

func (s *Server) handleGetAgent(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	agent, err := s.repo.GetAgentByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": agent})
}

func (s *Server) handleCreateAgent(c *gin.Context) {
	var agent models.Agent
	if err := c.ShouldBindJSON(&agent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.repo.CreateAgent(&agent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": agent})
}

func (s *Server) handleUpdateAgent(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var agent models.Agent
	if err := c.ShouldBindJSON(&agent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	agent.ID = uint(id)
	if err := s.repo.UpdateAgent(&agent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": agent})
}

func (s *Server) handleDeleteAgent(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := s.repo.DeleteAgent(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

// ==================== 大模型提供商API ====================

func (s *Server) handleGetProviders(c *gin.Context) {
	providers, err := s.repo.GetProviders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": providers})
}

func (s *Server) handleCreateProvider(c *gin.Context) {
	var provider models.LLMProvider
	if err := c.ShouldBindJSON(&provider); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.repo.CreateProvider(&provider); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": provider})
}

func (s *Server) handleUpdateProvider(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var provider models.LLMProvider
	if err := c.ShouldBindJSON(&provider); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	provider.ID = uint(id)
	if err := s.repo.UpdateProvider(&provider); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": provider})
}

func (s *Server) handleDeleteProvider(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := s.repo.DeleteProvider(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

// ==================== 提示词模板API ====================

func (s *Server) handleGetTemplates(c *gin.Context) {
	templates, err := s.repo.GetTemplates()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": templates})
}

func (s *Server) handleCreateTemplate(c *gin.Context) {
	var template models.PromptTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.repo.CreateTemplate(&template); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": template})
}

func (s *Server) handleUpdateTemplate(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var template models.PromptTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	template.ID = uint(id)
	if err := s.repo.UpdateTemplate(&template); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": template})
}

func (s *Server) handleDeleteTemplate(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := s.repo.DeleteTemplate(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

// ==================== 商店API ====================

func (s *Server) handleGetShops(c *gin.Context) {
	shops, err := s.repo.GetShops()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": shops})
}

func (s *Server) handleGetShop(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	shop, err := s.repo.GetShopByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shop not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": shop})
}

func (s *Server) handleCreateShop(c *gin.Context) {
	var shop models.Shop
	if err := c.ShouldBindJSON(&shop); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.repo.CreateShop(&shop); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": shop})
}

func (s *Server) handleUpdateShop(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var shop models.Shop
	if err := c.ShouldBindJSON(&shop); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	shop.ID = uint(id)
	if err := s.repo.UpdateShop(&shop); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": shop})
}

func (s *Server) handleDeleteShop(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := s.repo.DeleteShop(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

// ==================== 道具API ====================

func (s *Server) handleGetItems(c *gin.Context) {
	items, err := s.repo.GetItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (s *Server) handleCreateItem(c *gin.Context) {
	var item models.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.repo.CreateItem(&item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": item})
}

func (s *Server) handleUpdateItem(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var item models.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item.ID = uint(id)
	if err := s.repo.UpdateItem(&item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (s *Server) handleDeleteItem(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := s.repo.DeleteItem(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

// ==================== 任务API ====================

func (s *Server) handleGetTasks(c *gin.Context) {
	tasks, err := s.repo.GetTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tasks})
}

func (s *Server) handleGetTask(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	task, err := s.repo.GetTaskByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": task})
}

func (s *Server) handleCreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.repo.CreateTask(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": task})
}

func (s *Server) handleUpdateTask(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task.ID = uint(id)
	if err := s.repo.UpdateTask(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": task})
}

func (s *Server) handleDeleteTask(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := s.repo.DeleteTask(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

// ==================== 流程API ====================

func (s *Server) handleGetFlows(c *gin.Context) {
	flows, err := s.repo.GetFlows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": flows})
}

func (s *Server) handleCreateFlow(c *gin.Context) {
	var flow models.Flow
	if err := c.ShouldBindJSON(&flow); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.repo.CreateFlow(&flow); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": flow})
}

func (s *Server) handleUpdateFlow(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var flow models.Flow
	if err := c.ShouldBindJSON(&flow); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	flow.ID = uint(id)
	if err := s.repo.UpdateFlow(&flow); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": flow})
}

func (s *Server) handleDeleteFlow(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := s.repo.DeleteFlow(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

// ==================== 玩家API ====================

func (s *Server) handleGetPlayers(c *gin.Context) {
	players, err := s.repo.GetPlayers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": players})
}

func (s *Server) handleCreatePlayer(c *gin.Context) {
	var player models.Player
	if err := c.ShouldBindJSON(&player); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.repo.CreatePlayer(&player); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": player})
}

func (s *Server) handleUpdatePlayer(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var player models.Player
	if err := c.ShouldBindJSON(&player); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	player.ID = uint(id)
	if err := s.repo.UpdatePlayer(&player); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": conversations})
}

func (s *Server) handleCreateConversation(c *gin.Context) {
	var conv models.Conversation
	if err := c.ShouldBindJSON(&conv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.repo.CreateConversation(&conv); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO: 实现数据导入逻辑
	c.JSON(http.StatusOK, gin.H{"message": "Import successful"})
}

package network

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/agent"
	"github.com/ddc-111/agentGame/server/internal/config"
	"github.com/ddc-111/agentGame/server/internal/database"
	"github.com/ddc-111/agentGame/server/internal/database/models"
	"github.com/ddc-111/agentGame/server/internal/database/repository"
	"github.com/ddc-111/agentGame/server/internal/game"
	"github.com/ddc-111/agentGame/server/internal/generator"
	"github.com/ddc-111/agentGame/server/internal/mcp"
)

type Server struct {
	cfg           *config.Config
	router        *gin.Engine
	http          *http.Server
	db            *database.Database
	repo          *repository.Repository
	generator     *generator.Generator
	mcp           *mcp.Server
	chatMgr       *agent.ChatManager
	hub           *Hub
	behaviorMgr   *game.NPCBehaviorManager
	behaviorStore *game.NPCBehaviorStore
	cancel        context.CancelFunc
}

func NewServer(cfg *config.Config) *Server {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化数据库
	db, err := database.New(database.Config{
		Driver:   cfg.Database.Driver,
		DSN:      cfg.Database.DSN,
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
	})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// 自动迁移
	err = db.AutoMigrate(
		&models.Scene{},
		&models.SceneNPC{},
		&models.Portal{},
		&models.NPC{},
		&models.Agent{},
		&models.LLMProvider{},
		&models.PromptTemplate{},
		&models.Shop{},
		&models.ShopItem{},
		&models.Item{},
		&models.Task{},
		&models.Flow{},
		&models.GameConfig{},
		&models.Player{},
		&models.Conversation{},
		&models.SaveGame{},
		&models.Skill{},
		&models.Achievement{},
		&models.PlayerAchievement{},
		&models.GMUser{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 初始化种子数据
	if err := database.SeedData(db.DB); err != nil {
		log.Printf("Warning: Failed to seed data: %v", err)
	}

	repo := repository.New(db.DB)

	// 初始化生成智能体
	gen, err := generator.New(cfg.Generator)
	if err != nil {
		log.Printf("Warning: Failed to create generator: %v", err)
	}

	// 初始化AI对话管理器
	chatMgr := agent.NewChatManager(cfg.AI)
	if chatMgr.IsEnabled() {
		log.Println("AI chat manager enabled")
	} else {
		log.Println("AI chat manager disabled, using simple replies")
	}

	hub := NewHub()
	go hub.Run()

	mcpServer := mcp.New(repo, gen)

	// 演示场景在seed.go中已包含

	router := gin.Default()

	s := &Server{
		cfg:           cfg,
		router:        router,
		db:            db,
		repo:          repo,
		generator:     gen,
		mcp:           mcpServer,
		chatMgr:       chatMgr,
		hub:           hub,
		behaviorMgr:   game.NewNPCBehaviorManager(),
		behaviorStore: game.NewNPCBehaviorStore(),
	}

	s.setupRoutes()
	s.initNPCBehaviors()
	s.startGameLoop()
	return s
}

func (s *Server) setupRoutes() {
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	s.router.Use(CORSMiddleware(s.cfg.CORS.AllowedOrigins))

	// MCP端点
	s.router.POST("/mcp", func(c *gin.Context) {
		s.mcp.HandleHTTP(c.Writer, c.Request)
	})

	// MCP SSE端点（用于流式响应）
	s.router.GET("/mcp/sse", func(c *gin.Context) {
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.SSEvent("message", map[string]string{"status": "connected"})
	})

	api := s.router.Group("/api")
	api.Use(RateLimitMiddleware(100, 200))
	{
		// WebSocket
		api.GET("/ws", s.handleWebSocket)

		// GM登录
		api.POST("/gm/login", s.handleGMLogin)

		// GM受保护路由
		gm := api.Group("/gm")
		gm.Use(AuthMiddleware(s.cfg.Auth.JWTSecret))
		{
			gm.GET("/me", s.handleGMMe)
		}

		// 生成智能体API
		api.POST("/generator/generate", s.handleGenerate)
		api.GET("/generator/status", s.handleGeneratorStatus)
		api.POST("/generator/test", s.handleGeneratorTest)

		// MCP工具列表（兼容REST访问）
		api.GET("/mcp/tools", s.handleMCPTools)
		api.POST("/mcp/call", s.handleMCPCall)

		// 场景API
		api.GET("/scenes", s.handleGetScenes)
		api.GET("/scenes/:id", s.handleGetScene)
		api.POST("/scenes", s.handleCreateScene)
		api.PUT("/scenes/:id", s.handleUpdateScene)
		api.DELETE("/scenes/:id", s.handleDeleteScene)

		// NPC API
		api.GET("/npcs", s.handleGetNPCs)
		api.GET("/npcs/:id", s.handleGetNPC)
		api.POST("/npcs", s.handleCreateNPC)
		api.PUT("/npcs/:id", s.handleUpdateNPC)
		api.DELETE("/npcs/:id", s.handleDeleteNPC)

		// 智能体API
		api.GET("/agents", s.handleGetAgents)
		api.GET("/agents/:id", s.handleGetAgent)
		api.POST("/agents", s.handleCreateAgent)
		api.PUT("/agents/:id", s.handleUpdateAgent)
		api.DELETE("/agents/:id", s.handleDeleteAgent)

		// 大模型提供商API
		api.GET("/llm/providers", s.handleGetProviders)
		api.POST("/llm/providers", s.handleCreateProvider)
		api.PUT("/llm/providers/:id", s.handleUpdateProvider)
		api.DELETE("/llm/providers/:id", s.handleDeleteProvider)

		// 提示词模板API
		api.GET("/prompts", s.handleGetTemplates)
		api.POST("/prompts", s.handleCreateTemplate)
		api.PUT("/prompts/:id", s.handleUpdateTemplate)
		api.DELETE("/prompts/:id", s.handleDeleteTemplate)

		// 商店API
		api.GET("/shops", s.handleGetShops)
		api.GET("/shops/:id", s.handleGetShop)
		api.POST("/shops", s.handleCreateShop)
		api.PUT("/shops/:id", s.handleUpdateShop)
		api.DELETE("/shops/:id", s.handleDeleteShop)

		// 道具API
		api.GET("/items", s.handleGetItems)
		api.POST("/items", s.handleCreateItem)
		api.PUT("/items/:id", s.handleUpdateItem)
		api.DELETE("/items/:id", s.handleDeleteItem)

		// 任务API
		api.GET("/tasks", s.handleGetTasks)
		api.GET("/tasks/:id", s.handleGetTask)
		api.POST("/tasks", s.handleCreateTask)
		api.PUT("/tasks/:id", s.handleUpdateTask)
		api.DELETE("/tasks/:id", s.handleDeleteTask)

		// 流程API
		api.GET("/flows", s.handleGetFlows)
		api.GET("/flows/:id", s.handleGetFlow)
		api.POST("/flows", s.handleCreateFlow)
		api.PUT("/flows/:id", s.handleUpdateFlow)
		api.DELETE("/flows/:id", s.handleDeleteFlow)

		// 玩家API
		api.GET("/players", s.handleGetPlayers)
		api.POST("/players", s.handleCreatePlayer)
		api.PUT("/players/:id", s.handleUpdatePlayer)

		// 对话记录API
		api.GET("/conversations", s.handleGetConversations)
		api.POST("/conversations", s.handleCreateConversation)

		// 数据导出导入
		api.GET("/export", s.handleExport)
		api.POST("/import", s.handleImport)

		// 游戏API
		api.GET("/game/init", s.handleGetGameInit)
		api.GET("/game/scene/:code", s.handleGetSceneByCode)
		api.GET("/game/npc/:code", s.handleGetNPCByCode)
		api.GET("/game/shop/:code/items", s.handleGetShopItems)
		api.POST("/game/tick", s.handleGameTick)

		// NPC行为API
		api.GET("/npc/:code/behavior", s.handleGetNPCBehavior)
		api.POST("/npc/:code/behavior/event", s.handleNPCBehaviorEvent)

		// 玩家API
		api.POST("/player/create", s.handleCreatePlayer)
		api.GET("/player/:id", s.handleGetPlayer)
		api.PUT("/player/:id/pos", s.handleUpdatePlayerPos)
		api.GET("/player/:id/tasks", s.handleGetPlayerTasks)

		// NPC对话
		chat := api.Group("/npc/chat")
		chat.Use(RateLimitMiddleware(10, 20))
		chat.POST("", s.handleNPCChat)

		// 商店购买
		api.POST("/shop/buy", s.handleBuyItem)

		s.registerCombatRoutes(api)
		s.registerInventoryRoutes(api)
		s.registerSavegameRoutes(api)
		s.registerSkillRoutes(api)
		s.registerAchievementRoutes(api)
	}
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.cfg.Server.Port)
	s.http = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	log.Printf("Starting server on %s", addr)
	log.Printf("Generator enabled: %v", s.generator.IsEnabled())
	log.Printf("MCP endpoint: http://localhost:%d/mcp", s.cfg.Server.Port)
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown() error {
	if s.cancel != nil {
		s.cancel()
	}
	if s.db != nil {
		s.db.Close()
	}
	return s.http.Shutdown(context.Background())
}

func (s *Server) initNPCBehaviors() {
	npcs, err := s.repo.GetNPCs()
	if err != nil {
		log.Printf("Failed to load NPCs for behavior init: %v", err)
		return
	}
	for _, npc := range npcs {
		s.behaviorStore.GetOrCreate(npc.Code, npc.Schedule)
	}
	log.Printf("Initialized behaviors for %d NPCs", len(npcs))
}

func (s *Server) startGameLoop() {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		hour := time.Now().Hour()
		s.behaviorMgr.UpdateAllBehaviors(s.behaviorStore, hour)

		for {
			select {
			case <-ctx.Done():
				return
			case t := <-ticker.C:
				newHour := t.Hour()
				if newHour != hour {
					hour = newHour
					s.behaviorMgr.UpdateAllBehaviors(s.behaviorStore, hour)
					s.broadcastAllNPCStates()
				}
			}
		}
	}()
}

func (s *Server) broadcastAllNPCStates() {
	npcs, err := s.repo.GetNPCs()
	if err != nil {
		return
	}
	for _, npc := range npcs {
		behavior := s.behaviorStore.Get(npc.Code)
		if behavior == nil {
			continue
		}
		scenes, _ := s.repo.GetScenesByNPCID(npc.ID)
		if len(scenes) > 0 {
			s.BroadcastNPCState(npc.ID, npc.Code, npc.Name, scenes[0].Code, behavior.State, 0, 0)
		}
	}
}

// GetRouter 获取路由器（用于测试）
func (s *Server) GetRouter() *gin.Engine {
	return s.router
}

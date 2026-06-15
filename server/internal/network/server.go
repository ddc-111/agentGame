package network

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/config"
	"github.com/ddc-111/agentGame/server/internal/database"
	"github.com/ddc-111/agentGame/server/internal/database/models"
	"github.com/ddc-111/agentGame/server/internal/database/repository"
)

type Server struct {
	cfg    *config.Config
	router *gin.Engine
	http   *http.Server
	db     *database.Database
	repo   *repository.Repository
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
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 初始化种子数据
	if err := database.SeedData(db.DB); err != nil {
		log.Printf("Warning: Failed to seed data: %v", err)
	}

	repo := repository.New(db.DB)

	router := gin.Default()

	s := &Server{
		cfg:    cfg,
		router: router,
		db:     db,
		repo:   repo,
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := s.router.Group("/api")
	{
		// WebSocket
		api.GET("/ws", s.handleWebSocket)

		// GM登录
		api.POST("/gm/login", s.handleGMLogin)

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
	}
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.cfg.Server.Port)
	s.http = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	log.Printf("Starting server on %s", addr)
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown() error {
	if s.db != nil {
		s.db.Close()
	}
	return s.http.Shutdown(context.Background())
}

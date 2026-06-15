package network

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/config"
)

type Server struct {
	cfg    *config.Config
	router *gin.Engine
	http   *http.Server
}

func NewServer(cfg *config.Config) *Server {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	s := &Server{
		cfg:    cfg,
		router: router,
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
		api.GET("/ws", s.handleWebSocket)
		api.POST("/gm/login", s.handleGMLogin)
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
	return s.http.Shutdown(context.Background())
}

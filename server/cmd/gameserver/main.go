package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ddc-111/agentGame/server/internal/config"
	"github.com/ddc-111/agentGame/server/internal/network"
)

// @title           AgentGame API
// @version         1.0
// @description     AgentGame server API for game management, NPCs, players, combat, and more.
// @host            localhost:8080
// @BasePath        /api
// @schemes         http
// @produce         json
// @consumes        json
// @securityDefinitions.apikey BearerAuth
// @in                         header
// @name                       Authorization
// @description                JWT token. Format: "Bearer {token}"

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Printf("Failed to load config: %v, using defaults", err)
		cfg = config.Default()
	}

	server := network.NewServer(cfg)

	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	log.Printf("Game server started on port %d", cfg.Server.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	_ = server.Shutdown()
}

package main

import (
	"errors"
	"log"
	"net/http"
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
		if !os.IsNotExist(err) {
			log.Fatalf("Failed to load config: %v", err)
		}
		log.Printf("Config file not found, using development defaults")
		cfg = config.Default()
		if err := cfg.Validate(); err != nil {
			log.Fatalf("Invalid default config: %v", err)
		}
	}

	server := network.NewServer(cfg)
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- server.Start()
	}()

	log.Printf("Game server started on port %d", cfg.Server.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	select {
	case sig := <-quit:
		log.Printf("Received %s, shutting down server...", sig)
	case err := <-serverErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %v", err)
		}
		return
	}

	if err := server.Shutdown(); err != nil {
		log.Printf("Server shutdown completed with error: %v", err)
	}
	if err := <-serverErr; err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("Server stopped with error: %v", err)
	}
}

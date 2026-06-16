package network

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/database/models"
)

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
	errs := validateRequired(map[string]interface{}{"name": player.Name})
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}
	player.ID = uint(id)
	if err := s.repo.UpdatePlayer(&player); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": player})
}

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
	if err := s.repo.CreateConversation(&conv); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": conv})
}

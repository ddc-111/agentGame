package network

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/database/models"
)

func (s *Server) handleGetPlayers(c *gin.Context) {
	p := parsePagination(c)
	players, total, err := s.repo.GetPlayersPaginated(p.Offset, p.PageSize)
	if err != nil {
		respondInternalError(c, err)
		return
	}
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
	c.JSON(http.StatusOK, gin.H{"data": players, "total": total})
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
	p := parsePagination(c)

	conversations, total, err := s.repo.GetConversationsPaginated(uint(playerID), uint(npcID), p.Offset, p.PageSize)
	if err != nil {
		respondInternalError(c, err)
		return
	}
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
	c.JSON(http.StatusOK, gin.H{"data": conversations, "total": total})
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

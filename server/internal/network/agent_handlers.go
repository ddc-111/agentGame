package network

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/database/models"
)

func (s *Server) handleGetAgents(c *gin.Context) {
	ctx := c.Request.Context()
	p := parsePagination(c)
	agents, total, err := s.repo.GetAgentsPaginated(ctx, p.Offset, p.PageSize)
	if err != nil {
		respondInternalError(c, err)
		return
	}
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
	c.JSON(http.StatusOK, gin.H{"data": agents, "total": total})
}

func (s *Server) handleGetAgent(c *gin.Context) {
	ctx := c.Request.Context()
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	agent, err := s.repo.GetAgentByID(ctx, id)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Agent"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": agent})
}

func (s *Server) handleCreateAgent(c *gin.Context) {
	ctx := c.Request.Context()
	var agent models.Agent
	if err := c.ShouldBindJSON(&agent); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	errs := mergeErrors(
		validateRequired(map[string]interface{}{"name": agent.Name, "code": agent.Code}),
		validateIntRange("max_tokens", agent.MaxTokens, 1, 100000),
		validateIntRange("max_messages", agent.MaxMessages, 1, 1000),
	)
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}
	if err := s.repo.CreateAgent(ctx, &agent); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": agent})
}

func (s *Server) handleUpdateAgent(c *gin.Context) {
	ctx := c.Request.Context()
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var agent models.Agent
	if err := c.ShouldBindJSON(&agent); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	errs := mergeErrors(
		validateRequired(map[string]interface{}{"name": agent.Name, "code": agent.Code}),
		validateIntRange("max_tokens", agent.MaxTokens, 1, 100000),
		validateIntRange("max_messages", agent.MaxMessages, 1, 1000),
	)
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}
	agent.ID = id
	if err := s.repo.UpdateAgent(ctx, &agent); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": agent})
}

func (s *Server) handleDeleteAgent(c *gin.Context) {
	ctx := c.Request.Context()
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := s.repo.DeleteAgent(ctx, id); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

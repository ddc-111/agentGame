package network

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/database/models"
)

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
	errs := mergeErrors(
		validateRequired(map[string]interface{}{"name": agent.Name, "code": agent.Code}),
		validateIntRange("max_tokens", agent.MaxTokens, 1, 100000),
		validateIntRange("max_messages", agent.MaxMessages, 1, 1000),
	)
	if len(errs) > 0 {
		respondValidation(c, errs)
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
	errs := mergeErrors(
		validateRequired(map[string]interface{}{"name": agent.Name, "code": agent.Code}),
		validateIntRange("max_tokens", agent.MaxTokens, 1, 100000),
		validateIntRange("max_messages", agent.MaxMessages, 1, 1000),
	)
	if len(errs) > 0 {
		respondValidation(c, errs)
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

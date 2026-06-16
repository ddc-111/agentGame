package network

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/database/models"
)

func (s *Server) handleGetFlows(c *gin.Context) {
	flows, err := s.repo.GetFlows()
	if err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": flows})
}

func (s *Server) handleCreateFlow(c *gin.Context) {
	var flow models.Flow
	if err := c.ShouldBindJSON(&flow); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	errs := validateRequired(map[string]interface{}{"name": flow.Name, "code": flow.Code})
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}
	if err := s.repo.CreateFlow(&flow); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": flow})
}

func (s *Server) handleUpdateFlow(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var flow models.Flow
	if err := c.ShouldBindJSON(&flow); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	errs := validateRequired(map[string]interface{}{"name": flow.Name, "code": flow.Code})
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}
	flow.ID = uint(id)
	if err := s.repo.UpdateFlow(&flow); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": flow})
}

func (s *Server) handleDeleteFlow(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := s.repo.DeleteFlow(uint(id)); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

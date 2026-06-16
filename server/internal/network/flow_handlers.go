package network

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/database/models"
)

func (s *Server) handleGetFlow(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	flow, err := s.repo.GetFlowByID(uint(id))
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Flow"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": flow})
}

func (s *Server) handleGetFlows(c *gin.Context) {
	p := parsePagination(c)
	flows, total, err := s.repo.GetFlowsPaginated(p.Offset, p.PageSize)
	if err != nil {
		respondInternalError(c, err)
		return
	}
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
	c.JSON(http.StatusOK, gin.H{"data": flows, "total": total})
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

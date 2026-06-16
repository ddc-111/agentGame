package network

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/database/models"
)

func (s *Server) handleGetProviders(c *gin.Context) {
	p := parsePagination(c)
	providers, total, err := s.repo.GetProvidersPaginated(p.Offset, p.PageSize)
	if err != nil {
		respondInternalError(c, err)
		return
	}
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
	c.JSON(http.StatusOK, gin.H{"data": providers, "total": total})
}

func (s *Server) handleCreateProvider(c *gin.Context) {
	var provider models.LLMProvider
	if err := c.ShouldBindJSON(&provider); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	errs := validateRequired(map[string]interface{}{"name": provider.Name, "code": provider.Code})
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}
	if err := s.repo.CreateProvider(&provider); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": provider})
}

func (s *Server) handleUpdateProvider(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var provider models.LLMProvider
	if err := c.ShouldBindJSON(&provider); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	errs := validateRequired(map[string]interface{}{"name": provider.Name, "code": provider.Code})
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}
	provider.ID = uint(id)
	if err := s.repo.UpdateProvider(&provider); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": provider})
}

func (s *Server) handleDeleteProvider(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := s.repo.DeleteProvider(uint(id)); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

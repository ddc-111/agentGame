package network

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/database/models"
)

func (s *Server) handleGetTemplates(c *gin.Context) {
	templates, err := s.repo.GetTemplates()
	if err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": templates})
}

func (s *Server) handleCreateTemplate(c *gin.Context) {
	var template models.PromptTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	errs := validateRequired(map[string]interface{}{"name": template.Name, "code": template.Code, "content": template.Content})
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}
	if err := s.repo.CreateTemplate(&template); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": template})
}

func (s *Server) handleUpdateTemplate(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var template models.PromptTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	errs := validateRequired(map[string]interface{}{"name": template.Name, "code": template.Code, "content": template.Content})
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}
	template.ID = uint(id)
	if err := s.repo.UpdateTemplate(&template); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": template})
}

func (s *Server) handleDeleteTemplate(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := s.repo.DeleteTemplate(uint(id)); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

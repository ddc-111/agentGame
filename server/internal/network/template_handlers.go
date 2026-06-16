package network

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/database/models"
)

func (s *Server) handleGetTemplates(c *gin.Context) {
	ctx := c.Request.Context()
	p := parsePagination(c)
	templates, total, err := s.repo.GetTemplatesPaginated(ctx, p.Offset, p.PageSize)
	if err != nil {
		respondInternalError(c, err)
		return
	}
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
	c.JSON(http.StatusOK, gin.H{"data": templates, "total": total})
}

func (s *Server) handleCreateTemplate(c *gin.Context) {
	ctx := c.Request.Context()
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
	if err := s.repo.CreateTemplate(ctx, &template); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": template})
}

func (s *Server) handleUpdateTemplate(c *gin.Context) {
	ctx := c.Request.Context()
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
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
	template.ID = id
	if err := s.repo.UpdateTemplate(ctx, &template); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": template})
}

func (s *Server) handleDeleteTemplate(c *gin.Context) {
	ctx := c.Request.Context()
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := s.repo.DeleteTemplate(ctx, id); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

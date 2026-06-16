package network

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/database/models"
)

// handleGetTemplates godoc
// @Summary      List prompt templates
// @Description  Get paginated list of prompt templates
// @Tags         templates
// @Accept       json
// @Produce      json
// @Param        page    query  int  false  "Page number"  default(1)
// @Param        page_size  query  int  false  "Page size"  default(20)
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /prompts [get]
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

// handleCreateTemplate godoc
// @Summary      Create a prompt template
// @Description  Create a new prompt template
// @Tags         templates
// @Accept       json
// @Produce      json
// @Param        template  body  models.PromptTemplate  true  "Template data"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /prompts [post]
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

// handleUpdateTemplate godoc
// @Summary      Update a prompt template
// @Description  Update a prompt template by ID
// @Tags         templates
// @Accept       json
// @Produce      json
// @Param        id       path  int                    true  "Template ID"
// @Param        template  body  models.PromptTemplate  true  "Template data"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /prompts/{id} [put]
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

// handleDeleteTemplate godoc
// @Summary      Delete a prompt template
// @Description  Delete a prompt template by ID
// @Tags         templates
// @Accept       json
// @Produce      json
// @Param        id   path  int  true  "Template ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /prompts/{id} [delete]
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

package network

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/database/models"
)

// handleGetProviders godoc
// @Summary      List LLM providers
// @Description  Get paginated list of LLM providers
// @Tags         providers
// @Accept       json
// @Produce      json
// @Param        page    query  int  false  "Page number"  default(1)
// @Param        page_size  query  int  false  "Page size"  default(20)
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /llm/providers [get]
func (s *Server) handleGetProviders(c *gin.Context) {
	ctx := c.Request.Context()
	p := parsePagination(c)
	providers, total, err := s.repo.GetProvidersPaginated(ctx, p.Offset, p.PageSize)
	if err != nil {
		respondInternalError(c, err)
		return
	}
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
	c.JSON(http.StatusOK, gin.H{"data": providers, "total": total})
}

// handleCreateProvider godoc
// @Summary      Create an LLM provider
// @Description  Create a new LLM provider
// @Tags         providers
// @Accept       json
// @Produce      json
// @Param        provider  body  models.LLMProvider  true  "Provider data"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /llm/providers [post]
func (s *Server) handleCreateProvider(c *gin.Context) {
	ctx := c.Request.Context()
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
	if err := s.repo.CreateProvider(ctx, &provider); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": provider})
}

// handleUpdateProvider godoc
// @Summary      Update an LLM provider
// @Description  Update an LLM provider by ID
// @Tags         providers
// @Accept       json
// @Produce      json
// @Param        id       path  int                 true  "Provider ID"
// @Param        provider  body  models.LLMProvider  true  "Provider data"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /llm/providers/{id} [put]
func (s *Server) handleUpdateProvider(c *gin.Context) {
	ctx := c.Request.Context()
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
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
	provider.ID = id
	if err := s.repo.UpdateProvider(ctx, &provider); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": provider})
}

// handleDeleteProvider godoc
// @Summary      Delete an LLM provider
// @Description  Delete an LLM provider by ID
// @Tags         providers
// @Accept       json
// @Produce      json
// @Param        id   path  int  true  "Provider ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /llm/providers/{id} [delete]
func (s *Server) handleDeleteProvider(c *gin.Context) {
	ctx := c.Request.Context()
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := s.repo.DeleteProvider(ctx, id); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

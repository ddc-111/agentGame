package network

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/database/models"
)

func (s *Server) handleGetItems(c *gin.Context) {
	ctx := c.Request.Context()
	p := parsePagination(c)
	items, total, err := s.repo.GetItemsPaginated(ctx, p.Offset, p.PageSize)
	if err != nil {
		respondInternalError(c, err)
		return
	}
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
	c.JSON(http.StatusOK, gin.H{"data": items, "total": total})
}

func (s *Server) handleCreateItem(c *gin.Context) {
	ctx := c.Request.Context()
	var item models.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	errs := mergeErrors(
		validateRequired(map[string]interface{}{"name": item.Name, "code": item.Code}),
		validateStringMaxLen("name", item.Name, 100),
		validateStringMaxLen("code", item.Code, 50),
		validateStringMaxLen("description", item.Description, 500),
		validateStringMaxLen("icon", item.Icon, 255),
		validateStringMaxLen("effect", item.Effect, 500),
		validateJSON("effect", item.Effect, false),
	)
	if item.Category != "" {
		errs = append(errs, validateStringIn("category", item.Category, []string{"medicine", "food", "tool", "weapon", "armor", "material"})...)
	}
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}
	if err := s.repo.CreateItem(ctx, &item); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": item})
}

func (s *Server) handleUpdateItem(c *gin.Context) {
	ctx := c.Request.Context()
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var item models.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	errs := mergeErrors(
		validateRequired(map[string]interface{}{"name": item.Name, "code": item.Code}),
		validateStringMaxLen("name", item.Name, 100),
		validateStringMaxLen("code", item.Code, 50),
		validateStringMaxLen("description", item.Description, 500),
		validateStringMaxLen("icon", item.Icon, 255),
		validateStringMaxLen("effect", item.Effect, 500),
		validateJSON("effect", item.Effect, false),
	)
	if item.Category != "" {
		errs = append(errs, validateStringIn("category", item.Category, []string{"medicine", "food", "tool", "weapon", "armor", "material"})...)
	}
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}
	item.ID = id
	if err := s.repo.UpdateItem(ctx, &item); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (s *Server) handleDeleteItem(c *gin.Context) {
	ctx := c.Request.Context()
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := s.repo.DeleteItem(ctx, id); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

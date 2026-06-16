package network

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/database/models"
)

// handleGetShops godoc
// @Summary      List shops
// @Description  Get paginated list of shops
// @Tags         shops
// @Accept       json
// @Produce      json
// @Param        page    query  int  false  "Page number"  default(1)
// @Param        page_size  query  int  false  "Page size"  default(20)
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /shops [get]
func (s *Server) handleGetShops(c *gin.Context) {
	ctx := c.Request.Context()
	p := parsePagination(c)
	shops, total, err := s.repo.GetShopsPaginated(ctx, p.Offset, p.PageSize)
	if err != nil {
		respondInternalError(c, err)
		return
	}
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
	c.JSON(http.StatusOK, gin.H{"data": shops, "total": total})
}

// handleGetShop godoc
// @Summary      Get a shop
// @Description  Get a shop by ID
// @Tags         shops
// @Accept       json
// @Produce      json
// @Param        id   path  int  true  "Shop ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /shops/{id} [get]
func (s *Server) handleGetShop(c *gin.Context) {
	ctx := c.Request.Context()
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	shop, err := s.repo.GetShopByID(ctx, id)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Shop"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": shop})
}

// handleCreateShop godoc
// @Summary      Create a shop
// @Description  Create a new shop
// @Tags         shops
// @Accept       json
// @Produce      json
// @Param        shop  body  models.Shop  true  "Shop data"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /shops [post]
func (s *Server) handleCreateShop(c *gin.Context) {
	ctx := c.Request.Context()
	var shop models.Shop
	if err := c.ShouldBindJSON(&shop); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	errs := mergeErrors(
		validateRequired(map[string]interface{}{"name": shop.Name, "code": shop.Code}),
		validateStringMaxLen("name", shop.Name, 100),
		validateStringMaxLen("code", shop.Code, 50),
		validateStringMaxLen("description", shop.Description, 500),
		validateStringMaxLen("owner_npc", shop.OwnerNPC, 50),
		validateStringMaxLen("scene_code", shop.SceneCode, 50),
		validateStringMaxLen("discount", shop.Discount, 200),
		validateJSON("discount", shop.Discount, false),
		validateTimeFormat("open_time", shop.OpenTime),
		validateTimeFormat("close_time", shop.CloseTime),
	)
	if shop.Type != "" {
		errs = append(errs, validateStringIn("type", shop.Type, []string{"general", "blacksmith", "armory", "potion", "food", "specialty"})...)
	}
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}
	if err := s.repo.CreateShop(ctx, &shop); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": shop})
}

// handleUpdateShop godoc
// @Summary      Update a shop
// @Description  Update a shop by ID
// @Tags         shops
// @Accept       json
// @Produce      json
// @Param        id   path  int         true  "Shop ID"
// @Param        shop  body  models.Shop  true  "Shop data"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /shops/{id} [put]
func (s *Server) handleUpdateShop(c *gin.Context) {
	ctx := c.Request.Context()
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var shop models.Shop
	if err := c.ShouldBindJSON(&shop); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	errs := mergeErrors(
		validateRequired(map[string]interface{}{"name": shop.Name, "code": shop.Code}),
		validateStringMaxLen("name", shop.Name, 100),
		validateStringMaxLen("code", shop.Code, 50),
		validateStringMaxLen("description", shop.Description, 500),
		validateStringMaxLen("owner_npc", shop.OwnerNPC, 50),
		validateStringMaxLen("scene_code", shop.SceneCode, 50),
		validateStringMaxLen("discount", shop.Discount, 200),
		validateJSON("discount", shop.Discount, false),
		validateTimeFormat("open_time", shop.OpenTime),
		validateTimeFormat("close_time", shop.CloseTime),
	)
	if shop.Type != "" {
		errs = append(errs, validateStringIn("type", shop.Type, []string{"general", "blacksmith", "armory", "potion", "food", "specialty"})...)
	}
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}
	shop.ID = id
	if err := s.repo.UpdateShop(ctx, &shop); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": shop})
}

// handleDeleteShop godoc
// @Summary      Delete a shop
// @Description  Delete a shop by ID
// @Tags         shops
// @Accept       json
// @Produce      json
// @Param        id   path  int  true  "Shop ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /shops/{id} [delete]
func (s *Server) handleDeleteShop(c *gin.Context) {
	ctx := c.Request.Context()
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := s.repo.DeleteShop(ctx, id); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

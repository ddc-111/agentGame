package network

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/database/models"
)

func (s *Server) handleGetShops(c *gin.Context) {
	p := parsePagination(c)
	shops, total, err := s.repo.GetShopsPaginated(p.Offset, p.PageSize)
	if err != nil {
		respondInternalError(c, err)
		return
	}
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
	c.JSON(http.StatusOK, gin.H{"data": shops, "total": total})
}

func (s *Server) handleGetShop(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	shop, err := s.repo.GetShopByID(uint(id))
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Shop"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": shop})
}

func (s *Server) handleCreateShop(c *gin.Context) {
	var shop models.Shop
	if err := c.ShouldBindJSON(&shop); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	errs := validateRequired(map[string]interface{}{"name": shop.Name, "code": shop.Code})
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}
	if err := s.repo.CreateShop(&shop); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": shop})
}

func (s *Server) handleUpdateShop(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var shop models.Shop
	if err := c.ShouldBindJSON(&shop); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	errs := validateRequired(map[string]interface{}{"name": shop.Name, "code": shop.Code})
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}
	shop.ID = uint(id)
	if err := s.repo.UpdateShop(&shop); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": shop})
}

func (s *Server) handleDeleteShop(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := s.repo.DeleteShop(uint(id)); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

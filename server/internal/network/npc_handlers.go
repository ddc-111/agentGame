package network

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/database/models"
)

func (s *Server) handleGetNPCs(c *gin.Context) {
	p := parsePagination(c)
	npcs, total, err := s.repo.GetNPCsPaginated(p.Offset, p.PageSize)
	if err != nil {
		respondInternalError(c, err)
		return
	}
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
	c.JSON(http.StatusOK, gin.H{"data": npcs, "total": total})
}

func (s *Server) handleGetNPC(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	npc, err := s.repo.GetNPCByID(uint(id))
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("NPC"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": npc})
}

func (s *Server) handleCreateNPC(c *gin.Context) {
	var npc models.NPC
	if err := c.ShouldBindJSON(&npc); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	errs := validateRequired(map[string]interface{}{"name": npc.Name, "code": npc.Code})
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}
	if err := s.repo.CreateNPC(&npc); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": npc})
}

func (s *Server) handleUpdateNPC(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var npc models.NPC
	if err := c.ShouldBindJSON(&npc); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	errs := validateRequired(map[string]interface{}{"name": npc.Name, "code": npc.Code})
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}
	npc.ID = uint(id)
	if err := s.repo.UpdateNPC(&npc); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": npc})
}

func (s *Server) handleDeleteNPC(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := s.repo.DeleteNPC(uint(id)); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

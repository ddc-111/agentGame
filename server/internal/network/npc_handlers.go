package network

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/database/models"
)

// handleGetNPCs godoc
// @Summary      List NPCs
// @Description  Get paginated list of NPCs
// @Tags         npcs
// @Accept       json
// @Produce      json
// @Param        page    query  int  false  "Page number"  default(1)
// @Param        page_size  query  int  false  "Page size"  default(20)
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /npcs [get]
func (s *Server) handleGetNPCs(c *gin.Context) {
	ctx := c.Request.Context()
	p := parsePagination(c)
	npcs, total, err := s.repo.GetNPCsPaginated(ctx, p.Offset, p.PageSize)
	if err != nil {
		respondInternalError(c, err)
		return
	}
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
	c.JSON(http.StatusOK, gin.H{"data": npcs, "total": total})
}

// handleGetNPC godoc
// @Summary      Get an NPC
// @Description  Get an NPC by ID
// @Tags         npcs
// @Accept       json
// @Produce      json
// @Param        id   path  int  true  "NPC ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /npcs/{id} [get]
func (s *Server) handleGetNPC(c *gin.Context) {
	ctx := c.Request.Context()
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	npc, err := s.repo.GetNPCByID(ctx, id)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("NPC"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": npc})
}

// handleCreateNPC godoc
// @Summary      Create an NPC
// @Description  Create a new NPC
// @Tags         npcs
// @Accept       json
// @Produce      json
// @Param        npc  body  models.NPC  true  "NPC data"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /npcs [post]
func (s *Server) handleCreateNPC(c *gin.Context) {
	ctx := c.Request.Context()
	var npc models.NPC
	if err := c.ShouldBindJSON(&npc); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	errs := mergeErrors(
		validateRequired(map[string]interface{}{"name": npc.Name, "code": npc.Code}),
		validateStringMaxLen("name", npc.Name, 100),
		validateStringMaxLen("code", npc.Code, 50),
		validateStringMaxLen("description", npc.Description, 500),
		validateStringMaxLen("avatar", npc.Avatar, 255),
		validateStringMaxLen("sprite", npc.Sprite, 255),
		validateStringMaxLen("behaviors", npc.Behaviors, 500),
		validateStringMaxLen("schedule", npc.Schedule, 2000),
		validateJSON("behaviors", npc.Behaviors, true),
		validateJSON("schedule", npc.Schedule, true),
	)
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}
	if err := s.repo.CreateNPC(ctx, &npc); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": npc})
}

// handleUpdateNPC godoc
// @Summary      Update an NPC
// @Description  Update an NPC by ID
// @Tags         npcs
// @Accept       json
// @Produce      json
// @Param        id   path  int       true  "NPC ID"
// @Param        npc  body  models.NPC  true  "NPC data"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /npcs/{id} [put]
func (s *Server) handleUpdateNPC(c *gin.Context) {
	ctx := c.Request.Context()
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var npc models.NPC
	if err := c.ShouldBindJSON(&npc); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	errs := mergeErrors(
		validateRequired(map[string]interface{}{"name": npc.Name, "code": npc.Code}),
		validateStringMaxLen("name", npc.Name, 100),
		validateStringMaxLen("code", npc.Code, 50),
		validateStringMaxLen("description", npc.Description, 500),
		validateStringMaxLen("avatar", npc.Avatar, 255),
		validateStringMaxLen("sprite", npc.Sprite, 255),
		validateStringMaxLen("behaviors", npc.Behaviors, 500),
		validateStringMaxLen("schedule", npc.Schedule, 2000),
		validateJSON("behaviors", npc.Behaviors, true),
		validateJSON("schedule", npc.Schedule, true),
	)
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}
	npc.ID = id
	if err := s.repo.UpdateNPC(ctx, &npc); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": npc})
}

// handleDeleteNPC godoc
// @Summary      Delete an NPC
// @Description  Delete an NPC by ID
// @Tags         npcs
// @Accept       json
// @Produce      json
// @Param        id   path  int  true  "NPC ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /npcs/{id} [delete]
func (s *Server) handleDeleteNPC(c *gin.Context) {
	ctx := c.Request.Context()
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := s.repo.DeleteNPC(ctx, id); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

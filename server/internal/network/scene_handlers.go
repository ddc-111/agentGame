package network

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/database/models"
)

// handleGetScenes godoc
// @Summary      List scenes
// @Description  Get paginated list of scenes
// @Tags         scenes
// @Accept       json
// @Produce      json
// @Param        page    query  int  false  "Page number"  default(1)
// @Param        page_size  query  int  false  "Page size"  default(20)
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /scenes [get]
func (s *Server) handleGetScenes(c *gin.Context) {
	ctx := c.Request.Context()
	p := parsePagination(c)
	scenes, total, err := s.repo.GetScenesPaginated(ctx, p.Offset, p.PageSize)
	if err != nil {
		respondInternalError(c, err)
		return
	}
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
	c.JSON(http.StatusOK, gin.H{"data": scenes, "total": total})
}

// handleGetScene godoc
// @Summary      Get a scene
// @Description  Get a scene by ID
// @Tags         scenes
// @Accept       json
// @Produce      json
// @Param        id   path  int  true  "Scene ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /scenes/{id} [get]
func (s *Server) handleGetScene(c *gin.Context) {
	ctx := c.Request.Context()
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	scene, err := s.repo.GetSceneByID(ctx, id)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Scene"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": scene})
}

// handleCreateScene godoc
// @Summary      Create a scene
// @Description  Create a new scene
// @Tags         scenes
// @Accept       json
// @Produce      json
// @Param        scene  body  models.Scene  true  "Scene data"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /scenes [post]
func (s *Server) handleCreateScene(c *gin.Context) {
	ctx := c.Request.Context()
	var scene models.Scene
	if err := c.ShouldBindJSON(&scene); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	errs := mergeErrors(
		validateRequired(map[string]interface{}{"name": scene.Name, "code": scene.Code}),
		validateIntRange("width", scene.Width, 100, 10000),
		validateIntRange("height", scene.Height, 100, 10000),
	)
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}
	if err := s.repo.CreateScene(ctx, &scene); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": scene})
}

// handleUpdateScene godoc
// @Summary      Update a scene
// @Description  Update a scene by ID
// @Tags         scenes
// @Accept       json
// @Produce      json
// @Param        id     path  int          true  "Scene ID"
// @Param        scene  body  models.Scene  true  "Scene data"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /scenes/{id} [put]
func (s *Server) handleUpdateScene(c *gin.Context) {
	ctx := c.Request.Context()
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var scene models.Scene
	if err := c.ShouldBindJSON(&scene); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	errs := mergeErrors(
		validateRequired(map[string]interface{}{"name": scene.Name, "code": scene.Code}),
		validateIntRange("width", scene.Width, 100, 10000),
		validateIntRange("height", scene.Height, 100, 10000),
	)
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}
	scene.ID = id
	if err := s.repo.UpdateScene(ctx, &scene); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": scene})
}

// handleDeleteScene godoc
// @Summary      Delete a scene
// @Description  Delete a scene by ID
// @Tags         scenes
// @Accept       json
// @Produce      json
// @Param        id   path  int  true  "Scene ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /scenes/{id} [delete]
func (s *Server) handleDeleteScene(c *gin.Context) {
	ctx := c.Request.Context()
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := s.repo.DeleteScene(ctx, id); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

package network

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/database/models"
)

func (s *Server) handleGetScenes(c *gin.Context) {
	scenes, err := s.repo.GetScenes()
	if err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": scenes})
}

func (s *Server) handleGetScene(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	scene, err := s.repo.GetSceneByID(uint(id))
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Scene"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": scene})
}

func (s *Server) handleCreateScene(c *gin.Context) {
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
	if err := s.repo.CreateScene(&scene); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": scene})
}

func (s *Server) handleUpdateScene(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
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
	scene.ID = uint(id)
	if err := s.repo.UpdateScene(&scene); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": scene})
}

func (s *Server) handleDeleteScene(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := s.repo.DeleteScene(uint(id)); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

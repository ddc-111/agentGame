package network

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/database/models"
)

func (s *Server) handleGetTasks(c *gin.Context) {
	ctx := c.Request.Context()
	p := parsePagination(c)
	tasks, total, err := s.repo.GetTasksPaginated(ctx, p.Offset, p.PageSize)
	if err != nil {
		respondInternalError(c, err)
		return
	}
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
	c.JSON(http.StatusOK, gin.H{"data": tasks, "total": total})
}

func (s *Server) handleGetTask(c *gin.Context) {
	ctx := c.Request.Context()
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	task, err := s.repo.GetTaskByID(ctx, id)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Task"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": task})
}

func (s *Server) handleCreateTask(c *gin.Context) {
	ctx := c.Request.Context()
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	errs := mergeErrors(
		validateRequired(map[string]interface{}{"name": task.Name, "code": task.Code}),
		validateStringMaxLen("name", task.Name, 100),
		validateStringMaxLen("code", task.Code, 50),
		validateStringMaxLen("description", task.Description, 500),
		validateStringMaxLen("next_task", task.NextTask, 50),
		validateStringMaxLen("dialogue", task.Dialogue, 50),
		validateJSON("trigger", task.Trigger, false),
		validateJSON("objectives", task.Objectives, true),
		validateJSON("rewards", task.Rewards, false),
	)
	if task.Type != "" {
		errs = append(errs, validateStringIn("type", task.Type, []string{"main", "side", "daily", "event"})...)
	}
	if task.Status != "" {
		errs = append(errs, validateStringIn("status", task.Status, []string{"active", "locked", "completed", "failed"})...)
	}
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}
	if err := s.repo.CreateTask(ctx, &task); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": task})
}

func (s *Server) handleUpdateTask(c *gin.Context) {
	ctx := c.Request.Context()
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	errs := mergeErrors(
		validateRequired(map[string]interface{}{"name": task.Name, "code": task.Code}),
		validateStringMaxLen("name", task.Name, 100),
		validateStringMaxLen("code", task.Code, 50),
		validateStringMaxLen("description", task.Description, 500),
		validateStringMaxLen("next_task", task.NextTask, 50),
		validateStringMaxLen("dialogue", task.Dialogue, 50),
		validateJSON("trigger", task.Trigger, false),
		validateJSON("objectives", task.Objectives, true),
		validateJSON("rewards", task.Rewards, false),
	)
	if task.Type != "" {
		errs = append(errs, validateStringIn("type", task.Type, []string{"main", "side", "daily", "event"})...)
	}
	if task.Status != "" {
		errs = append(errs, validateStringIn("status", task.Status, []string{"active", "locked", "completed", "failed"})...)
	}
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}
	task.ID = id
	if err := s.repo.UpdateTask(ctx, &task); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": task})
}

func (s *Server) handleDeleteTask(c *gin.Context) {
	ctx := c.Request.Context()
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := s.repo.DeleteTask(ctx, id); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

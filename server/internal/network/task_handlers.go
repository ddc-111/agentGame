package network

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/database/models"
)

// handleGetTasks godoc
// @Summary      List tasks
// @Description  Get paginated list of tasks
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        page    query  int  false  "Page number"  default(1)
// @Param        page_size  query  int  false  "Page size"  default(20)
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /tasks [get]
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

// handleGetTask godoc
// @Summary      Get a task
// @Description  Get a task by ID
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id   path  int  true  "Task ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /tasks/{id} [get]
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

// handleCreateTask godoc
// @Summary      Create a task
// @Description  Create a new task
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        task  body  models.Task  true  "Task data"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /tasks [post]
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

// handleUpdateTask godoc
// @Summary      Update a task
// @Description  Update a task by ID
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id   path  int         true  "Task ID"
// @Param        task  body  models.Task  true  "Task data"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /tasks/{id} [put]
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

// handleDeleteTask godoc
// @Summary      Delete a task
// @Description  Delete a task by ID
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id   path  int  true  "Task ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /tasks/{id} [delete]
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

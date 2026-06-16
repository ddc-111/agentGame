package network

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/database/models"
)

// handleGetFlow godoc
// @Summary      Get a flow
// @Description  Get a flow by ID
// @Tags         flows
// @Accept       json
// @Produce      json
// @Param        id   path  int  true  "Flow ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /flows/{id} [get]
func (s *Server) handleGetFlow(c *gin.Context) {
	ctx := c.Request.Context()
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	flow, err := s.repo.GetFlowByID(ctx, id)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Flow"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": flow})
}

// handleGetFlows godoc
// @Summary      List flows
// @Description  Get paginated list of flows
// @Tags         flows
// @Accept       json
// @Produce      json
// @Param        page    query  int  false  "Page number"  default(1)
// @Param        page_size  query  int  false  "Page size"  default(20)
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /flows [get]
func (s *Server) handleGetFlows(c *gin.Context) {
	ctx := c.Request.Context()
	p := parsePagination(c)
	flows, total, err := s.repo.GetFlowsPaginated(ctx, p.Offset, p.PageSize)
	if err != nil {
		respondInternalError(c, err)
		return
	}
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
	c.JSON(http.StatusOK, gin.H{"data": flows, "total": total})
}

// handleCreateFlow godoc
// @Summary      Create a flow
// @Description  Create a new flow
// @Tags         flows
// @Accept       json
// @Produce      json
// @Param        flow  body  models.Flow  true  "Flow data"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /flows [post]
func (s *Server) handleCreateFlow(c *gin.Context) {
	ctx := c.Request.Context()
	var flow models.Flow
	if err := c.ShouldBindJSON(&flow); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	errs := validateRequired(map[string]interface{}{"name": flow.Name, "code": flow.Code})
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}
	if err := s.repo.CreateFlow(ctx, &flow); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": flow})
}

// handleUpdateFlow godoc
// @Summary      Update a flow
// @Description  Update a flow by ID
// @Tags         flows
// @Accept       json
// @Produce      json
// @Param        id   path  int         true  "Flow ID"
// @Param        flow  body  models.Flow  true  "Flow data"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /flows/{id} [put]
func (s *Server) handleUpdateFlow(c *gin.Context) {
	ctx := c.Request.Context()
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var flow models.Flow
	if err := c.ShouldBindJSON(&flow); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	errs := validateRequired(map[string]interface{}{"name": flow.Name, "code": flow.Code})
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}
	flow.ID = id
	if err := s.repo.UpdateFlow(ctx, &flow); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": flow})
}

// handleDeleteFlow godoc
// @Summary      Delete a flow
// @Description  Delete a flow by ID
// @Tags         flows
// @Accept       json
// @Produce      json
// @Param        id   path  int  true  "Flow ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /flows/{id} [delete]
func (s *Server) handleDeleteFlow(c *gin.Context) {
	ctx := c.Request.Context()
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := s.repo.DeleteFlow(ctx, id); err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

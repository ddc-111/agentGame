package network

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/mcp"
)

func (s *Server) handleMCPTools(c *gin.Context) {
	tools := s.mcp.GetTools()
	c.JSON(http.StatusOK, gin.H{"tools": tools})
}

func (s *Server) handleMCPCall(c *gin.Context) {
	var req struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}
	errs := validateRequired(map[string]interface{}{"name": req.Name})
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}

	result := s.mcp.HandleRequest(c.Request.Context(), mcp.MCPRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "tools/call",
		Params: map[string]interface{}{
			"name":      req.Name,
			"arguments": req.Arguments,
		},
	})

	c.JSON(http.StatusOK, result)
}

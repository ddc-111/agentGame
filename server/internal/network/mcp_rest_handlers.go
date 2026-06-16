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

func (s *Server) handleMCPResources(c *gin.Context) {
	resources := s.mcp.GetResources()
	c.JSON(http.StatusOK, gin.H{"resources": resources})
}

func (s *Server) handleMCPResourceRead(c *gin.Context) {
	uri := c.Query("uri")
	if uri == "" {
		respondError(c, http.StatusBadRequest, BadRequest("uri parameter is required"))
		return
	}

	result := s.mcp.HandleRequest(c.Request.Context(), mcp.MCPRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "resources/read",
		Params: map[string]interface{}{
			"uri": uri,
		},
	})

	c.JSON(http.StatusOK, result)
}

func (s *Server) handleMCPPrompts(c *gin.Context) {
	prompts := s.mcp.GetPrompts()
	c.JSON(http.StatusOK, gin.H{"prompts": prompts})
}

func (s *Server) handleMCPPromptGet(c *gin.Context) {
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
		Method:  "prompts/get",
		Params: map[string]interface{}{
			"name":      req.Name,
			"arguments": req.Arguments,
		},
	})

	c.JSON(http.StatusOK, result)
}

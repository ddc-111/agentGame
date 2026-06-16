package network

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/mcp"
)

// handleMCPTools godoc
// @Summary      List MCP tools
// @Description  Get available MCP tools
// @Tags         mcp
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /mcp/tools [get]
func (s *Server) handleMCPTools(c *gin.Context) {
	tools := s.mcp.GetTools()
	c.JSON(http.StatusOK, gin.H{"tools": tools})
}

// handleMCPCall godoc
// @Summary      Call an MCP tool
// @Description  Call an MCP tool by name with arguments
// @Tags         mcp
// @Accept       json
// @Produce      json
// @Param        request  body  object  true  "MCP tool call request"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /mcp/call [post]
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

// handleMCPResources godoc
// @Summary      List MCP resources
// @Description  Get available MCP resources
// @Tags         mcp
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /mcp/resources [get]
func (s *Server) handleMCPResources(c *gin.Context) {
	resources := s.mcp.GetResources()
	c.JSON(http.StatusOK, gin.H{"resources": resources})
}

// handleMCPResourceRead godoc
// @Summary      Read an MCP resource
// @Description  Read an MCP resource by URI
// @Tags         mcp
// @Accept       json
// @Produce      json
// @Param        uri  query  string  true  "Resource URI"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /mcp/resources/read [get]
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

// handleMCPPrompts godoc
// @Summary      List MCP prompts
// @Description  Get available MCP prompts
// @Tags         mcp
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /mcp/prompts [get]
func (s *Server) handleMCPPrompts(c *gin.Context) {
	prompts := s.mcp.GetPrompts()
	c.JSON(http.StatusOK, gin.H{"prompts": prompts})
}

// handleMCPPromptGet godoc
// @Summary      Get an MCP prompt
// @Description  Get an MCP prompt by name with arguments
// @Tags         mcp
// @Accept       json
// @Produce      json
// @Param        request  body  object  true  "MCP prompt get request"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /mcp/prompts/get [post]
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

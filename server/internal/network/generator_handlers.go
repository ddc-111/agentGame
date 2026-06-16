package network

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/generator"
)

// handleGenerate godoc
// @Summary      Generate content
// @Description  Generate game content using AI
// @Tags         generator
// @Accept       json
// @Produce      json
// @Param        request  body  generator.GenerateRequest  true  "Generate request"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      503  {object}  map[string]interface{}
// @Router       /generator/generate [post]
func (s *Server) handleGenerate(c *gin.Context) {
	var req generator.GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}

	errs := mergeErrors(
		validateRequired(map[string]interface{}{"type": req.Type, "action": req.Action}),
		validateStringIn("type", req.Type, []string{"npc", "scene", "task", "shop", "item", "agent", "dialogue", "flow"}),
		validateStringIn("action", req.Action, []string{"create", "complete", "expand", "translate"}),
	)
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}

	if s.generator == nil || !s.generator.IsEnabled() {
		respondError(c, http.StatusServiceUnavailable, BadRequest("生成智能体未启用，请检查配置"))
		return
	}

	resp, err := s.generator.Generate(c.Request.Context(), req)
	if err != nil {
		respondInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// handleGeneratorStatus godoc
// @Summary      Get generator status
// @Description  Get the current status and configuration of the generator
// @Tags         generator
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /generator/status [get]
func (s *Server) handleGeneratorStatus(c *gin.Context) {
	enabled := s.generator != nil && s.generator.IsEnabled()
	cfg := s.generator.GetConfig()

	c.JSON(http.StatusOK, gin.H{
		"enabled":  enabled,
		"provider": cfg.Provider,
		"model":    cfg.Model,
		"base_url": cfg.BaseURL,
	})
}

// handleGeneratorTest godoc
// @Summary      Test generator
// @Description  Test the generator with a sample NPC generation
// @Tags         generator
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      503  {object}  map[string]interface{}
// @Router       /generator/test [post]
func (s *Server) handleGeneratorTest(c *gin.Context) {
	if s.generator == nil || !s.generator.IsEnabled() {
		respondError(c, http.StatusServiceUnavailable, BadRequest("生成智能体未启用"))
		return
	}

	req := generator.GenerateRequest{
		Type:   "npc",
		Action: "create",
		Params: map[string]interface{}{
			"description": "一个卖包子的老大爷",
			"theme":       "古风小镇",
		},
	}

	resp, err := s.generator.Generate(c.Request.Context(), req)
	if err != nil {
		respondInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

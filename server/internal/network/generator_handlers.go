package network

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/generator"
)

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

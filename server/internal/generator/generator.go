package generator

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
	"github.com/ddc-111/agentGame/server/internal/config"
)

// Generator 生成智能体
type Generator struct {
	client *openai.Client
	cfg    config.GeneratorConfig
}

// GenerateRequest 生成请求
type GenerateRequest struct {
	Type    string                 `json:"type"`    // npc, scene, task, shop, item, agent, dialogue, flow
	Action  string                 `json:"action"`  // create, complete, expand, translate
	Context map[string]interface{} `json:"context"` // 上下文信息
	Params  map[string]interface{} `json:"params"`  // 参数
}

// GenerateResponse 生成响应
type GenerateResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Error   string      `json:"error,omitempty"`
}

// New 创建生成智能体
func New(cfg config.GeneratorConfig) (*Generator, error) {
	if !cfg.Enabled {
		return &Generator{cfg: cfg}, nil
	}

	clientConfig := openai.DefaultConfig(cfg.APIKey)
	if cfg.BaseURL != "" {
		clientConfig.BaseURL = cfg.BaseURL
	}

	client := openai.NewClientWithConfig(clientConfig)

	return &Generator{
		client: client,
		cfg:    cfg,
	}, nil
}

// Generate 生成内容
func (g *Generator) Generate(ctx context.Context, req GenerateRequest) (*GenerateResponse, error) {
	if !g.cfg.Enabled {
		return &GenerateResponse{
			Success: false,
			Error:   "生成智能体未启用",
		}, nil
	}

	prompt := g.buildPrompt(req)

	timeout := time.Duration(g.cfg.Timeout) * time.Second
	if timeout == 0 {
		timeout = 60 * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	resp, err := g.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: g.cfg.Model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: g.getSystemPrompt(req.Type),
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: float32(g.cfg.Temperature),
		MaxTokens:   g.cfg.MaxTokens,
	})
	if err != nil {
		return &GenerateResponse{
			Success: false,
			Error:   fmt.Sprintf("生成失败: %v", err),
		}, nil
	}

	if len(resp.Choices) == 0 {
		return &GenerateResponse{
			Success: false,
			Error:   "未生成任何内容",
		}, nil
	}

	content := resp.Choices[0].Message.Content
	result, err := g.parseResponse(req.Type, content)
	if err != nil {
		return &GenerateResponse{
			Success: false,
			Error:   fmt.Sprintf("解析响应失败: %v", err),
		}, nil
	}

	return &GenerateResponse{
		Success: true,
		Data:    result,
		Message: "生成成功",
	}, nil
}

// getSystemPrompt 获取系统提示词
func (g *Generator) getSystemPrompt(genType string) string {
	base := `你是一个古风RPG游戏配置生成助手。你需要根据用户的需求生成游戏配置数据。
请以JSON格式返回数据，确保JSON格式正确，可以直接被解析。
不要返回任何额外的解释文字，只返回JSON数据。`

	switch genType {
	case "npc":
		return base + `
你需要生成NPC配置，包含以下字段：
- name: NPC名称（古风名字）
- title: NPC称号
- description: NPC描述（50-100字）
- behaviors: 行为模式数组
- schedule: 日程安排数组

示例格式：
{
  "name": "李掌柜",
  "title": "杂货铺老板",
  "description": "一位精明的中年商人...",
  "behaviors": ["idle", "greet", "sell"],
  "schedule": [{"time": "06:00", "action": "open_shop"}]
}`

	case "scene":
		return base + `
你需要生成场景配置，包含以下字段：
- name: 场景名称
- description: 场景描述
- width: 宽度（1280-3840）
- height: 高度（720-2160）
- atmosphere: 氛围描述

示例格式：
{
  "name": "月老祠",
  "description": "香火缭绕的古老祠堂...",
  "width": 1920,
  "height": 1080,
  "atmosphere": "神秘浪漫"
}`

	case "task":
		return base + `
你需要生成任务配置，包含以下字段：
- name: 任务名称
- type: 任务类型(main/side/daily)
- description: 任务描述
- objectives: 任务目标数组
- rewards: 奖励对象

示例格式：
{
  "name": "初来乍到",
  "type": "main",
  "description": "新来的冒险者...",
  "objectives": [{"type": "dialogue", "target": "npc_001", "description": "..."}],
  "rewards": {"exp": 100, "gold": 500}
}`

	case "shop":
		return base + `
你需要生成商店配置，包含以下字段：
- name: 商店名称
- type: 商店类型(general/blacksmith/pharmacy/restaurant)
- description: 商店描述
- openTime: 开门时间
- closeTime: 关门时间
- suggestedItems: 建议商品数组

示例格式：
{
  "name": "回春堂",
  "type": "pharmacy",
  "description": "老字号药店...",
  "openTime": "08:00",
  "closeTime": "20:00",
  "suggestedItems": [{"name": "金疮药", "price": 200, "description": "..."}]
}`

	case "item":
		return base + `
你需要生成道具配置，包含以下字段：
- name: 道具名称
- category: 分类(medicine/food/weapon/armor/tool/material)
- description: 道具描述
- effect: 效果对象

示例格式：
{
  "name": "九转还魂丹",
  "category": "medicine",
  "description": "传说中的仙丹...",
  "effect": {"hp": 500, "mp": 200}
}`

	case "agent":
		return base + `
你需要生成智能体配置，包含以下字段：
- name: 智能体名称
- description: 描述
- systemPrompt: 系统提示词（详细的人设和行为规则）
- suggestedLLM: 建议使用的模型

示例格式：
{
  "name": "酒馆老板智能体",
  "description": "负责酒馆经营的AI...",
  "systemPrompt": "你是醉仙楼的老板...",
  "suggestedLLM": "gpt-4"
}`

	case "dialogue":
		return base + `
你需要生成对话树配置，包含以下字段：
- startNode: 起始节点ID
- nodes: 节点数组，每个节点包含id, type, content, choices

示例格式：
{
  "startNode": "node_1",
  "nodes": [
    {"id": "node_1", "type": "npc_say", "content": "客官好！...", "nextNode": "node_2"},
    {"id": "node_2", "type": "player_choice", "choices": [{"text": "...", "nextNode": "..."}]}
  ]
}`

	case "flow":
		return base + `
你需要生成流程配置，包含以下字段：
- name: 流程名称
- description: 描述
- nodes: 流程节点数组
- edges: 连接数组

示例格式：
{
  "name": "NPC巡逻流程",
  "description": "...",
  "nodes": [{"id": "node_1", "type": "start", "data": {"label": "开始"}}],
  "edges": [{"source": "node_1", "target": "node_2"}]
}`

	default:
		return base
	}
}

// buildPrompt 构建提示词
func (g *Generator) buildPrompt(req GenerateRequest) string {
	var parts []string

	switch req.Action {
	case "create":
		parts = append(parts, fmt.Sprintf("请生成一个新的%s配置。", req.Type))
	case "complete":
		parts = append(parts, fmt.Sprintf("请补全以下%s配置。", req.Type))
	case "expand":
		parts = append(parts, fmt.Sprintf("请扩展以下%s配置，添加更多细节。", req.Type))
	case "translate":
		parts = append(parts, fmt.Sprintf("请将以下%s配置翻译成古风风格。", req.Type))
	}

	if desc, ok := req.Params["description"].(string); ok && desc != "" {
		parts = append(parts, fmt.Sprintf("描述：%s", desc))
	}

	if theme, ok := req.Params["theme"].(string); ok && theme != "" {
		parts = append(parts, fmt.Sprintf("主题：%s", theme))
	}

	if dynasty, ok := req.Params["dynasty"].(string); ok && dynasty != "" {
		parts = append(parts, fmt.Sprintf("朝代背景：%s", dynasty))
	}

	if style, ok := req.Params["style"].(string); ok && style != "" {
		parts = append(parts, fmt.Sprintf("风格：%s", style))
	}

	if count, ok := req.Params["count"].(float64); ok && count > 1 {
		parts = append(parts, fmt.Sprintf("请生成%d个不同的配置。", int(count)))
	}

	// 添加上下文
	if len(req.Context) > 0 {
		contextJSON, _ := json.Marshal(req.Context)
		parts = append(parts, fmt.Sprintf("上下文信息：%s", string(contextJSON)))
	}

	// 如果是补全，添加现有数据
	if req.Action == "complete" || req.Action == "expand" {
		if existing, ok := req.Params["existing"]; ok {
			existingJSON, _ := json.Marshal(existing)
			parts = append(parts, fmt.Sprintf("现有数据：%s", string(existingJSON)))
		}
	}

	return strings.Join(parts, "\n")
}

// parseResponse 解析响应
func (g *Generator) parseResponse(genType, content string) (interface{}, error) {
	// 尝试提取JSON
	content = extractJSON(content)

	var result interface{}
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("JSON解析失败: %v, 原始内容: %s", err, content)
	}

	return result, nil
}

// extractJSON 提取JSON内容
func extractJSON(content string) string {
	content = strings.TrimSpace(content)

	// 如果内容被```json包围
	if strings.HasPrefix(content, "```json") {
		lines := strings.Split(content, "\n")
		var jsonLines []string
		inBlock := false
		for _, line := range lines {
			if strings.TrimSpace(line) == "```json" {
				inBlock = true
				continue
			}
			if strings.TrimSpace(line) == "```" {
				inBlock = false
				continue
			}
			if inBlock {
				jsonLines = append(jsonLines, line)
			}
		}
		return strings.Join(jsonLines, "\n")
	}

	// 如果内容被```包围
	if strings.HasPrefix(content, "```") {
		lines := strings.Split(content, "\n")
		var jsonLines []string
		inBlock := false
		for _, line := range lines {
			if strings.HasPrefix(strings.TrimSpace(line), "```") {
				if inBlock {
					inBlock = false
				} else {
					inBlock = true
				}
				continue
			}
			if inBlock {
				jsonLines = append(jsonLines, line)
			}
		}
		return strings.Join(jsonLines, "\n")
	}

	// 尝试找到JSON开始和结束
	start := strings.Index(content, "{")
	if start == -1 {
		start = strings.Index(content, "[")
	}
	if start == -1 {
		return content
	}

	end := strings.LastIndex(content, "}")
	if end == -1 {
		end = strings.LastIndex(content, "]")
	}
	if end == -1 {
		return content[start:]
	}

	return content[start : end+1]
}

// IsEnabled 是否启用
func (g *Generator) IsEnabled() bool {
	return g.cfg.Enabled
}

// GetConfig 获取配置
func (g *Generator) GetConfig() config.GeneratorConfig {
	return g.cfg
}

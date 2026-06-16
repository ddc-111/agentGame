package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ddc-111/agentGame/server/internal/database/models"
	"github.com/ddc-111/agentGame/server/internal/database/repository"
	"github.com/ddc-111/agentGame/server/internal/generator"
)

// Server MCP服务器
type Server struct {
	repo      *repository.Repository
	generator *generator.Generator
	tools     []Tool
}

// Tool MCP工具定义
type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

// ToolCall 工具调用
type ToolCall struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

// ToolResult 工具结果
type ToolResult struct {
	Content []Content `json:"content"`
	IsError bool      `json:"isError,omitempty"`
}

// Content 内容
type Content struct {
	Type string      `json:"type"`
	Text string      `json:"text,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

// MCPRequest MCP请求
type MCPRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// MCPResponse MCP响应
type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

// MCPError MCP错误
type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// New 创建MCP服务器
func New(repo *repository.Repository, gen *generator.Generator) *Server {
	s := &Server{
		repo:      repo,
		generator: gen,
	}
	s.tools = s.initTools()
	return s
}

// initTools 初始化工具集
func (s *Server) initTools() []Tool {
	return []Tool{
		// 场景工具
		{
			Name:        "list_scenes",
			Description: "获取所有场景列表",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
		},
		{
			Name:        "get_scene",
			Description: "根据ID获取场景详情",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"id": map[string]interface{}{"type": "number", "description": "场景ID"}}, "required": []string{"id"}},
		},
		{
			Name:        "create_scene",
			Description: "创建新场景",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"name": map[string]interface{}{"type": "string", "description": "场景名称"}, "code": map[string]interface{}{"type": "string", "description": "场景代码"}, "description": map[string]interface{}{"type": "string", "description": "场景描述"}, "width": map[string]interface{}{"type": "number", "description": "宽度"}, "height": map[string]interface{}{"type": "number", "description": "高度"}}, "required": []string{"name", "code"}},
		},
		{
			Name:        "update_scene",
			Description: "更新场景信息",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"id": map[string]interface{}{"type": "number", "description": "场景ID"}, "name": map[string]interface{}{"type": "string"}, "description": map[string]interface{}{"type": "string"}, "width": map[string]interface{}{"type": "number"}, "height": map[string]interface{}{"type": "number"}}, "required": []string{"id"}},
		},
		{
			Name:        "delete_scene",
			Description: "删除场景",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"id": map[string]interface{}{"type": "number", "description": "场景ID"}}, "required": []string{"id"}},
		},
		// NPC工具
		{
			Name:        "list_npcs",
			Description: "获取所有NPC列表",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
		},
		{
			Name:        "get_npc",
			Description: "根据ID获取NPC详情",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"id": map[string]interface{}{"type": "number", "description": "NPC ID"}}, "required": []string{"id"}},
		},
		{
			Name:        "create_npc",
			Description: "创建新NPC",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"name": map[string]interface{}{"type": "string", "description": "NPC名称"}, "code": map[string]interface{}{"type": "string", "description": "NPC代码"}, "title": map[string]interface{}{"type": "string", "description": "称号"}, "description": map[string]interface{}{"type": "string", "description": "描述"}}, "required": []string{"name", "code"}},
		},
		{
			Name:        "update_npc",
			Description: "更新NPC信息",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"id": map[string]interface{}{"type": "number", "description": "NPC ID"}, "name": map[string]interface{}{"type": "string"}, "title": map[string]interface{}{"type": "string"}, "description": map[string]interface{}{"type": "string"}, "agent_id": map[string]interface{}{"type": "number"}, "shop_id": map[string]interface{}{"type": "number"}}, "required": []string{"id"}},
		},
		{
			Name:        "delete_npc",
			Description: "删除NPC",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"id": map[string]interface{}{"type": "number", "description": "NPC ID"}}, "required": []string{"id"}},
		},
		// 智能体工具
		{
			Name:        "list_agents",
			Description: "获取所有智能体列表",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
		},
		{
			Name:        "get_agent",
			Description: "根据ID获取智能体详情",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"id": map[string]interface{}{"type": "number", "description": "智能体ID"}}, "required": []string{"id"}},
		},
		{
			Name:        "create_agent",
			Description: "创建新智能体",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"name": map[string]interface{}{"type": "string", "description": "智能体名称"}, "code": map[string]interface{}{"type": "string", "description": "智能体代码"}, "system_prompt": map[string]interface{}{"type": "string", "description": "系统提示词"}, "llm_model": map[string]interface{}{"type": "string", "description": "使用的模型"}}, "required": []string{"name", "code", "system_prompt"}},
		},
		{
			Name:        "update_agent",
			Description: "更新智能体信息",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"id": map[string]interface{}{"type": "number", "description": "智能体ID"}, "name": map[string]interface{}{"type": "string"}, "system_prompt": map[string]interface{}{"type": "string"}, "llm_model": map[string]interface{}{"type": "string"}}, "required": []string{"id"}},
		},
		{
			Name:        "delete_agent",
			Description: "删除智能体",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"id": map[string]interface{}{"type": "number", "description": "智能体ID"}}, "required": []string{"id"}},
		},
		// 商店工具
		{
			Name:        "list_shops",
			Description: "获取所有商店列表",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
		},
		{
			Name:        "get_shop",
			Description: "根据ID获取商店详情",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"id": map[string]interface{}{"type": "number", "description": "商店ID"}}, "required": []string{"id"}},
		},
		{
			Name:        "create_shop",
			Description: "创建新商店",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"name": map[string]interface{}{"type": "string", "description": "商店名称"}, "code": map[string]interface{}{"type": "string", "description": "商店代码"}, "type": map[string]interface{}{"type": "string", "description": "商店类型"}, "owner_npc": map[string]interface{}{"type": "string", "description": "店主NPC代码"}}, "required": []string{"name", "code"}},
		},
		{
			Name:        "update_shop",
			Description: "更新商店信息",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"id": map[string]interface{}{"type": "number", "description": "商店ID"}, "name": map[string]interface{}{"type": "string"}, "description": map[string]interface{}{"type": "string"}}, "required": []string{"id"}},
		},
		{
			Name:        "delete_shop",
			Description: "删除商店",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"id": map[string]interface{}{"type": "number", "description": "商店ID"}}, "required": []string{"id"}},
		},
		// 道具工具
		{
			Name:        "list_items",
			Description: "获取所有道具列表",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
		},
		{
			Name:        "get_item",
			Description: "根据ID获取道具详情",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"id": map[string]interface{}{"type": "number", "description": "道具ID"}}, "required": []string{"id"}},
		},
		{
			Name:        "create_item",
			Description: "创建新道具",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"name": map[string]interface{}{"type": "string", "description": "道具名称"}, "code": map[string]interface{}{"type": "string", "description": "道具代码"}, "category": map[string]interface{}{"type": "string", "description": "分类"}, "description": map[string]interface{}{"type": "string", "description": "描述"}}, "required": []string{"name", "code"}},
		},
		{
			Name:        "update_item",
			Description: "更新道具信息",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"id": map[string]interface{}{"type": "number", "description": "道具ID"}, "name": map[string]interface{}{"type": "string"}, "description": map[string]interface{}{"type": "string"}, "effect": map[string]interface{}{"type": "string"}}, "required": []string{"id"}},
		},
		{
			Name:        "delete_item",
			Description: "删除道具",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"id": map[string]interface{}{"type": "number", "description": "道具ID"}}, "required": []string{"id"}},
		},
		// 任务工具
		{
			Name:        "list_tasks",
			Description: "获取所有任务列表",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
		},
		{
			Name:        "get_task",
			Description: "根据ID获取任务详情",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"id": map[string]interface{}{"type": "number", "description": "任务ID"}}, "required": []string{"id"}},
		},
		{
			Name:        "create_task",
			Description: "创建新任务",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"name": map[string]interface{}{"type": "string", "description": "任务名称"}, "code": map[string]interface{}{"type": "string", "description": "任务代码"}, "type": map[string]interface{}{"type": "string", "description": "任务类型"}, "description": map[string]interface{}{"type": "string", "description": "描述"}}, "required": []string{"name", "code"}},
		},
		{
			Name:        "update_task",
			Description: "更新任务信息",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"id": map[string]interface{}{"type": "number", "description": "任务ID"}, "name": map[string]interface{}{"type": "string"}, "description": map[string]interface{}{"type": "string"}, "status": map[string]interface{}{"type": "string"}}, "required": []string{"id"}},
		},
		{
			Name:        "delete_task",
			Description: "删除任务",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"id": map[string]interface{}{"type": "number", "description": "任务ID"}}, "required": []string{"id"}},
		},
		// 流程工具
		{
			Name:        "list_flows",
			Description: "获取所有流程列表",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
		},
		{
			Name:        "create_flow",
			Description: "创建新流程",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"name": map[string]interface{}{"type": "string", "description": "流程名称"}, "code": map[string]interface{}{"type": "string", "description": "流程代码"}, "description": map[string]interface{}{"type": "string", "description": "描述"}}, "required": []string{"name", "code"}},
		},
		// 提示词模板工具
		{
			Name:        "list_templates",
			Description: "获取所有提示词模板列表",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
		},
		{
			Name:        "create_template",
			Description: "创建新提示词模板",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"name": map[string]interface{}{"type": "string", "description": "模板名称"}, "code": map[string]interface{}{"type": "string", "description": "模板代码"}, "content": map[string]interface{}{"type": "string", "description": "模板内容"}}, "required": []string{"name", "code", "content"}},
		},
		// 生成工具
		{
			Name:        "generate_config",
			Description: "使用AI生成游戏配置（NPC、场景、任务等）",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"type": map[string]interface{}{"type": "string", "description": "生成类型: npc, scene, task, shop, item, agent, dialogue, flow"}, "description": map[string]interface{}{"type": "string", "description": "描述/需求"}, "action": map[string]interface{}{"type": "string", "description": "操作: create, complete, expand"}}, "required": []string{"type", "description"}},
		},
		// 数据工具
		{
			Name:        "export_data",
			Description: "导出所有游戏数据",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
		},
		{
			Name:        "get_game_stats",
			Description: "获取游戏数据统计",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
		},
	}
}

// GetTools 获取工具列表
func (s *Server) GetTools() []Tool {
	return s.tools
}

// HandleRequest 处理MCP请求
func (s *Server) HandleRequest(ctx context.Context, req MCPRequest) MCPResponse {
	switch req.Method {
	case "initialize":
		return s.handleInitialize(req)
	case "tools/list":
		return s.handleToolsList(req)
	case "tools/call":
		return s.handleToolsCall(ctx, req)
	default:
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32601,
				Message: fmt.Sprintf("Method not found: %s", req.Method),
			},
		}
	}
}

// handleInitialize 处理初始化
func (s *Server) handleInitialize(req MCPRequest) MCPResponse {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{},
			},
			"serverInfo": map[string]interface{}{
				"name":    "agentgame-mcp",
				"version": "1.0.0",
			},
		},
	}
}

// handleToolsList 处理工具列表
func (s *Server) handleToolsList(req MCPRequest) MCPResponse {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"tools": s.tools,
		},
	}
}

// handleToolsCall 处理工具调用
func (s *Server) handleToolsCall(ctx context.Context, req MCPRequest) MCPResponse {
	params, ok := req.Params.(map[string]interface{})
	if !ok {
		return s.errorResponse(req.ID, -32602, "Invalid params")
	}

	name, _ := params["name"].(string)
	args, _ := params["arguments"].(map[string]interface{})

	result, err := s.callTool(ctx, name, args)
	if err != nil {
		return s.errorResponse(req.ID, -32000, err.Error())
	}

	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}
}

// callTool 调用工具
func (s *Server) callTool(ctx context.Context, name string, args map[string]interface{}) (*ToolResult, error) {
	switch name {
	case "list_scenes":
		return s.listScenes(ctx)
	case "get_scene":
		id, _ := args["id"].(float64)
		return s.getScene(ctx, uint(id))
	case "create_scene":
		return s.createScene(ctx, args)
	case "update_scene":
		id, _ := args["id"].(float64)
		return s.updateScene(ctx, uint(id), args)
	case "delete_scene":
		id, _ := args["id"].(float64)
		return s.deleteScene(ctx, uint(id))
	case "list_npcs":
		return s.listNPCs(ctx)
	case "get_npc":
		id, _ := args["id"].(float64)
		return s.getNPC(ctx, uint(id))
	case "create_npc":
		return s.createNPC(ctx, args)
	case "update_npc":
		id, _ := args["id"].(float64)
		return s.updateNPC(ctx, uint(id), args)
	case "delete_npc":
		id, _ := args["id"].(float64)
		return s.deleteNPC(ctx, uint(id))
	case "list_agents":
		return s.listAgents(ctx)
	case "get_agent":
		id, _ := args["id"].(float64)
		return s.getAgent(ctx, uint(id))
	case "create_agent":
		return s.createAgent(ctx, args)
	case "update_agent":
		id, _ := args["id"].(float64)
		return s.updateAgent(ctx, uint(id), args)
	case "delete_agent":
		id, _ := args["id"].(float64)
		return s.deleteAgent(ctx, uint(id))
	case "list_shops":
		return s.listShops(ctx)
	case "get_shop":
		id, _ := args["id"].(float64)
		return s.getShop(ctx, uint(id))
	case "create_shop":
		return s.createShop(ctx, args)
	case "update_shop":
		id, _ := args["id"].(float64)
		return s.updateShop(ctx, uint(id), args)
	case "delete_shop":
		id, _ := args["id"].(float64)
		return s.deleteShop(ctx, uint(id))
	case "list_items":
		return s.listItems(ctx)
	case "get_item":
		id, _ := args["id"].(float64)
		return s.getItem(ctx, uint(id))
	case "create_item":
		return s.createItem(ctx, args)
	case "update_item":
		id, _ := args["id"].(float64)
		return s.updateItem(ctx, uint(id), args)
	case "delete_item":
		id, _ := args["id"].(float64)
		return s.deleteItem(ctx, uint(id))
	case "list_tasks":
		return s.listTasks(ctx)
	case "get_task":
		id, _ := args["id"].(float64)
		return s.getTask(ctx, uint(id))
	case "create_task":
		return s.createTask(ctx, args)
	case "update_task":
		id, _ := args["id"].(float64)
		return s.updateTask(ctx, uint(id), args)
	case "delete_task":
		id, _ := args["id"].(float64)
		return s.deleteTask(ctx, uint(id))
	case "list_flows":
		return s.listFlows(ctx)
	case "create_flow":
		return s.createFlow(ctx, args)
	case "list_templates":
		return s.listTemplates(ctx)
	case "create_template":
		return s.createTemplate(ctx, args)
	case "generate_config":
		return s.generateConfig(ctx, args)
	case "export_data":
		return s.exportData(ctx)
	case "get_game_stats":
		return s.getGameStats(ctx)
	default:
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
}

// errorResponse 创建错误响应
func (s *Server) errorResponse(id interface{}, code int, message string) MCPResponse {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &MCPError{
			Code:    code,
			Message: message,
		},
	}
}

// successResult 创建成功结果
func (s *Server) successResult(data interface{}) *ToolResult {
	jsonData, _ := json.MarshalIndent(data, "", "  ")
	return &ToolResult{
		Content: []Content{
			{Type: "text", Text: string(jsonData)},
		},
	}
}

// errorResult 创建错误结果
func (s *Server) errorResult(err error) *ToolResult {
	return &ToolResult{
		Content: []Content{
			{Type: "text", Text: fmt.Sprintf("Error: %v", err)},
		},
		IsError: true,
	}
}

// ==================== 场景操作 ====================

func (s *Server) listScenes(ctx context.Context) (*ToolResult, error) {
	scenes, err := s.repo.GetScenes(ctx)
	if err != nil {
		return s.errorResult(err), nil
	}
	return s.successResult(scenes), nil
}

func (s *Server) getScene(ctx context.Context, id uint) (*ToolResult, error) {
	scene, err := s.repo.GetSceneByID(ctx, id)
	if err != nil {
		return s.errorResult(err), nil
	}
	return s.successResult(scene), nil
}

func (s *Server) createScene(ctx context.Context, args map[string]interface{}) (*ToolResult, error) {
	scene := &models.Scene{
		Name:        getString(args, "name"),
		Code:        getString(args, "code"),
		Description: getString(args, "description"),
		Width:       getInt(args, "width", 1920),
		Height:      getInt(args, "height", 1080),
	}
	if err := s.repo.CreateScene(ctx, scene); err != nil {
		return s.errorResult(err), nil
	}
	return s.successResult(scene), nil
}

func (s *Server) updateScene(ctx context.Context, id uint, args map[string]interface{}) (*ToolResult, error) {
	scene, err := s.repo.GetSceneByID(ctx, id)
	if err != nil {
		return s.errorResult(err), nil
	}
	if v, ok := args["name"].(string); ok {
		scene.Name = v
	}
	if v, ok := args["description"].(string); ok {
		scene.Description = v
	}
	if v, ok := args["width"].(float64); ok {
		scene.Width = int(v)
	}
	if v, ok := args["height"].(float64); ok {
		scene.Height = int(v)
	}
	if err := s.repo.UpdateScene(ctx, scene); err != nil {
		return s.errorResult(err), nil
	}
	return s.successResult(scene), nil
}

func (s *Server) deleteScene(ctx context.Context, id uint) (*ToolResult, error) {
	if err := s.repo.DeleteScene(ctx, id); err != nil {
		return s.errorResult(err), nil
	}
	return s.successResult(map[string]string{"message": "Scene deleted"}), nil
}

// ==================== NPC操作 ====================

func (s *Server) listNPCs(ctx context.Context) (*ToolResult, error) {
	npcs, err := s.repo.GetNPCs(ctx)
	if err != nil {
		return s.errorResult(err), nil
	}
	return s.successResult(npcs), nil
}

func (s *Server) getNPC(ctx context.Context, id uint) (*ToolResult, error) {
	npc, err := s.repo.GetNPCByID(ctx, id)
	if err != nil {
		return s.errorResult(err), nil
	}
	return s.successResult(npc), nil
}

func (s *Server) createNPC(ctx context.Context, args map[string]interface{}) (*ToolResult, error) {
	npc := &models.NPC{
		Name:        getString(args, "name"),
		Code:        getString(args, "code"),
		Title:       getString(args, "title"),
		Description: getString(args, "description"),
	}
	if err := s.repo.CreateNPC(ctx, npc); err != nil {
		return s.errorResult(err), nil
	}
	return s.successResult(npc), nil
}

func (s *Server) updateNPC(ctx context.Context, id uint, args map[string]interface{}) (*ToolResult, error) {
	npc, err := s.repo.GetNPCByID(ctx, id)
	if err != nil {
		return s.errorResult(err), nil
	}
	if v, ok := args["name"].(string); ok {
		npc.Name = v
	}
	if v, ok := args["title"].(string); ok {
		npc.Title = v
	}
	if v, ok := args["description"].(string); ok {
		npc.Description = v
	}
	if err := s.repo.UpdateNPC(ctx, npc); err != nil {
		return s.errorResult(err), nil
	}
	return s.successResult(npc), nil
}

func (s *Server) deleteNPC(ctx context.Context, id uint) (*ToolResult, error) {
	if err := s.repo.DeleteNPC(ctx, id); err != nil {
		return s.errorResult(err), nil
	}
	return s.successResult(map[string]string{"message": "NPC deleted"}), nil
}

// ==================== 智能体操作 ====================

func (s *Server) listAgents(ctx context.Context) (*ToolResult, error) {
	agents, err := s.repo.GetAgents(ctx)
	if err != nil {
		return s.errorResult(err), nil
	}
	return s.successResult(agents), nil
}

func (s *Server) getAgent(ctx context.Context, id uint) (*ToolResult, error) {
	agent, err := s.repo.GetAgentByID(ctx, id)
	if err != nil {
		return s.errorResult(err), nil
	}
	return s.successResult(agent), nil
}

func (s *Server) createAgent(ctx context.Context, args map[string]interface{}) (*ToolResult, error) {
	agent := &models.Agent{
		Name:         getString(args, "name"),
		Code:         getString(args, "code"),
		SystemPrompt: getString(args, "system_prompt"),
		LLMModel:     getString(args, "llm_model"),
	}
	if err := s.repo.CreateAgent(ctx, agent); err != nil {
		return s.errorResult(err), nil
	}
	return s.successResult(agent), nil
}

func (s *Server) updateAgent(ctx context.Context, id uint, args map[string]interface{}) (*ToolResult, error) {
	agent, err := s.repo.GetAgentByID(ctx, id)
	if err != nil {
		return s.errorResult(err), nil
	}
	if v, ok := args["name"].(string); ok {
		agent.Name = v
	}
	if v, ok := args["system_prompt"].(string); ok {
		agent.SystemPrompt = v
	}
	if v, ok := args["llm_model"].(string); ok {
		agent.LLMModel = v
	}
	if err := s.repo.UpdateAgent(ctx, agent); err != nil {
		return s.errorResult(err), nil
	}
	return s.successResult(agent), nil
}

func (s *Server) deleteAgent(ctx context.Context, id uint) (*ToolResult, error) {
	if err := s.repo.DeleteAgent(ctx, id); err != nil {
		return s.errorResult(err), nil
	}
	return s.successResult(map[string]string{"message": "Agent deleted"}), nil
}

// ==================== 商店操作 ====================

func (s *Server) listShops(ctx context.Context) (*ToolResult, error) {
	shops, err := s.repo.GetShops(ctx)
	if err != nil {
		return s.errorResult(err), nil
	}
	return s.successResult(shops), nil
}

func (s *Server) getShop(ctx context.Context, id uint) (*ToolResult, error) {
	shop, err := s.repo.GetShopByID(ctx, id)
	if err != nil {
		return s.errorResult(err), nil
	}
	return s.successResult(shop), nil
}

func (s *Server) createShop(ctx context.Context, args map[string]interface{}) (*ToolResult, error) {
	shop := &models.Shop{
		Name:     getString(args, "name"),
		Code:     getString(args, "code"),
		Type:     getString(args, "type"),
		OwnerNPC: getString(args, "owner_npc"),
	}
	if err := s.repo.CreateShop(ctx, shop); err != nil {
		return s.errorResult(err), nil
	}
	return s.successResult(shop), nil
}

func (s *Server) updateShop(ctx context.Context, id uint, args map[string]interface{}) (*ToolResult, error) {
	shop, err := s.repo.GetShopByID(ctx, id)
	if err != nil {
		return s.errorResult(err), nil
	}
	if v, ok := args["name"].(string); ok {
		shop.Name = v
	}
	if v, ok := args["description"].(string); ok {
		shop.Description = v
	}
	if err := s.repo.UpdateShop(ctx, shop); err != nil {
		return s.errorResult(err), nil
	}
	return s.successResult(shop), nil
}

func (s *Server) deleteShop(ctx context.Context, id uint) (*ToolResult, error) {
	if err := s.repo.DeleteShop(ctx, id); err != nil {
		return s.errorResult(err), nil
	}
	return s.successResult(map[string]string{"message": "Shop deleted"}), nil
}

// ==================== 道具操作 ====================

func (s *Server) listItems(ctx context.Context) (*ToolResult, error) {
	items, err := s.repo.GetItems(ctx)
	if err != nil {
		return s.errorResult(err), nil
	}
	return s.successResult(items), nil
}

func (s *Server) getItem(ctx context.Context, id uint) (*ToolResult, error) {
	item, err := s.repo.GetItemByID(ctx, id)
	if err != nil {
		return s.errorResult(err), nil
	}
	return s.successResult(item), nil
}

func (s *Server) createItem(ctx context.Context, args map[string]interface{}) (*ToolResult, error) {
	item := &models.Item{
		Name:        getString(args, "name"),
		Code:        getString(args, "code"),
		Category:    getString(args, "category"),
		Description: getString(args, "description"),
	}
	if err := s.repo.CreateItem(ctx, item); err != nil {
		return s.errorResult(err), nil
	}
	return s.successResult(item), nil
}

func (s *Server) updateItem(ctx context.Context, id uint, args map[string]interface{}) (*ToolResult, error) {
	item, err := s.repo.GetItemByID(ctx, id)
	if err != nil {
		return s.errorResult(err), nil
	}
	if v, ok := args["name"].(string); ok {
		item.Name = v
	}
	if v, ok := args["description"].(string); ok {
		item.Description = v
	}
	if v, ok := args["effect"].(string); ok {
		item.Effect = v
	}
	if err := s.repo.UpdateItem(ctx, item); err != nil {
		return s.errorResult(err), nil
	}
	return s.successResult(item), nil
}

func (s *Server) deleteItem(ctx context.Context, id uint) (*ToolResult, error) {
	if err := s.repo.DeleteItem(ctx, id); err != nil {
		return s.errorResult(err), nil
	}
	return s.successResult(map[string]string{"message": "Item deleted"}), nil
}

// ==================== 任务操作 ====================

func (s *Server) listTasks(ctx context.Context) (*ToolResult, error) {
	tasks, err := s.repo.GetTasks(ctx)
	if err != nil {
		return s.errorResult(err), nil
	}
	return s.successResult(tasks), nil
}

func (s *Server) getTask(ctx context.Context, id uint) (*ToolResult, error) {
	task, err := s.repo.GetTaskByID(ctx, id)
	if err != nil {
		return s.errorResult(err), nil
	}
	return s.successResult(task), nil
}

func (s *Server) createTask(ctx context.Context, args map[string]interface{}) (*ToolResult, error) {
	task := &models.Task{
		Name:        getString(args, "name"),
		Code:        getString(args, "code"),
		Type:        getString(args, "type"),
		Description: getString(args, "description"),
	}
	if err := s.repo.CreateTask(ctx, task); err != nil {
		return s.errorResult(err), nil
	}
	return s.successResult(task), nil
}

func (s *Server) updateTask(ctx context.Context, id uint, args map[string]interface{}) (*ToolResult, error) {
	task, err := s.repo.GetTaskByID(ctx, id)
	if err != nil {
		return s.errorResult(err), nil
	}
	if v, ok := args["name"].(string); ok {
		task.Name = v
	}
	if v, ok := args["description"].(string); ok {
		task.Description = v
	}
	if v, ok := args["status"].(string); ok {
		task.Status = v
	}
	if err := s.repo.UpdateTask(ctx, task); err != nil {
		return s.errorResult(err), nil
	}
	return s.successResult(task), nil
}

func (s *Server) deleteTask(ctx context.Context, id uint) (*ToolResult, error) {
	if err := s.repo.DeleteTask(ctx, id); err != nil {
		return s.errorResult(err), nil
	}
	return s.successResult(map[string]string{"message": "Task deleted"}), nil
}

// ==================== 流程操作 ====================

func (s *Server) listFlows(ctx context.Context) (*ToolResult, error) {
	flows, err := s.repo.GetFlows(ctx)
	if err != nil {
		return s.errorResult(err), nil
	}
	return s.successResult(flows), nil
}

func (s *Server) createFlow(ctx context.Context, args map[string]interface{}) (*ToolResult, error) {
	flow := &models.Flow{
		Name:        getString(args, "name"),
		Code:        getString(args, "code"),
		Description: getString(args, "description"),
	}
	if err := s.repo.CreateFlow(ctx, flow); err != nil {
		return s.errorResult(err), nil
	}
	return s.successResult(flow), nil
}

// ==================== 模板操作 ====================

func (s *Server) listTemplates(ctx context.Context) (*ToolResult, error) {
	templates, err := s.repo.GetTemplates(ctx)
	if err != nil {
		return s.errorResult(err), nil
	}
	return s.successResult(templates), nil
}

func (s *Server) createTemplate(ctx context.Context, args map[string]interface{}) (*ToolResult, error) {
	template := &models.PromptTemplate{
		Name:    getString(args, "name"),
		Code:    getString(args, "code"),
		Content: getString(args, "content"),
	}
	if err := s.repo.CreateTemplate(ctx, template); err != nil {
		return s.errorResult(err), nil
	}
	return s.successResult(template), nil
}

// ==================== 生成操作 ====================

func (s *Server) generateConfig(ctx context.Context, args map[string]interface{}) (*ToolResult, error) {
	if s.generator == nil || !s.generator.IsEnabled() {
		return s.errorResult(fmt.Errorf("generator not enabled")), nil
	}

	req := generator.GenerateRequest{
		Type:   getString(args, "type"),
		Action: getString(args, "action"),
		Params: map[string]interface{}{
			"description": getString(args, "description"),
		},
	}

	resp, err := s.generator.Generate(ctx, req)
	if err != nil {
		return s.errorResult(err), nil
	}

	return s.successResult(resp), nil
}

// ==================== 数据操作 ====================

func (s *Server) exportData(ctx context.Context) (*ToolResult, error) {
	data := make(map[string]interface{})

	scenes, _ := s.repo.GetScenes(ctx)
	data["scenes"] = scenes

	npcs, _ := s.repo.GetNPCs(ctx)
	data["npcs"] = npcs

	agents, _ := s.repo.GetAgents(ctx)
	data["agents"] = agents

	shops, _ := s.repo.GetShops(ctx)
	data["shops"] = shops

	items, _ := s.repo.GetItems(ctx)
	data["items"] = items

	tasks, _ := s.repo.GetTasks(ctx)
	data["tasks"] = tasks

	return s.successResult(data), nil
}

func (s *Server) getGameStats(ctx context.Context) (*ToolResult, error) {
	scenes, _ := s.repo.GetScenes(ctx)
	npcs, _ := s.repo.GetNPCs(ctx)
	agents, _ := s.repo.GetAgents(ctx)
	shops, _ := s.repo.GetShops(ctx)
	items, _ := s.repo.GetItems(ctx)
	tasks, _ := s.repo.GetTasks(ctx)
	flows, _ := s.repo.GetFlows(ctx)

	stats := map[string]interface{}{
		"scenes":  len(scenes),
		"npcs":    len(npcs),
		"agents":  len(agents),
		"shops":   len(shops),
		"items":   len(items),
		"tasks":   len(tasks),
		"flows":   len(flows),
	}

	return s.successResult(stats), nil
}

// ==================== 辅助函数 ====================

func getString(args map[string]interface{}, key string) string {
	if v, ok := args[key].(string); ok {
		return v
	}
	return ""
}

func getInt(args map[string]interface{}, key string, defaultVal int) int {
	if v, ok := args[key].(float64); ok {
		return int(v)
	}
	return defaultVal
}

// HandleHTTP 处理HTTP请求
func (s *Server) HandleHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(204)
		return
	}

	var req MCPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = json.NewEncoder(w).Encode(MCPResponse{
			JSONRPC: "2.0",
			Error: &MCPError{
				Code:    -32700,
				Message: "Parse error",
			},
		})
		return
	}

	resp := s.HandleRequest(r.Context(), req)
	_ = json.NewEncoder(w).Encode(resp)
}

// Log 工具调用日志
func (s *Server) Log(format string, args ...interface{}) {
	log.Printf("[MCP] "+format, args...)
}

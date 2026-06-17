```go
package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ddc-111/agentGame/server/internal/database/models"
	"github.com/ddc-111/agentGame/server/internal/database/repository"
	"github.com/ddc-111/agentGame/server/internal/generator"
)

// MockRepository implements the repository.Repository interface for testing
type MockRepository struct {
	Scenes    map[uint]*models.Scene
	NPCs      map[uint]*models.NPC
	Agents    map[uint]*models.Agent
	Shops     map[uint]*models.Shop
	Items     map[uint]*models.Item
	Tasks     map[uint]*models.Task
	Flows     map[uint]*models.Flow
	Templates map[uint]*models.PromptTemplate
	Error     error
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		Scenes:    make(map[uint]*models.Scene),
		NPCs:      make(map[uint]*models.NPC),
		Agents:    make(map[uint]*models.Agent),
		Shops:     make(map[uint]*models.Shop),
		Items:     make(map[uint]*models.Item),
		Tasks:     make(map[uint]*models.Task),
		Flows:     make(map[uint]*models.Flow),
		Templates: make(map[uint]*models.PromptTemplate),
	}
}

func (m *MockRepository) GetScenes(ctx context.Context) ([]models.Scene, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	var scenes []models.Scene
	for _, scene := range m.Scenes {
		scenes = append(scenes, *scene)
	}
	return scenes, nil
}

func (m *MockRepository) GetSceneByID(ctx context.Context, id uint) (*models.Scene, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	scene, exists := m.Scenes[id]
	if !exists {
		return nil, fmt.Errorf("scene not found: %d", id)
	}
	return scene, nil
}

func (m *MockRepository) CreateScene(ctx context.Context, scene *models.Scene) error {
	if m.Error != nil {
		return m.Error
	}
	scene.ID = uint(len(m.Scenes) + 1)
	m.Scenes[scene.ID] = scene
	return nil
}

func (m *MockRepository) UpdateScene(ctx context.Context, scene *models.Scene) error {
	if m.Error != nil {
		return m.Error
	}
	m.Scenes[scene.ID] = scene
	return nil
}

func (m *MockRepository) DeleteScene(ctx context.Context, id uint) error {
	if m.Error != nil {
		return m.Error
	}
	delete(m.Scenes, id)
	return nil
}

func (m *MockRepository) GetNPCs(ctx context.Context) ([]models.NPC, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	var npcs []models.NPC
	for _, npc := range m.NPCs {
		npcs = append(npcs, *npc)
	}
	return npcs, nil
}

func (m *MockRepository) GetNPCByID(ctx context.Context, id uint) (*models.NPC, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	npc, exists := m.NPCs[id]
	if !exists {
		return nil, fmt.Errorf("npc not found: %d", id)
	}
	return npc, nil
}

func (m *MockRepository) CreateNPC(ctx context.Context, npc *models.NPC) error {
	if m.Error != nil {
		return m.Error
	}
	npc.ID = uint(len(m.NPCs) + 1)
	m.NPCs[npc.ID] = npc
	return nil
}

func (m *MockRepository) UpdateNPC(ctx context.Context, npc *models.NPC) error {
	if m.Error != nil {
		return m.Error
	}
	m.NPCs[npc.ID] = npc
	return nil
}

func (m *MockRepository) DeleteNPC(ctx context.Context, id uint) error {
	if m.Error != nil {
		return m.Error
	}
	delete(m.NPCs, id)
	return nil
}

func (m *MockRepository) GetAgents(ctx context.Context) ([]models.Agent, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	var agents []models.Agent
	for _, agent := range m.Agents {
		agents = append(agents, *agent)
	}
	return agents, nil
}

func (m *MockRepository) GetAgentByID(ctx context.Context, id uint) (*models.Agent, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	agent, exists := m.Agents[id]
	if !exists {
		return nil, fmt.Errorf("agent not found: %d", id)
	}
	return agent, nil
}

func (m *MockRepository) CreateAgent(ctx context.Context, agent *models.Agent) error {
	if m.Error != nil {
		return m.Error
	}
	agent.ID = uint(len(m.Agents) + 1)
	m.Agents[agent.ID] = agent
	return nil
}

func (m *MockRepository) UpdateAgent(ctx context.Context, agent *models.Agent) error {
	if m.Error != nil {
		return m.Error
	}
	m.Agents[agent.ID] = agent
	return nil
}

func (m *MockRepository) DeleteAgent(ctx context.Context, id uint) error {
	if m.Error != nil {
		return m.Error
	}
	delete(m.Agents, id)
	return nil
}

func (m *MockRepository) GetShops(ctx context.Context) ([]models.Shop, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	var shops []models.Shop
	for _, shop := range m.Shops {
		shops = append(shops, *shop)
	}
	return shops, nil
}

func (m *MockRepository) GetShopByID(ctx context.Context, id uint) (*models.Shop, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	shop, exists := m.Shops[id]
	if !exists {
		return nil, fmt.Errorf("shop not found: %d", id)
	}
	return shop, nil
}

func (m *MockRepository) CreateShop(ctx context.Context, shop *models.Shop) error {
	if m.Error != nil {
		return m.Error
	}
	shop.ID = uint(len(m.Shops) + 1)
	m.Shops[shop.ID] = shop
	return nil
}

func (m *MockRepository) UpdateShop(ctx context.Context, shop *models.Shop) error {
	if m.Error != nil {
		return m.Error
	}
	m.Shops[shop.ID] = shop
	return nil
}

func (m *MockRepository) DeleteShop(ctx context.Context, id uint) error {
	if m.Error != nil {
		return m.Error
	}
	delete(m.Shops, id)
	return nil
}

func (m *MockRepository) GetItems(ctx context.Context) ([]models.Item, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	var items []models.Item
	for _, item := range m.Items {
		items = append(items, *item)
	}
	return items, nil
}

func (m *MockRepository) GetItemByID(ctx context.Context, id uint) (*models.Item, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	item, exists := m.Items[id]
	if !exists {
		return nil, fmt.Errorf("item not found: %d", id)
	}
	return item, nil
}

func (m *MockRepository) CreateItem(ctx context.Context, item *models.Item) error {
	if m.Error != nil {
		return m.Error
	}
	item.ID = uint(len(m.Items) + 1)
	m.Items[item.ID] = item
	return nil
}

func (m *MockRepository) UpdateItem(ctx context.Context, item *models.Item) error {
	if m.Error != nil {
		return m.Error
	}
	m.Items[item.ID] = item
	return nil
}

func (m *MockRepository) DeleteItem(ctx context.Context, id uint) error {
	if m.Error != nil {
		return m.Error
	}
	delete(m.Items, id)
	return nil
}

func (m *MockRepository) GetTasks(ctx context.Context) ([]models.Task, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	var tasks []models.Task
	for _, task := range m.Tasks {
		tasks = append(tasks, *task)
	}
	return tasks, nil
}

func (m *MockRepository) GetTaskByID(ctx context.Context, id uint) (*models.Task, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	task, exists := m.Tasks[id]
	if !exists {
		return nil, fmt.Errorf("task not found: %d", id)
	}
	return task, nil
}

func (m *MockRepository) CreateTask(ctx context.Context, task *models.Task) error {
	if m.Error != nil {
		return m.Error
	}
	task.ID = uint(len(m.Tasks) + 1)
	m.Tasks[task.ID] = task
	return nil
}

func (m *MockRepository) UpdateTask(ctx context.Context, task *models.Task) error {
	if m.Error != nil {
		return m.Error
	}
	m.Tasks[task.ID] = task
	return nil
}

func (m *MockRepository) DeleteTask(ctx context.Context, id uint) error {
	if m.Error != nil {
		return m.Error
	}
	delete(m.Tasks, id)
	return nil
}

func (m *MockRepository) GetFlows(ctx context.Context) ([]models.Flow, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	var flows []models.Flow
	for _, flow := range m.Flows {
		flows = append(flows, *flow)
	}
	return flows, nil
}

func (m *MockRepository) CreateFlow(ctx context.Context, flow *models.Flow) error {
	if m.Error != nil {
		return m.Error
	}
	flow.ID = uint(len(m.Flows) + 1)
	m.Flows[flow.ID] = flow
	return nil
}

func (m *MockRepository) GetTemplates(ctx context.Context) ([]models.PromptTemplate, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	var templates []models.PromptTemplate
	for _, template := range m.Templates {
		templates = append(templates, *template)
	}
	return templates, nil
}

func (m *MockRepository) CreateTemplate(ctx context.Context, template *models.PromptTemplate) error {
	if m.Error != nil {
		return m.Error
	}
	template.ID = uint(len(m.Templates) + 1)
	m.Templates[template.ID] = template
	return nil
}

// MockGenerator implements the generator.Generator interface for testing
type MockGenerator struct {
	Enabled bool
	Response *generator.GenerateResponse
	Error    error
}

func NewMockGenerator() *MockGenerator {
	return &MockGenerator{
		Enabled: true,
	}
}

func (m *MockGenerator) IsEnabled() bool {
	return m.Enabled
}

func (m *MockGenerator) Generate(ctx context.Context, req generator.GenerateRequest) (*generator.GenerateResponse, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	if m.Response != nil {
		return m.Response, nil
	}
	return &generator.GenerateResponse{
		Result: "Generated content",
	}, nil
}

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		repo      *MockRepository
		gen       *MockGenerator
		wantTools int
	}{
		{
			name:      "Valid initialization",
			repo:      NewMockRepository(),
			gen:       NewMockGenerator(),
			wantTools: 38, // Expected number of tools from initTools()
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := New(tt.repo, tt.gen)
			if server == nil {
				t.Fatal("Expected server to be created, got nil")
			}
			if len(server.tools) != tt.wantTools {
				t.Errorf("Expected %d tools, got %d", tt.wantTools, len(server.tools))
			}
			if len(server.resources) == 0 {
				t.Error("Expected resources to be initialized")
			}
			if len(server.prompts) == 0 {
				t.Error("Expected prompts to be initialized")
			}
		})
	}
}

func TestGetTools(t *testing.T) {
	server := New(NewMockRepository(), NewMockGenerator())
	tools := server.GetTools()

	if len(tools) == 0 {
		t.Error("Expected non-empty tools list")
	}

	// Verify some expected tool names
	expectedTools := []string{
		"list_scenes",
		"get_scene",
		"create_scene",
		"update_scene",
		"delete_scene",
	}

	toolNames := make(map[string]bool)
	for _, tool := range tools {
		toolNames[tool.Name] = true
	}

	for _, expected := range expectedTools {
		if !toolNames[expected] {
			t.Errorf("Expected tool %s not found", expected)
		}
	}
}

func TestGetResources(t *testing.T) {
	server := New(NewMockRepository(), NewMockGenerator())
	resources := server.GetResources()

	if len(resources) == 0 {
		t.Error("Expected non-empty resources list")
	}

	expectedResources := []string{
		"game_state://scenes",
		"game_state://npcs",
		"game_state://agents",
	}

	resourceURIs := make(map[string]bool)
	for _, res := range resources {
		resourceURIs[res.URI] = true
	}

	for _, expected := range expectedResources {
		if !resourceURIs[expected] {
			t.Errorf("Expected resource %s not found", expected)
		}
	}
}

func TestGetPrompts(t *testing.T) {
	server := New(NewMockRepository(), NewMockGenerator())
	prompts := server.GetPrompts()

	if len(prompts) == 0 {
		t.Error("Expected non-empty prompts list")
	}

	expectedPrompts := []string{
		"npc_personality",
		"npc_dialogue",
		"scene_description",
		"quest_design",
	}

	promptNames := make(map[string]bool)
	for _, prompt := range prompts {
		promptNames[prompt.Name] = true
	}

	for _, expected := range expectedPrompts {
		if !promptNames[expected] {
			t.Errorf("Expected prompt %s not found", expected)
		}
	}
}

func TestHandleRequest(t *testing.T) {
	server := New(NewMockRepository(), NewMockGenerator())

	tests := []struct {
		name       string
		method     string
		id         interface{}
		params     interface{}
		wantError  bool
		wantResult bool
	}{
		{
			name:       "Initialize request",
			method:     "initialize",
			id:         1,
			params:     nil,
			wantError:  false,
			wantResult: true,
		},
		{
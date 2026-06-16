package repository

import (
	"github.com/ddc-111/agentGame/server/internal/database/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Transaction(fn func(repo *Repository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(&Repository{db: tx})
	})
}

// Scene 场景相关操作
func (r *Repository) GetScenes() ([]models.Scene, error) {
	var scenes []models.Scene
	err := r.db.Preload("SceneNPCs").Preload("Portals").Find(&scenes).Error
	return scenes, err
}

func (r *Repository) GetScenesPaginated(offset, limit int) ([]models.Scene, int64, error) {
	var scenes []models.Scene
	var total int64
	r.db.Model(&models.Scene{}).Count(&total)
	err := r.db.Preload("SceneNPCs").Preload("Portals").Offset(offset).Limit(limit).Find(&scenes).Error
	return scenes, total, err
}

func (r *Repository) GetSceneByID(id uint) (*models.Scene, error) {
	var scene models.Scene
	err := r.db.Preload("SceneNPCs").Preload("Portals").First(&scene, id).Error
	return &scene, err
}

func (r *Repository) GetSceneByCode(code string) (*models.Scene, error) {
	var scene models.Scene
	err := r.db.Preload("SceneNPCs").Preload("Portals").Where("code = ?", code).First(&scene).Error
	return &scene, err
}

func (r *Repository) CreateScene(scene *models.Scene) error {
	return r.db.Create(scene).Error
}

func (r *Repository) UpdateScene(scene *models.Scene) error {
	return r.db.Save(scene).Error
}

func (r *Repository) DeleteScene(id uint) error {
	return r.db.Delete(&models.Scene{}, id).Error
}

// NPC 相关操作
func (r *Repository) GetNPCs() ([]models.NPC, error) {
	var npcs []models.NPC
	err := r.db.Preload("Agent").Preload("Shop").Find(&npcs).Error
	return npcs, err
}

func (r *Repository) GetNPCsPaginated(offset, limit int) ([]models.NPC, int64, error) {
	var npcs []models.NPC
	var total int64
	r.db.Model(&models.NPC{}).Count(&total)
	err := r.db.Preload("Agent").Preload("Shop").Offset(offset).Limit(limit).Find(&npcs).Error
	return npcs, total, err
}

func (r *Repository) GetNPCByID(id uint) (*models.NPC, error) {
	var npc models.NPC
	err := r.db.Preload("Agent").Preload("Shop").First(&npc, id).Error
	return &npc, err
}

func (r *Repository) GetNPCByCode(code string) (*models.NPC, error) {
	var npc models.NPC
	err := r.db.Preload("Agent").Preload("Shop").Where("code = ?", code).First(&npc).Error
	return &npc, err
}

func (r *Repository) CreateNPC(npc *models.NPC) error {
	return r.db.Create(npc).Error
}

func (r *Repository) UpdateNPC(npc *models.NPC) error {
	return r.db.Save(npc).Error
}

func (r *Repository) DeleteNPC(id uint) error {
	return r.db.Delete(&models.NPC{}, id).Error
}

// Agent 智能体相关操作
func (r *Repository) GetAgents() ([]models.Agent, error) {
	var agents []models.Agent
	err := r.db.Find(&agents).Error
	return agents, err
}

func (r *Repository) GetAgentsPaginated(offset, limit int) ([]models.Agent, int64, error) {
	var agents []models.Agent
	var total int64
	r.db.Model(&models.Agent{}).Count(&total)
	err := r.db.Offset(offset).Limit(limit).Find(&agents).Error
	return agents, total, err
}

func (r *Repository) GetAgentByID(id uint) (*models.Agent, error) {
	var agent models.Agent
	err := r.db.First(&agent, id).Error
	return &agent, err
}

func (r *Repository) GetAgentByCode(code string) (*models.Agent, error) {
	var agent models.Agent
	err := r.db.Where("code = ?", code).First(&agent).Error
	return &agent, err
}

func (r *Repository) CreateAgent(agent *models.Agent) error {
	return r.db.Create(agent).Error
}

func (r *Repository) UpdateAgent(agent *models.Agent) error {
	return r.db.Save(agent).Error
}

func (r *Repository) DeleteAgent(id uint) error {
	return r.db.Delete(&models.Agent{}, id).Error
}

// LLMProvider 相关操作
func (r *Repository) GetProviders() ([]models.LLMProvider, error) {
	var providers []models.LLMProvider
	err := r.db.Find(&providers).Error
	return providers, err
}

func (r *Repository) GetProvidersPaginated(offset, limit int) ([]models.LLMProvider, int64, error) {
	var providers []models.LLMProvider
	var total int64
	r.db.Model(&models.LLMProvider{}).Count(&total)
	err := r.db.Offset(offset).Limit(limit).Find(&providers).Error
	return providers, total, err
}

func (r *Repository) GetProviderByID(id uint) (*models.LLMProvider, error) {
	var provider models.LLMProvider
	err := r.db.First(&provider, id).Error
	return &provider, err
}

func (r *Repository) CreateProvider(provider *models.LLMProvider) error {
	return r.db.Create(provider).Error
}

func (r *Repository) UpdateProvider(provider *models.LLMProvider) error {
	return r.db.Save(provider).Error
}

func (r *Repository) DeleteProvider(id uint) error {
	return r.db.Delete(&models.LLMProvider{}, id).Error
}

// PromptTemplate 相关操作
func (r *Repository) GetTemplates() ([]models.PromptTemplate, error) {
	var templates []models.PromptTemplate
	err := r.db.Find(&templates).Error
	return templates, err
}

func (r *Repository) GetTemplatesPaginated(offset, limit int) ([]models.PromptTemplate, int64, error) {
	var templates []models.PromptTemplate
	var total int64
	r.db.Model(&models.PromptTemplate{}).Count(&total)
	err := r.db.Offset(offset).Limit(limit).Find(&templates).Error
	return templates, total, err
}

func (r *Repository) GetTemplateByID(id uint) (*models.PromptTemplate, error) {
	var template models.PromptTemplate
	err := r.db.First(&template, id).Error
	return &template, err
}

func (r *Repository) GetTemplateByCode(code string) (*models.PromptTemplate, error) {
	var template models.PromptTemplate
	err := r.db.Where("code = ?", code).First(&template).Error
	return &template, err
}

func (r *Repository) CreateTemplate(template *models.PromptTemplate) error {
	return r.db.Create(template).Error
}

func (r *Repository) UpdateTemplate(template *models.PromptTemplate) error {
	return r.db.Save(template).Error
}

func (r *Repository) DeleteTemplate(id uint) error {
	return r.db.Delete(&models.PromptTemplate{}, id).Error
}

// Shop 商店相关操作
func (r *Repository) GetShops() ([]models.Shop, error) {
	var shops []models.Shop
	err := r.db.Preload("Items.Item").Find(&shops).Error
	return shops, err
}

func (r *Repository) GetShopsPaginated(offset, limit int) ([]models.Shop, int64, error) {
	var shops []models.Shop
	var total int64
	r.db.Model(&models.Shop{}).Count(&total)
	err := r.db.Preload("Items.Item").Offset(offset).Limit(limit).Find(&shops).Error
	return shops, total, err
}

func (r *Repository) GetShopByID(id uint) (*models.Shop, error) {
	var shop models.Shop
	err := r.db.Preload("Items.Item").First(&shop, id).Error
	return &shop, err
}

func (r *Repository) GetShopByCode(code string) (*models.Shop, error) {
	var shop models.Shop
	err := r.db.Preload("Items.Item").Where("code = ?", code).First(&shop).Error
	return &shop, err
}

func (r *Repository) CreateShop(shop *models.Shop) error {
	return r.db.Create(shop).Error
}

func (r *Repository) UpdateShop(shop *models.Shop) error {
	return r.db.Save(shop).Error
}

func (r *Repository) DeleteShop(id uint) error {
	return r.db.Delete(&models.Shop{}, id).Error
}

// Item 道具相关操作
func (r *Repository) GetItems() ([]models.Item, error) {
	var items []models.Item
	err := r.db.Find(&items).Error
	return items, err
}

func (r *Repository) GetItemsPaginated(offset, limit int) ([]models.Item, int64, error) {
	var items []models.Item
	var total int64
	r.db.Model(&models.Item{}).Count(&total)
	err := r.db.Offset(offset).Limit(limit).Find(&items).Error
	return items, total, err
}

func (r *Repository) GetItemByID(id uint) (*models.Item, error) {
	var item models.Item
	err := r.db.First(&item, id).Error
	return &item, err
}

func (r *Repository) GetItemByCode(code string) (*models.Item, error) {
	var item models.Item
	err := r.db.Where("code = ?", code).First(&item).Error
	return &item, err
}

func (r *Repository) CreateItem(item *models.Item) error {
	return r.db.Create(item).Error
}

func (r *Repository) UpdateItem(item *models.Item) error {
	return r.db.Save(item).Error
}

func (r *Repository) DeleteItem(id uint) error {
	return r.db.Delete(&models.Item{}, id).Error
}

// Task 任务相关操作
func (r *Repository) GetTasks() ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.Find(&tasks).Error
	return tasks, err
}

func (r *Repository) GetTasksPaginated(offset, limit int) ([]models.Task, int64, error) {
	var tasks []models.Task
	var total int64
	r.db.Model(&models.Task{}).Count(&total)
	err := r.db.Offset(offset).Limit(limit).Find(&tasks).Error
	return tasks, total, err
}

func (r *Repository) GetTaskByID(id uint) (*models.Task, error) {
	var task models.Task
	err := r.db.First(&task, id).Error
	return &task, err
}

func (r *Repository) GetTaskByCode(code string) (*models.Task, error) {
	var task models.Task
	err := r.db.Where("code = ?", code).First(&task).Error
	return &task, err
}

func (r *Repository) CreateTask(task *models.Task) error {
	return r.db.Create(task).Error
}

func (r *Repository) UpdateTask(task *models.Task) error {
	return r.db.Save(task).Error
}

func (r *Repository) DeleteTask(id uint) error {
	return r.db.Delete(&models.Task{}, id).Error
}

// Flow 流程相关操作
func (r *Repository) GetFlows() ([]models.Flow, error) {
	var flows []models.Flow
	err := r.db.Find(&flows).Error
	return flows, err
}

func (r *Repository) GetFlowsPaginated(offset, limit int) ([]models.Flow, int64, error) {
	var flows []models.Flow
	var total int64
	r.db.Model(&models.Flow{}).Count(&total)
	err := r.db.Offset(offset).Limit(limit).Find(&flows).Error
	return flows, total, err
}

func (r *Repository) GetFlowByID(id uint) (*models.Flow, error) {
	var flow models.Flow
	err := r.db.First(&flow, id).Error
	return &flow, err
}

func (r *Repository) GetFlowByCode(code string) (*models.Flow, error) {
	var flow models.Flow
	err := r.db.Where("code = ?", code).First(&flow).Error
	return &flow, err
}

func (r *Repository) CreateFlow(flow *models.Flow) error {
	return r.db.Create(flow).Error
}

func (r *Repository) UpdateFlow(flow *models.Flow) error {
	return r.db.Save(flow).Error
}

func (r *Repository) DeleteFlow(id uint) error {
	return r.db.Delete(&models.Flow{}, id).Error
}

// GameConfig 游戏配置相关操作
func (r *Repository) GetConfig(key string) (*models.GameConfig, error) {
	var config models.GameConfig
	err := r.db.Where("key = ?", key).First(&config).Error
	return &config, err
}

func (r *Repository) SetConfig(key, value string) error {
	var config models.GameConfig
	result := r.db.Where("key = ?", key).First(&config)
	if result.Error != nil {
		config = models.GameConfig{Key: key, Value: value}
		return r.db.Create(&config).Error
	}
	config.Value = value
	return r.db.Save(&config).Error
}

// Player 玩家相关操作
func (r *Repository) GetPlayers() ([]models.Player, error) {
	var players []models.Player
	err := r.db.Find(&players).Error
	return players, err
}

func (r *Repository) GetPlayersPaginated(offset, limit int) ([]models.Player, int64, error) {
	var players []models.Player
	var total int64
	r.db.Model(&models.Player{}).Count(&total)
	err := r.db.Offset(offset).Limit(limit).Find(&players).Error
	return players, total, err
}

func (r *Repository) GetPlayerByID(id uint) (*models.Player, error) {
	var player models.Player
	err := r.db.First(&player, id).Error
	return &player, err
}

func (r *Repository) CreatePlayer(player *models.Player) error {
	return r.db.Create(player).Error
}

func (r *Repository) UpdatePlayer(player *models.Player) error {
	return r.db.Save(player).Error
}

// Conversation 对话记录相关操作
func (r *Repository) GetConversations(playerID, npcID uint, limit int) ([]models.Conversation, error) {
	var conversations []models.Conversation
	query := r.db.Order("created_at desc")
	if playerID > 0 {
		query = query.Where("player_id = ?", playerID)
	}
	if npcID > 0 {
		query = query.Where("npc_id = ?", npcID)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&conversations).Error
	return conversations, err
}

func (r *Repository) GetConversationsPaginated(playerID, npcID uint, offset, limit int) ([]models.Conversation, int64, error) {
	var conversations []models.Conversation
	var total int64
	countQuery := r.db.Model(&models.Conversation{})
	if playerID > 0 {
		countQuery = countQuery.Where("player_id = ?", playerID)
	}
	if npcID > 0 {
		countQuery = countQuery.Where("npc_id = ?", npcID)
	}
	countQuery.Count(&total)
	query := r.db.Order("created_at desc")
	if playerID > 0 {
		query = query.Where("player_id = ?", playerID)
	}
	if npcID > 0 {
		query = query.Where("npc_id = ?", npcID)
	}
	err := query.Offset(offset).Limit(limit).Find(&conversations).Error
	return conversations, total, err
}

func (r *Repository) CreateConversation(conv *models.Conversation) error {
	return r.db.Create(conv).Error
}

func (r *Repository) GetConversationSummary(playerID, npcID uint) (string, error) {
	var conv models.Conversation
	err := r.db.Where("player_id = ? AND npc_id = ? AND summary != ''", playerID, npcID).
		Order("created_at desc").First(&conv).Error
	if err != nil {
		return "", err
	}
	return conv.Summary, nil
}

// Player 玩家相关操作
func (r *Repository) GetPlayerByAccount(account string) (*models.Player, error) {
	var player models.Player
	err := r.db.Where("account = ?", account).First(&player).Error
	return &player, err
}

// SaveShopItem 保存商店商品
func (r *Repository) SaveShopItem(item *models.ShopItem) error {
	return r.db.Save(item).Error
}

// GetShopItemByID 获取商店商品
func (r *Repository) GetShopItemByID(id uint) (*models.ShopItem, error) {
	var item models.ShopItem
	err := r.db.First(&item, id).Error
	return &item, err
}

// Skill 技能相关操作
func (r *Repository) GetSkills() ([]models.Skill, error) {
	var skills []models.Skill
	err := r.db.Find(&skills).Error
	return skills, err
}

func (r *Repository) GetSkillsPaginated(offset, limit int) ([]models.Skill, int64, error) {
	var skills []models.Skill
	var total int64
	r.db.Model(&models.Skill{}).Count(&total)
	err := r.db.Offset(offset).Limit(limit).Find(&skills).Error
	return skills, total, err
}

func (r *Repository) GetSkillByID(id uint) (*models.Skill, error) {
	var skill models.Skill
	err := r.db.First(&skill, id).Error
	return &skill, err
}

func (r *Repository) GetSkillByCode(code string) (*models.Skill, error) {
	var skill models.Skill
	err := r.db.Where("code = ?", code).First(&skill).Error
	return &skill, err
}

func (r *Repository) CreateSkill(skill *models.Skill) error {
	return r.db.Create(skill).Error
}

func (r *Repository) UpdateSkill(skill *models.Skill) error {
	return r.db.Save(skill).Error
}

func (r *Repository) DeleteSkill(id uint) error {
	return r.db.Delete(&models.Skill{}, id).Error
}

// Achievement 成就相关操作
func (r *Repository) GetAchievements() ([]models.Achievement, error) {
	var achievements []models.Achievement
	err := r.db.Find(&achievements).Error
	return achievements, err
}

func (r *Repository) GetAchievementsPaginated(offset, limit int) ([]models.Achievement, int64, error) {
	var achievements []models.Achievement
	var total int64
	r.db.Model(&models.Achievement{}).Count(&total)
	err := r.db.Offset(offset).Limit(limit).Find(&achievements).Error
	return achievements, total, err
}

func (r *Repository) GetAchievementByID(id uint) (*models.Achievement, error) {
	var achievement models.Achievement
	err := r.db.First(&achievement, id).Error
	return &achievement, err
}

func (r *Repository) GetAchievementByCode(code string) (*models.Achievement, error) {
	var achievement models.Achievement
	err := r.db.Where("code = ?", code).First(&achievement).Error
	return &achievement, err
}

func (r *Repository) CreateAchievement(achievement *models.Achievement) error {
	return r.db.Create(achievement).Error
}

func (r *Repository) GetPlayerAchievements(playerID uint) ([]models.PlayerAchievement, error) {
	var achievements []models.PlayerAchievement
	err := r.db.Preload("Achievement").Where("player_id = ?", playerID).Find(&achievements).Error
	return achievements, err
}

func (r *Repository) CreatePlayerAchievement(pa *models.PlayerAchievement) error {
	return r.db.Create(pa).Error
}

func (r *Repository) HasAchievement(playerID, achievementID uint) bool {
	var count int64
	r.db.Model(&models.PlayerAchievement{}).Where("player_id = ? AND achievement_id = ?", playerID, achievementID).Count(&count)
	return count > 0
}

// SaveGame 存档相关操作
func (r *Repository) GetSaveGames(playerID uint) ([]models.SaveGame, error) {
	var saves []models.SaveGame
	err := r.db.Where("player_id = ?", playerID).Order("slot").Find(&saves).Error
	return saves, err
}

func (r *Repository) GetSaveGame(playerID uint, slot int) (*models.SaveGame, error) {
	var save models.SaveGame
	err := r.db.Where("player_id = ? AND slot = ?", playerID, slot).First(&save).Error
	return &save, err
}

func (r *Repository) CreateSaveGame(save *models.SaveGame) error {
	return r.db.Create(save).Error
}

func (r *Repository) UpdateSaveGame(save *models.SaveGame) error {
	return r.db.Save(save).Error
}

func (r *Repository) DeleteSaveGame(playerID uint, slot int) error {
	return r.db.Where("player_id = ? AND slot = ?", playerID, slot).Delete(&models.SaveGame{}).Error
}

// PlayerConversationContext 对话上下文相关操作
func (r *Repository) GetConversationContext(playerID, npcID uint) (*models.PlayerConversationContext, error) {
	var ctx models.PlayerConversationContext
	err := r.db.Where("player_id = ? AND npc_id = ?", playerID, npcID).First(&ctx).Error
	return &ctx, err
}

func (r *Repository) CreateConversationContext(ctx *models.PlayerConversationContext) error {
	return r.db.Create(ctx).Error
}

func (r *Repository) UpdateConversationContext(ctx *models.PlayerConversationContext) error {
	return r.db.Save(ctx).Error
}

func (r *Repository) UpsertConversationContext(ctx *models.PlayerConversationContext) error {
	var existing models.PlayerConversationContext
	result := r.db.Where("player_id = ? AND npc_id = ?", ctx.PlayerID, ctx.NPCID).First(&existing)
	if result.Error != nil {
		return r.db.Create(ctx).Error
	}
	existing.PlayerName = ctx.PlayerName
	existing.PlayerLevel = ctx.PlayerLevel
	existing.TalkCount = ctx.TalkCount
	existing.Summary = ctx.Summary
	existing.Extra = ctx.Extra
	return r.db.Save(&existing).Error
}

func (r *Repository) GetConversationsByPair(playerID, npcID uint, limit int) ([]models.Conversation, error) {
	var conversations []models.Conversation
	query := r.db.Where("player_id = ? AND npc_id = ?", playerID, npcID).Order("created_at desc")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&conversations).Error
	return conversations, err
}

func (r *Repository) DeleteConversationsByPair(playerID, npcID uint) error {
	return r.db.Where("player_id = ? AND npc_id = ?", playerID, npcID).Delete(&models.Conversation{}).Error
}

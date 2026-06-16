package repository

import (
	"context"

	"github.com/ddc-111/agentGame/server/internal/database/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Transaction(ctx context.Context, fn func(repo *Repository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(&Repository{db: tx})
	})
}

func (r *Repository) GetScenes(ctx context.Context) ([]models.Scene, error) {
	var scenes []models.Scene
	err := r.db.WithContext(ctx).Preload("SceneNPCs").Preload("Portals").Find(&scenes).Error
	return scenes, err
}

func (r *Repository) GetScenesPaginated(ctx context.Context, offset, limit int) ([]models.Scene, int64, error) {
	var scenes []models.Scene
	var total int64
	r.db.WithContext(ctx).Model(&models.Scene{}).Count(&total)
	err := r.db.WithContext(ctx).Preload("SceneNPCs").Preload("Portals").Offset(offset).Limit(limit).Find(&scenes).Error
	return scenes, total, err
}

func (r *Repository) GetSceneByID(ctx context.Context, id uint) (*models.Scene, error) {
	var scene models.Scene
	err := r.db.WithContext(ctx).Preload("SceneNPCs").Preload("Portals").First(&scene, id).Error
	return &scene, err
}

func (r *Repository) GetSceneByCode(ctx context.Context, code string) (*models.Scene, error) {
	var scene models.Scene
	err := r.db.WithContext(ctx).Preload("SceneNPCs").Preload("Portals").Where("code = ?", code).First(&scene).Error
	return &scene, err
}

func (r *Repository) CreateScene(ctx context.Context, scene *models.Scene) error {
	return r.db.WithContext(ctx).Create(scene).Error
}

func (r *Repository) UpdateScene(ctx context.Context, scene *models.Scene) error {
	return r.db.WithContext(ctx).Save(scene).Error
}

func (r *Repository) DeleteScene(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Scene{}, id).Error
}

func (r *Repository) GetNPCs(ctx context.Context) ([]models.NPC, error) {
	var npcs []models.NPC
	err := r.db.WithContext(ctx).Preload("Agent").Preload("Shop").Find(&npcs).Error
	return npcs, err
}

func (r *Repository) GetNPCsPaginated(ctx context.Context, offset, limit int) ([]models.NPC, int64, error) {
	var npcs []models.NPC
	var total int64
	r.db.WithContext(ctx).Model(&models.NPC{}).Count(&total)
	err := r.db.WithContext(ctx).Preload("Agent").Preload("Shop").Offset(offset).Limit(limit).Find(&npcs).Error
	return npcs, total, err
}

func (r *Repository) GetNPCByID(ctx context.Context, id uint) (*models.NPC, error) {
	var npc models.NPC
	err := r.db.WithContext(ctx).Preload("Agent").Preload("Shop").First(&npc, id).Error
	return &npc, err
}

func (r *Repository) GetNPCByCode(ctx context.Context, code string) (*models.NPC, error) {
	var npc models.NPC
	err := r.db.WithContext(ctx).Preload("Agent").Preload("Shop").Where("code = ?", code).First(&npc).Error
	return &npc, err
}

func (r *Repository) CreateNPC(ctx context.Context, npc *models.NPC) error {
	return r.db.WithContext(ctx).Create(npc).Error
}

func (r *Repository) UpdateNPC(ctx context.Context, npc *models.NPC) error {
	return r.db.WithContext(ctx).Save(npc).Error
}

func (r *Repository) DeleteNPC(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.NPC{}, id).Error
}

func (r *Repository) GetAgents(ctx context.Context) ([]models.Agent, error) {
	var agents []models.Agent
	err := r.db.WithContext(ctx).Find(&agents).Error
	return agents, err
}

func (r *Repository) GetAgentsPaginated(ctx context.Context, offset, limit int) ([]models.Agent, int64, error) {
	var agents []models.Agent
	var total int64
	r.db.WithContext(ctx).Model(&models.Agent{}).Count(&total)
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&agents).Error
	return agents, total, err
}

func (r *Repository) GetAgentByID(ctx context.Context, id uint) (*models.Agent, error) {
	var agent models.Agent
	err := r.db.WithContext(ctx).First(&agent, id).Error
	return &agent, err
}

func (r *Repository) GetAgentByCode(ctx context.Context, code string) (*models.Agent, error) {
	var agent models.Agent
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&agent).Error
	return &agent, err
}

func (r *Repository) CreateAgent(ctx context.Context, agent *models.Agent) error {
	return r.db.WithContext(ctx).Create(agent).Error
}

func (r *Repository) UpdateAgent(ctx context.Context, agent *models.Agent) error {
	return r.db.WithContext(ctx).Save(agent).Error
}

func (r *Repository) DeleteAgent(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Agent{}, id).Error
}

func (r *Repository) GetProviders(ctx context.Context) ([]models.LLMProvider, error) {
	var providers []models.LLMProvider
	err := r.db.WithContext(ctx).Find(&providers).Error
	return providers, err
}

func (r *Repository) GetProvidersPaginated(ctx context.Context, offset, limit int) ([]models.LLMProvider, int64, error) {
	var providers []models.LLMProvider
	var total int64
	r.db.WithContext(ctx).Model(&models.LLMProvider{}).Count(&total)
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&providers).Error
	return providers, total, err
}

func (r *Repository) GetProviderByID(ctx context.Context, id uint) (*models.LLMProvider, error) {
	var provider models.LLMProvider
	err := r.db.WithContext(ctx).First(&provider, id).Error
	return &provider, err
}

func (r *Repository) CreateProvider(ctx context.Context, provider *models.LLMProvider) error {
	return r.db.WithContext(ctx).Create(provider).Error
}

func (r *Repository) UpdateProvider(ctx context.Context, provider *models.LLMProvider) error {
	return r.db.WithContext(ctx).Save(provider).Error
}

func (r *Repository) DeleteProvider(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.LLMProvider{}, id).Error
}

func (r *Repository) GetTemplates(ctx context.Context) ([]models.PromptTemplate, error) {
	var templates []models.PromptTemplate
	err := r.db.WithContext(ctx).Find(&templates).Error
	return templates, err
}

func (r *Repository) GetTemplatesPaginated(ctx context.Context, offset, limit int) ([]models.PromptTemplate, int64, error) {
	var templates []models.PromptTemplate
	var total int64
	r.db.WithContext(ctx).Model(&models.PromptTemplate{}).Count(&total)
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&templates).Error
	return templates, total, err
}

func (r *Repository) GetTemplateByID(ctx context.Context, id uint) (*models.PromptTemplate, error) {
	var template models.PromptTemplate
	err := r.db.WithContext(ctx).First(&template, id).Error
	return &template, err
}

func (r *Repository) GetTemplateByCode(ctx context.Context, code string) (*models.PromptTemplate, error) {
	var template models.PromptTemplate
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&template).Error
	return &template, err
}

func (r *Repository) CreateTemplate(ctx context.Context, template *models.PromptTemplate) error {
	return r.db.WithContext(ctx).Create(template).Error
}

func (r *Repository) UpdateTemplate(ctx context.Context, template *models.PromptTemplate) error {
	return r.db.WithContext(ctx).Save(template).Error
}

func (r *Repository) DeleteTemplate(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.PromptTemplate{}, id).Error
}

func (r *Repository) GetShops(ctx context.Context) ([]models.Shop, error) {
	var shops []models.Shop
	err := r.db.WithContext(ctx).Preload("Items.Item").Find(&shops).Error
	return shops, err
}

func (r *Repository) GetShopsPaginated(ctx context.Context, offset, limit int) ([]models.Shop, int64, error) {
	var shops []models.Shop
	var total int64
	r.db.WithContext(ctx).Model(&models.Shop{}).Count(&total)
	err := r.db.WithContext(ctx).Preload("Items.Item").Offset(offset).Limit(limit).Find(&shops).Error
	return shops, total, err
}

func (r *Repository) GetShopByID(ctx context.Context, id uint) (*models.Shop, error) {
	var shop models.Shop
	err := r.db.WithContext(ctx).Preload("Items.Item").First(&shop, id).Error
	return &shop, err
}

func (r *Repository) GetShopByCode(ctx context.Context, code string) (*models.Shop, error) {
	var shop models.Shop
	err := r.db.WithContext(ctx).Preload("Items.Item").Where("code = ?", code).First(&shop).Error
	return &shop, err
}

func (r *Repository) CreateShop(ctx context.Context, shop *models.Shop) error {
	return r.db.WithContext(ctx).Create(shop).Error
}

func (r *Repository) UpdateShop(ctx context.Context, shop *models.Shop) error {
	return r.db.WithContext(ctx).Save(shop).Error
}

func (r *Repository) DeleteShop(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Shop{}, id).Error
}

func (r *Repository) GetItems(ctx context.Context) ([]models.Item, error) {
	var items []models.Item
	err := r.db.WithContext(ctx).Find(&items).Error
	return items, err
}

func (r *Repository) GetItemsPaginated(ctx context.Context, offset, limit int) ([]models.Item, int64, error) {
	var items []models.Item
	var total int64
	r.db.WithContext(ctx).Model(&models.Item{}).Count(&total)
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&items).Error
	return items, total, err
}

func (r *Repository) GetItemByID(ctx context.Context, id uint) (*models.Item, error) {
	var item models.Item
	err := r.db.WithContext(ctx).First(&item, id).Error
	return &item, err
}

func (r *Repository) GetItemByCode(ctx context.Context, code string) (*models.Item, error) {
	var item models.Item
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&item).Error
	return &item, err
}

func (r *Repository) CreateItem(ctx context.Context, item *models.Item) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *Repository) UpdateItem(ctx context.Context, item *models.Item) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *Repository) DeleteItem(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Item{}, id).Error
}

func (r *Repository) GetTasks(ctx context.Context) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.WithContext(ctx).Find(&tasks).Error
	return tasks, err
}

func (r *Repository) GetTasksPaginated(ctx context.Context, offset, limit int) ([]models.Task, int64, error) {
	var tasks []models.Task
	var total int64
	r.db.WithContext(ctx).Model(&models.Task{}).Count(&total)
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&tasks).Error
	return tasks, total, err
}

func (r *Repository) GetTaskByID(ctx context.Context, id uint) (*models.Task, error) {
	var task models.Task
	err := r.db.WithContext(ctx).First(&task, id).Error
	return &task, err
}

func (r *Repository) GetTaskByCode(ctx context.Context, code string) (*models.Task, error) {
	var task models.Task
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&task).Error
	return &task, err
}

func (r *Repository) CreateTask(ctx context.Context, task *models.Task) error {
	return r.db.WithContext(ctx).Create(task).Error
}

func (r *Repository) UpdateTask(ctx context.Context, task *models.Task) error {
	return r.db.WithContext(ctx).Save(task).Error
}

func (r *Repository) DeleteTask(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Task{}, id).Error
}

func (r *Repository) GetFlows(ctx context.Context) ([]models.Flow, error) {
	var flows []models.Flow
	err := r.db.WithContext(ctx).Find(&flows).Error
	return flows, err
}

func (r *Repository) GetFlowsPaginated(ctx context.Context, offset, limit int) ([]models.Flow, int64, error) {
	var flows []models.Flow
	var total int64
	r.db.WithContext(ctx).Model(&models.Flow{}).Count(&total)
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&flows).Error
	return flows, total, err
}

func (r *Repository) GetFlowByID(ctx context.Context, id uint) (*models.Flow, error) {
	var flow models.Flow
	err := r.db.WithContext(ctx).First(&flow, id).Error
	return &flow, err
}

func (r *Repository) GetFlowByCode(ctx context.Context, code string) (*models.Flow, error) {
	var flow models.Flow
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&flow).Error
	return &flow, err
}

func (r *Repository) CreateFlow(ctx context.Context, flow *models.Flow) error {
	return r.db.WithContext(ctx).Create(flow).Error
}

func (r *Repository) UpdateFlow(ctx context.Context, flow *models.Flow) error {
	return r.db.WithContext(ctx).Save(flow).Error
}

func (r *Repository) DeleteFlow(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Flow{}, id).Error
}

func (r *Repository) GetConfig(ctx context.Context, key string) (*models.GameConfig, error) {
	var config models.GameConfig
	err := r.db.WithContext(ctx).Where("key = ?", key).First(&config).Error
	return &config, err
}

func (r *Repository) SetConfig(ctx context.Context, key, value string) error {
	var config models.GameConfig
	result := r.db.WithContext(ctx).Where("key = ?", key).First(&config)
	if result.Error != nil {
		config = models.GameConfig{Key: key, Value: value}
		return r.db.WithContext(ctx).Create(&config).Error
	}
	config.Value = value
	return r.db.WithContext(ctx).Save(&config).Error
}

func (r *Repository) GetPlayers(ctx context.Context) ([]models.Player, error) {
	var players []models.Player
	err := r.db.WithContext(ctx).Find(&players).Error
	return players, err
}

func (r *Repository) GetPlayersPaginated(ctx context.Context, offset, limit int) ([]models.Player, int64, error) {
	var players []models.Player
	var total int64
	r.db.WithContext(ctx).Model(&models.Player{}).Count(&total)
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&players).Error
	return players, total, err
}

func (r *Repository) GetPlayerByID(ctx context.Context, id uint) (*models.Player, error) {
	var player models.Player
	err := r.db.WithContext(ctx).First(&player, id).Error
	return &player, err
}

func (r *Repository) CreatePlayer(ctx context.Context, player *models.Player) error {
	return r.db.WithContext(ctx).Create(player).Error
}

func (r *Repository) UpdatePlayer(ctx context.Context, player *models.Player) error {
	return r.db.WithContext(ctx).Save(player).Error
}

func (r *Repository) GetConversations(ctx context.Context, playerID, npcID uint, limit int) ([]models.Conversation, error) {
	var conversations []models.Conversation
	query := r.db.WithContext(ctx).Order("created_at desc")
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

func (r *Repository) GetConversationsPaginated(ctx context.Context, playerID, npcID uint, offset, limit int) ([]models.Conversation, int64, error) {
	var conversations []models.Conversation
	var total int64
	countQuery := r.db.WithContext(ctx).Model(&models.Conversation{})
	if playerID > 0 {
		countQuery = countQuery.Where("player_id = ?", playerID)
	}
	if npcID > 0 {
		countQuery = countQuery.Where("npc_id = ?", npcID)
	}
	countQuery.Count(&total)
	query := r.db.WithContext(ctx).Order("created_at desc")
	if playerID > 0 {
		query = query.Where("player_id = ?", playerID)
	}
	if npcID > 0 {
		query = query.Where("npc_id = ?", npcID)
	}
	err := query.Offset(offset).Limit(limit).Find(&conversations).Error
	return conversations, total, err
}

func (r *Repository) CreateConversation(ctx context.Context, conv *models.Conversation) error {
	return r.db.WithContext(ctx).Create(conv).Error
}

func (r *Repository) GetConversationSummary(ctx context.Context, playerID, npcID uint) (string, error) {
	var conv models.Conversation
	err := r.db.WithContext(ctx).Where("player_id = ? AND npc_id = ? AND summary != ''", playerID, npcID).
		Order("created_at desc").First(&conv).Error
	if err != nil {
		return "", err
	}
	return conv.Summary, nil
}

func (r *Repository) GetPlayerByAccount(ctx context.Context, account string) (*models.Player, error) {
	var player models.Player
	err := r.db.WithContext(ctx).Where("account = ?", account).First(&player).Error
	return &player, err
}

func (r *Repository) SaveShopItem(ctx context.Context, item *models.ShopItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *Repository) GetShopItemByID(ctx context.Context, id uint) (*models.ShopItem, error) {
	var item models.ShopItem
	err := r.db.WithContext(ctx).First(&item, id).Error
	return &item, err
}

func (r *Repository) GetSkills(ctx context.Context) ([]models.Skill, error) {
	var skills []models.Skill
	err := r.db.WithContext(ctx).Find(&skills).Error
	return skills, err
}

func (r *Repository) GetSkillsPaginated(ctx context.Context, offset, limit int) ([]models.Skill, int64, error) {
	var skills []models.Skill
	var total int64
	r.db.WithContext(ctx).Model(&models.Skill{}).Count(&total)
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&skills).Error
	return skills, total, err
}

func (r *Repository) GetSkillByID(ctx context.Context, id uint) (*models.Skill, error) {
	var skill models.Skill
	err := r.db.WithContext(ctx).First(&skill, id).Error
	return &skill, err
}

func (r *Repository) GetSkillByCode(ctx context.Context, code string) (*models.Skill, error) {
	var skill models.Skill
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&skill).Error
	return &skill, err
}

func (r *Repository) CreateSkill(ctx context.Context, skill *models.Skill) error {
	return r.db.WithContext(ctx).Create(skill).Error
}

func (r *Repository) UpdateSkill(ctx context.Context, skill *models.Skill) error {
	return r.db.WithContext(ctx).Save(skill).Error
}

func (r *Repository) DeleteSkill(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Skill{}, id).Error
}

func (r *Repository) GetAchievements(ctx context.Context) ([]models.Achievement, error) {
	var achievements []models.Achievement
	err := r.db.WithContext(ctx).Find(&achievements).Error
	return achievements, err
}

func (r *Repository) GetAchievementsPaginated(ctx context.Context, offset, limit int) ([]models.Achievement, int64, error) {
	var achievements []models.Achievement
	var total int64
	r.db.WithContext(ctx).Model(&models.Achievement{}).Count(&total)
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&achievements).Error
	return achievements, total, err
}

func (r *Repository) GetAchievementByID(ctx context.Context, id uint) (*models.Achievement, error) {
	var achievement models.Achievement
	err := r.db.WithContext(ctx).First(&achievement, id).Error
	return &achievement, err
}

func (r *Repository) GetAchievementByCode(ctx context.Context, code string) (*models.Achievement, error) {
	var achievement models.Achievement
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&achievement).Error
	return &achievement, err
}

func (r *Repository) CreateAchievement(ctx context.Context, achievement *models.Achievement) error {
	return r.db.WithContext(ctx).Create(achievement).Error
}

func (r *Repository) GetPlayerAchievements(ctx context.Context, playerID uint) ([]models.PlayerAchievement, error) {
	var achievements []models.PlayerAchievement
	err := r.db.WithContext(ctx).Preload("Achievement").Where("player_id = ?", playerID).Find(&achievements).Error
	return achievements, err
}

func (r *Repository) CreatePlayerAchievement(ctx context.Context, pa *models.PlayerAchievement) error {
	return r.db.WithContext(ctx).Create(pa).Error
}

func (r *Repository) HasAchievement(ctx context.Context, playerID, achievementID uint) bool {
	var count int64
	r.db.WithContext(ctx).Model(&models.PlayerAchievement{}).Where("player_id = ? AND achievement_id = ?", playerID, achievementID).Count(&count)
	return count > 0
}

func (r *Repository) GetSaveGames(ctx context.Context, playerID uint) ([]models.SaveGame, error) {
	var saves []models.SaveGame
	err := r.db.WithContext(ctx).Where("player_id = ?", playerID).Order("slot").Find(&saves).Error
	return saves, err
}

func (r *Repository) GetSaveGame(ctx context.Context, playerID uint, slot int) (*models.SaveGame, error) {
	var save models.SaveGame
	err := r.db.WithContext(ctx).Where("player_id = ? AND slot = ?", playerID, slot).First(&save).Error
	return &save, err
}

func (r *Repository) CreateSaveGame(ctx context.Context, save *models.SaveGame) error {
	return r.db.WithContext(ctx).Create(save).Error
}

func (r *Repository) UpdateSaveGame(ctx context.Context, save *models.SaveGame) error {
	return r.db.WithContext(ctx).Save(save).Error
}

func (r *Repository) DeleteSaveGame(ctx context.Context, playerID uint, slot int) error {
	return r.db.WithContext(ctx).Where("player_id = ? AND slot = ?", playerID, slot).Delete(&models.SaveGame{}).Error
}

func (r *Repository) GetConversationContext(ctx context.Context, playerID, npcID uint) (*models.PlayerConversationContext, error) {
	var pctx models.PlayerConversationContext
	err := r.db.WithContext(ctx).Where("player_id = ? AND npc_id = ?", playerID, npcID).First(&pctx).Error
	return &pctx, err
}

func (r *Repository) CreateConversationContext(ctx context.Context, pctx *models.PlayerConversationContext) error {
	return r.db.WithContext(ctx).Create(pctx).Error
}

func (r *Repository) UpdateConversationContext(ctx context.Context, pctx *models.PlayerConversationContext) error {
	return r.db.WithContext(ctx).Save(pctx).Error
}

func (r *Repository) UpsertConversationContext(ctx context.Context, pctx *models.PlayerConversationContext) error {
	var existing models.PlayerConversationContext
	result := r.db.WithContext(ctx).Where("player_id = ? AND npc_id = ?", pctx.PlayerID, pctx.NPCID).First(&existing)
	if result.Error != nil {
		return r.db.WithContext(ctx).Create(pctx).Error
	}
	existing.PlayerName = pctx.PlayerName
	existing.PlayerLevel = pctx.PlayerLevel
	existing.TalkCount = pctx.TalkCount
	existing.Summary = pctx.Summary
	existing.Extra = pctx.Extra
	return r.db.WithContext(ctx).Save(&existing).Error
}

func (r *Repository) GetConversationsByPair(ctx context.Context, playerID, npcID uint, limit int) ([]models.Conversation, error) {
	var conversations []models.Conversation
	query := r.db.WithContext(ctx).Where("player_id = ? AND npc_id = ?", playerID, npcID).Order("created_at desc")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&conversations).Error
	return conversations, err
}

func (r *Repository) DeleteConversationsByPair(ctx context.Context, playerID, npcID uint) error {
	return r.db.WithContext(ctx).Where("player_id = ? AND npc_id = ?", playerID, npcID).Delete(&models.Conversation{}).Error
}

func (r *Repository) GetScenesByNPCID(ctx context.Context, npcID uint) ([]models.Scene, error) {
	var scenes []models.Scene
	err := r.db.WithContext(ctx).Joins("JOIN scene_npcs ON scene_npcs.scene_id = scenes.id").
		Where("scene_npcs.npc_id = ?", npcID).
		Find(&scenes).Error
	return scenes, err
}

package migrations

import (
	"github.com/ddc-111/agentGame/server/internal/database/models"
	"gorm.io/gorm"
)

var _001 = models.Migration{
	Version: "001",
	Name:    "initial_schema",
	Up: func(db *gorm.DB) error {
		return db.AutoMigrate(
			&models.Scene{},
			&models.SceneNPC{},
			&models.Portal{},
			&models.NPC{},
			&models.Agent{},
			&models.LLMProvider{},
			&models.PromptTemplate{},
			&models.Shop{},
			&models.ShopItem{},
			&models.Item{},
			&models.Task{},
			&models.Flow{},
			&models.GameConfig{},
			&models.Player{},
			&models.Conversation{},
			&models.SaveGame{},
			&models.Skill{},
			&models.Achievement{},
			&models.PlayerAchievement{},
			&models.PlayerConversationContext{},
			&models.GMUser{},
		)
	},
	Down: func(db *gorm.DB) error {
		return db.Migrator().DropTable(
			&models.GMUser{},
			&models.PlayerConversationContext{},
			&models.PlayerAchievement{},
			&models.Achievement{},
			&models.Skill{},
			&models.SaveGame{},
			&models.Conversation{},
			&models.Player{},
			&models.GameConfig{},
			&models.Flow{},
			&models.Task{},
			&models.Item{},
			&models.ShopItem{},
			&models.Shop{},
			&models.PromptTemplate{},
			&models.LLMProvider{},
			&models.Agent{},
			&models.NPC{},
			&models.Portal{},
			&models.SceneNPC{},
			&models.Scene{},
		)
	},
}

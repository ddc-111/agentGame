package models

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// Scene 场景
type Scene struct {
	BaseModel
	Name        string    `gorm:"size:100;not null" json:"name"`
	Code        string    `gorm:"size:50;uniqueIndex" json:"code"`
	Description string    `gorm:"size:500" json:"description"`
	Background  string    `gorm:"size:255" json:"background"`
	Width       int       `gorm:"default:1920" json:"width"`
	Height      int       `gorm:"default:1080" json:"height"`
	SceneNPCs   []SceneNPC   `gorm:"foreignKey:SceneID" json:"scene_npcs,omitempty"`
	Portals     []Portal     `gorm:"foreignKey:SceneID" json:"portals,omitempty"`
}

// SceneNPC 场景NPC关联
type SceneNPC struct {
	BaseModel
	SceneID uint   `gorm:"index" json:"scene_id"`
	NPCID   uint   `gorm:"index" json:"npc_id"`
	X       int    `gorm:"default:0" json:"x"`
	Y       int    `gorm:"default:0" json:"y"`
	NPC     *NPC   `gorm:"foreignKey:NPCID" json:"npc,omitempty"`
	Scene   *Scene `gorm:"foreignKey:SceneID" json:"scene,omitempty"`
}

// Portal 传送点
type Portal struct {
	BaseModel
	SceneID      uint   `gorm:"index" json:"scene_id"`
	X            int    `json:"x"`
	Y            int    `json:"y"`
	TargetScene  string `gorm:"size:50" json:"target_scene"`
	TargetX      int    `json:"target_x"`
	TargetY      int    `json:"target_y"`
}

// NPC 非玩家角色
type NPC struct {
	BaseModel
	Name        string  `gorm:"size:100;not null" json:"name"`
	Code        string  `gorm:"size:50;uniqueIndex" json:"code"`
	Title       string  `gorm:"size:100" json:"title"`
	Description string  `gorm:"size:500" json:"description"`
	Avatar      string  `gorm:"size:255" json:"avatar"`
	Sprite      string  `gorm:"size:255" json:"sprite"`
	AgentID     *uint   `gorm:"index" json:"agent_id,omitempty"`
	ShopID      *uint   `gorm:"index" json:"shop_id,omitempty"`
	Behaviors   string  `gorm:"size:500" json:"behaviors"` // JSON数组
	Schedule    string  `gorm:"size:2000" json:"schedule"` // JSON数组
	Agent       *Agent  `gorm:"foreignKey:AgentID" json:"agent,omitempty"`
	Shop        *Shop   `gorm:"foreignKey:ShopID" json:"shop,omitempty"`
}

// Agent 智能体
type Agent struct {
	BaseModel
	Name           string `gorm:"size:100;not null" json:"name"`
	Code           string `gorm:"size:50;uniqueIndex" json:"code"`
	Description    string `gorm:"size:500" json:"description"`
	LLMProvider    string `gorm:"size:50" json:"llm_provider"`
	LLMModel       string `gorm:"size:100" json:"llm_model"`
	Temperature    float64 `gorm:"default:0.7" json:"temperature"`
	MaxTokens      int    `gorm:"default:500" json:"max_tokens"`
	SystemPrompt   string `gorm:"type:text" json:"system_prompt"`
	MemoryType     string `gorm:"size:50;default:sliding_window" json:"memory_type"`
	MaxMessages    int    `gorm:"default:20" json:"max_messages"`
	SummaryEnabled bool   `gorm:"default:true" json:"summary_enabled"`
	KnowledgeBase  string `gorm:"type:text" json:"knowledge_base"` // JSON数组
	Tools          string `gorm:"type:text" json:"tools"`          // JSON数组
}

// LLMProvider 大模型提供商
type LLMProvider struct {
	BaseModel
	Name    string `gorm:"size:100;not null" json:"name"`
	Code    string `gorm:"size:50;uniqueIndex" json:"code"`
	BaseURL string `gorm:"size:255" json:"base_url"`
	APIKey  string `gorm:"size:255" json:"api_key"`
	Models  string `gorm:"type:text" json:"models"` // JSON数组
}

// PromptTemplate 提示词模板
type PromptTemplate struct {
	BaseModel
	Name       string `gorm:"size:100;not null" json:"name"`
	Code       string `gorm:"size:50;uniqueIndex" json:"code"`
	Category   string `gorm:"size:50" json:"category"`
	Content    string `gorm:"type:text;not null" json:"content"`
	Variables  string `gorm:"type:text" json:"variables"` // JSON数组
}

// Shop 商店
type Shop struct {
	BaseModel
	Name        string `gorm:"size:100;not null" json:"name"`
	Code        string `gorm:"size:50;uniqueIndex" json:"code"`
	Type        string `gorm:"size:50" json:"type"`
	Description string `gorm:"size:500" json:"description"`
	OwnerNPC    string `gorm:"size:50" json:"owner_npc"`
	SceneCode   string `gorm:"size:50" json:"scene_code"`
	OpenTime    string `gorm:"size:10" json:"open_time"`
	CloseTime   string `gorm:"size:10" json:"close_time"`
	Discount    string `gorm:"size:200" json:"discount"` // JSON对象
	Items       []ShopItem `gorm:"foreignKey:ShopID" json:"items,omitempty"`
}

// ShopItem 商店商品
type ShopItem struct {
	BaseModel
	ShopID  uint   `gorm:"index" json:"shop_id"`
	ItemID  uint   `gorm:"index" json:"item_id"`
	Price   int    `gorm:"default:0" json:"price"`
	Stock   int    `gorm:"default:0" json:"stock"`
	Shop    *Shop  `gorm:"foreignKey:ShopID" json:"shop,omitempty"`
	Item    *Item  `gorm:"foreignKey:ItemID" json:"item,omitempty"`
}

// Item 道具
type Item struct {
	BaseModel
	Name        string `gorm:"size:100;not null" json:"name"`
	Code        string `gorm:"size:50;uniqueIndex" json:"code"`
	Category    string `gorm:"size:50" json:"category"`
	Description string `gorm:"size:500" json:"description"`
	Icon        string `gorm:"size:255" json:"icon"`
	Effect      string `gorm:"size:500" json:"effect"` // JSON对象
}

// Task 任务
type Task struct {
	BaseModel
	Name        string `gorm:"size:100;not null" json:"name"`
	Code        string `gorm:"size:50;uniqueIndex" json:"code"`
	Type        string `gorm:"size:50" json:"type"`
	Description string `gorm:"size:500" json:"description"`
	Status      string `gorm:"size:50;default:active" json:"status"`
	Trigger     string `gorm:"type:text" json:"trigger"`    // JSON对象
	Objectives  string `gorm:"type:text" json:"objectives"` // JSON数组
	Rewards     string `gorm:"type:text" json:"rewards"`    // JSON对象
	NextTask    string `gorm:"size:50" json:"next_task"`
	Dialogue    string `gorm:"size:50" json:"dialogue"`
}

// Flow 流程
type Flow struct {
	BaseModel
	Name        string `gorm:"size:100;not null" json:"name"`
	Code        string `gorm:"size:50;uniqueIndex" json:"code"`
	Description string `gorm:"size:500" json:"description"`
	Nodes       string `gorm:"type:text" json:"nodes"` // JSON数组
	Edges       string `gorm:"type:text" json:"edges"` // JSON数组
}

// GameConfig 游戏配置
type GameConfig struct {
	BaseModel
	Key   string `gorm:"size:100;uniqueIndex" json:"key"`
	Value string `gorm:"type:text" json:"value"`
}

// Player 玩家
type Player struct {
	BaseModel
	Name         string `gorm:"size:100;not null" json:"name"`
	Account      string `gorm:"size:100;uniqueIndex" json:"account"`
	Level        int    `gorm:"default:1" json:"level"`
	Exp          int    `gorm:"default:0" json:"exp"`
	Gold         int    `gorm:"default:1000" json:"gold"`
	HP           int    `gorm:"default:100" json:"hp"`
	MP           int    `gorm:"default:50" json:"mp"`
	Attack       int    `gorm:"default:10" json:"attack"`
	Defense      int    `gorm:"default:5" json:"defense"`
	SceneID      string `gorm:"size:50" json:"scene_id"`
	PosX         int    `gorm:"default:0" json:"pos_x"`
	PosY         int    `gorm:"default:0" json:"pos_y"`
	Items        string `gorm:"type:text" json:"items"`        // JSON对象 {item_id: count}
	Equipment    string `gorm:"type:text" json:"equipment"`    // JSON对象 {weapon_id, armor_id}
	CombatWins   int    `gorm:"default:0" json:"combat_wins"`
	SkillsUsed   int    `gorm:"default:0" json:"skills_used"`
	VisitedScenes string `gorm:"type:text" json:"visited_scenes"` // JSON数组 ["scene_code", ...]
}

// Conversation 对话记录
type Conversation struct {
	BaseModel
	PlayerID uint   `gorm:"index" json:"player_id"`
	NPCID    uint   `gorm:"index" json:"npc_id"`
	AgentID  uint   `gorm:"index" json:"agent_id"`
	Role     string `gorm:"size:20" json:"role"`
	Content  string `gorm:"type:text" json:"content"`
	Summary  string `gorm:"type:text" json:"summary"`
}

// SaveGame 存档
type SaveGame struct {
	BaseModel
	PlayerID uint   `gorm:"index" json:"player_id"`
	Slot     int    `gorm:"index" json:"slot"`
	Name     string `gorm:"size:100" json:"name"`
	Snapshot string `gorm:"type:text" json:"snapshot"` // JSON快照数据
}

// Skill 技能
type Skill struct {
	BaseModel
	Name        string `gorm:"size:100;not null" json:"name"`
	Code        string `gorm:"size:50;uniqueIndex" json:"code"`
	Description string `gorm:"size:500" json:"description"`
	Type        string `gorm:"size:20" json:"type"`        // attack/heal/buff/debuff
	MPCost      int    `gorm:"default:0" json:"mp_cost"`
	Damage      int    `gorm:"default:0" json:"damage"`
	Heal        int    `gorm:"default:0" json:"heal"`
	Cooldown    int    `gorm:"default:0" json:"cooldown"`
	Level       int    `gorm:"default:1" json:"level"`
	Effect      string `gorm:"size:500" json:"effect"` // JSON effect
}

// Achievement 成就
type Achievement struct {
	BaseModel
	Name        string `gorm:"size:100;not null" json:"name"`
	Code        string `gorm:"size:50;uniqueIndex" json:"code"`
	Description string `gorm:"size:500" json:"description"`
	Condition   string `gorm:"type:text" json:"condition"` // JSON condition
	Reward      string `gorm:"type:text" json:"reward"`    // JSON reward
	Icon        string `gorm:"size:20" json:"icon"`
}

// PlayerAchievement 玩家成就记录
type PlayerAchievement struct {
	BaseModel
	PlayerID      uint   `gorm:"index" json:"player_id"`
	AchievementID uint   `gorm:"index" json:"achievement_id"`
	Achievement   *Achievement `gorm:"foreignKey:AchievementID" json:"achievement,omitempty"`
}

// PlayerConversationContext 玩家NPC对话上下文
type PlayerConversationContext struct {
	BaseModel
	PlayerID    uint   `gorm:"index" json:"player_id"`
	NPCID       uint   `gorm:"index" json:"npc_id"`
	PlayerName  string `gorm:"size:100" json:"player_name"`
	PlayerLevel int    `gorm:"default:1" json:"player_level"`
	TalkCount   int    `gorm:"default:0" json:"talk_count"`
	Summary     string `gorm:"type:text" json:"summary"`
	Extra       string `gorm:"type:text" json:"extra"`
}

// GMUser GM管理员用户
type GMUser struct {
	BaseModel
	Username string `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Password string `gorm:"size:255;not null" json:"-"`
	Role     string `gorm:"size:20;default:gm" json:"role"`
}

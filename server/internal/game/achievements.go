package game

import (
	"encoding/json"
	"fmt"
)

// Achievement tracks player accomplishments
type Achievement struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Condition   string `json:"condition"` // JSON condition
	Reward      string `json:"reward"`    // JSON reward
	Icon        string `json:"icon"`
}

// AchievementCondition defines when an achievement is unlocked
type AchievementCondition struct {
	Type  string `json:"type"`  // combat_win, gold, level, quest, explore, collect
	Value int    `json:"value"` // target value
	Extra string `json:"extra"` // additional filter (e.g., quest code)
}

// AchievementReward defines what the player gets
type AchievementReward struct {
	Exp  int `json:"exp"`
	Gold int `json:"gold"`
}

// PlayerAchievement tracks which achievements a player has unlocked
type PlayerAchievement struct {
	AchievementID uint   `json:"achievement_id"`
	UnlockedAt    string `json:"unlocked_at"`
}

// AchievementManager handles achievement checking and awarding
type AchievementManager struct{}

// NewAchievementManager creates a new achievement manager
func NewAchievementManager() *AchievementManager {
	return &AchievementManager{}
}

// CheckAchievements checks if a player has unlocked new achievements
// Returns list of newly unlocked achievement IDs
func (am *AchievementManager) CheckAchievements(achievements []*Achievement, playerData *PlayerAchievementData, unlocked map[uint]bool) []*Achievement {
	var newlyUnlocked []*Achievement

	for _, ach := range achievements {
		if unlocked[ach.ID] {
			continue // already unlocked
		}

		if am.checkCondition(ach.Condition, playerData) {
			newlyUnlocked = append(newlyUnlocked, ach)
		}
	}

	return newlyUnlocked
}

// checkCondition evaluates a single achievement condition
func (am *AchievementManager) checkCondition(conditionJSON string, data *PlayerAchievementData) bool {
	if conditionJSON == "" {
		return false
	}

	var cond AchievementCondition
	if err := json.Unmarshal([]byte(conditionJSON), &cond); err != nil {
		return false
	}

	switch cond.Type {
	case "combat_win":
		return data.CombatWins >= cond.Value
	case "gold":
		return data.TotalGold >= cond.Value
	case "level":
		return data.Level >= cond.Value
	case "quest_complete":
		if cond.Extra != "" {
			return data.CompletedQuests[cond.Extra]
		}
		return data.QuestCount >= cond.Value
	case "explore":
		return data.VisitedScenes >= cond.Value
	case "collect":
		return data.UniqueItems >= cond.Value
	case "talk_all_npcs":
		return data.TalkedToAllNPCs
	case "skill_use":
		return data.SkillsUsed >= cond.Value
	}

	return false
}

// GetReward parses the reward JSON
func (am *AchievementManager) GetReward(rewardJSON string) *AchievementReward {
	if rewardJSON == "" || rewardJSON == "{}" {
		return &AchievementReward{}
	}
	var reward AchievementReward
	if err := json.Unmarshal([]byte(rewardJSON), &reward); err != nil {
		return &AchievementReward{}
	}
	return &reward
}

// PlayerAchievementData holds player stats for achievement checking
type PlayerAchievementData struct {
	Level            int              `json:"level"`
	TotalGold        int              `json:"total_gold"`
	CombatWins       int              `json:"combat_win"`
	QuestCount       int              `json:"quest_count"`
	CompletedQuests  map[string]bool  `json:"completed_quests"`
	VisitedScenes    int              `json:"visited_scenes"`
	UniqueItems      int              `json:"unique_items"`
	TalkedToAllNPCs  bool             `json:"talked_to_all_npcs"`
	SkillsUsed       int              `json:"skills_used"`
}

// PredefinedAchievements returns the built-in achievements
func PredefinedAchievements() []*Achievement {
	return []*Achievement{
		{
			Name:        "初来乍到",
			Code:        "ach_first_quest",
			Description: "完成第一个任务",
			Condition:   `{"type":"quest_complete","value":1}`,
			Reward:      `{"exp":50,"gold":100}`,
			Icon:        "⭐",
		},
		{
			Name:        "村庄之友",
			Code:        "ach_talk_all_npcs",
			Description: "与所有NPC对话过",
			Condition:   `{"type":"talk_all_npcs","value":1}`,
			Reward:      `{"exp":200,"gold":300}`,
			Icon:        "👥",
		},
		{
			Name:        "富甲一方",
			Code:        "ach_rich",
			Description: "累计获得10000金币",
			Condition:   `{"type":"gold","value":10000}`,
			Reward:      `{"exp":500,"gold":1000}`,
			Icon:        "💰",
		},
		{
			Name:        "百战百胜",
			Code:        "ach_combat_100",
			Description: "赢得100场战斗",
			Condition:   `{"type":"combat_win","value":100}`,
			Reward:      `{"exp":1000,"gold":2000}`,
			Icon:        "⚔️",
		},
		{
			Name:        "探索者",
			Code:        "ach_explorer",
			Description: "探索所有场景",
			Condition:   `{"type":"explore","value":6}`,
			Reward:      `{"exp":300,"gold":500}`,
			Icon:        "🗺️",
		},
		{
			Name:        "收藏家",
			Code:        "ach_collector",
			Description: "拥有50种不同的道具",
			Condition:   `{"type":"collect","value":50}`,
			Reward:      `{"exp":800,"gold":1500}`,
			Icon:        "🎒",
		},
		{
			Name:        "初试牛刀",
			Code:        "ach_first_combat",
			Description: "赢得第一场战斗",
			Condition:   `{"type":"combat_win","value":1}`,
			Reward:      `{"exp":30,"gold":50}`,
			Icon:        "🗡️",
		},
		{
			Name:        "技能大师",
			Code:        "ach_skill_master",
			Description: "使用100次技能",
			Condition:   `{"type":"skill_use","value":100}`,
			Reward:      `{"exp":600,"gold":800}`,
			Icon:        "✨",
		},
		{
			Name:        "等级10",
			Code:        "ach_level_10",
			Description: "达到10级",
			Condition:   `{"type":"level","value":10}`,
			Reward:      `{"exp":200,"gold":500}`,
			Icon:        "🏆",
		},
		{
			Name:        "等级20",
			Code:        "ach_level_20",
			Description: "达到20级",
			Condition:   `{"type":"level","value":20}`,
			Reward:      `{"exp":500,"gold":1000}`,
			Icon:        "👑",
		},
	}
}

// SerializeAchievements serializes achievements to JSON
func SerializeAchievements(achievements []*Achievement) (string, error) {
	data, err := json.Marshal(achievements)
	if err != nil {
		return "", fmt.Errorf("序列化成就失败: %v", err)
	}
	return string(data), nil
}

// DeserializeAchievements deserializes achievements from JSON
func DeserializeAchievements(data string) ([]*Achievement, error) {
	var achievements []*Achievement
	err := json.Unmarshal([]byte(data), &achievements)
	if err != nil {
		return nil, fmt.Errorf("反序列化成就失败: %v", err)
	}
	return achievements, nil
}

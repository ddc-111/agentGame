package game

import (
	"encoding/json"
	"fmt"
	"time"
)

// SaveGameManager 存档管理器
type SaveGameManager struct{}

// SaveGame 存档数据
type SaveGame struct {
	ID        uint       `json:"id"`
	PlayerID  uint       `json:"player_id"`
	Slot      int        `json:"slot"`
	Name      string     `json:"name"`
	Snapshot  *Snapshot  `json:"snapshot"`
	CreatedAt time.Time  `json:"created_at"`
}

// Snapshot 玩家状态快照
type Snapshot struct {
	Name          string `json:"name"`
	Level         int    `json:"level"`
	Exp           int    `json:"exp"`
	Gold          int    `json:"gold"`
	HP            int    `json:"hp"`
	MP            int    `json:"mp"`
	Attack        int    `json:"attack"`
	Defense       int    `json:"defense"`
	SceneID       string `json:"scene_id"`
	PosX          int    `json:"pos_x"`
	PosY          int    `json:"pos_y"`
	Items         string `json:"items"`
	Equipment     string `json:"equipment"`
	CombatWins    int    `json:"combat_wins"`
	SkillsUsed    int    `json:"skills_used"`
	VisitedScenes string `json:"visited_scenes"`
}

// SaveSlotInfo 存档槽信息
type SaveSlotInfo struct {
	SaveID    uint      `json:"save_id"`
	Slot      int       `json:"slot"`
	Name      string    `json:"name"`
	Level     int       `json:"level"`
	SceneID   string    `json:"scene_id"`
	CreatedAt time.Time `json:"created_at"`
	IsEmpty   bool      `json:"is_empty"`
}

// NewSaveGameManager 创建存档管理器
func NewSaveGameManager() *SaveGameManager {
	return &SaveGameManager{}
}

// CreateSnapshot 创建玩家状态快照
func (sgm *SaveGameManager) CreateSnapshot(name string, level, exp, gold, hp, mp, attack, defense int, sceneID string, posX, posY int, items, equipment string, combatWins, skillsUsed int, visitedScenes string) *Snapshot {
	return &Snapshot{
		Name:          name,
		Level:         level,
		Exp:           exp,
		Gold:          gold,
		HP:            hp,
		MP:            mp,
		Attack:        attack,
		Defense:       defense,
		SceneID:       sceneID,
		PosX:          posX,
		PosY:          posY,
		Items:         items,
		Equipment:     equipment,
		CombatWins:    combatWins,
		SkillsUsed:    skillsUsed,
		VisitedScenes: visitedScenes,
	}
}

// SerializeSnapshot 序列化快照
func (sgm *SaveGameManager) SerializeSnapshot(snapshot *Snapshot) (string, error) {
	data, err := json.Marshal(snapshot)
	if err != nil {
		return "", fmt.Errorf("序列化快照失败: %v", err)
	}
	return string(data), nil
}

// DeserializeSnapshot 反序列化快照
func (sgm *SaveGameManager) DeserializeSnapshot(snapshotJSON string) (*Snapshot, error) {
	var snapshot Snapshot
	err := json.Unmarshal([]byte(snapshotJSON), &snapshot)
	if err != nil {
		return nil, fmt.Errorf("反序列化快照失败: %v", err)
	}
	return &snapshot, nil
}

// AutoSaveTrigger 自动存档触发条件
type AutoSaveTrigger string

const (
	AutoSaveSceneChange   AutoSaveTrigger = "scene_change"   // 场景切换
	AutoSaveQuestComplete AutoSaveTrigger = "quest_complete" // 任务完成
	AutoSaveLevelUp       AutoSaveTrigger = "level_up"       // 升级
	AutoSaveCombatWin     AutoSaveTrigger = "combat_win"     // 战斗胜利
)

// ShouldAutoSave 检查是否应该自动存档
func (sgm *SaveGameManager) ShouldAutoSave(trigger AutoSaveTrigger) bool {
	switch trigger {
	case AutoSaveSceneChange, AutoSaveQuestComplete, AutoSaveLevelUp, AutoSaveCombatWin:
		return true
	default:
		return false
	}
}

// GetAutoSaveSlot 获取自动存档槽位（使用槽位0）
func (sgm *SaveGameManager) GetAutoSaveSlot() int {
	return 0
}

// ValidateSlot 验证存档槽位
func (sgm *SaveGameManager) ValidateSlot(slot int) error {
	if slot < 0 || slot > 10 {
		return fmt.Errorf("无效的存档槽位: %d，有效范围 0-10", slot)
	}
	return nil
}

// FormatSaveName 格式化存档名称
func (sgm *SaveGameManager) FormatSaveName(slot int, playerName string) string {
	if slot == 0 {
		return "自动存档"
	}
	return fmt.Sprintf("存档%d - %s", slot, playerName)
}

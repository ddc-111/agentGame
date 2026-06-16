package game

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

// CombatSystem 处理回合制战斗
type CombatSystem struct{}

// CombatState 战斗状态
type CombatState struct {
	PlayerID  uint     `json:"player_id"`
	EnemyType string   `json:"enemy_type"`
	EnemyName string   `json:"enemy_name"`
	PlayerHP  int      `json:"player_hp"`
	PlayerMP  int      `json:"player_mp"`
	EnemyHP   int      `json:"enemy_hp"`
	EnemyAtk  int      `json:"enemy_atk"`
	EnemyDef  int      `json:"enemy_def"`
	Turn      int      `json:"turn"`
	Log       []string `json:"log"`
	IsActive  bool     `json:"is_active"`
	Rewards   Rewards  `json:"rewards"`
}

// Rewards 战斗奖励
type Rewards struct {
	Exp   int    `json:"exp"`
	Gold  int    `json:"gold"`
	Items []uint `json:"items,omitempty"`
}

// EnemyConfig 敌人配置
type EnemyConfig struct {
	Name   string `json:"name"`
	HP     int    `json:"hp"`
	Attack int    `json:"attack"`
	Defense int   `json:"defense"`
	Exp    int    `json:"exp"`
	Gold   int    `json:"gold"`
	Items  []uint `json:"items,omitempty"`
}

// 敌人模板
var enemyTemplates = map[string]EnemyConfig{
	"wolf": {
		Name:   "野狼",
		HP:     50,
		Attack: 12,
		Defense: 3,
		Exp:    30,
		Gold:   20,
	},
	"bandit": {
		Name:   "山贼",
		HP:     80,
		Attack: 15,
		Defense: 5,
		Exp:    50,
		Gold:   50,
	},
	"bear": {
		Name:   "黑熊",
		HP:     120,
		Attack: 20,
		Defense: 8,
		Exp:    80,
		Gold:   30,
	},
	"tiger": {
		Name:   "猛虎",
		HP:     150,
		Attack: 25,
		Defense: 10,
		Exp:    100,
		Gold:   60,
	},
	"ghost": {
		Name:   "厉鬼",
		HP:     100,
		Attack: 18,
		Defense: 2,
		Exp:    70,
		Gold:   40,
	},
}

// NewCombatSystem 创建战斗系统
func NewCombatSystem() *CombatSystem {
	return &CombatSystem{}
}

// StartCombat 发起战斗
func (cs *CombatSystem) StartCombat(playerID uint, enemyType string, playerHP, playerMP int) *CombatState {
	enemy, exists := enemyTemplates[enemyType]
	if !exists {
		enemy = enemyTemplates["wolf"] // 默认敌人
	}

	state := &CombatState{
		PlayerID:  playerID,
		EnemyType: enemyType,
		EnemyName: enemy.Name,
		PlayerHP:  playerHP,
		PlayerMP:  playerMP,
		EnemyHP:   enemy.HP,
		EnemyAtk:  enemy.Attack,
		EnemyDef:  enemy.Defense,
		Turn:      1,
		Log:       []string{"战斗开始！遭遇了" + enemy.Name + "！"},
		IsActive:  true,
		Rewards: Rewards{
			Exp:   enemy.Exp,
			Gold:  enemy.Gold,
			Items: enemy.Items,
		},
	}

	return state
}

// Attack 玩家攻击
func (cs *CombatSystem) Attack(state *CombatState, playerAtk int) *CombatState {
	if !state.IsActive {
		return state
	}

	// 玩家攻击
	damage := cs.calculateDamage(playerAtk, state.EnemyDef)
	state.EnemyHP -= damage
	state.Log = append(state.Log, "你对"+state.EnemyName+"造成了"+formatInt(damage)+"点伤害")

	// 检查敌人是否死亡
	if state.EnemyHP <= 0 {
		state.EnemyHP = 0
		state.IsActive = false
		state.Log = append(state.Log, state.EnemyName+"被击败了！战斗胜利！")
		return state
	}

	// 敌人反击
	enemyDamage := cs.calculateDamage(state.EnemyAtk, 5) // 假设玩家防御为5
	state.PlayerHP -= enemyDamage
	state.Log = append(state.Log, state.EnemyName+"对你造成了"+formatInt(enemyDamage)+"点伤害")

	// 检查玩家是否死亡
	if state.PlayerHP <= 0 {
		state.PlayerHP = 0
		state.IsActive = false
		state.Log = append(state.Log, "你被击败了！战斗失败...")
	}

	state.Turn++
	return state
}

// UseItem 在战斗中使用道具
func (cs *CombatSystem) UseItem(state *CombatState, itemEffect map[string]int) *CombatState {
	if !state.IsActive {
		return state
	}

	// 应用道具效果
	if hp, ok := itemEffect["hp"]; ok {
		state.PlayerHP += hp
		state.Log = append(state.Log, "恢复了"+formatInt(hp)+"点生命")
	}
	if mp, ok := itemEffect["mp"]; ok {
		state.PlayerMP += mp
		state.Log = append(state.Log, "恢复了"+formatInt(mp)+"点法力")
	}

	// 敌人回合
	enemyDamage := cs.calculateDamage(state.EnemyAtk, 5)
	state.PlayerHP -= enemyDamage
	state.Log = append(state.Log, state.EnemyName+"对你造成了"+formatInt(enemyDamage)+"点伤害")

	if state.PlayerHP <= 0 {
		state.PlayerHP = 0
		state.IsActive = false
		state.Log = append(state.Log, "你被击败了！战斗失败...")
	}

	state.Turn++
	return state
}

// Flee 逃跑（基于等级差异）
func (cs *CombatSystem) Flee(state *CombatState, playerLevel int) (bool, *CombatState) {
	if !state.IsActive {
		return false, state
	}

	// 基础逃跑概率50%，每级加5%
	baseChance := 50 + (playerLevel-1)*5
	if baseChance > 90 {
		baseChance = 90
	}

	rand.Seed(time.Now().UnixNano())
	if rand.Intn(100) < baseChance {
		state.IsActive = false
		state.Log = append(state.Log, "逃跑成功！")
		return true, state
	}

	// 逃跑失败，敌人攻击
	enemyDamage := cs.calculateDamage(state.EnemyAtk, 5)
	state.PlayerHP -= enemyDamage
	state.Log = append(state.Log, "逃跑失败！"+state.EnemyName+"对你造成了"+formatInt(enemyDamage)+"点伤害")

	if state.PlayerHP <= 0 {
		state.PlayerHP = 0
		state.IsActive = false
		state.Log = append(state.Log, "你被击败了！战斗失败...")
	}

	state.Turn++
	return false, state
}

// GetRewards 获取战斗奖励
func (cs *CombatSystem) GetRewards(state *CombatState) Rewards {
	return state.Rewards
}

// calculateDamage 计算伤害
func (cs *CombatSystem) calculateDamage(atk, def int) int {
	baseDamage := atk - def
	if baseDamage < 1 {
		baseDamage = 1
	}

	// 添加随机浮动 ±20%
	variation := float64(baseDamage) * 0.2
	min := float64(baseDamage) - variation
	max := float64(baseDamage) + variation

	rand.Seed(time.Now().UnixNano())
	damage := int(min + rand.Float64()*(max-min))

	if damage < 1 {
		damage = 1
	}

	return damage
}

// formatInt 格式化整数
func formatInt(n int) string {
	return fmt.Sprintf("%d", n)
}

// SerializeCombatState 序列化战斗状态
func SerializeCombatState(state *CombatState) (string, error) {
	data, err := json.Marshal(state)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// DeserializeCombatState 反序列化战斗状态
func DeserializeCombatState(data string) (*CombatState, error) {
	var state CombatState
	err := json.Unmarshal([]byte(data), &state)
	if err != nil {
		return nil, err
	}
	return &state, nil
}

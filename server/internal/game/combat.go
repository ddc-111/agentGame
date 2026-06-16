package game

import (
	"encoding/json"
	"fmt"
	"math/rand"
)

// CombatSystem 处理回合制战斗
type CombatSystem struct{}

// CombatState 战斗状态
type CombatState struct {
	PlayerID     uint           `json:"player_id"`
	EnemyType    string         `json:"enemy_type"`
	EnemyName    string         `json:"enemy_name"`
	PlayerHP     int            `json:"player_hp"`
	PlayerMP     int            `json:"player_mp"`
	PlayerDef    int            `json:"player_def"`
	EnemyHP      int            `json:"enemy_hp"`
	EnemyAtk     int            `json:"enemy_atk"`
	EnemyDef     int            `json:"enemy_def"`
	Turn         int            `json:"turn"`
	Log          []string       `json:"log"`
	IsActive     bool           `json:"is_active"`
	Rewards      Rewards        `json:"rewards"`
	ActiveEffects []StatusEffect `json:"active_effects,omitempty"`
}

// Rewards 战斗奖励
type Rewards struct {
	Exp   int    `json:"exp"`
	Gold  int    `json:"gold"`
	Items []uint `json:"items,omitempty"`
}

// StatusEffect 战斗中的状态效果
type StatusEffect struct {
	Type      string `json:"type"`      // stun, poison, defense_up, attack_up
	Duration  int    `json:"duration"`  // remaining turns
	Value     int    `json:"value"`     // effect magnitude
	Source    string `json:"source"`    // "player" or "enemy"
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
		PlayerDef: 5,
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

	// Process enemy effects (poison, stun)
	enemyStunned := cs.ProcessEffects(state, "enemy")

	// Apply attack_up bonus
	atkBonus := cs.GetEffectBonus(state, "player", "attack_up")

	// 玩家攻击
	damage := cs.calculateDamage(playerAtk+atkBonus, state.EnemyDef)
	state.EnemyHP -= damage
	state.Log = append(state.Log, "你对"+state.EnemyName+"造成了"+formatInt(damage)+"点伤害")

	// 检查敌人是否死亡
	if state.EnemyHP <= 0 {
		state.EnemyHP = 0
		state.IsActive = false
		state.Log = append(state.Log, state.EnemyName+"被击败了！战斗胜利！")
		return state
	}

	// 敌人反击（眩晕则跳过）
	if !enemyStunned {
		defBonus := cs.GetEffectBonus(state, "player", "defense_up")
		enemyDamage := cs.calculateDamage(state.EnemyAtk, state.PlayerDef+defBonus)
		state.PlayerHP -= enemyDamage
		state.Log = append(state.Log, state.EnemyName+"对你造成了"+formatInt(enemyDamage)+"点伤害")
	}

	// Process player effects (poison)
	cs.ProcessEffects(state, "player")

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
	enemyDamage := cs.calculateDamage(state.EnemyAtk, state.PlayerDef)
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

	if rand.Intn(100) < baseChance {
		state.IsActive = false
		state.Log = append(state.Log, "逃跑成功！")
		return true, state
	}

	// 逃跑失败，敌人攻击
	enemyDamage := cs.calculateDamage(state.EnemyAtk, state.PlayerDef)
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

// AddEffect 添加状态效果
func (cs *CombatSystem) AddEffect(state *CombatState, effectType string, duration, value int, source string) {
	// Remove existing effect of same type from same source
	var filtered []StatusEffect
	for _, e := range state.ActiveEffects {
		if !(e.Type == effectType && e.Source == source) {
			filtered = append(filtered, e)
		}
	}
	filtered = append(filtered, StatusEffect{
		Type:     effectType,
		Duration: duration,
		Value:    value,
		Source:   source,
	})
	state.ActiveEffects = filtered
}

// ProcessEffects 处理回合开始时的状态效果，返回true表示目标被眩晕跳过回合
func (cs *CombatSystem) ProcessEffects(state *CombatState, target string) bool {
	stunned := false
	var remaining []StatusEffect

	for _, effect := range state.ActiveEffects {
		if effect.Source != target {
			remaining = append(remaining, effect)
			continue
		}

		switch effect.Type {
		case "stun":
			if target == "enemy" {
				state.Log = append(state.Log, state.EnemyName+"被眩晕，无法行动！")
			} else {
				state.Log = append(state.Log, "你被眩晕，无法行动！")
			}
			stunned = true
		case "poison":
			damage := effect.Value
			if target == "enemy" {
				state.EnemyHP -= damage
				if state.EnemyHP < 0 {
					state.EnemyHP = 0
				}
				state.Log = append(state.Log, state.EnemyName+"受到中毒伤害 "+formatInt(damage)+" 点")
			} else {
				state.PlayerHP -= damage
				if state.PlayerHP < 0 {
					state.PlayerHP = 0
				}
				state.Log = append(state.Log, "你受到中毒伤害 "+formatInt(damage)+" 点")
			}
		case "defense_up":
			// Applied during damage calculation, just tick here
		case "attack_up":
			// Applied during damage calculation, just tick here
		}

		effect.Duration--
		if effect.Duration > 0 {
			remaining = append(remaining, effect)
		} else {
			state.Log = append(state.Log, effect.Type+" 效果结束了")
		}
	}

	state.ActiveEffects = remaining
	return stunned
}

// GetEffectBonus 计算效果加成
func (cs *CombatSystem) GetEffectBonus(state *CombatState, source string, effectType string) int {
	bonus := 0
	for _, effect := range state.ActiveEffects {
		if effect.Source == source && effect.Type == effectType {
			bonus += effect.Value
		}
	}
	return bonus
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

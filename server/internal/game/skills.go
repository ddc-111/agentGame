package game

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

// Skill represents a player ability
type Skill struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Type        string `json:"type"` // attack/heal/buff/debuff
	MPCost      int    `json:"mp_cost"`
	Damage      int    `json:"damage"`
	Heal        int    `json:"heal"`
	Cooldown    int    `json:"cooldown"`
	Level       int    `json:"level"`
	Effect      string `json:"effect"` // JSON: {"type":"stun","duration":2}
}

// SkillEffect represents a parsed skill effect
type SkillEffect struct {
	Type     string `json:"type"`
	Duration int    `json:"duration"`
	Value    int    `json:"value"`
}

// PlayerSkill tracks a player's skill cooldowns
type PlayerSkill struct {
	SkillID    uint `json:"skill_id"`
	CooldownCD int  `json:"cooldown_cd"` // remaining cooldown turns
}

// SkillManager handles skill usage
type SkillManager struct{}

// NewSkillManager creates a new skill manager
func NewSkillManager() *SkillManager {
	return &SkillManager{}
}

// UseSkill validates and applies a skill effect in combat
func (sm *SkillManager) UseSkill(skill *Skill, state *CombatState, playerAtk int) (*CombatState, string, error) {
	if !state.IsActive {
		return state, "", fmt.Errorf("战斗已结束")
	}

	if state.PlayerMP < skill.MPCost {
		return state, "", fmt.Errorf("法力不足，需要 %d MP，当前 %d MP", skill.MPCost, state.PlayerMP)
	}

	// Deduct MP
	state.PlayerMP -= skill.MPCost

	var logMsg string

	switch skill.Type {
	case "attack":
		damage := sm.CalculateDamage(skill, playerAtk, state.EnemyDef)
		state.EnemyHP -= damage
		logMsg = fmt.Sprintf("使用【%s】对 %s 造成了 %d 点伤害！(消耗 %d MP)", skill.Name, state.EnemyName, damage, skill.MPCost)

	case "heal":
		healAmount := skill.Heal
		state.PlayerHP += healAmount
		logMsg = fmt.Sprintf("使用【%s】恢复了 %d 点生命！(消耗 %d MP)", skill.Name, healAmount, skill.MPCost)

	case "buff":
		effect := sm.ParseEffect(skill.Effect)
		if effect != nil {
			logMsg = fmt.Sprintf("使用【%s】！%s效果持续 %d 回合！(消耗 %d MP)", skill.Name, effect.Type, effect.Duration, skill.MPCost)
		} else {
			logMsg = fmt.Sprintf("使用【%s】！(消耗 %d MP)", skill.Name, skill.MPCost)
		}

	case "debuff":
		effect := sm.ParseEffect(skill.Effect)
		if effect != nil {
			logMsg = fmt.Sprintf("对 %s 使用【%s】！%s效果持续 %d 回合！(消耗 %d MP)", state.EnemyName, skill.Name, effect.Type, effect.Duration, skill.MPCost)
		} else {
			logMsg = fmt.Sprintf("对 %s 使用【%s】！(消耗 %d MP)", state.EnemyName, skill.Name, skill.MPCost)
		}

	default:
		return state, "", fmt.Errorf("未知技能类型: %s", skill.Type)
	}

	state.Log = append(state.Log, logMsg)

	// Check enemy death
	if state.EnemyHP <= 0 {
		state.EnemyHP = 0
		state.IsActive = false
		state.Log = append(state.Log, state.EnemyName+"被击败了！战斗胜利！")
		return state, logMsg, nil
	}

	// Enemy counterattack
	enemyDamage := sm.calculateEnemyDamage(state.EnemyAtk, state.PlayerDef)
	state.PlayerHP -= enemyDamage
	state.Log = append(state.Log, state.EnemyName+"对你造成了"+formatInt(enemyDamage)+"点伤害")

	if state.PlayerHP <= 0 {
		state.PlayerHP = 0
		state.IsActive = false
		state.Log = append(state.Log, "你被击败了！战斗失败...")
	}

	state.Turn++
	return state, logMsg, nil
}

// CalculateDamage computes final damage with skill multiplier + weapon + stats
func (sm *SkillManager) CalculateDamage(skill *Skill, playerAtk, enemyDef int) int {
	// Base damage = skill damage + player attack - enemy defense
	baseDamage := skill.Damage + playerAtk - enemyDef
	if baseDamage < 1 {
		baseDamage = 1
	}

	// Add random variance ±20%
	variation := float64(baseDamage) * 0.2
	min := float64(baseDamage) - variation
	max := float64(baseDamage) + variation

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	damage := int(min + rng.Float64()*(max-min))

	if damage < 1 {
		damage = 1
	}
	return damage
}

// calculateEnemyDamage calculates enemy counterattack damage
func (sm *SkillManager) calculateEnemyDamage(enemyAtk, playerDef int) int {
	baseDamage := enemyAtk - playerDef
	if baseDamage < 1 {
		baseDamage = 1
	}

	variation := float64(baseDamage) * 0.2
	min := float64(baseDamage) - variation
	max := float64(baseDamage) + variation

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	damage := int(min + rng.Float64()*(max-min))

	if damage < 1 {
		damage = 1
	}
	return damage
}

// GetAvailableSkills returns skills the player can use based on level
func (sm *SkillManager) GetAvailableSkills(allSkills []*Skill, playerLevel int) []*Skill {
	var available []*Skill
	for _, skill := range allSkills {
		if playerLevel >= skill.Level {
			available = append(available, skill)
		}
	}
	return available
}

// ParseEffect parses the JSON effect string
func (sm *SkillManager) ParseEffect(effectJSON string) *SkillEffect {
	if effectJSON == "" || effectJSON == "{}" {
		return nil
	}
	var effect SkillEffect
	if err := json.Unmarshal([]byte(effectJSON), &effect); err != nil {
		return nil
	}
	return &effect
}

// SerializeSkills serializes skills to JSON
func SerializeSkills(skills []*Skill) (string, error) {
	data, err := json.Marshal(skills)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// DeserializeSkills deserializes skills from JSON
func DeserializeSkills(data string) ([]*Skill, error) {
	var skills []*Skill
	err := json.Unmarshal([]byte(data), &skills)
	if err != nil {
		return nil, err
	}
	return skills, nil
}

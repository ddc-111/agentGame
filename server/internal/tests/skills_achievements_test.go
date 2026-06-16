package tests

import (
	"testing"

	"github.com/ddc-111/agentGame/server/internal/game"
)

func TestSkillManager_AttackSkill(t *testing.T) {
	sm := game.NewSkillManager()

	skill := &game.Skill{
		ID:       1,
		Name:     "基础斩击",
		Code:     "skill_basic_slash",
		Type:     "attack",
		MPCost:   5,
		Damage:   15,
		Level:    1,
	}

	state := &game.CombatState{
		PlayerID:  1,
		EnemyName: "野狼",
		PlayerHP:  100,
		PlayerMP:  50,
		EnemyHP:   50,
		EnemyAtk:  12,
		EnemyDef:  3,
		Turn:      1,
		Log:       []string{"战斗开始！"},
		IsActive:  true,
	}

	newState, logMsg, err := sm.UseSkill(skill, state, 10)
	if err != nil {
		t.Fatalf("UseSkill failed: %v", err)
	}

	if newState.PlayerMP != 45 {
		t.Errorf("Expected MP=45, got %d", newState.PlayerMP)
	}

	if logMsg == "" {
		t.Error("Expected non-empty log message")
	}

	t.Logf("Attack skill result: %s", logMsg)
	t.Logf("Enemy HP: %d -> %d", 50, newState.EnemyHP)
}

func TestSkillManager_HealSkill(t *testing.T) {
	sm := game.NewSkillManager()

	skill := &game.Skill{
		ID:     2,
		Name:   "治疗术",
		Code:   "skill_heal",
		Type:   "heal",
		MPCost: 8,
		Heal:   30,
		Level:  3,
	}

	state := &game.CombatState{
		PlayerID:  1,
		EnemyName: "野狼",
		PlayerHP:  50,
		PlayerMP:  30,
		EnemyHP:   50,
		EnemyAtk:  12,
		EnemyDef:  3,
		Turn:      1,
		Log:       []string{},
		IsActive:  true,
	}

	newState, logMsg, err := sm.UseSkill(skill, state, 10)
	if err != nil {
		t.Fatalf("UseSkill failed: %v", err)
	}

	// HP should be 50 + 30 (heal) - enemy counterattack damage
	if newState.PlayerHP < 50 || newState.PlayerHP > 80 {
		t.Errorf("Expected HP between 50-80 (after heal + enemy counter), got %d", newState.PlayerHP)
	}

	if newState.PlayerMP != 22 {
		t.Errorf("Expected MP=22, got %d", newState.PlayerMP)
	}

	t.Logf("Heal skill result: %s", logMsg)
	t.Logf("Player HP after heal + enemy counter: %d", newState.PlayerHP)
}

func TestSkillManager_InsufficientMP(t *testing.T) {
	sm := game.NewSkillManager()

	skill := &game.Skill{
		ID:     3,
		Name:   "火球术",
		Code:   "skill_fireball",
		Type:   "attack",
		MPCost: 15,
		Damage: 25,
		Level:  5,
	}

	state := &game.CombatState{
		PlayerID:  1,
		EnemyName: "野狼",
		PlayerHP:  100,
		PlayerMP:  10, // Not enough MP
		EnemyHP:   50,
		IsActive:  true,
	}

	_, _, err := sm.UseSkill(skill, state, 10)
	if err == nil {
		t.Error("Expected error for insufficient MP")
	}
	t.Logf("Expected error: %v", err)
}

func TestSkillManager_BuffSkill(t *testing.T) {
	sm := game.NewSkillManager()

	skill := &game.Skill{
		ID:       4,
		Name:     "铁壁",
		Code:     "skill_iron_wall",
		Type:     "buff",
		MPCost:   10,
		Level:    4,
		Effect:   `{"type":"defense_up","duration":3,"value":50}`,
	}

	state := &game.CombatState{
		PlayerID:  1,
		EnemyName: "野狼",
		PlayerHP:  100,
		PlayerMP:  50,
		EnemyHP:   50,
		EnemyAtk:  12,
		EnemyDef:  3,
		Turn:      1,
		Log:       []string{},
		IsActive:  true,
	}

	newState, logMsg, err := sm.UseSkill(skill, state, 10)
	if err != nil {
		t.Fatalf("UseSkill failed: %v", err)
	}

	if newState.PlayerMP != 40 {
		t.Errorf("Expected MP=40, got %d", newState.PlayerMP)
	}

	t.Logf("Buff skill result: %s", logMsg)
}

func TestSkillManager_GetAvailableSkills(t *testing.T) {
	sm := game.NewSkillManager()

	allSkills := []*game.Skill{
		{ID: 1, Name: "基础斩击", Level: 1},
		{ID: 2, Name: "治疗术", Level: 3},
		{ID: 3, Name: "火球术", Level: 5},
		{ID: 4, Name: "铁壁", Level: 4},
	}

	// Level 1 player
	available := sm.GetAvailableSkills(allSkills, 1)
	if len(available) != 1 {
		t.Errorf("Expected 1 skill at level 1, got %d", len(available))
	}

	// Level 4 player
	available = sm.GetAvailableSkills(allSkills, 4)
	if len(available) != 3 {
		t.Errorf("Expected 3 skills at level 4, got %d", len(available))
	}

	// Level 10 player
	available = sm.GetAvailableSkills(allSkills, 10)
	if len(available) != 4 {
		t.Errorf("Expected 4 skills at level 10, got %d", len(available))
	}
}

func TestAchievementManager_CheckConditions(t *testing.T) {
	am := game.NewAchievementManager()

	achievements := game.PredefinedAchievements()

	data := &game.PlayerAchievementData{
		Level:           1,
		TotalGold:       100,
		CombatWins:      0,
		QuestCount:      0,
		VisitedScenes:   1,
		UniqueItems:     5,
		TalkedToAllNPCs: false,
		SkillsUsed:      0,
	}

	unlocked := make(map[uint]bool)
	newAchs := am.CheckAchievements(achievements, data, unlocked)
	if len(newAchs) != 0 {
		t.Errorf("Expected 0 achievements for new player, got %d", len(newAchs))
	}

	data.QuestCount = 1
	data.CompletedQuests = map[string]bool{"task_first_arrival": true}
	newAchs = am.CheckAchievements(achievements, data, unlocked)
	if len(newAchs) != 1 {
		t.Errorf("Expected 1 achievement, got %d", len(newAchs))
	}
	if len(newAchs) > 0 && newAchs[0].Code != "ach_first_quest" {
		t.Errorf("Expected ach_first_quest, got %s", newAchs[0].Code)
	}

	data.Level = 10
	newAchs = am.CheckAchievements(achievements, data, unlocked)
	found := false
	for _, a := range newAchs {
		if a.Code == "ach_level_10" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected ach_level_10 achievement")
	}
}

func TestAchievementManager_CombatWinCondition(t *testing.T) {
	am := game.NewAchievementManager()
	achievements := game.PredefinedAchievements()

	unlocked := make(map[uint]bool)
	data := &game.PlayerAchievementData{
		Level:         1,
		TotalGold:     100,
		CombatWins:    1,
		QuestCount:    0,
		VisitedScenes: 1,
		UniqueItems:   0,
		SkillsUsed:    0,
	}

	newAchs := am.CheckAchievements(achievements, data, unlocked)
	found := false
	for _, a := range newAchs {
		if a.Code == "ach_first_combat" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected ach_first_combat achievement when CombatWins=1")
	}
}

func TestAchievementManager_ExploreCondition(t *testing.T) {
	am := game.NewAchievementManager()
	achievements := game.PredefinedAchievements()

	unlocked := make(map[uint]bool)
	data := &game.PlayerAchievementData{
		Level:         1,
		TotalGold:     100,
		CombatWins:    0,
		QuestCount:    0,
		VisitedScenes: 6,
		UniqueItems:   0,
		SkillsUsed:    0,
	}

	newAchs := am.CheckAchievements(achievements, data, unlocked)
	found := false
	for _, a := range newAchs {
		if a.Code == "ach_explorer" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected ach_explorer achievement when VisitedScenes=6")
	}
}

func TestAchievementManager_SkillUseCondition(t *testing.T) {
	am := game.NewAchievementManager()
	achievements := game.PredefinedAchievements()

	unlocked := make(map[uint]bool)
	data := &game.PlayerAchievementData{
		Level:         1,
		TotalGold:     100,
		CombatWins:    0,
		QuestCount:    0,
		VisitedScenes: 1,
		UniqueItems:   0,
		SkillsUsed:    100,
	}

	newAchs := am.CheckAchievements(achievements, data, unlocked)
	found := false
	for _, a := range newAchs {
		if a.Code == "ach_skill_master" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected ach_skill_master achievement when SkillsUsed=100")
	}
}

func TestAchievementManager_GetReward(t *testing.T) {
	am := game.NewAchievementManager()

	reward := am.GetReward(`{"exp":50,"gold":100}`)
	if reward.Exp != 50 {
		t.Errorf("Expected exp=50, got %d", reward.Exp)
	}
	if reward.Gold != 100 {
		t.Errorf("Expected gold=100, got %d", reward.Gold)
	}
}

func TestNPCBehaviorManager_UpdateBehavior(t *testing.T) {
	bm := game.NewNPCBehaviorManager()

	behavior := &game.NPCBehavior{
		NPCCode: "npc_merchant_li",
		State:   "idle",
		Schedule: []game.ScheduleEntry{
			{Time: "06:00", Action: "open_shop", Scene: "scene_general_store"},
			{Time: "22:00", Action: "close_shop", Scene: "scene_village_center"},
		},
		Mood: "neutral",
	}

	// At 6am, should be trading
	bm.UpdateBehavior(behavior, 6)
	if behavior.State != "trading" {
		t.Errorf("Expected state=trading at 6am, got %s", behavior.State)
	}
	if behavior.Target != "scene_general_store" {
		t.Errorf("Expected target=scene_general_store, got %s", behavior.Target)
	}

	// At 10pm, should be idle
	bm.UpdateBehavior(behavior, 22)
	if behavior.State != "idle" {
		t.Errorf("Expected state=idle at 10pm, got %s", behavior.State)
	}
}

func TestNPCBehaviorManager_ReactToPlayer(t *testing.T) {
	bm := game.NewNPCBehaviorManager()

	behavior := &game.NPCBehavior{
		NPCCode: "npc_chief_chen",
		State:   "idle",
		Mood:    "neutral",
		Memory:  []game.NPCEvent{},
	}

	// Talk to NPC
	bm.ReactToPlayer(behavior, 1, "talk")
	if behavior.State != "talking" {
		t.Errorf("Expected state=talking, got %s", behavior.State)
	}
	if len(behavior.Memory) != 1 {
		t.Errorf("Expected 1 memory entry, got %d", len(behavior.Memory))
	}

	// Gift to NPC - mood should become happy
	bm.ReactToPlayer(behavior, 1, "gift")
	if behavior.Mood != "happy" {
		t.Errorf("Expected mood=happy, got %s", behavior.Mood)
	}

	// Attack NPC - mood should become angry
	bm.ReactToPlayer(behavior, 1, "attack")
	if behavior.Mood != "angry" {
		t.Errorf("Expected mood=angry, got %s", behavior.Mood)
	}
	if behavior.State != "fleeing" {
		t.Errorf("Expected state=fleeing, got %s", behavior.State)
	}
}

func TestNPCBehaviorManager_GetDialogMood(t *testing.T) {
	bm := game.NewNPCBehaviorManager()

	behavior := &game.NPCBehavior{
		Mood:  "happy",
		State: "idle",
	}

	mood := bm.GetDialogMood(behavior)
	if mood != "npc_happy" {
		t.Errorf("Expected npc_happy, got %s", mood)
	}

	behavior.Mood = "angry"
	mood = bm.GetDialogMood(behavior)
	if mood != "npc_angry" {
		t.Errorf("Expected npc_angry, got %s", mood)
	}

	behavior.Mood = "neutral"
	behavior.State = "trading"
	mood = bm.GetDialogMood(behavior)
	if mood != "npc_merchant" {
		t.Errorf("Expected npc_merchant, got %s", mood)
	}
}

func TestSkillSerialization(t *testing.T) {
	skills := []*game.Skill{
		{ID: 1, Name: "基础斩击", Type: "attack", MPCost: 5},
		{ID: 2, Name: "治疗术", Type: "heal", MPCost: 8},
	}

	data, err := game.SerializeSkills(skills)
	if err != nil {
		t.Fatalf("SerializeSkills failed: %v", err)
	}

	deserialized, err := game.DeserializeSkills(data)
	if err != nil {
		t.Fatalf("DeserializeSkills failed: %v", err)
	}

	if len(deserialized) != 2 {
		t.Errorf("Expected 2 skills, got %d", len(deserialized))
	}

	if deserialized[0].Name != "基础斩击" {
		t.Errorf("Expected 基础斩击, got %s", deserialized[0].Name)
	}
}

func TestAchievementSerialization(t *testing.T) {
	achievements := game.PredefinedAchievements()

	data, err := game.SerializeAchievements(achievements)
	if err != nil {
		t.Fatalf("SerializeAchievements failed: %v", err)
	}

	deserialized, err := game.DeserializeAchievements(data)
	if err != nil {
		t.Fatalf("DeserializeAchievements failed: %v", err)
	}

	if len(deserialized) < 5 {
		t.Errorf("Expected at least 5 achievements, got %d", len(deserialized))
	}
}

func TestNPCBehaviorSerialization(t *testing.T) {
	behavior := &game.NPCBehavior{
		NPCCode: "npc_test",
		State:   "idle",
		Mood:    "happy",
		Schedule: []game.ScheduleEntry{
			{Time: "08:00", Action: "patrol"},
		},
	}

	data, err := game.SerializeNPCBehavior(behavior)
	if err != nil {
		t.Fatalf("SerializeNPCBehavior failed: %v", err)
	}

	deserialized, err := game.DeserializeNPCBehavior(data)
	if err != nil {
		t.Fatalf("DeserializeNPCBehavior failed: %v", err)
	}

	if deserialized.NPCCode != "npc_test" {
		t.Errorf("Expected npc_test, got %s", deserialized.NPCCode)
	}
	if deserialized.Mood != "happy" {
		t.Errorf("Expected happy, got %s", deserialized.Mood)
	}
	if len(deserialized.Schedule) != 1 {
		t.Errorf("Expected 1 schedule entry, got %d", len(deserialized.Schedule))
	}
}

func TestCombatSystem_SharedRNG_Attack(t *testing.T) {
	cs := game.NewCombatSystem()

	state := cs.StartCombat(1, "wolf", 100, 50)

	for i := 0; i < 10 && state.IsActive; i++ {
		cs.Attack(state, 20)
	}

	if len(state.Log) == 0 {
		t.Error("Expected combat log entries")
	}

	for _, entry := range state.Log {
		if entry == "" {
			t.Error("Expected non-empty log entry")
		}
	}
}

func TestCombatSystem_SharedRNG_DamageVariation(t *testing.T) {
	cs := game.NewCombatSystem()

	damages := make(map[int]bool)
	for i := 0; i < 50; i++ {
		state := cs.StartCombat(1, "wolf", 1000, 50)
		prevHP := state.EnemyHP
		cs.Attack(state, 20)
		damage := prevHP - state.EnemyHP
		if damage > 0 {
			damages[damage] = true
		}
	}

	if len(damages) < 2 {
		t.Errorf("Expected damage variation, got only %d distinct values", len(damages))
	}
}

func TestCombatSystem_SharedRNG_Flee(t *testing.T) {
	cs := game.NewCombatSystem()

	successCount := 0
	for i := 0; i < 100; i++ {
		state := cs.StartCombat(1, "wolf", 100, 50)
		fled, _ := cs.Flee(state, 10)
		if fled {
			successCount++
		}
	}

	if successCount == 0 || successCount == 100 {
		t.Errorf("Expected mixed flee results, got %d/100 successes", successCount)
	}
}

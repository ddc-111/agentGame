package tests

import (
	"encoding/json"
	"testing"

	"github.com/ddc-111/agentGame/server/internal/game"
)

func TestCreateDefaultBehavior(t *testing.T) {
	b := game.CreateDefaultBehavior("npc_test", "")
	if b.NPCCode != "npc_test" {
		t.Errorf("expected npc_test, got %s", b.NPCCode)
	}
	if b.State != "idle" {
		t.Errorf("expected idle, got %s", b.State)
	}
	if b.Mood != "neutral" {
		t.Errorf("expected neutral, got %s", b.Mood)
	}
}

func TestCreateDefaultBehavior_WithSchedule(t *testing.T) {
	scheduleJSON := `[{"time":"08:00","action":"open_shop","scene":"scene_shop"},{"time":"18:00","action":"close_shop","scene":"scene_home"}]`
	b := game.CreateDefaultBehavior("npc_merchant", scheduleJSON)
	if len(b.Schedule) != 2 {
		t.Fatalf("expected 2 schedule entries, got %d", len(b.Schedule))
	}
	if b.Schedule[0].Action != "open_shop" {
		t.Errorf("expected open_shop, got %s", b.Schedule[0].Action)
	}
}

func TestUpdateBehavior_ScheduleMatch(t *testing.T) {
	bm := game.NewNPCBehaviorManager()
	scheduleJSON := `[{"time":"08:00","action":"open_shop","scene":"scene_shop"},{"time":"18:00","action":"close_shop","scene":"scene_home"}]`
	b := game.CreateDefaultBehavior("npc_merchant", scheduleJSON)

	bm.UpdateBehavior(b, 8)
	if b.State != "trading" {
		t.Errorf("expected trading at hour 8, got %s", b.State)
	}
	if b.Target != "scene_shop" {
		t.Errorf("expected scene_shop, got %s", b.Target)
	}
}

func TestUpdateBehavior_NoSchedule(t *testing.T) {
	bm := game.NewNPCBehaviorManager()
	b := game.CreateDefaultBehavior("npc_idle", "")

	bm.UpdateBehavior(b, 12)
	if b.State != "idle" {
		t.Errorf("expected idle with no schedule, got %s", b.State)
	}
}

func TestUpdateBehavior_NoMatch(t *testing.T) {
	bm := game.NewNPCBehaviorManager()
	scheduleJSON := `[{"time":"08:00","action":"open_shop","scene":"scene_shop"}]`
	b := game.CreateDefaultBehavior("npc_merchant", scheduleJSON)
	b.State = "trading"

	bm.UpdateBehavior(b, 15)
	if b.State != "trading" {
		t.Errorf("expected state unchanged when no schedule match, got %s", b.State)
	}
}

func TestReactToPlayer_Talk(t *testing.T) {
	bm := game.NewNPCBehaviorManager()
	b := game.CreateDefaultBehavior("npc_test", "")

	bm.ReactToPlayer(b, 1, "talk")
	if b.State != "talking" {
		t.Errorf("expected talking, got %s", b.State)
	}
	if len(b.Memory) != 1 {
		t.Fatalf("expected 1 memory event, got %d", len(b.Memory))
	}
	if b.Memory[0].Type != "talk" {
		t.Errorf("expected talk event, got %s", b.Memory[0].Type)
	}
}

func TestReactToPlayer_Gift(t *testing.T) {
	bm := game.NewNPCBehaviorManager()
	b := game.CreateDefaultBehavior("npc_test", "")
	b.Mood = "neutral"

	bm.ReactToPlayer(b, 1, "gift")
	if b.Mood != "happy" {
		t.Errorf("expected happy after gift, got %s", b.Mood)
	}
}

func TestReactToPlayer_Attack(t *testing.T) {
	bm := game.NewNPCBehaviorManager()
	b := game.CreateDefaultBehavior("npc_test", "")

	bm.ReactToPlayer(b, 1, "attack")
	if b.Mood != "angry" {
		t.Errorf("expected angry after attack, got %s", b.Mood)
	}
	if b.State != "fleeing" {
		t.Errorf("expected fleeing after attack, got %s", b.State)
	}
}

func TestReactToPlayer_MemoryLimit(t *testing.T) {
	bm := game.NewNPCBehaviorManager()
	b := game.CreateDefaultBehavior("npc_test", "")

	for i := 0; i < 15; i++ {
		bm.ReactToPlayer(b, 1, "talk")
	}
	if len(b.Memory) != 10 {
		t.Errorf("expected memory capped at 10, got %d", len(b.Memory))
	}
}

func TestGetDialogMood_Happy(t *testing.T) {
	bm := game.NewNPCBehaviorManager()
	b := game.CreateDefaultBehavior("npc_test", "")
	b.Mood = "happy"

	mood := bm.GetDialogMood(b)
	if mood != "npc_happy" {
		t.Errorf("expected npc_happy, got %s", mood)
	}
}

func TestGetDialogMood_Angry(t *testing.T) {
	bm := game.NewNPCBehaviorManager()
	b := game.CreateDefaultBehavior("npc_test", "")
	b.Mood = "angry"

	mood := bm.GetDialogMood(b)
	if mood != "npc_angry" {
		t.Errorf("expected npc_angry, got %s", mood)
	}
}

func TestGetDialogMood_NeutralTrading(t *testing.T) {
	bm := game.NewNPCBehaviorManager()
	b := game.CreateDefaultBehavior("npc_test", "")
	b.Mood = "neutral"
	b.State = "trading"

	mood := bm.GetDialogMood(b)
	if mood != "npc_merchant" {
		t.Errorf("expected npc_merchant, got %s", mood)
	}
}

func TestGetDialogContext(t *testing.T) {
	bm := game.NewNPCBehaviorManager()
	b := game.CreateDefaultBehavior("npc_test", "")
	b.State = "trading"
	b.Mood = "happy"
	b.Target = "scene_shop"

	ctx := bm.GetDialogContext(b)
	if ctx == "" {
		t.Error("expected non-empty dialog context")
	}
}

func TestBehaviorStore_GetOrCreate(t *testing.T) {
	store := game.NewNPCBehaviorStore()

	b1 := store.GetOrCreate("npc_test", "")
	if b1.NPCCode != "npc_test" {
		t.Errorf("expected npc_test, got %s", b1.NPCCode)
	}

	b1.State = "trading"
	b2 := store.GetOrCreate("npc_test", "")
	if b2.State != "trading" {
		t.Errorf("expected cached trading state, got %s", b2.State)
	}
}

func TestBehaviorStore_Get(t *testing.T) {
	store := game.NewNPCBehaviorStore()

	b := store.Get("nonexistent")
	if b != nil {
		t.Error("expected nil for nonexistent NPC")
	}

	store.Set("npc_test", game.CreateDefaultBehavior("npc_test", ""))
	b = store.Get("npc_test")
	if b == nil {
		t.Error("expected behavior after Set")
	}
}

func TestBehaviorStore_All(t *testing.T) {
	store := game.NewNPCBehaviorStore()
	store.Set("npc_a", game.CreateDefaultBehavior("npc_a", ""))
	store.Set("npc_b", game.CreateDefaultBehavior("npc_b", ""))

	all := store.All()
	if len(all) != 2 {
		t.Errorf("expected 2 behaviors, got %d", len(all))
	}
}

func TestUpdateAllBehaviors(t *testing.T) {
	bm := game.NewNPCBehaviorManager()
	store := game.NewNPCBehaviorStore()

	scheduleJSON := `[{"time":"08:00","action":"open_shop","scene":"scene_shop"}]`
	store.Set("npc_shop", game.CreateDefaultBehavior("npc_shop", scheduleJSON))
	store.Set("npc_idle", game.CreateDefaultBehavior("npc_idle", ""))

	bm.UpdateAllBehaviors(store, 8)

	shop := store.Get("npc_shop")
	if shop.State != "trading" {
		t.Errorf("expected trading, got %s", shop.State)
	}

	idle := store.Get("npc_idle")
	if idle.State != "idle" {
		t.Errorf("expected idle, got %s", idle.State)
	}
}

func TestSerializeDeserialize(t *testing.T) {
	b := game.CreateDefaultBehavior("npc_test", `[{"time":"10:00","action":"patrol","scene":"scene_forest"}]`)
	b.Mood = "happy"
	b.Memory = []game.NPCEvent{{Time: "10:30", Type: "talk", PlayerID: 1, Detail: "hello"}}

	data, err := game.SerializeNPCBehavior(b)
	if err != nil {
		t.Fatalf("serialize error: %v", err)
	}

	b2, err := game.DeserializeNPCBehavior(data)
	if err != nil {
		t.Fatalf("deserialize error: %v", err)
	}

	if b2.NPCCode != b.NPCCode {
		t.Errorf("expected %s, got %s", b.NPCCode, b2.NPCCode)
	}
	if b2.Mood != "happy" {
		t.Errorf("expected happy, got %s", b2.Mood)
	}
	if len(b2.Memory) != 1 {
		t.Fatalf("expected 1 memory, got %d", len(b2.Memory))
	}
	if b2.Memory[0].Type != "talk" {
		t.Errorf("expected talk, got %s", b2.Memory[0].Type)
	}
}

func TestActionToState(t *testing.T) {
	bm := game.NewNPCBehaviorManager()

	tests := []struct {
		action string
		state  string
	}{
		{"open_shop", "trading"},
		{"close_shop", "idle"},
		{"patrol", "patrolling"},
		{"go_hunt", "patrolling"},
		{"rest", "idle"},
		{"stand_at_tree", "idle"},
		{"unknown", "idle"},
	}

	for _, tt := range tests {
		b := game.CreateDefaultBehavior("npc_test", "")
		b.Schedule = []game.ScheduleEntry{{Time: "12:00", Action: tt.action}}
		bm.UpdateBehavior(b, 12)
		if b.State != tt.state {
			t.Errorf("action %s: expected %s, got %s", tt.action, tt.state, b.State)
		}
	}
}

func TestReactToPlayer_GiftWhenHappy(t *testing.T) {
	bm := game.NewNPCBehaviorManager()
	b := game.CreateDefaultBehavior("npc_test", "")
	b.Mood = "happy"

	bm.ReactToPlayer(b, 1, "gift")
	if b.Mood != "happy" {
		t.Errorf("expected happy to stay happy after gift, got %s", b.Mood)
	}
}

func TestReactToPlayer_GiftWhenAngry(t *testing.T) {
	bm := game.NewNPCBehaviorManager()
	b := game.CreateDefaultBehavior("npc_test", "")
	b.Mood = "angry"

	bm.ReactToPlayer(b, 1, "gift")
	if b.Mood != "angry" {
		t.Errorf("expected angry to stay angry after gift, got %s", b.Mood)
	}
}

func TestGetDialogMood_Scared(t *testing.T) {
	bm := game.NewNPCBehaviorManager()
	b := game.CreateDefaultBehavior("npc_test", "")
	b.Mood = "scared"

	mood := bm.GetDialogMood(b)
	if mood != "npc_scared" {
		t.Errorf("expected npc_scared, got %s", mood)
	}
}

func TestGetDialogMood_NeutralTalking(t *testing.T) {
	bm := game.NewNPCBehaviorManager()
	b := game.CreateDefaultBehavior("npc_test", "")
	b.Mood = "neutral"
	b.State = "talking"

	mood := bm.GetDialogMood(b)
	if mood != "npc_friendly" {
		t.Errorf("expected npc_friendly, got %s", mood)
	}
}

func TestGetDialogMood_NeutralPatrolling(t *testing.T) {
	bm := game.NewNPCBehaviorManager()
	b := game.CreateDefaultBehavior("npc_test", "")
	b.Mood = "neutral"
	b.State = "patrolling"

	mood := bm.GetDialogMood(b)
	if mood != "npc_alert" {
		t.Errorf("expected npc_alert, got %s", mood)
	}
}

func TestGetDialogMood_NeutralIdle(t *testing.T) {
	bm := game.NewNPCBehaviorManager()
	b := game.CreateDefaultBehavior("npc_test", "")
	b.Mood = "neutral"
	b.State = "idle"

	mood := bm.GetDialogMood(b)
	if mood != "npc_neutral" {
		t.Errorf("expected npc_neutral, got %s", mood)
	}
}

func TestGetDialogContext_WithMemory(t *testing.T) {
	bm := game.NewNPCBehaviorManager()
	b := game.CreateDefaultBehavior("npc_test", "")
	b.State = "talking"
	b.Mood = "happy"
	b.Memory = []game.NPCEvent{{Time: "10:30", Type: "talk", PlayerID: 1, Detail: "友好交谈"}}

	ctx := bm.GetDialogContext(b)
	if ctx == "" {
		t.Error("expected non-empty dialog context")
	}
}

func TestReactToPlayer_TalkWhenAngry(t *testing.T) {
	bm := game.NewNPCBehaviorManager()
	b := game.CreateDefaultBehavior("npc_test", "")
	b.Mood = "angry"

	bm.ReactToPlayer(b, 1, "talk")
	if b.State != "talking" {
		t.Errorf("expected talking, got %s", b.State)
	}
	if len(b.Memory) != 1 {
		t.Fatalf("expected 1 memory event, got %d", len(b.Memory))
	}
	if b.Memory[0].Detail != "NPC不情愿地交谈" {
		t.Errorf("expected reluctant talk detail, got %s", b.Memory[0].Detail)
	}
}

func TestBehaviorStore_GetOrCreatePreservesState(t *testing.T) {
	store := game.NewNPCBehaviorStore()

	b := store.GetOrCreate("npc_test", "")
	b.State = "talking"
	b.Mood = "happy"
	b.Memory = []game.NPCEvent{{Time: "10:00", Type: "talk", PlayerID: 1, Detail: "hello"}}

	b2 := store.GetOrCreate("npc_test", `[{"time":"08:00","action":"patrol"}]`)
	if b2.State != "talking" {
		t.Errorf("expected state preserved, got %s", b2.State)
	}
	if b2.Mood != "happy" {
		t.Errorf("expected mood preserved, got %s", b2.Mood)
	}
	if len(b2.Memory) != 1 {
		t.Errorf("expected memory preserved, got %d", len(b2.Memory))
	}
}

func TestSerializeDeserialize_EmptyMemory(t *testing.T) {
	b := game.CreateDefaultBehavior("npc_test", "")
	b.Memory = []game.NPCEvent{}

	data, err := game.SerializeNPCBehavior(b)
	if err != nil {
		t.Fatalf("serialize error: %v", err)
	}

	b2, err := game.DeserializeNPCBehavior(data)
	if err != nil {
		t.Fatalf("deserialize error: %v", err)
	}

	if b2.NPCCode != "npc_test" {
		t.Errorf("expected npc_test, got %s", b2.NPCCode)
	}
	if b2.Memory == nil {
		t.Error("expected non-nil memory after deserialize")
	}
}

func TestDeserializeNPCBehavior_InvalidJSON(t *testing.T) {
	_, err := game.DeserializeNPCBehavior("not json")
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestUpdateBehavior_MultipleScheduleEntries(t *testing.T) {
	bm := game.NewNPCBehaviorManager()
	scheduleJSON := `[{"time":"06:00","action":"patrol","scene":"scene_forest"},{"time":"12:00","action":"rest","scene":"scene_home"},{"time":"18:00","action":"open_shop","scene":"scene_shop"}]`
	b := game.CreateDefaultBehavior("npc_merchant", scheduleJSON)

	bm.UpdateBehavior(b, 12)
	if b.State != "idle" {
		t.Errorf("expected idle at 12 (rest), got %s", b.State)
	}
	if b.Target != "scene_home" {
		t.Errorf("expected scene_home at 12, got %s", b.Target)
	}

	bm.UpdateBehavior(b, 18)
	if b.State != "trading" {
		t.Errorf("expected trading at 18, got %s", b.State)
	}
	if b.Target != "scene_shop" {
		t.Errorf("expected scene_shop at 18, got %s", b.Target)
	}
}

func TestReactToPlayer_MultipleEvents(t *testing.T) {
	bm := game.NewNPCBehaviorManager()
	b := game.CreateDefaultBehavior("npc_test", "")

	bm.ReactToPlayer(b, 1, "talk")
	bm.ReactToPlayer(b, 2, "gift")
	bm.ReactToPlayer(b, 3, "attack")

	if len(b.Memory) != 3 {
		t.Fatalf("expected 3 memory events, got %d", len(b.Memory))
	}
	if b.Memory[0].Type != "talk" {
		t.Errorf("expected talk, got %s", b.Memory[0].Type)
	}
	if b.Memory[1].Type != "gift" {
		t.Errorf("expected gift, got %s", b.Memory[1].Type)
	}
	if b.Memory[2].Type != "attack" {
		t.Errorf("expected attack, got %s", b.Memory[2].Type)
	}
	if b.Mood != "angry" {
		t.Errorf("expected angry after attack, got %s", b.Mood)
	}
	if b.State != "fleeing" {
		t.Errorf("expected fleeing after attack, got %s", b.State)
	}
}

func TestGetDialogContext_EmptyMemory(t *testing.T) {
	bm := game.NewNPCBehaviorManager()
	b := game.CreateDefaultBehavior("npc_test", "")
	b.State = "idle"
	b.Mood = "neutral"
	b.Target = ""

	ctx := bm.GetDialogContext(b)
	if ctx == "" {
		t.Error("expected non-empty dialog context even with no memory")
	}
}

func TestNPCBehaviorJSON(t *testing.T) {
	b := game.CreateDefaultBehavior("npc_test", `[{"time":"08:00","action":"open_shop","scene":"scene_shop"}]`)
	b.Mood = "happy"
	b.State = "trading"
	b.Target = "scene_shop"
	b.Memory = []game.NPCEvent{
		{Time: "08:05", Type: "talk", PlayerID: 1, Detail: "hello"},
	}

	data, err := json.Marshal(b)
	if err != nil {
		t.Fatalf("json marshal error: %v", err)
	}

	var b2 game.NPCBehavior
	if err := json.Unmarshal(data, &b2); err != nil {
		t.Fatalf("json unmarshal error: %v", err)
	}

	if b2.NPCCode != "npc_test" {
		t.Errorf("expected npc_test, got %s", b2.NPCCode)
	}
	if b2.Mood != "happy" {
		t.Errorf("expected happy, got %s", b2.Mood)
	}
	if b2.State != "trading" {
		t.Errorf("expected trading, got %s", b2.State)
	}
	if len(b2.Schedule) != 1 {
		t.Fatalf("expected 1 schedule, got %d", len(b2.Schedule))
	}
	if b2.Schedule[0].Action != "open_shop" {
		t.Errorf("expected open_shop, got %s", b2.Schedule[0].Action)
	}
}

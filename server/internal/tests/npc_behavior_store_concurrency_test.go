package tests

import (
	"sync"
	"testing"

	"github.com/ddc-111/agentGame/server/internal/game"
)

func TestBehaviorStoreSnapshotsAreIndependent(t *testing.T) {
	store := game.NewNPCBehaviorStore()
	original := game.CreateDefaultBehavior("npc_test", `[{"time":"08:00","action":"open_shop","scene":"scene_shop"}]`)
	original.Memory = []game.NPCEvent{{Type: "talk", Detail: "original"}}
	store.Set("npc_test", original)

	first := store.Get("npc_test")
	second := store.Get("npc_test")
	first.Schedule[0].Scene = "mutated_scene"
	first.Memory[0].Detail = "mutated_memory"

	if second.Schedule[0].Scene != "scene_shop" {
		t.Fatalf("schedule snapshot shares backing storage: %q", second.Schedule[0].Scene)
	}
	if second.Memory[0].Detail != "original" {
		t.Fatalf("memory snapshot shares backing storage: %q", second.Memory[0].Detail)
	}

	stored := store.Get("npc_test")
	if stored.Schedule[0].Scene != "scene_shop" || stored.Memory[0].Detail != "original" {
		t.Fatalf("caller mutation leaked into store: %+v", stored)
	}
}

func TestBehaviorStoreSetCopiesInput(t *testing.T) {
	store := game.NewNPCBehaviorStore()
	behavior := game.CreateDefaultBehavior("npc_test", "")
	behavior.Memory = []game.NPCEvent{{Type: "talk", Detail: "before"}}
	store.Set("npc_test", behavior)

	behavior.State = "fleeing"
	behavior.Memory[0].Detail = "after"

	stored := store.Get("npc_test")
	if stored.State != "idle" {
		t.Fatalf("state mutation leaked into store: %q", stored.State)
	}
	if stored.Memory[0].Detail != "before" {
		t.Fatalf("slice mutation leaked into store: %q", stored.Memory[0].Detail)
	}
}

func TestBehaviorStoreConcurrentSnapshotsAndTicks(t *testing.T) {
	store := game.NewNPCBehaviorStore()
	manager := game.NewNPCBehaviorManager()
	store.Set("npc_shop", game.CreateDefaultBehavior("npc_shop", `[{"time":"08:00","action":"open_shop","scene":"scene_shop"}]`))

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 500; j++ {
				snapshot := store.GetOrCreateCopy("npc_shop", "")
				_ = snapshot.State
				_ = snapshot.Schedule
				_ = store.All()
			}
		}()
	}

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 500; j++ {
				manager.UpdateAllBehaviors(store, 8)
				manager.UpdateAllBehaviors(store, 12)
			}
		}()
	}

	wg.Wait()
	if got := store.Get("npc_shop"); got == nil {
		t.Fatal("behavior disappeared during concurrent access")
	}
}

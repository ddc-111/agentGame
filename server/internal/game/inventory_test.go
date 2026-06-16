package game

import (
	"fmt"
	"testing"
)

func TestEquipmentStatsFromEquip_EmptyEquipment(t *testing.T) {
	im := NewInventoryManager()
	lookup := func(itemID uint) (map[string]int, error) {
		return nil, nil
	}

	stats, err := im.EquipmentStatsFromEquip("", lookup)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.Attack != 0 || stats.Defense != 0 || stats.HP != 0 || stats.MP != 0 {
		t.Errorf("expected zero stats, got %+v", stats)
	}

	stats, err = im.EquipmentStatsFromEquip("{}", lookup)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.Attack != 0 || stats.Defense != 0 || stats.HP != 0 || stats.MP != 0 {
		t.Errorf("expected zero stats, got %+v", stats)
	}
}

func TestEquipmentStatsFromEquip_WeaponOnly(t *testing.T) {
	im := NewInventoryManager()
	items := map[uint]map[string]int{
		1: {"attack": 10},
	}
	lookup := func(itemID uint) (map[string]int, error) {
		if eff, ok := items[itemID]; ok {
			return eff, nil
		}
		return nil, fmt.Errorf("item not found")
	}

	stats, err := im.EquipmentStatsFromEquip(`{"weapon_id":1,"armor_id":0}`, lookup)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.Attack != 10 {
		t.Errorf("expected attack=10, got %d", stats.Attack)
	}
	if stats.Defense != 0 {
		t.Errorf("expected defense=0, got %d", stats.Defense)
	}
}

func TestEquipmentStatsFromEquip_FullEquipment(t *testing.T) {
	im := NewInventoryManager()
	items := map[uint]map[string]int{
		1: {"attack": 10},
		2: {"defense": 15, "hp": 50},
	}
	lookup := func(itemID uint) (map[string]int, error) {
		if eff, ok := items[itemID]; ok {
			return eff, nil
		}
		return nil, fmt.Errorf("item not found")
	}

	stats, err := im.EquipmentStatsFromEquip(`{"weapon_id":1,"armor_id":2}`, lookup)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.Attack != 10 {
		t.Errorf("expected attack=10, got %d", stats.Attack)
	}
	if stats.Defense != 15 {
		t.Errorf("expected defense=15, got %d", stats.Defense)
	}
	if stats.HP != 50 {
		t.Errorf("expected hp=50, got %d", stats.HP)
	}
	if stats.MP != 0 {
		t.Errorf("expected mp=0, got %d", stats.MP)
	}
}

func TestEquipmentStatsFromEquip_LookupError(t *testing.T) {
	im := NewInventoryManager()
	lookup := func(itemID uint) (map[string]int, error) {
		return nil, fmt.Errorf("item not found")
	}

	stats, err := im.EquipmentStatsFromEquip(`{"weapon_id":999,"armor_id":999}`, lookup)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.Attack != 0 || stats.Defense != 0 {
		t.Errorf("expected zero stats on lookup error, got %+v", stats)
	}
}

func TestEquipmentStatsFromEquip_InvalidJSON(t *testing.T) {
	im := NewInventoryManager()
	lookup := func(itemID uint) (map[string]int, error) {
		return nil, nil
	}

	_, err := im.EquipmentStatsFromEquip(`{invalid json`, lookup)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestCalculateStats_WithEquipment(t *testing.T) {
	im := NewInventoryManager()
	equipStats := EquipmentStats{Attack: 10, Defense: 15, HP: 50, MP: 20}
	stats := im.CalculateStats(100, 50, 200, 100, equipStats)

	if stats.TotalAttack != 110 {
		t.Errorf("expected total_attack=110, got %d", stats.TotalAttack)
	}
	if stats.TotalDefense != 65 {
		t.Errorf("expected total_defense=65, got %d", stats.TotalDefense)
	}
	if stats.TotalHP != 250 {
		t.Errorf("expected total_hp=250, got %d", stats.TotalHP)
	}
	if stats.TotalMP != 120 {
		t.Errorf("expected total_mp=120, got %d", stats.TotalMP)
	}
}

package network

import (
	"testing"

	"github.com/ddc-111/agentGame/server/internal/database/models"
)

func TestProcessLevelUp(t *testing.T) {
	player := &models.Player{
		Level:   1,
		Exp:     0,
		HP:      100,
		MP:      50,
		Attack:  10,
		Defense: 5,
	}

	gained := processLevelUp(player)
	if gained != 0 {
		t.Errorf("expected 0 levels gained, got %d", gained)
	}
	if player.Level != 1 {
		t.Errorf("expected level 1, got %d", player.Level)
	}

	player.Exp = 100
	gained = processLevelUp(player)
	if gained != 1 {
		t.Errorf("expected 1 level gained, got %d", gained)
	}
	if player.Level != 2 {
		t.Errorf("expected level 2, got %d", player.Level)
	}
	if player.Exp != 0 {
		t.Errorf("expected exp 0, got %d", player.Exp)
	}
	if player.HP != 110 {
		t.Errorf("expected hp 110, got %d", player.HP)
	}
	if player.MP != 55 {
		t.Errorf("expected mp 55, got %d", player.MP)
	}
	if player.Attack != 12 {
		t.Errorf("expected attack 12, got %d", player.Attack)
	}
	if player.Defense != 6 {
		t.Errorf("expected defense 6, got %d", player.Defense)
	}

	player.Exp = 500
	gained = processLevelUp(player)
	if gained != 2 {
		t.Errorf("expected 2 levels gained, got %d", gained)
	}
	if player.Level != 4 {
		t.Errorf("expected level 4, got %d", player.Level)
	}
	if player.Exp != 0 {
		t.Errorf("expected exp 0, got %d", player.Exp)
	}
}

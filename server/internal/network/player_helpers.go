package network

import (
	"github.com/ddc-111/agentGame/server/internal/database/models"
)

func processLevelUp(player *models.Player) int {
	levelsGained := 0
	for {
		expNeeded := player.Level * 100
		if player.Exp < expNeeded {
			break
		}
		player.Exp -= expNeeded
		player.Level++
		player.HP += 10
		player.MP += 5
		player.Attack += 2
		player.Defense += 1
		levelsGained++
	}
	return levelsGained
}

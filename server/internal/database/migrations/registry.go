package migrations

import "github.com/ddc-111/agentGame/server/internal/database/models"

func All() []models.Migration {
	return []models.Migration{
		_001,
	}
}

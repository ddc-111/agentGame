package models

import "time"

type SchemaMigration struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Version   string    `gorm:"size:50;uniqueIndex;not null" json:"version"`
	Name      string    `gorm:"size:200;not null" json:"name"`
	AppliedAt time.Time `gorm:"not null" json:"applied_at"`
}

func (SchemaMigration) TableName() string {
	return "schema_migrations"
}

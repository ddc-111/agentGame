package database

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/ddc-111/agentGame/server/internal/database/models"
	"gorm.io/gorm"
)

type Migrator struct {
	db         *gorm.DB
	migrations []models.Migration
}

func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{db: db}
}

func (m *Migrator) Register(migrations ...models.Migration) {
	m.migrations = append(m.migrations, migrations...)
	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].Version < m.migrations[j].Version
	})
}

func (m *Migrator) ensureTable() error {
	return m.db.AutoMigrate(&models.SchemaMigration{})
}

func (m *Migrator) appliedVersions() (map[string]bool, error) {
	var records []models.SchemaMigration
	if err := m.db.Find(&records).Error; err != nil {
		return nil, err
	}
	versions := make(map[string]bool, len(records))
	for _, r := range records {
		versions[r.Version] = true
	}
	return versions, nil
}

func (m *Migrator) Up() error {
	if err := m.ensureTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	applied, err := m.appliedVersions()
	if err != nil {
		return fmt.Errorf("failed to query applied migrations: %w", err)
	}

	for _, mig := range m.migrations {
		if applied[mig.Version] {
			continue
		}

		log.Printf("Running migration %s_%s", mig.Version, mig.Name)

		err := m.db.Transaction(func(tx *gorm.DB) error {
			if err := mig.Up(tx); err != nil {
				return err
			}
			return tx.Create(&models.SchemaMigration{
				Version:   mig.Version,
				Name:      mig.Name,
				AppliedAt: time.Now(),
			}).Error
		})
		if err != nil {
			return fmt.Errorf("migration %s failed: %w", mig.Version, err)
		}

		log.Printf("Migration %s_%s applied", mig.Version, mig.Name)
	}

	return nil
}

func (m *Migrator) Down(steps int) error {
	if err := m.ensureTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	if steps <= 0 {
		return nil
	}

	var records []models.SchemaMigration
	if err := m.db.Order("version DESC").Limit(steps).Find(&records).Error; err != nil {
		return fmt.Errorf("failed to query applied migrations: %w", err)
	}

	migMap := make(map[string]models.Migration)
	for _, mig := range m.migrations {
		migMap[mig.Version] = mig
	}

	for _, record := range records {
		mig, ok := migMap[record.Version]
		if !ok {
			return fmt.Errorf("migration %s not found in registered migrations", record.Version)
		}

		log.Printf("Rolling back migration %s_%s", mig.Version, mig.Name)

		err := m.db.Transaction(func(tx *gorm.DB) error {
			if err := mig.Down(tx); err != nil {
				return err
			}
			return tx.Where("version = ?", mig.Version).Delete(&models.SchemaMigration{}).Error
		})
		if err != nil {
			return fmt.Errorf("rollback %s failed: %w", mig.Version, err)
		}

		log.Printf("Migration %s_%s rolled back", mig.Version, mig.Name)
	}

	return nil
}

type MigrationStatus struct {
	Version   string
	Name      string
	Applied   bool
	AppliedAt *time.Time
}

func (m *Migrator) Status() ([]MigrationStatus, error) {
	if err := m.ensureTable(); err != nil {
		return nil, err
	}

	var records []models.SchemaMigration
	if err := m.db.Find(&records).Error; err != nil {
		return nil, err
	}

	applied := make(map[string]time.Time, len(records))
	for _, r := range records {
		applied[r.Version] = r.AppliedAt
	}

	statuses := make([]MigrationStatus, 0, len(m.migrations))
	for _, mig := range m.migrations {
		s := MigrationStatus{
			Version: mig.Version,
			Name:    mig.Name,
		}
		if ts, ok := applied[mig.Version]; ok {
			s.Applied = true
			s.AppliedAt = &ts
		}
		statuses = append(statuses, s)
	}

	return statuses, nil
}

func (m *Migrator) CurrentVersion() (string, error) {
	if err := m.ensureTable(); err != nil {
		return "", err
	}

	var record models.SchemaMigration
	if err := m.db.Order("version DESC").First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil
		}
		return "", err
	}

	return record.Version, nil
}

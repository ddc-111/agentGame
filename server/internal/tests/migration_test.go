package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ddc-111/agentGame/server/internal/database"
	"github.com/ddc-111/agentGame/server/internal/database/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var migrationDBCounter int64

func setupMigrationTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	id := atomic.AddInt64(&migrationDBCounter, 1)
	tmpDir := t.TempDir()
	dsn := filepath.Join(tmpDir, fmt.Sprintf("migration_test_%d.db", id))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
		os.Remove(dsn)
	})
	return db
}

func TestMigratorUp(t *testing.T) {
	db := setupMigrationTestDB(t)
	migrator := database.NewMigrator(db)

	migrator.Register(
		models.Migration{
			Version: "001",
			Name:    "create_scenes",
			Up: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&models.Scene{})
			},
			Down: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(&models.Scene{})
			},
		},
		models.Migration{
			Version: "002",
			Name:    "create_npcs",
			Up: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&models.NPC{})
			},
			Down: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(&models.NPC{})
			},
		},
	)

	if err := migrator.Up(); err != nil {
		t.Fatalf("Up() failed: %v", err)
	}

	if !db.Migrator().HasTable(&models.Scene{}) {
		t.Error("scenes table should exist after Up")
	}
	if !db.Migrator().HasTable(&models.NPC{}) {
		t.Error("npcs table should exist after Up")
	}
	if !db.Migrator().HasTable(&models.SchemaMigration{}) {
		t.Error("schema_migrations table should exist after Up")
	}

	var count int64
	db.Model(&models.SchemaMigration{}).Count(&count)
	if count != 2 {
		t.Errorf("expected 2 migration records, got %d", count)
	}
}

func TestMigratorUpIdempotent(t *testing.T) {
	db := setupMigrationTestDB(t)
	migrator := database.NewMigrator(db)

	runCount := 0
	migrator.Register(
		models.Migration{
			Version: "001",
			Name:    "test_migration",
			Up: func(tx *gorm.DB) error {
				runCount++
				return tx.AutoMigrate(&models.Scene{})
			},
			Down: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(&models.Scene{})
			},
		},
	)

	if err := migrator.Up(); err != nil {
		t.Fatalf("first Up() failed: %v", err)
	}
	if runCount != 1 {
		t.Errorf("expected migration to run once, ran %d times", runCount)
	}

	if err := migrator.Up(); err != nil {
		t.Fatalf("second Up() failed: %v", err)
	}
	if runCount != 1 {
		t.Errorf("expected migration to not run again, ran %d times", runCount)
	}
}

func TestMigratorDown(t *testing.T) {
	db := setupMigrationTestDB(t)
	migrator := database.NewMigrator(db)

	migrator.Register(
		models.Migration{
			Version: "001",
			Name:    "create_scenes",
			Up: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&models.Scene{})
			},
			Down: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(&models.Scene{})
			},
		},
		models.Migration{
			Version: "002",
			Name:    "create_npcs",
			Up: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&models.NPC{})
			},
			Down: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(&models.NPC{})
			},
		},
	)

	if err := migrator.Up(); err != nil {
		t.Fatalf("Up() failed: %v", err)
	}

	if err := migrator.Down(1); err != nil {
		t.Fatalf("Down(1) failed: %v", err)
	}

	if !db.Migrator().HasTable(&models.Scene{}) {
		t.Error("scenes table should still exist after rolling back only NPCs")
	}
	if db.Migrator().HasTable(&models.NPC{}) {
		t.Error("npcs table should not exist after Down(1)")
	}

	var count int64
	db.Model(&models.SchemaMigration{}).Count(&count)
	if count != 1 {
		t.Errorf("expected 1 migration record after Down(1), got %d", count)
	}
}

func TestMigratorDownAll(t *testing.T) {
	db := setupMigrationTestDB(t)
	migrator := database.NewMigrator(db)

	migrator.Register(
		models.Migration{
			Version: "001",
			Name:    "create_scenes",
			Up: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&models.Scene{})
			},
			Down: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(&models.Scene{})
			},
		},
	)

	if err := migrator.Up(); err != nil {
		t.Fatalf("Up() failed: %v", err)
	}

	if err := migrator.Down(1); err != nil {
		t.Fatalf("Down(1) failed: %v", err)
	}

	if db.Migrator().HasTable(&models.Scene{}) {
		t.Error("scenes table should not exist after Down(1)")
	}

	var count int64
	db.Model(&models.SchemaMigration{}).Count(&count)
	if count != 0 {
		t.Errorf("expected 0 migration records after full rollback, got %d", count)
	}
}

func TestMigratorDownZeroSteps(t *testing.T) {
	db := setupMigrationTestDB(t)
	migrator := database.NewMigrator(db)

	migrator.Register(
		models.Migration{
			Version: "001",
			Name:    "create_scenes",
			Up: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&models.Scene{})
			},
			Down: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(&models.Scene{})
			},
		},
	)

	if err := migrator.Up(); err != nil {
		t.Fatalf("Up() failed: %v", err)
	}

	if err := migrator.Down(0); err != nil {
		t.Fatalf("Down(0) should not fail: %v", err)
	}

	var count int64
	db.Model(&models.SchemaMigration{}).Count(&count)
	if count != 1 {
		t.Errorf("Down(0) should not remove any records, got %d", count)
	}
}

func TestMigratorStatus(t *testing.T) {
	db := setupMigrationTestDB(t)
	migrator := database.NewMigrator(db)

	migrator.Register(
		models.Migration{
			Version: "001",
			Name:    "first",
			Up:      func(tx *gorm.DB) error { return nil },
			Down:    func(tx *gorm.DB) error { return nil },
		},
		models.Migration{
			Version: "002",
			Name:    "second",
			Up:      func(tx *gorm.DB) error { return nil },
			Down:    func(tx *gorm.DB) error { return nil },
		},
	)

	if err := migrator.Up(); err != nil {
		t.Fatalf("Up() failed: %v", err)
	}

	statuses, err := migrator.Status()
	if err != nil {
		t.Fatalf("Status() failed: %v", err)
	}

	if len(statuses) != 2 {
		t.Fatalf("expected 2 statuses, got %d", len(statuses))
	}

	if !statuses[0].Applied {
		t.Error("first migration should be applied")
	}
	if statuses[0].AppliedAt == nil {
		t.Error("first migration should have AppliedAt")
	}
	if !statuses[1].Applied {
		t.Error("second migration should be applied")
	}
}

func TestMigratorStatusPartial(t *testing.T) {
	db := setupMigrationTestDB(t)
	migrator := database.NewMigrator(db)

	migrator.Register(
		models.Migration{
			Version: "001",
			Name:    "first",
			Up:      func(tx *gorm.DB) error { return tx.AutoMigrate(&models.Scene{}) },
			Down:    func(tx *gorm.DB) error { return tx.Migrator().DropTable(&models.Scene{}) },
		},
		models.Migration{
			Version: "002",
			Name:    "second",
			Up:      func(tx *gorm.DB) error { return nil },
			Down:    func(tx *gorm.DB) error { return nil },
		},
	)

	if err := migrator.Up(); err != nil {
		t.Fatalf("Up() failed: %v", err)
	}

	if err := migrator.Down(1); err != nil {
		t.Fatalf("Down(1) failed: %v", err)
	}

	statuses, err := migrator.Status()
	if err != nil {
		t.Fatalf("Status() failed: %v", err)
	}

	if !statuses[0].Applied {
		t.Error("first migration should still be applied")
	}
	if statuses[1].Applied {
		t.Error("second migration should not be applied after rollback")
	}
}

func TestMigratorCurrentVersion(t *testing.T) {
	db := setupMigrationTestDB(t)
	migrator := database.NewMigrator(db)

	migrator.Register(
		models.Migration{
			Version: "001",
			Name:    "first",
			Up:      func(tx *gorm.DB) error { return nil },
			Down:    func(tx *gorm.DB) error { return nil },
		},
		models.Migration{
			Version: "002",
			Name:    "second",
			Up:      func(tx *gorm.DB) error { return nil },
			Down:    func(tx *gorm.DB) error { return nil },
		},
	)

	version, err := migrator.CurrentVersion()
	if err != nil {
		t.Fatalf("CurrentVersion() failed: %v", err)
	}
	if version != "" {
		t.Errorf("expected empty version before any migration, got %s", version)
	}

	if err := migrator.Up(); err != nil {
		t.Fatalf("Up() failed: %v", err)
	}

	version, err = migrator.CurrentVersion()
	if err != nil {
		t.Fatalf("CurrentVersion() failed: %v", err)
	}
	if version != "002" {
		t.Errorf("expected version 002, got %s", version)
	}
}

func TestMigratorDownMoreThanApplied(t *testing.T) {
	db := setupMigrationTestDB(t)
	migrator := database.NewMigrator(db)

	migrator.Register(
		models.Migration{
			Version: "001",
			Name:    "first",
			Up:      func(tx *gorm.DB) error { return tx.AutoMigrate(&models.Scene{}) },
			Down:    func(tx *gorm.DB) error { return tx.Migrator().DropTable(&models.Scene{}) },
		},
	)

	if err := migrator.Up(); err != nil {
		t.Fatalf("Up() failed: %v", err)
	}

	if err := migrator.Down(5); err != nil {
		t.Fatalf("Down(5) should not fail with more steps than applied: %v", err)
	}

	var count int64
	db.Model(&models.SchemaMigration{}).Count(&count)
	if count != 0 {
		t.Errorf("expected 0 records after rolling back all, got %d", count)
	}
}

func TestMigratorVersionOrdering(t *testing.T) {
	db := setupMigrationTestDB(t)
	migrator := database.NewMigrator(db)

	order := []string{}
	migrator.Register(
		models.Migration{
			Version: "003",
			Name:    "third",
			Up: func(tx *gorm.DB) error {
				order = append(order, "003")
				return nil
			},
			Down: func(tx *gorm.DB) error { return nil },
		},
		models.Migration{
			Version: "001",
			Name:    "first",
			Up: func(tx *gorm.DB) error {
				order = append(order, "001")
				return nil
			},
			Down: func(tx *gorm.DB) error { return nil },
		},
		models.Migration{
			Version: "002",
			Name:    "second",
			Up: func(tx *gorm.DB) error {
				order = append(order, "002")
				return nil
			},
			Down: func(tx *gorm.DB) error { return nil },
		},
	)

	if err := migrator.Up(); err != nil {
		t.Fatalf("Up() failed: %v", err)
	}

	if len(order) != 3 || order[0] != "001" || order[1] != "002" || order[2] != "003" {
		t.Errorf("migrations should run in version order, got %v", order)
	}
}

func TestMigratorTransactionalRollback(t *testing.T) {
	db := setupMigrationTestDB(t)
	migrator := database.NewMigrator(db)

	migrator.Register(
		models.Migration{
			Version: "001",
			Name:    "good_migration",
			Up: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&models.Scene{})
			},
			Down: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(&models.Scene{})
			},
		},
		models.Migration{
			Version: "002",
			Name:    "bad_migration",
			Up: func(tx *gorm.DB) error {
				return fmt.Errorf("intentional failure")
			},
			Down: func(tx *gorm.DB) error { return nil },
		},
	)

	err := migrator.Up()
	if err == nil {
		t.Fatal("Up() should fail when a migration fails")
	}

	var count int64
	db.Model(&models.SchemaMigration{}).Count(&count)
	if count != 1 {
		t.Errorf("expected 1 applied migration after failed second, got %d", count)
	}
}

func TestSchemaMigrationTableName(t *testing.T) {
	m := models.SchemaMigration{}
	if m.TableName() != "schema_migrations" {
		t.Errorf("expected table name schema_migrations, got %s", m.TableName())
	}
}

func TestSchemaMigrationFields(t *testing.T) {
	now := time.Now()
	m := models.SchemaMigration{
		ID:        1,
		Version:   "001",
		Name:      "test",
		AppliedAt: now,
	}

	if m.Version != "001" {
		t.Errorf("expected version 001, got %s", m.Version)
	}
	if m.Name != "test" {
		t.Errorf("expected name test, got %s", m.Name)
	}
	if m.AppliedAt != now {
		t.Error("AppliedAt mismatch")
	}
}

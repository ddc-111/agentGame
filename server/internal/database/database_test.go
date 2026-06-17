package database

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"gorm.io/gorm"
)

type TestModel struct {
	gorm.Model
	Name string
}

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		cfg       Config
		wantErr   bool
		checkFunc func(*Database) error
	}{
		{
			name: "SQLite with default DSN",
			cfg: Config{
				Driver: "sqlite",
				DSN:    "",
			},
			wantErr: false,
			checkFunc: func(db *Database) error {
				if db.DB == nil {
					return fmt.Errorf("expected non-nil DB")
				}
				return nil
			},
		},
		{
			name: "SQLite with custom DSN",
			cfg: Config{
				Driver: "sqlite",
				DSN:    filepath.Join(os.TempDir(), "test.db"),
			},
			wantErr: false,
			checkFunc: func(db *Database) error {
				if db.DB == nil {
					return fmt.Errorf("expected non-nil DB")
				}
				return nil
			},
		},
		{
			name: "SQLite memory database",
			cfg: Config{
				Driver: "sqlite",
				DSN:    ":memory:",
			},
			wantErr: false,
			checkFunc: func(db *Database) error {
				if db.DB == nil {
					return fmt.Errorf("expected non-nil DB")
				}
				return nil
			},
		},
		{
			name: "Unsupported driver",
			cfg: Config{
				Driver: "unsupported",
				DSN:    "some_dsn",
			},
			wantErr: true,
			checkFunc: func(db *Database) error {
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := New(tt.cfg)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("New() expected error, got nil")
					return
				}
				return
			}

			if err != nil {
				t.Errorf("New() unexpected error: %v", err)
				return
			}

			if err := tt.checkFunc(db); err != nil {
				t.Errorf("checkFunc failed: %v", err)
			}

			if tt.cfg.Driver == "sqlite" && tt.cfg.DSN != "" && tt.cfg.DSN != ":memory:" {
				defer os.Remove(tt.cfg.DSN)
			}
		})
	}
}

func TestDatabase_Close(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "Close SQLite database",
			cfg: Config{
				Driver: "sqlite",
				DSN:    ":memory:",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := New(tt.cfg)
			if err != nil {
				t.Fatalf("New() failed: %v", err)
			}

			err = db.Close()
			if tt.wantErr && err == nil {
				t.Errorf("Close() expected error, got nil")
			} else if !tt.wantErr && err != nil {
				t.Errorf("Close() unexpected error: %v", err)
			}
		})
	}
}

func TestDatabase_AutoMigrate(t *testing.T) {
	tests := []struct {
		name    string
		models  []interface{}
		wantErr bool
	}{
		{
			name:    "AutoMigrate single model",
			models:  []interface{}{&TestModel{}},
			wantErr: false,
		},
		{
			name:    "AutoMigrate empty models",
			models:  []interface{}{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := New(Config{
				Driver: "sqlite",
				DSN:    ":memory:",
			})
			if err != nil {
				t.Fatalf("New() failed: %v", err)
			}
			defer db.Close()

			err = db.AutoMigrate(tt.models...)
			if tt.wantErr && err == nil {
				t.Errorf("AutoMigrate() expected error, got nil")
			} else if !tt.wantErr && err != nil {
				t.Errorf("AutoMigrate() unexpected error: %v", err)
			}
		})
	}
}

func TestDatabase_DB_Field(t *testing.T) {
	db, err := New(Config{
		Driver: "sqlite",
		DSN:    ":memory:",
	})
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer db.Close()

	if db.DB == nil {
		t.Errorf("DB field should not be nil")
	}

	sqlDB, err := db.DB.DB()
	if err != nil {
		t.Errorf("Failed to get underlying sql.DB: %v", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		t.Errorf("Database ping failed: %v", err)
	}
}

func TestConfigStruct(t *testing.T) {
	cfg := Config{
		Driver:   "test_driver",
		DSN:      "test_dsn",
		Host:     "test_host",
		Port:     3306,
		User:     "test_user",
		Password: "test_password",
		DBName:   "test_db",
	}

	if cfg.Driver != "test_driver" {
		t.Errorf("Expected Driver to be 'test_driver', got '%s'", cfg.Driver)
	}
	if cfg.DSN != "test_dsn" {
		t.Errorf("Expected DSN to be 'test_dsn', got '%s'", cfg.DSN)
	}
	if cfg.Host != "test_host" {
		t.Errorf("Expected Host to be 'test_host', got '%s'", cfg.Host)
	}
	if cfg.Port != 3306 {
		t.Errorf("Expected Port to be 3306, got %d", cfg.Port)
	}
}

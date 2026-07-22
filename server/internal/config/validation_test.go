package config

import (
	"strings"
	"testing"
)

func TestDefaultConfigValidatesForDevelopment(t *testing.T) {
	if err := Default().Validate(); err != nil {
		t.Fatalf("Default().Validate() error = %v", err)
	}
}

func TestReleaseConfigRejectsDefaultCredentials(t *testing.T) {
	cfg := Default()
	cfg.Server.Mode = "release"

	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected release config with default credentials to fail")
	}
	if !strings.Contains(err.Error(), "default GM credentials") {
		t.Fatalf("error %q does not mention default credentials", err)
	}
}

func TestValidReleaseConfig(t *testing.T) {
	cfg := Default()
	cfg.Server.Mode = "release"
	cfg.Auth.JWTSecret = "01234567890123456789012345678901"
	cfg.Auth.GMUsername = "operator"
	cfg.Auth.GMPassword = "correct-horse-battery-staple"
	cfg.CORS.AllowedOrigins = []string{"https://gm.example.com"}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestMySQLConfigRequiresConnectionFieldsWithoutDSN(t *testing.T) {
	cfg := Default()
	cfg.Database = DatabaseConfig{Driver: "mysql"}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected invalid mysql config to fail")
	}
	for _, field := range []string{"database.host", "database.port", "database.user", "database.dbname"} {
		if !strings.Contains(err.Error(), field) {
			t.Fatalf("error %q does not mention %s", err, field)
		}
	}
}

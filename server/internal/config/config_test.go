package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestDefault(t *testing.T) {
	tests := []struct {
		name     string
		expected *Config
	}{
		{
			name: "default config values",
			expected: &Config{
				Server: ServerConfig{
					Port:     8080,
					Mode:     "debug",
					LogLevel: "info",
				},
				Database: DatabaseConfig{
					Driver: "sqlite",
					DSN:    "game.db",
				},
				AI: AIConfig{
					Provider:    "openai",
					BaseURL:     "https://api.openai.com/v1",
					Model:       "gpt-4",
					Temperature: 0.7,
					MaxTokens:   500,
				},
				Generator: GeneratorConfig{
					Enabled:     true,
					Provider:    "openai",
					BaseURL:     "https://api.openai.com/v1",
					Model:       "gpt-4-turbo",
					Temperature: 0.7,
					MaxTokens:   4000,
					Timeout:     60,
				},
				Game: GameConfig{
					MaxPlayers: 100,
					TickRate:   20,
				},
				Auth: AuthConfig{
					JWTSecret:   "change-me-in-production",
					TokenExpiry: 24,
					GMUsername:  "admin",
					GMPassword:  "admin123",
				},
				CORS: CORSConfig{
					AllowedOrigins: []string{"http://localhost:5173", "http://localhost:5174"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := Default()
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("Default() = %v, want %v", actual, tt.expected)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	createTestFile := func(filename, content string) string {
		filePath := filepath.Join(tmpDir, filename)
		err := os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			t.Fatal(err)
		}
		return filePath
	}

	validConfig := `
server:
  port: 9090
  mode: "release"
  log_level: "warn"
database:
  driver: "mysql"
  dsn: ""
  host: "localhost"
  port: 3306
  user: "testuser"
  password: "testpass"
  dbname: "testdb"
ai:
  provider: "anthropic"
  base_url: "https://api.anthropic.com/v1"
  api_key: "test-key"
  model: "claude-3"
  temperature: 0.8
  max_tokens: 1000
generator:
  enabled: false
  provider: "custom"
  base_url: "http://localhost:8000/v1"
  api_key: "custom-key"
  model: "custom-model"
  temperature: 0.5
  max_tokens: 2000
  timeout: 30
game:
  max_players: 50
  tick_rate: 30
auth:
  jwt_secret: "custom-secret"
  token_expiry: 12
  gm_username: "gmuser"
  gm_password: "gmpass123"
cors:
  allowed_origins:
    - "http://example.com"
    - "https://example.com"
`

	invalidConfig := `
server:
  port: "invalid-port"
  mode: 123
`

	tests := []struct {
		name       string
		filename   string
		content    string
		wantErr    bool
		wantConfig *Config
	}{
		{
			name:     "valid config file",
			filename: "valid_config.yaml",
			content:  validConfig,
			wantErr:  false,
			wantConfig: &Config{
				Server: ServerConfig{
					Port:     9090,
					Mode:     "release",
					LogLevel: "warn",
				},
				Database: DatabaseConfig{
					Driver:   "mysql",
					DSN:      "",
					Host:     "localhost",
					Port:     3306,
					User:     "testuser",
					Password: "testpass",
					DBName:   "testdb",
				},
				AI: AIConfig{
					Provider:    "anthropic",
					BaseURL:     "https://api.anthropic.com/v1",
					APIKey:      "test-key",
					Model:       "claude-3",
					Temperature: 0.8,
					MaxTokens:   1000,
				},
				Generator: GeneratorConfig{
					Enabled:     false,
					Provider:    "custom",
					BaseURL:     "http://localhost:8000/v1",
					APIKey:      "custom-key",
					Model:       "custom-model",
					Temperature: 0.5,
					MaxTokens:   2000,
					Timeout:     30,
				},
				Game: GameConfig{
					MaxPlayers: 50,
					TickRate:   30,
				},
				Auth: AuthConfig{
					JWTSecret:   "custom-secret",
					TokenExpiry: 12,
					GMUsername:  "gmuser",
					GMPassword:  "gmpass123",
				},
				CORS: CORSConfig{
					AllowedOrigins: []string{"http://example.com", "https://example.com"},
				},
			},
		},
		{
			name:     "non-existent file",
			filename: "nonexistent.yaml",
			wantErr:  true,
		},
		{
			name:     "invalid yaml content",
			filename: "invalid_config.yaml",
			content:  invalidConfig,
			wantErr:  true,
		},
		{
			name:     "empty config file",
			filename: "empty_config.yaml",
			content:  "",
			wantErr:  false,
			wantConfig: &Config{
				Server: ServerConfig{
					Port:     8080,
					Mode:     "debug",
					LogLevel: "info",
				},
				Database: DatabaseConfig{
					Driver: "sqlite",
					DSN:    "game.db",
				},
				AI: AIConfig{
					Provider:    "openai",
					BaseURL:     "https://api.openai.com/v1",
					Model:       "gpt-4",
					Temperature: 0.7,
					MaxTokens:   500,
				},
				Generator: GeneratorConfig{
					Enabled:     true,
					Provider:    "openai",
					BaseURL:     "https://api.openai.com/v1",
					Model:       "gpt-4-turbo",
					Temperature: 0.7,
					MaxTokens:   4000,
					Timeout:     60,
				},
				Game: GameConfig{
					MaxPlayers: 100,
					TickRate:   20,
				},
				Auth: AuthConfig{
					JWTSecret:   "change-me-in-production",
					TokenExpiry: 24,
					GMUsername:  "admin",
					GMPassword:  "admin123",
				},
				CORS: CORSConfig{
					AllowedOrigins: []string{"http://localhost:5173", "http://localhost:5174"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var filePath string
			if tt.filename == "nonexistent.yaml" {
				filePath = filepath.Join(tmpDir, tt.filename)
			} else {
				filePath = createTestFile(tt.filename, tt.content)
			}

			config, err := Load(filePath)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Load() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Load() unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(config, tt.wantConfig) {
				t.Errorf("Load() got = %v, want %v", config, tt.wantConfig)
			}
		})
	}
}

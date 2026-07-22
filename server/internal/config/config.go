package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	defaultJWTSecret  = "change-me-in-production"
	defaultGMUsername = "admin"
	defaultGMPassword = "admin123"
)

type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Database  DatabaseConfig  `yaml:"database"`
	AI        AIConfig        `yaml:"ai"`
	Generator GeneratorConfig `yaml:"generator"`
	Game      GameConfig      `yaml:"game"`
	Auth      AuthConfig      `yaml:"auth"`
	CORS      CORSConfig      `yaml:"cors"`
}

type ServerConfig struct {
	Port     int    `yaml:"port"`
	Mode     string `yaml:"mode"`
	LogLevel string `yaml:"log_level"`
}

type AuthConfig struct {
	JWTSecret   string `yaml:"jwt_secret"`
	TokenExpiry int    `yaml:"token_expiry"`
	GMUsername  string `yaml:"gm_username"`
	GMPassword  string `yaml:"gm_password"`
}

type CORSConfig struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
}

type DatabaseConfig struct {
	Driver   string `yaml:"driver"` // sqlite 或 mysql
	DSN      string `yaml:"dsn"`    // 数据库连接字符串
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type AIConfig struct {
	Provider    string  `yaml:"provider"`
	BaseURL     string  `yaml:"base_url"`
	APIKey      string  `yaml:"api_key"`
	Model       string  `yaml:"model"`
	Temperature float64 `yaml:"temperature"`
	MaxTokens   int     `yaml:"max_tokens"`
}

// GeneratorConfig 生成智能体配置（独立大模型地址）
type GeneratorConfig struct {
	Enabled     bool    `yaml:"enabled"`
	Provider    string  `yaml:"provider"`    // openai, anthropic, custom
	BaseURL     string  `yaml:"base_url"`    // 大模型API地址
	APIKey      string  `yaml:"api_key"`     // API密钥
	Model       string  `yaml:"model"`       // 模型名称
	Temperature float64 `yaml:"temperature"` // 温度
	MaxTokens   int     `yaml:"max_tokens"`  // 最大token
	Timeout     int     `yaml:"timeout"`     // 超时时间(秒)
}

type GameConfig struct {
	MaxPlayers int `yaml:"max_players"`
	TickRate   int `yaml:"tick_rate"`
}

func Default() *Config {
	return &Config{
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
			JWTSecret:   defaultJWTSecret,
			TokenExpiry: 24,
			GMUsername:  defaultGMUsername,
			GMPassword:  defaultGMPassword,
		},
		CORS: CORSConfig{
			AllowedOrigins: []string{"http://localhost:5173", "http://localhost:5174"},
		},
	}
}

func Load(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	cfg := Default()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return cfg, nil
}

func (c *Config) Validate() error {
	if c == nil {
		return errors.New("config is nil")
	}

	var errs []error
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		errs = append(errs, fmt.Errorf("server.port must be between 1 and 65535"))
	}

	mode := strings.ToLower(strings.TrimSpace(c.Server.Mode))
	switch mode {
	case "debug", "release", "test":
		c.Server.Mode = mode
	default:
		errs = append(errs, fmt.Errorf("server.mode must be debug, release, or test"))
	}

	if strings.TrimSpace(c.Auth.JWTSecret) == "" {
		errs = append(errs, errors.New("auth.jwt_secret is required"))
	}
	if c.Auth.TokenExpiry <= 0 {
		errs = append(errs, errors.New("auth.token_expiry must be greater than zero"))
	}
	if strings.TrimSpace(c.Auth.GMUsername) == "" {
		errs = append(errs, errors.New("auth.gm_username is required"))
	}
	if strings.TrimSpace(c.Auth.GMPassword) == "" {
		errs = append(errs, errors.New("auth.gm_password is required"))
	}

	if mode == "release" {
		if c.Auth.JWTSecret == defaultJWTSecret || len(c.Auth.JWTSecret) < 32 {
			errs = append(errs, errors.New("auth.jwt_secret must be changed and contain at least 32 characters in release mode"))
		}
		if c.Auth.GMUsername == defaultGMUsername && c.Auth.GMPassword == defaultGMPassword {
			errs = append(errs, errors.New("default GM credentials are forbidden in release mode"))
		}
		if len(c.Auth.GMPassword) < 12 {
			errs = append(errs, errors.New("auth.gm_password must contain at least 12 characters in release mode"))
		}
		for _, origin := range c.CORS.AllowedOrigins {
			if strings.TrimSpace(origin) == "*" {
				errs = append(errs, errors.New("cors.allowed_origins cannot contain '*' in release mode"))
				break
			}
		}
	}

	switch strings.ToLower(strings.TrimSpace(c.Database.Driver)) {
	case "sqlite":
		if strings.TrimSpace(c.Database.DSN) == "" {
			errs = append(errs, errors.New("database.dsn is required for sqlite"))
		}
	case "mysql":
		if strings.TrimSpace(c.Database.DSN) == "" {
			if strings.TrimSpace(c.Database.Host) == "" {
				errs = append(errs, errors.New("database.host is required for mysql when dsn is empty"))
			}
			if c.Database.Port < 1 || c.Database.Port > 65535 {
				errs = append(errs, errors.New("database.port must be between 1 and 65535 for mysql"))
			}
			if strings.TrimSpace(c.Database.User) == "" {
				errs = append(errs, errors.New("database.user is required for mysql when dsn is empty"))
			}
			if strings.TrimSpace(c.Database.DBName) == "" {
				errs = append(errs, errors.New("database.dbname is required for mysql when dsn is empty"))
			}
		}
	default:
		errs = append(errs, fmt.Errorf("database.driver must be sqlite or mysql"))
	}

	if c.Game.MaxPlayers <= 0 {
		errs = append(errs, errors.New("game.max_players must be greater than zero"))
	}
	if c.Game.TickRate <= 0 {
		errs = append(errs, errors.New("game.tick_rate must be greater than zero"))
	}
	if c.Generator.Enabled {
		if strings.TrimSpace(c.Generator.BaseURL) == "" {
			errs = append(errs, errors.New("generator.base_url is required when generator is enabled"))
		}
		if strings.TrimSpace(c.Generator.Model) == "" {
			errs = append(errs, errors.New("generator.model is required when generator is enabled"))
		}
		if c.Generator.Timeout <= 0 {
			errs = append(errs, errors.New("generator.timeout must be greater than zero when generator is enabled"))
		}
	}

	return errors.Join(errs...)
}

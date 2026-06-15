package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	AI       AIConfig       `yaml:"ai"`
	Game     GameConfig     `yaml:"game"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

type DatabaseConfig struct {
	Driver   string `yaml:"driver"`   // sqlite 或 mysql
	DSN      string `yaml:"dsn"`      // 数据库连接字符串
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type AIConfig struct {
	Provider string `yaml:"provider"`
	APIKey   string `yaml:"api_key"`
	Model    string `yaml:"model"`
}

type GameConfig struct {
	MaxPlayers int `yaml:"max_players"`
	TickRate   int `yaml:"tick_rate"`
}

func Default() *Config {
	return &Config{
		Server: ServerConfig{
			Port: 8080,
			Mode: "debug",
		},
		Database: DatabaseConfig{
			Driver: "sqlite",
			DSN:    "game.db",
		},
		AI: AIConfig{
			Provider: "openai",
			Model:    "gpt-4",
		},
		Game: GameConfig{
			MaxPlayers: 100,
			TickRate:   20,
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
		return nil, err
	}

	return cfg, nil
}

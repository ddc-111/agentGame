package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Database  DatabaseConfig  `yaml:"database"`
	AI        AIConfig        `yaml:"ai"`
	Generator GeneratorConfig `yaml:"generator"`
	Game      GameConfig      `yaml:"game"`
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
			Port: 8080,
			Mode: "debug",
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

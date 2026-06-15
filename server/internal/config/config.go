package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	AI     AIConfig     `yaml:"ai"`
	Game   GameConfig   `yaml:"game"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
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

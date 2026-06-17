package generator

import (
	"context"
	"strings"
	"testing"

	"github.com/ddc-111/agentGame/server/internal/config"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.GeneratorConfig
		wantErr bool
	}{
		{
			name: "disabled config",
			cfg: config.GeneratorConfig{
				Enabled: false,
			},
			wantErr: false,
		},
		{
			name: "enabled config with valid API key",
			cfg: config.GeneratorConfig{
				Enabled:     true,
				APIKey:      "test-api-key",
				Model:       "gpt-3.5-turbo",
				Temperature: 0.7,
				MaxTokens:   100,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen, err := New(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gen != nil && gen.cfg.Enabled != tt.cfg.Enabled {
				t.Errorf("New() cfg.Enabled = %v, want %v", gen.cfg.Enabled, tt.cfg.Enabled)
			}
		})
	}
}

func TestGenerator_IsEnabled(t *testing.T) {
	tests := []struct {
		name string
		cfg  config.GeneratorConfig
		want bool
	}{
		{
			name: "enabled",
			cfg:  config.GeneratorConfig{Enabled: true},
			want: true,
		},
		{
			name: "disabled",
			cfg:  config.GeneratorConfig{Enabled: false},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen, _ := New(tt.cfg)
			if got := gen.IsEnabled(); got != tt.want {
				t.Errorf("Generator.IsEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerator_GetConfig(t *testing.T) {
	cfg := config.GeneratorConfig{
		Enabled:     true,
		APIKey:      "test-key",
		Model:       "gpt-4",
		Temperature: 0.8,
		MaxTokens:   150,
		Timeout:     30,
	}
	gen, _ := New(cfg)
	got := gen.GetConfig()

	if got.APIKey != cfg.APIKey {
		t.Errorf("Generator.GetConfig() APIKey = %v, want %v", got.APIKey, cfg.APIKey)
	}
	if got.Model != cfg.Model {
		t.Errorf("Generator.GetConfig() Model = %v, want %v", got.Model, cfg.Model)
	}
}

func TestGenerator_Generate(t *testing.T) {
	tests := []struct {
		name        string
		cfg         config.GeneratorConfig
		req         GenerateRequest
		wantSuccess bool
		wantErr     bool
	}{
		{
			name: "disabled generator",
			cfg:  config.GeneratorConfig{Enabled: false},
			req: GenerateRequest{
				Type:   "npc",
				Action: "create",
			},
			wantSuccess: false,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen, err := New(tt.cfg)
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}
			got, err := gen.Generate(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generator.Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Success != tt.wantSuccess {
				t.Errorf("Generator.Generate() success = %v, wantSuccess %v", got.Success, tt.wantSuccess)
			}
		})
	}
}

func TestGenerator_getSystemPrompt(t *testing.T) {
	gen, _ := New(config.GeneratorConfig{})
	tests := []struct {
		name     string
		genType  string
		contains string
	}{
		{
			name:     "npc type",
			genType:  "npc",
			contains: "NPC配置",
		},
		{
			name:     "scene type",
			genType:  "scene",
			contains: "场景配置",
		},
		{
			name:     "task type",
			genType:  "task",
			contains: "任务配置",
		},
		{
			name:     "unknown type",
			genType:  "unknown",
			contains: "古风RPG游戏配置生成助手",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gen.getSystemPrompt(tt.genType)
			if !strings.Contains(got, tt.contains) {
				t.Errorf("getSystemPrompt(%s) = %v, want contains %v", tt.genType, got, tt.contains)
			}
		})
	}
}

func Test_extractJSON(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "json in code block with json marker",
			content: "```json\n{\"name\": \"test\"}\n```",
			want:    "{\"name\": \"test\"}",
		},
		{
			name:    "json in code block without json marker",
			content: "```\n{\"name\": \"test\"}\n```",
			want:    "{\"name\": \"test\"}",
		},
		{
			name:    "json with explanation",
			content: "Here is the JSON response:\n{\"name\": \"test\"}\nEnd of response",
			want:    "{\"name\": \"test\"}",
		},
		{
			name:    "empty string",
			content: "",
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractJSON(tt.content); got != tt.want {
				t.Errorf("extractJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

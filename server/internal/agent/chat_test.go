package agent

import (
	"context"
	"testing"

	"github.com/ddc-111/agentGame/server/internal/config"
)

func TestNewChatManager(t *testing.T) {
	tests := []struct {
		name string
		cfg  config.AIConfig
	}{
		{
			name: "with API key",
			cfg: config.AIConfig{
				APIKey:  "test-key",
				BaseURL: "https://api.example.com/v1",
				Model:   "gpt-4",
			},
		},
		{
			name: "without base URL",
			cfg: config.AIConfig{
				APIKey: "test-key",
				Model:  "gpt-4",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := NewChatManager(tt.cfg)
			if cm == nil {
				t.Error("NewChatManager() returned nil")
			}
			if cm.cfg.APIKey != tt.cfg.APIKey {
				t.Errorf("NewChatManager() APIKey = %v, want %v", cm.cfg.APIKey, tt.cfg.APIKey)
			}
		})
	}
}

func TestChatManager_IsEnabled(t *testing.T) {
	tests := []struct {
		name string
		cfg  config.AIConfig
		want bool
	}{
		{
			name: "enabled with API key",
			cfg:  config.AIConfig{APIKey: "test-key"},
			want: true,
		},
		{
			name: "disabled without API key",
			cfg:  config.AIConfig{APIKey: ""},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := NewChatManager(tt.cfg)
			if got := cm.IsEnabled(); got != tt.want {
				t.Errorf("ChatManager.IsEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChatManager_buildSystemPrompt(t *testing.T) {
	cm := NewChatManager(config.AIConfig{})

	tests := []struct {
		name    string
		persona *NPCPersona
		want    string
	}{
		{
			name: "with custom system prompt",
			persona: &NPCPersona{
				Name:         "张三",
				Title:        "铁匠",
				Description:  "一位经验丰富的铁匠",
				SystemPrompt: "你是铁匠张三，专门打造武器。",
			},
			want: "你是铁匠张三，专门打造武器。",
		},
		{
			name: "with default system prompt",
			persona: &NPCPersona{
				Name:        "李四",
				Title:       "客栈老板",
				Description: "经营一家客栈",
			},
			want: "李四",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cm.buildSystemPrompt(tt.persona)
			if tt.name == "with custom system prompt" {
				if got != tt.want {
					t.Errorf("buildSystemPrompt() = %v, want %v", got, tt.want)
				}
			} else {
				if !contains(got, tt.want) {
					t.Errorf("buildSystemPrompt() = %v, should contain %v", got, tt.want)
				}
			}
		})
	}
}

func TestChatManager_fallbackReply(t *testing.T) {
	cm := NewChatManager(config.AIConfig{})
	persona := &NPCPersona{
		Name: "王五",
	}

	tests := []struct {
		name    string
		userMsg string
		want    string
	}{
		{
			name:    "greeting",
			userMsg: "你好",
			want:    "王五",
		},
		{
			name:    "task related",
			userMsg: "有什么任务吗",
			want:    "忙",
		},
		{
			name:    "shop related",
			userMsg: "我想买东西",
			want:    "准备",
		},
		{
			name:    "default",
			userMsg: "今天天气不错",
			want:    "什么事",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cm.fallbackReply(persona, tt.userMsg)
			if !contains(got, tt.want) {
				t.Errorf("fallbackReply() = %v, should contain %v", got, tt.want)
			}
		})
	}
}

func TestChatManager_cleanReply(t *testing.T) {
	cm := NewChatManager(config.AIConfig{})

	tests := []struct {
		name  string
		reply string
		want  string
	}{
		{
			name:  "with quotes",
			reply: `"你好"`,
			want:  "你好",
		},
		{
			name:  "with whitespace",
			reply: "  你好  ",
			want:  "你好",
		},
		{
			name:  "normal",
			reply: "你好",
			want:  "你好",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cm.cleanReply(tt.reply)
			if got != tt.want {
				t.Errorf("cleanReply() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChatManager_buildMessages(t *testing.T) {
	cm := NewChatManager(config.AIConfig{})
	persona := &NPCPersona{
		Name:        "测试NPC",
		Title:       "测试",
		Description: "测试描述",
	}

	history := []ChatMessage{
		{Role: "user", Content: "你好"},
		{Role: "assistant", Content: "你好，客官！"},
	}

	messages := cm.buildMessages(persona, "玩家等级:10", history, "有什么任务?", 10)

	// 应该有: system + 2 history + 1 user = 4 messages
	if len(messages) != 4 {
		t.Errorf("buildMessages() returned %d messages, want 4", len(messages))
	}

	// 第一条应该是system
	if messages[0].Role != "system" {
		t.Errorf("buildMessages() first message role = %v, want system", messages[0].Role)
	}

	// 最后一条应该是user
	if messages[3].Role != "user" {
		t.Errorf("buildMessages() last message role = %v, want user", messages[3].Role)
	}
}

func TestMessage_Struct(t *testing.T) {
	msg := Message{
		Role:    "user",
		Content: "test content",
	}

	if msg.Role != "user" {
		t.Errorf("Message.Role = %v, want user", msg.Role)
	}
	if msg.Content != "test content" {
		t.Errorf("Message.Content = %v, want test content", msg.Content)
	}
}

func TestNPCPersona_Struct(t *testing.T) {
	persona := NPCPersona{
		Name:         "test",
		Title:        "test title",
		Description:  "test desc",
		SystemPrompt: "test prompt",
	}

	if persona.Name != "test" {
		t.Errorf("NPCPersona.Name = %v, want test", persona.Name)
	}
	if persona.Title != "test title" {
		t.Errorf("NPCPersona.Title = %v, want test title", persona.Title)
	}
}

func TestChatWithNPC_ContextTimeout(t *testing.T) {
	cm := NewChatManager(config.AIConfig{
		APIKey: "test-key",
		Model:  "gpt-4",
	})

	persona := &NPCPersona{
		Name:        "测试NPC",
		Title:       "测试",
		Description: "测试描述",
	}

	// 创建一个已取消的context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// 应该返回降级回复
	reply, err := cm.ChatWithNPC(ctx, persona, "", nil, "你好", 10)
	if err != nil {
		t.Errorf("ChatWithNPC() error = %v", err)
	}
	if reply == "" {
		t.Error("ChatWithNPC() returned empty reply")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && (s[0:len(substr)] == substr || contains(s[1:], substr)))
}

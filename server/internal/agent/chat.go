package agent

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
	"github.com/ddc-111/agentGame/server/internal/config"
)

// ChatManager NPC对话管理器
type ChatManager struct {
	client *openai.Client
	cfg    config.AIConfig
}

// Message 对话消息（通用格式，供memory和handlers使用）
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatMessage 对话消息
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// NPCPersona NPC人设
type NPCPersona struct {
	Name         string `json:"name"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	SystemPrompt string `json:"system_prompt"`
}

// NewChatManager 创建对话管理器
func NewChatManager(cfg config.AIConfig) *ChatManager {
	clientConfig := openai.DefaultConfig(cfg.APIKey)
	if cfg.BaseURL != "" {
		clientConfig.BaseURL = cfg.BaseURL
	}

	client := openai.NewClientWithConfig(clientConfig)

	return &ChatManager{
		client: client,
		cfg:    cfg,
	}
}

// ChatWithNPC 与NPC对话（调用OpenAI）
func (cm *ChatManager) ChatWithNPC(ctx context.Context, persona *NPCPersona, playerContext string, history []ChatMessage, userMsg string, maxMessages int) (string, error) {
	// 构建消息列表
	messages := cm.buildMessages(persona, playerContext, history, userMsg, maxMessages)

	// 创建超时上下文
	timeout := 30 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// 调用OpenAI API
	resp, err := cm.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: cm.cfg.Model,
		Messages: messages,
		Temperature: 0.7,
		MaxTokens:   500,
	})
	if err != nil {
		log.Printf("OpenAI API调用失败: %v", err)
		return cm.fallbackReply(persona, userMsg), nil
	}

	if len(resp.Choices) == 0 {
		return cm.fallbackReply(persona, userMsg), nil
	}

	reply := resp.Choices[0].Message.Content
	reply = cm.cleanReply(reply)

	return reply, nil
}

// buildMessages 构建OpenAI消息列表
func (cm *ChatManager) buildMessages(persona *NPCPersona, playerContext string, history []ChatMessage, userMsg string, maxMessages int) []openai.ChatCompletionMessage {
	var messages []openai.ChatCompletionMessage

	// 系统提示词
	systemPrompt := cm.buildSystemPrompt(persona)
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: systemPrompt,
	})

	// 对话历史（滑动窗口）
	startIdx := 0
	if len(history) > maxMessages {
		startIdx = len(history) - maxMessages
	}
	for i := startIdx; i < len(history); i++ {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    history[i].Role,
			Content: history[i].Content,
		})
	}

	// 当前用户消息（带上下文）
	contextMsg := playerContext + "\n【玩家消息】" + userMsg
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: contextMsg,
	})

	return messages
}

// buildSystemPrompt 构建系统提示词
func (cm *ChatManager) buildSystemPrompt(persona *NPCPersona) string {
	if persona.SystemPrompt != "" {
		return persona.SystemPrompt
	}

	// 默认系统提示词
	return fmt.Sprintf(`你是%s，%s。
%s

请以角色的身份回复玩家的消息，保持角色的性格特点。
回复要简洁自然，符合古风RPG游戏的氛围。
不要打破角色，不要提及自己是AI。`,
		persona.Name,
		persona.Title,
		persona.Description)
}

// fallbackReply 降级回复（API失败时使用）
func (cm *ChatManager) fallbackReply(persona *NPCPersona, userMsg string) string {
	// 简单的关键词回复
	lowerMsg := strings.ToLower(userMsg)

	if strings.Contains(lowerMsg, "你好") || strings.Contains(lowerMsg, "hi") || strings.Contains(lowerMsg, "hello") {
		return fmt.Sprintf("你好，客官！我是%s。", persona.Name)
	}

	if strings.Contains(lowerMsg, "任务") || strings.Contains(lowerMsg, "帮忙") {
		return "抱歉，我现在有点忙，稍后再聊吧。"
	}

	if strings.Contains(lowerMsg, "买") || strings.Contains(lowerMsg, "卖") || strings.Contains(lowerMsg, "商品") {
		return "客官请稍等，我这就去准备。"
	}

	return fmt.Sprintf("嗯，客官有什么事吗？")
}

// cleanReply 清理回复内容
func (cm *ChatManager) cleanReply(reply string) string {
	// 移除可能的引号
	reply = strings.Trim(reply, "\"")
	// 移除多余的空白
	reply = strings.TrimSpace(reply)
	return reply
}

// IsEnabled 检查是否启用
func (cm *ChatManager) IsEnabled() bool {
	return cm.cfg.APIKey != ""
}

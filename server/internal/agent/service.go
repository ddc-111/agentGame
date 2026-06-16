package agent

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

// ChatService NPC对话服务（兼容旧接口）
type ChatService struct {
	client  *openai.Client
	apiKey  string
	baseURL string
}

// NewChatService 创建对话服务
func NewChatService(apiKey, baseURL string) *ChatService {
	if apiKey == "" {
		return &ChatService{apiKey: apiKey, baseURL: baseURL}
	}

	clientConfig := openai.DefaultConfig(apiKey)
	if baseURL != "" {
		clientConfig.BaseURL = baseURL
	}

	client := openai.NewClientWithConfig(clientConfig)

	return &ChatService{
		client:  client,
		apiKey:  apiKey,
		baseURL: baseURL,
	}
}

// Chat 与NPC对话
func (cs *ChatService) Chat(ctx context.Context, systemPrompt string, history []Message, userMsg string, model string, temperature float64, maxTokens int) (string, error) {
	if cs.client == nil {
		return "", fmt.Errorf("OpenAI client not initialized")
	}

	// 构建消息列表
	var messages []openai.ChatCompletionMessage

	// 系统提示词
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: systemPrompt,
	})

	// 对话历史
	for _, msg := range history {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// 用户消息
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userMsg,
	})

	// 设置默认值
	if model == "" {
		model = "gpt-4"
	}
	if temperature == 0 {
		temperature = 0.7
	}
	if maxTokens == 0 {
		maxTokens = 500
	}

	// 创建超时上下文
	timeout := 30 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// 调用OpenAI API
	resp, err := cs.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       model,
		Messages:    messages,
		Temperature: float32(temperature),
		MaxTokens:   maxTokens,
	})
	if err != nil {
		log.Printf("OpenAI API调用失败: %v", err)
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	reply := resp.Choices[0].Message.Content
	reply = strings.TrimSpace(reply)
	reply = strings.Trim(reply, "\"")

	return reply, nil
}

// IsEnabled 检查是否启用
func (cs *ChatService) IsEnabled() bool {
	return cs.apiKey != ""
}

package agent

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// ConversationMemory manages NPC memory of conversations
type ConversationMemory struct {
	mu       sync.RWMutex
	memories map[string]*PlayerMemory // key: "playerID:npcID"
}

// PlayerMemory stores memory for a specific player-NPC pair
type PlayerMemory struct {
	PlayerID    uint
	NPCID       uint
	Messages    []Message
	Summary     string
	LastTalkAt  time.Time
	TalkCount   int
	PlayerName  string
	PlayerLevel int
	Extra       map[string]string
}

// MemoryStore is the global conversation memory store
var MemoryStore = &ConversationMemory{
	memories: make(map[string]*PlayerMemory),
}

// GetMemoryKey generates a unique key for player-NPC pair
func GetMemoryKey(playerID, npcID uint) string {
	return fmt.Sprintf("%d:%d", playerID, npcID)
}

// GetRecentMessages returns last N messages for a player-NPC pair
func (cm *ConversationMemory) GetRecentMessages(playerID, npcID uint, limit int) []Message {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	key := GetMemoryKey(playerID, npcID)
	mem, ok := cm.memories[key]
	if !ok {
		return nil
	}

	if limit <= 0 {
		limit = 10
	}

	start := len(mem.Messages) - limit
	if start < 0 {
		start = 0
	}

	result := make([]Message, len(mem.Messages[start:]))
	copy(result, mem.Messages[start:])
	return result
}

// AddMessage adds a message to the conversation memory
func (cm *ConversationMemory) AddMessage(playerID, npcID uint, role, content string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	key := GetMemoryKey(playerID, npcID)
	mem, ok := cm.memories[key]
	if !ok {
		mem = &PlayerMemory{
			PlayerID: playerID,
			NPCID:    npcID,
			Extra:    make(map[string]string),
		}
		cm.memories[key] = mem
	}

	mem.Messages = append(mem.Messages, Message{
		Role:    role,
		Content: content,
	})
	mem.LastTalkAt = time.Now()
	mem.TalkCount++

	// 保持消息窗口大小
	maxMessages := 20
	if len(mem.Messages) > maxMessages {
		// 保留最近的消息
		mem.Messages = mem.Messages[len(mem.Messages)-maxMessages:]
	}
}

// UpdatePlayerInfo updates player context in memory
func (cm *ConversationMemory) UpdatePlayerInfo(playerID, npcID uint, name string, level int) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	key := GetMemoryKey(playerID, npcID)
	mem, ok := cm.memories[key]
	if !ok {
		mem = &PlayerMemory{
			PlayerID: playerID,
			NPCID:    npcID,
			Extra:    make(map[string]string),
		}
		cm.memories[key] = mem
	}

	mem.PlayerName = name
	mem.PlayerLevel = level
}

// SetExtra stores extra context information
func (cm *ConversationMemory) SetExtra(playerID, npcID uint, key, value string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	memKey := GetMemoryKey(playerID, npcID)
	mem, ok := cm.memories[memKey]
	if !ok {
		mem = &PlayerMemory{
			PlayerID: playerID,
			NPCID:    npcID,
			Extra:    make(map[string]string),
		}
		cm.memories[memKey] = mem
	}

	mem.Extra[key] = value
}

// SummarizeConversation creates a summary of old conversations
func (cm *ConversationMemory) SummarizeConversation(playerID, npcID uint) string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	key := GetMemoryKey(playerID, npcID)
	mem, ok := cm.memories[key]
	if !ok {
		return ""
	}

	if len(mem.Messages) < 5 {
		return mem.Summary
	}

	// 简单摘要：记录对话次数和最后话题
	summary := fmt.Sprintf("与该玩家对话%d次。", mem.TalkCount)
	if len(mem.Messages) > 0 {
		lastMsg := mem.Messages[len(mem.Messages)-1]
		summary += fmt.Sprintf("最后话题: %s", truncateString(lastMsg.Content, 50))
	}

	log.Printf("Memory summary for player %d - NPC %d: %s", playerID, npcID, summary)
	return summary
}

// GetNPCContext builds context from memory
func (cm *ConversationMemory) GetNPCContext(playerID, npcID uint) string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	key := GetMemoryKey(playerID, npcID)
	mem, ok := cm.memories[key]
	if !ok {
		return ""
	}

	context := ""
	if mem.PlayerName != "" {
		context += fmt.Sprintf("玩家姓名: %s, 等级: %d\n", mem.PlayerName, mem.PlayerLevel)
	}
	if mem.TalkCount > 0 {
		context += fmt.Sprintf("对话次数: %d\n", mem.TalkCount)
	}
	if mem.Summary != "" {
		context += fmt.Sprintf("历史摘要: %s\n", mem.Summary)
	}
	for k, v := range mem.Extra {
		context += fmt.Sprintf("%s: %s\n", k, v)
	}

	return context
}

// truncateString truncates a string to maxLen characters
func truncateString(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}

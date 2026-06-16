package agent

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ddc-111/agentGame/server/internal/database/models"
	"github.com/ddc-111/agentGame/server/internal/database/repository"
)

type MemoryStore interface {
	GetRecentMessages(playerID, npcID uint, limit int) []Message
	AddMessage(playerID, npcID uint, role, content string)
	UpdatePlayerInfo(playerID, npcID uint, name string, level int)
	SetExtra(playerID, npcID uint, key, value string)
	SummarizeConversation(playerID, npcID uint) string
	GetNPCContext(playerID, npcID uint) string
}

type ConversationMemory struct {
	mu       sync.RWMutex
	memories map[string]*PlayerMemory
}

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

var DefaultMemoryStore MemoryStore = &ConversationMemory{
	memories: make(map[string]*PlayerMemory),
}

func (cm *ConversationMemory) Init() {
	cm.memories = make(map[string]*PlayerMemory)
}

func GetMemoryKey(playerID, npcID uint) string {
	return fmt.Sprintf("%d:%d", playerID, npcID)
}

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

	maxMessages := 20
	if len(mem.Messages) > maxMessages {
		mem.Messages = mem.Messages[len(mem.Messages)-maxMessages:]
	}
}

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

	summary := fmt.Sprintf("与该玩家对话%d次。", mem.TalkCount)
	if len(mem.Messages) > 0 {
		lastMsg := mem.Messages[len(mem.Messages)-1]
		summary += fmt.Sprintf("最后话题: %s", truncateString(lastMsg.Content, 50))
	}

	log.Printf("Memory summary for player %d - NPC %d: %s", playerID, npcID, summary)
	return summary
}

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

// DBMemoryStore 数据库持久化对话记忆
type DBMemoryStore struct {
	repo      *repository.Repository
	cache     sync.Map
	maxWindow int
}

func NewDBMemoryStore(repo *repository.Repository) *DBMemoryStore {
	return &DBMemoryStore{
		repo:      repo,
		maxWindow: 20,
	}
}

func (s *DBMemoryStore) cacheKey(playerID, npcID uint) string {
	return GetMemoryKey(playerID, npcID)
}

func (s *DBMemoryStore) getOrLoadCache(playerID, npcID uint) *PlayerMemory {
	key := s.cacheKey(playerID, npcID)
	if cached, ok := s.cache.Load(key); ok {
		return cached.(*PlayerMemory)
	}

	mem := &PlayerMemory{
		PlayerID: playerID,
		NPCID:    npcID,
		Extra:    make(map[string]string),
	}

	ctx, err := s.repo.GetConversationContext(playerID, npcID)
	if err == nil {
		mem.PlayerName = ctx.PlayerName
		mem.PlayerLevel = ctx.PlayerLevel
		mem.TalkCount = ctx.TalkCount
		mem.Summary = ctx.Summary
		if ctx.Extra != "" {
			var extra map[string]string
			if json.Unmarshal([]byte(ctx.Extra), &extra) == nil {
				mem.Extra = extra
			}
		}
	}

	convs, err := s.repo.GetConversationsByPair(playerID, npcID, s.maxWindow)
	if err == nil {
		for i := len(convs) - 1; i >= 0; i-- {
			mem.Messages = append(mem.Messages, Message{
				Role:    convs[i].Role,
				Content: convs[i].Content,
			})
		}
	}

	s.cache.Store(key, mem)
	return mem
}

func (s *DBMemoryStore) saveContext(mem *PlayerMemory) {
	extraJSON := "{}"
	if mem.Extra != nil {
		if data, err := json.Marshal(mem.Extra); err == nil {
			extraJSON = string(data)
		}
	}

	ctx := &models.PlayerConversationContext{
		PlayerID:    mem.PlayerID,
		NPCID:       mem.NPCID,
		PlayerName:  mem.PlayerName,
		PlayerLevel: mem.PlayerLevel,
		TalkCount:   mem.TalkCount,
		Summary:     mem.Summary,
		Extra:       extraJSON,
	}

	if err := s.repo.UpsertConversationContext(ctx); err != nil {
		log.Printf("Failed to save conversation context: %v", err)
	}
}

func (s *DBMemoryStore) GetRecentMessages(playerID, npcID uint, limit int) []Message {
	mem := s.getOrLoadCache(playerID, npcID)

	if len(mem.Messages) == 0 {
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

func (s *DBMemoryStore) AddMessage(playerID, npcID uint, role, content string) {
	mem := s.getOrLoadCache(playerID, npcID)

	mem.Messages = append(mem.Messages, Message{
		Role:    role,
		Content: content,
	})
	mem.TalkCount++
	mem.LastTalkAt = time.Now()

	if len(mem.Messages) > s.maxWindow {
		mem.Messages = mem.Messages[len(mem.Messages)-s.maxWindow:]
	}

	conv := &models.Conversation{
		PlayerID: playerID,
		NPCID:    npcID,
		Role:     role,
		Content:  content,
	}
	if err := s.repo.CreateConversation(conv); err != nil {
		log.Printf("Failed to save conversation: %v", err)
	}

	s.saveContext(mem)
}

func (s *DBMemoryStore) UpdatePlayerInfo(playerID, npcID uint, name string, level int) {
	mem := s.getOrLoadCache(playerID, npcID)
	mem.PlayerName = name
	mem.PlayerLevel = level
	s.saveContext(mem)
}

func (s *DBMemoryStore) SetExtra(playerID, npcID uint, key, value string) {
	mem := s.getOrLoadCache(playerID, npcID)
	mem.Extra[key] = value
	s.saveContext(mem)
}

func (s *DBMemoryStore) SummarizeConversation(playerID, npcID uint) string {
	mem := s.getOrLoadCache(playerID, npcID)

	if len(mem.Messages) < 5 {
		return mem.Summary
	}

	summary := fmt.Sprintf("与该玩家对话%d次。", mem.TalkCount)
	if len(mem.Messages) > 0 {
		lastMsg := mem.Messages[len(mem.Messages)-1]
		summary += fmt.Sprintf("最后话题: %s", truncateString(lastMsg.Content, 50))
	}

	mem.Summary = summary
	s.saveContext(mem)

	return summary
}

func (s *DBMemoryStore) GetNPCContext(playerID, npcID uint) string {
	mem := s.getOrLoadCache(playerID, npcID)

	ctx := ""
	if mem.PlayerName != "" {
		ctx += fmt.Sprintf("玩家姓名: %s, 等级: %d\n", mem.PlayerName, mem.PlayerLevel)
	}
	if mem.TalkCount > 0 {
		ctx += fmt.Sprintf("对话次数: %d\n", mem.TalkCount)
	}
	if mem.Summary != "" {
		ctx += fmt.Sprintf("历史摘要: %s\n", mem.Summary)
	}
	for k, v := range mem.Extra {
		ctx += fmt.Sprintf("%s: %s\n", k, v)
	}

	return ctx
}

func (s *DBMemoryStore) InvalidateCache(playerID, npcID uint) {
	s.cache.Delete(s.cacheKey(playerID, npcID))
}

func truncateString(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}

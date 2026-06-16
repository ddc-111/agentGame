package tests

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"

	"github.com/ddc-111/agentGame/server/internal/agent"
	"github.com/ddc-111/agentGame/server/internal/database/models"
	"github.com/ddc-111/agentGame/server/internal/database/repository"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var dbCounter int64

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	id := atomic.AddInt64(&dbCounter, 1)
	tmpDir := t.TempDir()
	dsn := filepath.Join(tmpDir, fmt.Sprintf("test_%d.db", id))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	if err := db.AutoMigrate(
		&models.Conversation{},
		&models.PlayerConversationContext{},
	); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}
	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
		os.Remove(dsn)
	})
	return db
}

func TestDBMemoryStore_AddAndGetMessages(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.New(db)
	store := agent.NewDBMemoryStore(repo)

	store.AddMessage(1, 1, "user", "你好")
	store.AddMessage(1, 1, "assistant", "你好，冒险者！")

	msgs := store.GetRecentMessages(1, 1, 10)
	if len(msgs) != 2 {
		t.Fatalf("Expected 2 messages, got %d", len(msgs))
	}
	if msgs[0].Role != "user" || msgs[0].Content != "你好" {
		t.Errorf("First message mismatch: %+v", msgs[0])
	}
	if msgs[1].Role != "assistant" || msgs[1].Content != "你好，冒险者！" {
		t.Errorf("Second message mismatch: %+v", msgs[1])
	}
}

func TestDBMemoryStore_Persistence(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.New(db)

	store1 := agent.NewDBMemoryStore(repo)
	store1.AddMessage(1, 1, "user", "消息1")
	store1.AddMessage(1, 1, "assistant", "回复1")
	store1.UpdatePlayerInfo(1, 1, "测试玩家", 5)

	store2 := agent.NewDBMemoryStore(repo)
	msgs := store2.GetRecentMessages(1, 1, 10)
	if len(msgs) != 2 {
		t.Fatalf("Expected 2 persisted messages, got %d", len(msgs))
	}
	if msgs[0].Content != "消息1" {
		t.Errorf("Expected persisted message '消息1', got '%s'", msgs[0].Content)
	}

	ctx := store2.GetNPCContext(1, 1)
	if ctx == "" {
		t.Error("Expected non-empty NPC context after reload")
	}
}

func TestDBMemoryStore_PlayerInfo(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.New(db)
	store := agent.NewDBMemoryStore(repo)

	store.UpdatePlayerInfo(1, 1, "张三", 10)
	store.AddMessage(1, 1, "user", "测试")

	ctx := store.GetNPCContext(1, 1)
	if ctx == "" {
		t.Fatal("Expected non-empty context")
	}
}

func TestDBMemoryStore_Extra(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.New(db)
	store := agent.NewDBMemoryStore(repo)

	store.SetExtra(1, 1, "mood", "happy")
	store.SetExtra(1, 1, "last_quest", "quest_001")
	store.AddMessage(1, 1, "user", "你好")

	ctx := store.GetNPCContext(1, 1)
	if ctx == "" {
		t.Fatal("Expected non-empty context with extras")
	}
}

func TestDBMemoryStore_Summarize(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.New(db)
	store := agent.NewDBMemoryStore(repo)

	for i := 0; i < 6; i++ {
		store.AddMessage(1, 1, "user", "测试消息")
	}

	summary := store.SummarizeConversation(1, 1)
	if summary == "" {
		t.Error("Expected non-empty summary after 6 messages")
	}

	store2 := agent.NewDBMemoryStore(repo)
	summary2 := store2.SummarizeConversation(1, 1)
	if summary2 == "" {
		t.Error("Expected persisted summary to be available on new store")
	}
}

func TestDBMemoryStore_WindowLimit(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.New(db)
	store := agent.NewDBMemoryStore(repo)

	for i := 0; i < 25; i++ {
		store.AddMessage(1, 1, "user", "msg")
	}

	msgs := store.GetRecentMessages(1, 1, 100)
	if len(msgs) > 20 {
		t.Errorf("Expected window limit of 20, got %d messages", len(msgs))
	}
}

func TestDBMemoryStore_EmptyMemory(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.New(db)
	store := agent.NewDBMemoryStore(repo)

	msgs := store.GetRecentMessages(999, 999, 10)
	if msgs != nil {
		t.Errorf("Expected nil for non-existent pair, got %v", msgs)
	}

	ctx := store.GetNPCContext(999, 999)
	if ctx != "" {
		t.Errorf("Expected empty context for non-existent pair, got '%s'", ctx)
	}

	summary := store.SummarizeConversation(999, 999)
	if summary != "" {
		t.Errorf("Expected empty summary for non-existent pair, got '%s'", summary)
	}
}

func TestDBMemoryStore_InvalidateCache(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.New(db)
	store := agent.NewDBMemoryStore(repo)

	store.AddMessage(1, 1, "user", "hello")
	msgs := store.GetRecentMessages(1, 1, 10)
	if len(msgs) != 1 {
		t.Fatalf("Expected 1 message, got %d", len(msgs))
	}

	store.InvalidateCache(1, 1)
	msgs = store.GetRecentMessages(1, 1, 10)
	if len(msgs) != 1 {
		t.Errorf("Expected 1 message after cache invalidation (loaded from DB), got %d", len(msgs))
	}
}

func TestConversationMemory_InMemory(t *testing.T) {
	mem := &agent.ConversationMemory{}
	mem.Init()

	mem.AddMessage(1, 1, "user", "你好")
	mem.AddMessage(1, 1, "assistant", "你好！")

	msgs := mem.GetRecentMessages(1, 1, 10)
	if len(msgs) != 2 {
		t.Fatalf("Expected 2 messages, got %d", len(msgs))
	}

	mem.UpdatePlayerInfo(1, 1, "玩家A", 5)
	ctx := mem.GetNPCContext(1, 1)
	if ctx == "" {
		t.Error("Expected non-empty context")
	}
}

func TestDBMemoryStore_DifferentPairs(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.New(db)
	store := agent.NewDBMemoryStore(repo)

	store.AddMessage(1, 1, "user", "对NPC1说")
	store.AddMessage(1, 2, "user", "对NPC2说")
	store.AddMessage(2, 1, "user", "玩家2对NPC1说")

	msgs1 := store.GetRecentMessages(1, 1, 10)
	if len(msgs1) != 1 || msgs1[0].Content != "对NPC1说" {
		t.Errorf("Pair (1,1) mismatch: %+v", msgs1)
	}

	msgs2 := store.GetRecentMessages(1, 2, 10)
	if len(msgs2) != 1 || msgs2[0].Content != "对NPC2说" {
		t.Errorf("Pair (1,2) mismatch: %+v", msgs2)
	}

	msgs3 := store.GetRecentMessages(2, 1, 10)
	if len(msgs3) != 1 || msgs3[0].Content != "玩家2对NPC1说" {
		t.Errorf("Pair (2,1) mismatch: %+v", msgs3)
	}
}

func TestDBMemoryStore_RepositoryContext(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.New(db)

	ctx := &models.PlayerConversationContext{
		PlayerID:    1,
		NPCID:       1,
		PlayerName:  "测试玩家",
		PlayerLevel: 10,
		TalkCount:   5,
		Summary:     "测试摘要",
		Extra:       `{"key":"value"}`,
	}
	if err := repo.CreateConversationContext(context.Background(), ctx); err != nil {
		t.Fatalf("Failed to create context: %v", err)
	}

	got, err := repo.GetConversationContext(context.Background(), 1, 1)
	if err != nil {
		t.Fatalf("Failed to get context: %v", err)
	}
	if got.PlayerName != "测试玩家" {
		t.Errorf("Expected player name '测试玩家', got '%s'", got.PlayerName)
	}
	if got.TalkCount != 5 {
		t.Errorf("Expected talk count 5, got %d", got.TalkCount)
	}
	if got.Summary != "测试摘要" {
		t.Errorf("Expected summary '测试摘要', got '%s'", got.Summary)
	}

	ctx.ID = got.ID
	ctx.PlayerLevel = 20
	ctx.TalkCount = 10
	if err := repo.UpsertConversationContext(context.Background(), ctx); err != nil {
		t.Fatalf("Failed to upsert context: %v", err)
	}

	got2, err := repo.GetConversationContext(context.Background(), 1, 1)
	if err != nil {
		t.Fatalf("Failed to get context after upsert: %v", err)
	}
	if got2.PlayerLevel != 20 {
		t.Errorf("Expected level 20 after upsert, got %d", got2.PlayerLevel)
	}
	if got2.TalkCount != 10 {
		t.Errorf("Expected talk count 10 after upsert, got %d", got2.TalkCount)
	}
}

func TestDBMemoryStore_RepositoryConversations(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.New(db)

	conv1 := &models.Conversation{PlayerID: 1, NPCID: 1, Role: "user", Content: "msg1"}
	conv2 := &models.Conversation{PlayerID: 1, NPCID: 1, Role: "assistant", Content: "reply1"}
	conv3 := &models.Conversation{PlayerID: 1, NPCID: 2, Role: "user", Content: "other"}

	repo.CreateConversation(context.Background(), conv1)
	repo.CreateConversation(context.Background(), conv2)
	repo.CreateConversation(context.Background(), conv3)

	convs, err := repo.GetConversationsByPair(context.Background(), 1, 1, 10)
	if err != nil {
		t.Fatalf("Failed to get conversations: %v", err)
	}
	if len(convs) != 2 {
		t.Errorf("Expected 2 conversations for pair (1,1), got %d", len(convs))
	}

	err = repo.DeleteConversationsByPair(context.Background(), 1, 1)
	if err != nil {
		t.Fatalf("Failed to delete conversations: %v", err)
	}

	convs, _ = repo.GetConversationsByPair(context.Background(), 1, 1, 10)
	if len(convs) != 0 {
		t.Errorf("Expected 0 conversations after delete, got %d", len(convs))
	}

	convs, _ = repo.GetConversationsByPair(context.Background(), 1, 2, 10)
	if len(convs) != 1 {
		t.Errorf("Expected 1 conversation for pair (1,2) still present, got %d", len(convs))
	}
}

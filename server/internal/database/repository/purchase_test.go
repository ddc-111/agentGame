package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/ddc-111/agentGame/server/internal/database/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestPurchaseItemCommitsPlayerAndStockTogether(t *testing.T) {
	repo, db := newPurchaseTestRepository(t)
	player, shopItem := seedPurchaseData(t, db, 100, 5, `{}`)

	result, err := repo.PurchaseItem(context.Background(), player.ID, "shop_test", shopItem.ItemID, 2)
	if err != nil {
		t.Fatalf("PurchaseItem() error = %v", err)
	}
	if result.TotalPrice != 20 {
		t.Fatalf("TotalPrice = %d, want 20", result.TotalPrice)
	}
	if result.Player.Gold != 80 {
		t.Fatalf("player gold = %d, want 80", result.Player.Gold)
	}
	if result.ShopItem.Stock != 3 {
		t.Fatalf("shop stock = %d, want 3", result.ShopItem.Stock)
	}

	var inventory map[string]int
	if err := json.Unmarshal([]byte(result.Player.Items), &inventory); err != nil {
		t.Fatalf("decode inventory: %v", err)
	}
	if inventory[fmt.Sprint(shopItem.ItemID)] != 2 {
		t.Fatalf("inventory = %#v, want purchased item count 2", inventory)
	}
}

func TestPurchaseItemRollsBackWhenInventoryIsInvalid(t *testing.T) {
	repo, db := newPurchaseTestRepository(t)
	player, shopItem := seedPurchaseData(t, db, 100, 5, `{invalid`)

	_, err := repo.PurchaseItem(context.Background(), player.ID, "shop_test", shopItem.ItemID, 1)
	if err == nil || !strings.Contains(err.Error(), "decode player inventory") {
		t.Fatalf("PurchaseItem() error = %v, want inventory decode error", err)
	}

	assertPurchaseState(t, db, player.ID, shopItem.ID, 100, 5, `{invalid`)
}

func TestPurchaseItemRollsBackStockWhenPlayerUpdateFails(t *testing.T) {
	repo, db := newPurchaseTestRepository(t)
	player, shopItem := seedPurchaseData(t, db, 100, 5, `{}`)

	if err := db.Exec(`CREATE TRIGGER fail_player_purchase BEFORE UPDATE ON players BEGIN SELECT RAISE(ABORT, 'forced player update failure'); END;`).Error; err != nil {
		t.Fatalf("create trigger: %v", err)
	}

	_, err := repo.PurchaseItem(context.Background(), player.ID, "shop_test", shopItem.ItemID, 1)
	if err == nil || !strings.Contains(err.Error(), "update player purchase") {
		t.Fatalf("PurchaseItem() error = %v, want forced player update error", err)
	}

	assertPurchaseState(t, db, player.ID, shopItem.ID, 100, 5, `{}`)
}

func TestPurchaseItemRejectsInsufficientStockWithoutCharging(t *testing.T) {
	repo, db := newPurchaseTestRepository(t)
	player, shopItem := seedPurchaseData(t, db, 100, 1, `{}`)

	_, err := repo.PurchaseItem(context.Background(), player.ID, "shop_test", shopItem.ItemID, 2)
	if !errors.Is(err, ErrInsufficientStock) {
		t.Fatalf("PurchaseItem() error = %v, want ErrInsufficientStock", err)
	}

	assertPurchaseState(t, db, player.ID, shopItem.ID, 100, 1, `{}`)
}

func TestPurchaseItemRejectsInsufficientGoldWithoutReducingStock(t *testing.T) {
	repo, db := newPurchaseTestRepository(t)
	player, shopItem := seedPurchaseData(t, db, 5, 5, `{}`)

	_, err := repo.PurchaseItem(context.Background(), player.ID, "shop_test", shopItem.ItemID, 1)
	if !errors.Is(err, ErrInsufficientGold) {
		t.Fatalf("PurchaseItem() error = %v, want ErrInsufficientGold", err)
	}

	assertPurchaseState(t, db, player.ID, shopItem.ID, 5, 5, `{}`)
}

func newPurchaseTestRepository(t *testing.T) (*Repository, *gorm.DB) {
	t.Helper()
	dsnName := strings.NewReplacer("/", "_", " ", "_").Replace(t.Name())
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared&_pragma=busy_timeout(5000)", dsnName)
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.Player{}, &models.Shop{}, &models.Item{}, &models.ShopItem{}); err != nil {
		t.Fatalf("AutoMigrate() error = %v", err)
	}
	return New(db), db
}

func seedPurchaseData(t *testing.T, db *gorm.DB, gold, stock int, inventory string) (*models.Player, *models.ShopItem) {
	t.Helper()
	item := &models.Item{Name: "Potion", Code: "item_potion"}
	if err := db.Create(item).Error; err != nil {
		t.Fatalf("create item: %v", err)
	}
	shop := &models.Shop{Name: "Test Shop", Code: "shop_test"}
	if err := db.Create(shop).Error; err != nil {
		t.Fatalf("create shop: %v", err)
	}
	shopItem := &models.ShopItem{ShopID: shop.ID, ItemID: item.ID, Price: 10, Stock: stock}
	if err := db.Create(shopItem).Error; err != nil {
		t.Fatalf("create shop item: %v", err)
	}
	player := &models.Player{Name: "Tester", Account: "tester", Gold: gold, Items: inventory}
	if err := db.Create(player).Error; err != nil {
		t.Fatalf("create player: %v", err)
	}
	return player, shopItem
}

func assertPurchaseState(t *testing.T, db *gorm.DB, playerID, shopItemID uint, expectedGold, expectedStock int, expectedItems string) {
	t.Helper()
	var player models.Player
	if err := db.First(&player, playerID).Error; err != nil {
		t.Fatalf("reload player: %v", err)
	}
	var shopItem models.ShopItem
	if err := db.First(&shopItem, shopItemID).Error; err != nil {
		t.Fatalf("reload shop item: %v", err)
	}
	if player.Gold != expectedGold || player.Items != expectedItems || shopItem.Stock != expectedStock {
		t.Fatalf("state = gold:%d items:%q stock:%d, want gold:%d items:%q stock:%d", player.Gold, player.Items, shopItem.Stock, expectedGold, expectedItems, expectedStock)
	}
}

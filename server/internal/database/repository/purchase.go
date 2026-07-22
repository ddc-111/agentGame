package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/ddc-111/agentGame/server/internal/database/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrPlayerNotFound    = errors.New("player not found")
	ErrShopNotFound      = errors.New("shop not found")
	ErrShopItemNotFound  = errors.New("item not found in shop")
	ErrInsufficientStock = errors.New("not enough stock")
	ErrInsufficientGold  = errors.New("not enough gold")
	ErrInvalidPurchase   = errors.New("invalid purchase")
)

type PurchaseResult struct {
	Player     *models.Player
	Shop       *models.Shop
	ShopItem   *models.ShopItem
	TotalPrice int
}

func (r *Repository) PurchaseItem(ctx context.Context, playerID uint, shopCode string, itemID uint, count int) (*PurchaseResult, error) {
	if playerID == 0 || itemID == 0 || strings.TrimSpace(shopCode) == "" || count <= 0 {
		return nil, ErrInvalidPurchase
	}

	var result *PurchaseResult
	operation := func() error {
		result = nil
		return r.Transaction(ctx, func(tx *Repository) error {
			purchase, err := tx.purchaseItem(ctx, playerID, shopCode, itemID, count)
			if err != nil {
				return err
			}
			result = purchase
			return nil
		})
	}

	var err error
	attempts := 1
	if r.db.Dialector.Name() == "sqlite" {
		attempts = 4
	}
	for attempt := 0; attempt < attempts; attempt++ {
		err = operation()
		if err == nil || !isRetryableDatabaseLock(err) {
			break
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(time.Duration(attempt+1) * 10 * time.Millisecond):
		}
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Repository) purchaseItem(ctx context.Context, playerID uint, shopCode string, itemID uint, count int) (*PurchaseResult, error) {
	var shop models.Shop
	if err := r.lockForUpdate(r.db.WithContext(ctx).Where("code = ?", shopCode)).First(&shop).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrShopNotFound
		}
		return nil, fmt.Errorf("load shop: %w", err)
	}

	var shopItem models.ShopItem
	itemQuery := r.db.WithContext(ctx).
		Preload("Item").
		Where("shop_id = ? AND item_id = ?", shop.ID, itemID)
	if err := r.lockForUpdate(itemQuery).First(&shopItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrShopItemNotFound
		}
		return nil, fmt.Errorf("load shop item: %w", err)
	}
	if shopItem.Stock < count {
		return nil, ErrInsufficientStock
	}
	if shopItem.Price < 0 || shopItem.Price > math.MaxInt/count {
		return nil, ErrInvalidPurchase
	}
	totalPrice := shopItem.Price * count

	var player models.Player
	if err := r.lockForUpdate(r.db.WithContext(ctx).Where("id = ?", playerID)).First(&player).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPlayerNotFound
		}
		return nil, fmt.Errorf("load player: %w", err)
	}
	if player.Gold < totalPrice {
		return nil, ErrInsufficientGold
	}

	items := make(map[string]int)
	if strings.TrimSpace(player.Items) != "" {
		if err := json.Unmarshal([]byte(player.Items), &items); err != nil {
			return nil, fmt.Errorf("decode player inventory: %w", err)
		}
	}
	if items == nil {
		items = make(map[string]int)
	}
	itemKey := strconv.FormatUint(uint64(itemID), 10)
	if items[itemKey] > math.MaxInt-count {
		return nil, ErrInvalidPurchase
	}
	items[itemKey] += count
	itemsJSON, err := json.Marshal(items)
	if err != nil {
		return nil, fmt.Errorf("encode player inventory: %w", err)
	}

	stockUpdate := r.db.WithContext(ctx).
		Model(&models.ShopItem{}).
		Where("id = ? AND stock >= ?", shopItem.ID, count).
		UpdateColumn("stock", gorm.Expr("stock - ?", count))
	if stockUpdate.Error != nil {
		return nil, fmt.Errorf("decrement stock: %w", stockUpdate.Error)
	}
	if stockUpdate.RowsAffected != 1 {
		return nil, ErrInsufficientStock
	}

	playerUpdate := r.db.WithContext(ctx).
		Model(&models.Player{}).
		Where("id = ? AND gold >= ?", player.ID, totalPrice).
		Updates(map[string]interface{}{
			"gold":  gorm.Expr("gold - ?", totalPrice),
			"items": string(itemsJSON),
		})
	if playerUpdate.Error != nil {
		return nil, fmt.Errorf("update player purchase: %w", playerUpdate.Error)
	}
	if playerUpdate.RowsAffected != 1 {
		return nil, ErrInsufficientGold
	}

	if err := r.db.WithContext(ctx).First(&player, player.ID).Error; err != nil {
		return nil, fmt.Errorf("reload player: %w", err)
	}
	if err := r.db.WithContext(ctx).Preload("Item").First(&shopItem, shopItem.ID).Error; err != nil {
		return nil, fmt.Errorf("reload shop item: %w", err)
	}

	return &PurchaseResult{
		Player:     &player,
		Shop:       &shop,
		ShopItem:   &shopItem,
		TotalPrice: totalPrice,
	}, nil
}

func (r *Repository) lockForUpdate(query *gorm.DB) *gorm.DB {
	if r.db.Dialector.Name() == "mysql" {
		return query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	return query
}

func isRetryableDatabaseLock(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "database is locked") ||
		strings.Contains(message, "database table is locked") ||
		strings.Contains(message, "sqlite_busy") ||
		strings.Contains(message, "sqlite_locked")
}

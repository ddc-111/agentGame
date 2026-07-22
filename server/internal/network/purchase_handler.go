package network

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ddc-111/agentGame/server/internal/database/repository"
	"github.com/gin-gonic/gin"
)

func (s *Server) handlePurchaseItem(c *gin.Context) {
	ctx := c.Request.Context()
	var req struct {
		PlayerID uint   `json:"player_id"`
		ShopCode string `json:"shop_code"`
		ItemID   uint   `json:"item_id"`
		Count    int    `json:"count"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}

	errs := mergeErrors(
		validatePositiveInt("player_id", req.PlayerID),
		validateRequired(map[string]interface{}{"shop_code": req.ShopCode}),
		validatePositiveInt("item_id", req.ItemID),
	)
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}
	if req.Count <= 0 {
		req.Count = 1
	}
	if req.Count > 1000 {
		respondError(c, http.StatusBadRequest, BadRequest("count must be between 1 and 1000"))
		return
	}

	purchase, err := s.repo.PurchaseItem(ctx, req.PlayerID, req.ShopCode, req.ItemID, req.Count)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrPlayerNotFound):
			respondError(c, http.StatusNotFound, NotFound("Player"))
		case errors.Is(err, repository.ErrShopNotFound):
			respondError(c, http.StatusNotFound, NotFound("Shop"))
		case errors.Is(err, repository.ErrShopItemNotFound):
			respondError(c, http.StatusBadRequest, BadRequest("Item not found in shop"))
		case errors.Is(err, repository.ErrInsufficientStock):
			respondError(c, http.StatusConflict, BadRequest("Not enough stock"))
		case errors.Is(err, repository.ErrInsufficientGold):
			respondError(c, http.StatusConflict, BadRequest("Not enough gold"))
		case errors.Is(err, repository.ErrInvalidPurchase):
			respondError(c, http.StatusBadRequest, BadRequest("Invalid purchase"))
		default:
			respondInternalError(c, err)
		}
		return
	}

	if purchase.Shop.OwnerNPC != "" {
		npc, err := s.repo.GetNPCByCode(ctx, purchase.Shop.OwnerNPC)
		if err == nil {
			behavior := s.behaviorStore.GetOrCreateCopy(npc.Code, npc.Schedule)
			s.behaviorMgr.ReactToPlayer(behavior, req.PlayerID, "gift")
			s.behaviorStore.Set(npc.Code, behavior)
			scenes, _ := s.repo.GetScenesByNPCID(ctx, npc.ID)
			if len(scenes) > 0 {
				s.BroadcastNPCState(npc.ID, npc.Code, npc.Name, scenes[0].Code, behavior.State, 0, 0)
			}
		}
	}

	equipment := map[string]interface{}{
		"weapon_id": nil,
		"armor_id":  nil,
	}
	if purchase.Player.Equipment != "" {
		var equip struct {
			WeaponID uint `json:"weapon_id"`
			ArmorID  uint `json:"armor_id"`
		}
		if err := json.Unmarshal([]byte(purchase.Player.Equipment), &equip); err == nil {
			if equip.WeaponID > 0 {
				equipment["weapon_id"] = equip.WeaponID
			}
			if equip.ArmorID > 0 {
				equipment["armor_id"] = equip.ArmorID
			}
		}
	}

	itemName := ""
	if purchase.ShopItem.Item != nil {
		itemName = purchase.ShopItem.Item.Name
	}
	c.JSON(http.StatusOK, gin.H{
		"message":     "购买成功",
		"gold":        purchase.Player.Gold,
		"items":       purchase.Player.Items,
		"equipment":   equipment,
		"item_name":   itemName,
		"total_price": purchase.TotalPrice,
	})
}

package network

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/game"
)

func (s *Server) registerInventoryRoutes(api *gin.RouterGroup) {
	api.GET("/inventory/:player_id", s.handleGetInventory)
	api.POST("/inventory/equip", s.handleEquipItem)
	api.POST("/inventory/unequip", s.handleUnequipItem)
	api.POST("/inventory/use", s.handleUseItem)
}

func (s *Server) handleGetInventory(c *gin.Context) {
	ctx := c.Request.Context()
	playerID, ok := parseID(c, "player_id")
	if !ok {
		return
	}

	player, err := s.repo.GetPlayerByID(ctx, playerID)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Player"))
		return
	}

	im := game.NewInventoryManager()
	items, _ := im.GetItems(player.Items)

	type ItemDetail struct {
		game.InventoryItem
		Name        string `json:"name"`
		Code        string `json:"code"`
		Category    string `json:"category"`
		Description string `json:"description"`
		Effect      string `json:"effect"`
	}

	var itemDetails []ItemDetail
	for _, invItem := range items {
		item, err := s.repo.GetItemByID(ctx, invItem.ItemID)
		if err != nil {
			continue
		}
		itemDetails = append(itemDetails, ItemDetail{
			InventoryItem: invItem,
			Name:          item.Name,
			Code:          item.Code,
			Category:      item.Category,
			Description:   item.Description,
			Effect:        item.Effect,
		})
	}

	equipment := map[string]interface{}{
		"weapon_id": nil,
		"armor_id":  nil,
	}
	if player.Equipment != "" {
		var equip game.Equipment
		if err := json.Unmarshal([]byte(player.Equipment), &equip); err == nil {
			if equip.WeaponID > 0 {
				equipment["weapon_id"] = equip.WeaponID
			}
			if equip.ArmorID > 0 {
				equipment["armor_id"] = equip.ArmorID
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"player_id": uint(playerID),
		"items":     itemDetails,
		"equipment": equipment,
		"gold":      player.Gold,
		"stats": gin.H{
			"attack":  player.Attack,
			"defense": player.Defense,
			"hp":      player.HP,
			"mp":      player.MP,
		},
	})
}

func (s *Server) handleEquipItem(c *gin.Context) {
	ctx := c.Request.Context()
	var req struct {
		PlayerID uint `json:"player_id"`
		ItemID   uint `json:"item_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}

	errs := mergeErrors(
		validatePositiveInt("player_id", req.PlayerID),
		validatePositiveInt("item_id", req.ItemID),
	)
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}

	player, err := s.repo.GetPlayerByID(ctx, req.PlayerID)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Player"))
		return
	}

	item, err := s.repo.GetItemByID(ctx, req.ItemID)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Item"))
		return
	}

	im := game.NewInventoryManager()
	equipJSON := player.Equipment
	if equipJSON == "" {
		equipJSON = "{}"
	}
	newEquipJSON, newItemsJSON, err := im.EquipItem(equipJSON, player.Items, req.ItemID, item.Category)
	if err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}

	player.Items = newItemsJSON
	player.Equipment = newEquipJSON

	if err := s.repo.UpdatePlayer(ctx, player); err != nil {
		respondInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "装备成功",
		"items":     player.Items,
		"item_name": item.Name,
	})
}

func (s *Server) handleUnequipItem(c *gin.Context) {
	ctx := c.Request.Context()
	var req struct {
		PlayerID uint   `json:"player_id"`
		Slot     string `json:"slot"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}

	errs := mergeErrors(
		validatePositiveInt("player_id", req.PlayerID),
		validateRequired(map[string]interface{}{"slot": req.Slot}),
		validateStringIn("slot", req.Slot, []string{"weapon", "armor"}),
	)
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}

	player, err := s.repo.GetPlayerByID(ctx, req.PlayerID)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Player"))
		return
	}

	im := game.NewInventoryManager()
	equipJSON := player.Equipment
	if equipJSON == "" {
		equipJSON = "{}"
	}
	newEquipJSON, newItemsJSON, err := im.UnequipItem(equipJSON, player.Items, req.Slot)
	if err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}

	player.Items = newItemsJSON
	player.Equipment = newEquipJSON

	if err := s.repo.UpdatePlayer(ctx, player); err != nil {
		respondInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "卸下装备成功",
		"items":   player.Items,
	})
}

func (s *Server) handleUseItem(c *gin.Context) {
	ctx := c.Request.Context()
	var req struct {
		PlayerID uint `json:"player_id"`
		ItemID   uint `json:"item_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}

	errs := mergeErrors(
		validatePositiveInt("player_id", req.PlayerID),
		validatePositiveInt("item_id", req.ItemID),
	)
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}

	player, err := s.repo.GetPlayerByID(ctx, req.PlayerID)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Player"))
		return
	}

	item, err := s.repo.GetItemByID(ctx, req.ItemID)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Item"))
		return
	}

	var effect map[string]int
	err = json.Unmarshal([]byte(item.Effect), &effect)
	if err != nil {
		respondError(c, http.StatusBadRequest, BadRequest("Invalid item effect"))
		return
	}

	im := game.NewInventoryManager()
	newItemsJSON, appliedEffect, err := im.UseItem(player.Items, req.ItemID, effect)
	if err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}

	if hp, ok := appliedEffect["hp"]; ok {
		player.HP += hp
	}
	if mp, ok := appliedEffect["mp"]; ok {
		player.MP += mp
	}

	player.Items = newItemsJSON
	if err := s.repo.UpdatePlayer(ctx, player); err != nil {
		respondInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "使用道具成功",
		"item_name": item.Name,
		"effect":    appliedEffect,
		"player": gin.H{
			"hp":    player.HP,
			"mp":    player.MP,
			"items": player.Items,
		},
	})
}

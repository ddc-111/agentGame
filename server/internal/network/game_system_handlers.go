package network

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/database/models"
	"github.com/ddc-111/agentGame/server/internal/game"
)

// ==================== 战斗系统API ====================

// StartCombat 开始战斗
func (s *Server) handleStartCombat(c *gin.Context) {
	var req struct {
		PlayerID  uint   `json:"player_id"`
		EnemyType string `json:"enemy_type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}

	// 获取玩家信息
	player, err := s.repo.GetPlayerByID(req.PlayerID)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Player"))
		return
	}

	// 计算装备加成后的总属性
	im := game.NewInventoryManager()
	equipStats := game.EquipmentStats{}
	if player.Equipment != "" && player.Equipment != "{}" {
		var equip game.Equipment
		if err := json.Unmarshal([]byte(player.Equipment), &equip); err == nil {
			if equip.WeaponID > 0 {
				weapon, err := s.repo.GetItemByID(equip.WeaponID)
				if err == nil {
					var effect map[string]int
					json.Unmarshal([]byte(weapon.Effect), &effect)
					equipStats.Attack += effect["attack"]
					equipStats.Defense += effect["defense"]
					equipStats.HP += effect["hp"]
					equipStats.MP += effect["mp"]
				}
			}
			if equip.ArmorID > 0 {
				armor, err := s.repo.GetItemByID(equip.ArmorID)
				if err == nil {
					var effect map[string]int
					json.Unmarshal([]byte(armor.Effect), &effect)
					equipStats.Attack += effect["attack"]
					equipStats.Defense += effect["defense"]
					equipStats.HP += effect["hp"]
					equipStats.MP += effect["mp"]
				}
			}
		}
	}
	playerStats := im.CalculateStats(player.Attack, player.Defense, player.HP, player.MP, equipStats)

	// 创建战斗系统
	combatSys := game.NewCombatSystem()
	state := combatSys.StartCombat(req.PlayerID, req.EnemyType, playerStats.TotalHP, player.MP)
	state.PlayerDef = playerStats.TotalDefense

	c.JSON(http.StatusOK, gin.H{
		"data":    state,
		"message": "战斗开始",
	})
}

// CombatAction 战斗行动
func (s *Server) handleCombatAction(c *gin.Context) {
	var req struct {
		PlayerID uint   `json:"player_id"`
		Action   string `json:"action"` // attack, skill, item, flee
		ItemID   uint   `json:"item_id,omitempty"`
		SkillID  uint   `json:"skill_id,omitempty"`
		State    *game.CombatState `json:"state"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}

	if req.State == nil {
		respondError(c, http.StatusBadRequest, BadRequest("Combat state is required"))
		return
	}

	// 获取玩家信息
	player, err := s.repo.GetPlayerByID(req.PlayerID)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Player"))
		return
	}

	// 计算装备加成后的攻击力
	im := game.NewInventoryManager()
	totalAttack := player.Attack
	if player.Equipment != "" && player.Equipment != "{}" {
		var equip game.Equipment
		if err := json.Unmarshal([]byte(player.Equipment), &equip); err == nil {
			if equip.WeaponID > 0 {
				weapon, err := s.repo.GetItemByID(equip.WeaponID)
				if err == nil {
					var effect map[string]int
					json.Unmarshal([]byte(weapon.Effect), &effect)
					totalAttack += effect["attack"]
				}
			}
			if equip.ArmorID > 0 {
				armor, err := s.repo.GetItemByID(equip.ArmorID)
				if err == nil {
					var effect map[string]int
					json.Unmarshal([]byte(armor.Effect), &effect)
					totalAttack += effect["attack"]
				}
			}
		}
	}

	combatSys := game.NewCombatSystem()
	var newState *game.CombatState

	switch req.Action {
	case "attack":
		newState = combatSys.Attack(req.State, totalAttack)

	case "skill":
		if req.SkillID == 0 {
			respondError(c, http.StatusBadRequest, BadRequest("Skill ID is required"))
			return
		}
		skillModel, err := s.repo.GetSkillByID(req.SkillID)
		if err != nil {
			respondError(c, http.StatusNotFound, NotFound("Skill"))
			return
		}
		skill := &game.Skill{
			ID:       skillModel.ID,
			Name:     skillModel.Name,
			Code:     skillModel.Code,
			Type:     skillModel.Type,
			MPCost:   skillModel.MPCost,
			Damage:   skillModel.Damage,
			Heal:     skillModel.Heal,
			Cooldown: skillModel.Cooldown,
			Level:    skillModel.Level,
			Effect:   skillModel.Effect,
		}
		sm := game.NewSkillManager()
		newState, _, err = sm.UseSkill(skill, req.State, totalAttack)
		if err != nil {
			respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
			return
		}

	case "item":
		// 获取道具效果
		item, err := s.repo.GetItemByID(req.ItemID)
		if err != nil {
			respondError(c, http.StatusBadRequest, BadRequest("Item not found"))
			return
		}

		// 解析道具效果
		var effect map[string]int
		if err := json.Unmarshal([]byte(item.Effect), &effect); err != nil {
			respondError(c, http.StatusBadRequest, BadRequest("Invalid item effect"))
			return
		}

		newState = combatSys.UseItem(req.State, effect)

		// 消耗道具
		newItemsJSON, _, err := im.UseItem(player.Items, req.ItemID, effect)
		if err == nil {
			player.Items = newItemsJSON
			s.repo.UpdatePlayer(player)
		}

	case "flee":
		success, state := combatSys.Flee(req.State, player.Level)
		newState = state
		if success {
			c.JSON(http.StatusOK, gin.H{
				"data":    newState,
				"message": "逃跑成功",
			})
			return
		}

	default:
		respondError(c, http.StatusBadRequest, BadRequest("Invalid action"))
		return
	}

	// 检查战斗是否结束
	if !newState.IsActive {
		// 战斗结束，处理奖励
		rewards := combatSys.GetRewards(newState)
		if newState.PlayerHP > 0 {
			// 胜利，发放奖励
			player.Exp += rewards.Exp
			player.Gold += rewards.Gold

			// 检查升级（支持连续升级）
			levelsGained := 0
			for {
				expNeeded := player.Level * 100
				if player.Exp < expNeeded {
					break
				}
				player.Exp -= expNeeded
				player.Level++
				player.HP += 10
				player.MP += 5
				player.Attack += 2
				player.Defense += 1
				levelsGained++
			}

			// Sync HP/MP from combat state (combat damage applied)
			player.HP = newState.PlayerHP
			player.MP = newState.PlayerMP

			s.repo.UpdatePlayer(player)

			// Add level-up info to response
			if levelsGained > 0 {
				newState.Log = append(newState.Log, fmt.Sprintf("恭喜！升级到 %d 级！", player.Level))
			}
		} else {
			// Defeat - still sync HP/MP
			player.HP = newState.PlayerHP
			player.MP = newState.PlayerMP
			s.repo.UpdatePlayer(player)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    newState,
		"message": "战斗更新",
	})
}

// ==================== 背包系统API ====================

// GetInventory 获取玩家背包
func (s *Server) handleGetInventory(c *gin.Context) {
	playerID, _ := strconv.ParseUint(c.Param("player_id"), 10, 32)

	player, err := s.repo.GetPlayerByID(uint(playerID))
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Player"))
		return
	}

	im := game.NewInventoryManager()
	items, _ := im.GetItems(player.Items)

	// 获取道具详情
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
		item, err := s.repo.GetItemByID(invItem.ItemID)
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

	// 解析装备
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
		"player_id": playerID,
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

// EquipItem 装备道具
func (s *Server) handleEquipItem(c *gin.Context) {
	var req struct {
		PlayerID uint   `json:"player_id"`
		ItemID   uint   `json:"item_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}

	player, err := s.repo.GetPlayerByID(req.PlayerID)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Player"))
		return
	}

	item, err := s.repo.GetItemByID(req.ItemID)
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

	if err := s.repo.UpdatePlayer(player); err != nil {
		respondInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "装备成功",
		"items":    player.Items,
		"item_name": item.Name,
	})
}

// UnequipItem 卸下装备
func (s *Server) handleUnequipItem(c *gin.Context) {
	var req struct {
		PlayerID uint   `json:"player_id"`
		Slot     string `json:"slot"` // weapon, armor
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}

	player, err := s.repo.GetPlayerByID(req.PlayerID)
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

	if err := s.repo.UpdatePlayer(player); err != nil {
		respondInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "卸下装备成功",
		"items":   player.Items,
	})
}

// UseItem 使用道具
func (s *Server) handleUseItem(c *gin.Context) {
	var req struct {
		PlayerID uint `json:"player_id"`
		ItemID   uint `json:"item_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}

	player, err := s.repo.GetPlayerByID(req.PlayerID)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Player"))
		return
	}

	item, err := s.repo.GetItemByID(req.ItemID)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Item"))
		return
	}

	// 解析道具效果
	var effect map[string]int
	if err := json.Unmarshal([]byte(item.Effect), &effect); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest("Invalid item effect"))
		return
	}

	im := game.NewInventoryManager()
	newItemsJSON, appliedEffect, err := im.UseItem(player.Items, req.ItemID, effect)
	if err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}

	// 应用效果
	if hp, ok := appliedEffect["hp"]; ok {
		player.HP += hp
	}
	if mp, ok := appliedEffect["mp"]; ok {
		player.MP += mp
	}

	player.Items = newItemsJSON
	if err := s.repo.UpdatePlayer(player); err != nil {
		respondInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "使用道具成功",
		"item_name": item.Name,
		"effect":   appliedEffect,
		"player": gin.H{
			"hp":    player.HP,
			"mp":    player.MP,
			"items": player.Items,
		},
	})
}

// ==================== 存档系统API ====================

// SaveGame 保存游戏
func (s *Server) handleSaveGame(c *gin.Context) {
	var req struct {
		PlayerID uint   `json:"player_id"`
		Slot     int    `json:"slot"`
		Name     string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}

	// 验证槽位
	sgm := game.NewSaveGameManager()
	if err := sgm.ValidateSlot(req.Slot); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}

	// 获取玩家信息
	player, err := s.repo.GetPlayerByID(req.PlayerID)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Player"))
		return
	}

	// 创建快照
	snapshot := sgm.CreateSnapshot(
		player.Name,
		player.Level,
		player.Exp,
		player.Gold,
		player.HP,
		player.MP,
		player.Attack,
		player.Defense,
		player.SceneID,
		player.PosX,
		player.PosY,
		player.Items,
		player.Equipment,
	)

	snapshotJSON, err := sgm.SerializeSnapshot(snapshot)
	if err != nil {
		respondInternalError(c, err)
		return
	}

	// 检查是否已存在存档
	existingSave, _ := s.repo.GetSaveGame(req.PlayerID, req.Slot)

	saveName := req.Name
	if saveName == "" {
		saveName = sgm.FormatSaveName(req.Slot, player.Name)
	}

	if existingSave != nil && existingSave.ID > 0 {
		// 更新存档
		existingSave.Name = saveName
		existingSave.Snapshot = snapshotJSON
		if err := s.repo.UpdateSaveGame(existingSave); err != nil {
			respondInternalError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "存档更新成功",
			"data":    existingSave,
		})
	} else {
		// 创建新存档
		save := &models.SaveGame{
			PlayerID: req.PlayerID,
			Slot:     req.Slot,
			Name:     saveName,
			Snapshot: snapshotJSON,
		}
		if err := s.repo.CreateSaveGame(save); err != nil {
			respondInternalError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "存档保存成功",
			"data":    save,
		})
	}
}

// GetSaves 获取玩家存档列表
func (s *Server) handleGetSaves(c *gin.Context) {
	playerID, _ := strconv.ParseUint(c.Param("player_id"), 10, 32)

	saves, err := s.repo.GetSaveGames(uint(playerID))
	if err != nil {
		respondInternalError(c, err)
		return
	}

	sgm := game.NewSaveGameManager()
	var saveInfos []game.SaveSlotInfo

	// 创建10个槽位的信息
	for i := 0; i <= 10; i++ {
		info := game.SaveSlotInfo{
			Slot:    i,
			IsEmpty: true,
		}

		// 查找对应存档
		for _, save := range saves {
			if save.Slot == i {
				snapshot, _ := sgm.DeserializeSnapshot(save.Snapshot)
				info.SaveID = save.ID
				info.Name = save.Name
				info.IsEmpty = false
				info.CreatedAt = save.CreatedAt
				if snapshot != nil {
					info.Level = snapshot.Level
					info.SceneID = snapshot.SceneID
				}
				break
			}
		}

		saveInfos = append(saveInfos, info)
	}

	c.JSON(http.StatusOK, gin.H{
		"player_id": playerID,
		"saves":     saveInfos,
	})
}

// LoadGame 加载游戏
func (s *Server) handleLoadGame(c *gin.Context) {
	saveID, _ := strconv.ParseUint(c.Param("save_id"), 10, 32)

	var req struct {
		PlayerID uint `json:"player_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}

	// 获取存档
	saves, err := s.repo.GetSaveGames(req.PlayerID)
	if err != nil {
		respondInternalError(c, err)
		return
	}

	var targetSave *models.SaveGame
	for _, save := range saves {
		if save.ID == uint(saveID) {
			targetSave = &save
			break
		}
	}

	if targetSave == nil {
		respondError(c, http.StatusNotFound, NotFound("Save"))
		return
	}

	// 反序列化快照
	sgm := game.NewSaveGameManager()
	snapshot, err := sgm.DeserializeSnapshot(targetSave.Snapshot)
	if err != nil {
		respondInternalError(c, err)
		return
	}

	// 获取玩家
	player, err := s.repo.GetPlayerByID(req.PlayerID)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Player"))
		return
	}

	// 恢复玩家状态
	player.Name = snapshot.Name
	player.Level = snapshot.Level
	player.Exp = snapshot.Exp
	player.Gold = snapshot.Gold
	player.HP = snapshot.HP
	player.MP = snapshot.MP
	player.Attack = snapshot.Attack
	player.Defense = snapshot.Defense
	player.SceneID = snapshot.SceneID
	player.PosX = snapshot.PosX
	player.PosY = snapshot.PosY
	player.Items = snapshot.Items
	player.Equipment = snapshot.Equipment

	if err := s.repo.UpdatePlayer(player); err != nil {
		respondInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "读档成功",
		"data":    player,
	})
}

// ==================== 技能系统API ====================

// GetSkills 获取所有技能
func (s *Server) handleGetSkills(c *gin.Context) {
	skills, err := s.repo.GetSkills()
	if err != nil {
		respondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": skills})
}

// UseSkill 使用技能
func (s *Server) handleUseSkill(c *gin.Context) {
	var req struct {
		PlayerID uint   `json:"player_id"`
		SkillID  uint   `json:"skill_id"`
		State    *game.CombatState `json:"state"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}

	if req.State == nil {
		respondError(c, http.StatusBadRequest, BadRequest("Combat state is required"))
		return
	}

	// Get player
	player, err := s.repo.GetPlayerByID(req.PlayerID)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Player"))
		return
	}

	// Get skill
	skillModel, err := s.repo.GetSkillByID(req.SkillID)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Skill"))
		return
	}

	// Check level requirement
	if player.Level < skillModel.Level {
		respondError(c, http.StatusBadRequest, BadRequest(fmt.Sprintf("等级不足，需要 %d 级", skillModel.Level)))
		return
	}

	// Convert to game skill
	skill := &game.Skill{
		ID:          skillModel.ID,
		Name:        skillModel.Name,
		Code:        skillModel.Code,
		Description: skillModel.Description,
		Type:        skillModel.Type,
		MPCost:      skillModel.MPCost,
		Damage:      skillModel.Damage,
		Heal:        skillModel.Heal,
		Cooldown:    skillModel.Cooldown,
		Level:       skillModel.Level,
		Effect:      skillModel.Effect,
	}

	// Calculate total attack with equipment bonuses
	totalAttack := player.Attack
	if player.Equipment != "" && player.Equipment != "{}" {
		var equip game.Equipment
		if err := json.Unmarshal([]byte(player.Equipment), &equip); err == nil {
			if equip.WeaponID > 0 {
				weapon, err := s.repo.GetItemByID(equip.WeaponID)
				if err == nil {
					var effect map[string]int
					json.Unmarshal([]byte(weapon.Effect), &effect)
					totalAttack += effect["attack"]
				}
			}
			if equip.ArmorID > 0 {
				armor, err := s.repo.GetItemByID(equip.ArmorID)
				if err == nil {
					var effect map[string]int
					json.Unmarshal([]byte(armor.Effect), &effect)
					totalAttack += effect["attack"]
				}
			}
		}
	}

	sm := game.NewSkillManager()
	newState, logMsg, err := sm.UseSkill(skill, req.State, totalAttack)
	if err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}

	// Check combat end and update player
	if !newState.IsActive {
		combatSys := game.NewCombatSystem()
		rewards := combatSys.GetRewards(newState)
		if newState.PlayerHP > 0 {
			player.Exp += rewards.Exp
			player.Gold += rewards.Gold

			// Level up check (support multi-level)
			levelsGained := 0
			for {
				expNeeded := player.Level * 100
				if player.Exp < expNeeded {
					break
				}
				player.Exp -= expNeeded
				player.Level++
				player.HP += 10
				player.MP += 5
				player.Attack += 2
				player.Defense += 1
				levelsGained++
			}

			if levelsGained > 0 {
				newState.Log = append(newState.Log, fmt.Sprintf("恭喜！升级到 %d 级！", player.Level))
			}
		}
		// Update HP/MP from combat state
		player.HP = newState.PlayerHP
		player.MP = newState.PlayerMP
		s.repo.UpdatePlayer(player)
	} else {
		// Update HP/MP from combat state
		player.HP = newState.PlayerHP
		player.MP = newState.PlayerMP
		s.repo.UpdatePlayer(player)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    newState,
		"message": logMsg,
		"skill":   skillModel.Name,
	})
}

// ==================== 成就系统API ====================

// GetPlayerAchievements 获取玩家成就
func (s *Server) handleGetPlayerAchievements(c *gin.Context) {
	playerID, _ := strconv.ParseUint(c.Param("player_id"), 10, 32)

	// Get all achievements
	allAchievements, err := s.repo.GetAchievements()
	if err != nil {
		respondInternalError(c, err)
		return
	}

	// Get player's unlocked achievements
	playerAchievements, err := s.repo.GetPlayerAchievements(uint(playerID))
	if err != nil {
		respondInternalError(c, err)
		return
	}

	unlockedMap := make(map[uint]bool)
	for _, pa := range playerAchievements {
		unlockedMap[pa.AchievementID] = true
	}

	type AchievementStatus struct {
		models.Achievement
		Unlocked bool `json:"unlocked"`
	}

	var result []AchievementStatus
	for _, ach := range allAchievements {
		result = append(result, AchievementStatus{
			Achievement: ach,
			Unlocked:    unlockedMap[ach.ID],
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"player_id":   playerID,
		"achievements": result,
		"total":       len(allAchievements),
		"unlocked":    len(playerAchievements),
	})
}

// CheckAchievements 检查并解锁新成就
func (s *Server) handleCheckAchievements(c *gin.Context) {
	var req struct {
		PlayerID uint `json:"player_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}

	player, err := s.repo.GetPlayerByID(req.PlayerID)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Player"))
		return
	}

	// Get all achievements
	allAchievements, err := s.repo.GetAchievements()
	if err != nil {
		respondInternalError(c, err)
		return
	}

	// Get player's unlocked achievements
	playerAchievements, err := s.repo.GetPlayerAchievements(req.PlayerID)
	if err != nil {
		respondInternalError(c, err)
		return
	}

	unlockedMap := make(map[uint]bool)
	for _, pa := range playerAchievements {
		unlockedMap[pa.AchievementID] = true
	}

	// Build player achievement data
	am := game.NewAchievementManager()

	// Count unique items
	uniqueItems := 0
	if player.Items != "" {
		var items map[string]int
		json.Unmarshal([]byte(player.Items), &items)
		uniqueItems = len(items)
	}

	// Count completed quests (tasks with status "completed")
	completedQuests := make(map[string]bool)
	tasks, _ := s.repo.GetTasks()
	questCount := 0
	for _, t := range tasks {
		if t.Status == "completed" {
			completedQuests[t.Code] = true
			questCount++
		}
	}

	// Count visited scenes (simplified - based on player scene)
	visitedScenes := 1 // at least current scene

	playerData := &game.PlayerAchievementData{
		Level:           player.Level,
		TotalGold:       player.Gold + (player.Level-1)*500, // approximate total
		CombatWins:      0, // would need a counter field
		QuestCount:      questCount,
		CompletedQuests: completedQuests,
		VisitedScenes:   visitedScenes,
		UniqueItems:     uniqueItems,
		TalkedToAllNPCs: false,
		SkillsUsed:      0,
	}

	// Convert achievements
	var gameAchievements []*game.Achievement
	for _, ach := range allAchievements {
		gameAchievements = append(gameAchievements, &game.Achievement{
			ID:          ach.ID,
			Name:        ach.Name,
			Code:        ach.Code,
			Description: ach.Description,
			Condition:   ach.Condition,
			Reward:      ach.Reward,
			Icon:        ach.Icon,
		})
	}

	// Check for new achievements
	newAchievements := am.CheckAchievements(gameAchievements, playerData, unlockedMap)

	var unlockedNames []string
	for _, ach := range newAchievements {
		// Save to DB
		pa := &models.PlayerAchievement{
			PlayerID:      req.PlayerID,
			AchievementID: ach.ID,
		}
		s.repo.CreatePlayerAchievement(pa)

		// Apply reward
		reward := am.GetReward(ach.Reward)
		player.Exp += reward.Exp
		player.Gold += reward.Gold

		unlockedNames = append(unlockedNames, ach.Name)
	}

	if len(newAchievements) > 0 {
		s.repo.UpdatePlayer(player)
	}

	c.JSON(http.StatusOK, gin.H{
		"new_achievements": unlockedNames,
		"count":           len(newAchievements),
		"player_exp":      player.Exp,
		"player_gold":     player.Gold,
	})
}

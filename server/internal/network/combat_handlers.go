package network

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/database/models"
	"github.com/ddc-111/agentGame/server/internal/game"
)

func (s *Server) registerCombatRoutes(api *gin.RouterGroup) {
	combat := api.Group("/combat")
	combat.Use(RateLimitMiddleware(30, 50))
	{
		combat.POST("/start", s.handleStartCombat)
		combat.POST("/action", s.handleCombatAction)
	}
}

// handleStartCombat godoc
// @Summary      Start combat
// @Description  Start a new combat encounter
// @Tags         combat
// @Accept       json
// @Produce      json
// @Param        request  body  object  true  "Combat start request"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /combat/start [post]
func (s *Server) handleStartCombat(c *gin.Context) {
	ctx := c.Request.Context()
	var req struct {
		PlayerID  uint   `json:"player_id"`
		EnemyType string `json:"enemy_type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}

	errs := mergeErrors(
		validatePositiveInt("player_id", req.PlayerID),
		validateRequired(map[string]interface{}{"enemy_type": req.EnemyType}),
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

	playerStats := s.calcPlayerStats(ctx, player)

	combatSys := game.NewCombatSystem()
	state := combatSys.StartCombat(req.PlayerID, req.EnemyType, playerStats.TotalHP, player.MP)
	state.PlayerDef = playerStats.TotalDefense

	c.JSON(http.StatusOK, gin.H{
		"data":    state,
		"message": "战斗开始",
	})
}

// handleCombatAction godoc
// @Summary      Perform combat action
// @Description  Perform a combat action (attack, skill, item, flee)
// @Tags         combat
// @Accept       json
// @Produce      json
// @Param        request  body  object  true  "Combat action request"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /combat/action [post]
func (s *Server) handleCombatAction(c *gin.Context) {
	ctx := c.Request.Context()
	var req struct {
		PlayerID uint              `json:"player_id"`
		Action   string            `json:"action"`
		ItemID   uint              `json:"item_id,omitempty"`
		SkillID  uint              `json:"skill_id,omitempty"`
		State    *game.CombatState `json:"state"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}

	errs := mergeErrors(
		validatePositiveInt("player_id", req.PlayerID),
		validateRequired(map[string]interface{}{"action": req.Action}),
		validateStringIn("action", req.Action, []string{"attack", "skill", "item", "flee"}),
	)
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}

	if req.State == nil {
		respondError(c, http.StatusBadRequest, BadRequest("Combat state is required"))
		return
	}

	player, err := s.repo.GetPlayerByID(ctx, req.PlayerID)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Player"))
		return
	}

	playerStats := s.calcPlayerStats(ctx, player)
	totalAttack := playerStats.TotalAttack

	im := game.NewInventoryManager()
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
		skillModel, err := s.repo.GetSkillByID(ctx, req.SkillID)
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
		player.SkillsUsed++

	case "item":
		item, err := s.repo.GetItemByID(ctx, req.ItemID)
		if err != nil {
			respondError(c, http.StatusBadRequest, BadRequest("Item not found"))
			return
		}

		var effect map[string]int
		err = json.Unmarshal([]byte(item.Effect), &effect)
		if err != nil {
			respondError(c, http.StatusBadRequest, BadRequest("Invalid item effect"))
			return
		}

		newState = combatSys.UseItem(req.State, effect)

		newItemsJSON, _, err := im.UseItem(player.Items, req.ItemID, effect)
		if err == nil {
			player.Items = newItemsJSON
			if err := s.repo.UpdatePlayer(ctx, player); err != nil {
				respondInternalError(c, err)
				return
			}
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

	if !newState.IsActive {
		rewards := combatSys.GetRewards(newState)
		if newState.PlayerHP > 0 {
			player.Exp += rewards.Exp
			player.Gold += rewards.Gold
			player.CombatWins++

			levelsGained := processLevelUp(player)

			player.HP = newState.PlayerHP
			player.MP = newState.PlayerMP

			if err := s.repo.UpdatePlayer(ctx, player); err != nil {
				respondInternalError(c, err)
				return
			}

			if levelsGained > 0 {
				newState.Log = append(newState.Log, fmt.Sprintf("恭喜！升级到 %d 级！", player.Level))
			}
		} else {
			player.HP = newState.PlayerHP
			player.MP = newState.PlayerMP
			if err := s.repo.UpdatePlayer(ctx, player); err != nil {
				respondInternalError(c, err)
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    newState,
		"message": "战斗更新",
	})
}

func (s *Server) playerEquipStats(ctx context.Context, equipJSON string) game.EquipmentStats {
	im := game.NewInventoryManager()
	stats, _ := im.EquipmentStatsFromEquip(equipJSON, s.itemEffectLookup(ctx))
	return stats
}

func (s *Server) calcPlayerStats(ctx context.Context, p *models.Player) *game.PlayerStats {
	im := game.NewInventoryManager()
	equipStats := s.playerEquipStats(ctx, p.Equipment)
	return im.CalculateStats(p.Attack, p.Defense, p.HP, p.MP, equipStats)
}

func (s *Server) itemEffectLookup(ctx context.Context) game.ItemLookupFunc {
	return func(itemID uint) (map[string]int, error) {
		item, err := s.repo.GetItemByID(ctx, itemID)
		if err != nil {
			return nil, err
		}
		var effect map[string]int
		if err := json.Unmarshal([]byte(item.Effect), &effect); err != nil {
			return nil, err
		}
		return effect, nil
	}
}

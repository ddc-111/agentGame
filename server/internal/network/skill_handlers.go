package network

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/game"
)

func (s *Server) registerSkillRoutes(api *gin.RouterGroup) {
	api.GET("/skills", s.handleGetSkills)
	api.POST("/skills/use", s.handleUseSkill)
}

func (s *Server) handleGetSkills(c *gin.Context) {
	ctx := c.Request.Context()
	p := parsePagination(c)
	skills, total, err := s.repo.GetSkillsPaginated(ctx, p.Offset, p.PageSize)
	if err != nil {
		respondInternalError(c, err)
		return
	}
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
	c.JSON(http.StatusOK, gin.H{"data": skills, "total": total})
}

func (s *Server) handleUseSkill(c *gin.Context) {
	ctx := c.Request.Context()
	var req struct {
		PlayerID uint              `json:"player_id"`
		SkillID  uint              `json:"skill_id"`
		State    *game.CombatState `json:"state"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}

	errs := mergeErrors(
		validatePositiveInt("player_id", req.PlayerID),
		validatePositiveInt("skill_id", req.SkillID),
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

	skillModel, err := s.repo.GetSkillByID(ctx, req.SkillID)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Skill"))
		return
	}

	if player.Level < skillModel.Level {
		respondError(c, http.StatusBadRequest, BadRequest(fmt.Sprintf("等级不足，需要 %d 级", skillModel.Level)))
		return
	}

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

	playerStats := s.calcPlayerStats(ctx, player)
	totalAttack := playerStats.TotalAttack

	sm := game.NewSkillManager()
	newState, logMsg, err := sm.UseSkill(skill, req.State, totalAttack)
	if err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}

	if !newState.IsActive {
		combatSys := game.NewCombatSystem()
		rewards := combatSys.GetRewards(newState)
		if newState.PlayerHP > 0 {
			player.Exp += rewards.Exp
			player.Gold += rewards.Gold

			levelsGained := processLevelUp(player)

			if levelsGained > 0 {
				newState.Log = append(newState.Log, fmt.Sprintf("恭喜！升级到 %d 级！", player.Level))
			}
		}
		player.HP = newState.PlayerHP
		player.MP = newState.PlayerMP
		if err := s.repo.UpdatePlayer(ctx, player); err != nil {
			respondInternalError(c, err)
			return
		}
	} else {
		player.HP = newState.PlayerHP
		player.MP = newState.PlayerMP
		if err := s.repo.UpdatePlayer(ctx, player); err != nil {
			respondInternalError(c, err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    newState,
		"message": logMsg,
		"skill":   skillModel.Name,
	})
}

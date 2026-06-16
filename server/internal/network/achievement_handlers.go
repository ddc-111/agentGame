package network

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/database/models"
	"github.com/ddc-111/agentGame/server/internal/game"
)

func (s *Server) registerAchievementRoutes(api *gin.RouterGroup) {
	api.GET("/achievements/:player_id", s.handleGetPlayerAchievements)
	api.POST("/achievements/check", s.handleCheckAchievements)
}

func (s *Server) handleGetPlayerAchievements(c *gin.Context) {
	playerID, _ := strconv.ParseUint(c.Param("player_id"), 10, 32)
	p := parsePagination(c)

	allAchievements, total, err := s.repo.GetAchievementsPaginated(p.Offset, p.PageSize)
	if err != nil {
		respondInternalError(c, err)
		return
	}

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

	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
	c.JSON(http.StatusOK, gin.H{
		"player_id":    playerID,
		"achievements": result,
		"total":        total,
		"unlocked":     len(playerAchievements),
	})
}

func (s *Server) handleCheckAchievements(c *gin.Context) {
	var req struct {
		PlayerID uint `json:"player_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}

	errs := validatePositiveInt("player_id", req.PlayerID)
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}

	player, err := s.repo.GetPlayerByID(req.PlayerID)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Player"))
		return
	}

	allAchievements, err := s.repo.GetAchievements()
	if err != nil {
		respondInternalError(c, err)
		return
	}

	playerAchievements, err := s.repo.GetPlayerAchievements(req.PlayerID)
	if err != nil {
		respondInternalError(c, err)
		return
	}

	unlockedMap := make(map[uint]bool)
	for _, pa := range playerAchievements {
		unlockedMap[pa.AchievementID] = true
	}

	am := game.NewAchievementManager()

	uniqueItems := 0
	if player.Items != "" {
		var items map[string]int
		json.Unmarshal([]byte(player.Items), &items)
		uniqueItems = len(items)
	}

	completedQuests := make(map[string]bool)
	tasks, _ := s.repo.GetTasks()
	questCount := 0
	for _, t := range tasks {
		if t.Status == "completed" {
			completedQuests[t.Code] = true
			questCount++
		}
	}

	visitedSceneCount := 1
	if player.VisitedScenes != "" {
		var sceneCodes []string
		json.Unmarshal([]byte(player.VisitedScenes), &sceneCodes)
		if len(sceneCodes) > 0 {
			visitedSceneCount = len(sceneCodes)
		}
	}

	playerData := &game.PlayerAchievementData{
		Level:           player.Level,
		TotalGold:       player.Gold + (player.Level-1)*500,
		CombatWins:      player.CombatWins,
		QuestCount:      questCount,
		CompletedQuests: completedQuests,
		VisitedScenes:   visitedSceneCount,
		UniqueItems:     uniqueItems,
		TalkedToAllNPCs: false,
		SkillsUsed:      player.SkillsUsed,
	}

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

	newAchievements := am.CheckAchievements(gameAchievements, playerData, unlockedMap)

	var unlockedNames []string
	for _, ach := range newAchievements {
		pa := &models.PlayerAchievement{
			PlayerID:      req.PlayerID,
			AchievementID: ach.ID,
		}
		s.repo.CreatePlayerAchievement(pa)

		reward := am.GetReward(ach.Reward)
		player.Exp += reward.Exp
		player.Gold += reward.Gold

		unlockedNames = append(unlockedNames, ach.Name)
	}

	if len(newAchievements) > 0 {
		if err := s.repo.UpdatePlayer(player); err != nil {
			respondInternalError(c, err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"new_achievements": unlockedNames,
		"count":            len(newAchievements),
		"player_exp":       player.Exp,
		"player_gold":      player.Gold,
	})
}

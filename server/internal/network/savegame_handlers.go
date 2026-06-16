package network

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/database/models"
	"github.com/ddc-111/agentGame/server/internal/game"
)

func (s *Server) registerSavegameRoutes(api *gin.RouterGroup) {
	api.POST("/save", s.handleSaveGame)
	api.GET("/saves/:player_id", s.handleGetSaves)
	api.POST("/load/:save_id", s.handleLoadGame)
}

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

	errs := mergeErrors(
		validatePositiveInt("player_id", req.PlayerID),
		validateIntMin("slot", req.Slot, 0),
		validateIntRange("slot", req.Slot, 0, 10),
	)
	if len(errs) > 0 {
		respondValidation(c, errs)
		return
	}

	sgm := game.NewSaveGameManager()
	if err := sgm.ValidateSlot(req.Slot); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}

	player, err := s.repo.GetPlayerByID(req.PlayerID)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Player"))
		return
	}

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

	existingSave, _ := s.repo.GetSaveGame(req.PlayerID, req.Slot)

	saveName := req.Name
	if saveName == "" {
		saveName = sgm.FormatSaveName(req.Slot, player.Name)
	}

	if existingSave != nil && existingSave.ID > 0 {
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

func (s *Server) handleGetSaves(c *gin.Context) {
	playerID, _ := strconv.ParseUint(c.Param("player_id"), 10, 32)

	saves, err := s.repo.GetSaveGames(uint(playerID))
	if err != nil {
		respondInternalError(c, err)
		return
	}

	sgm := game.NewSaveGameManager()
	var saveInfos []game.SaveSlotInfo

	for i := 0; i <= 10; i++ {
		info := game.SaveSlotInfo{
			Slot:    i,
			IsEmpty: true,
		}

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

func (s *Server) handleLoadGame(c *gin.Context) {
	saveID, _ := strconv.ParseUint(c.Param("save_id"), 10, 32)

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

	sgm := game.NewSaveGameManager()
	snapshot, err := sgm.DeserializeSnapshot(targetSave.Snapshot)
	if err != nil {
		respondInternalError(c, err)
		return
	}

	player, err := s.repo.GetPlayerByID(req.PlayerID)
	if err != nil {
		respondError(c, http.StatusNotFound, NotFound("Player"))
		return
	}

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

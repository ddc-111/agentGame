package network

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *Server) handleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	playerIDStr := c.Query("player_id")
	sceneID := c.Query("scene_id")

	var playerID uint
	if playerIDStr != "" {
		if id, err := strconv.ParseUint(playerIDStr, 10, 32); err == nil {
			playerID = uint(id)
		}
	}

	client := &Client{
		hub:      s.hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		playerID: playerID,
		sceneID:  sceneID,
	}

	if sceneID != "" {
		s.hub.mu.Lock()
		s.hub.addToRoom(client, sceneID)
		s.hub.mu.Unlock()
	}

	s.hub.register <- client

	if playerID > 0 {
		s.sendInitialState(c.Request.Context(), client, playerID, sceneID)
	}

	go client.writePump()
	go client.readPump(s.hub)
}

func (s *Server) sendInitialState(ctx context.Context, client *Client, playerID uint, sceneID string) {
	player, err := s.repo.GetPlayerByID(ctx, playerID)
	if err != nil {
		log.Printf("Failed to get player for initial state: %v", err)
		return
	}

	stateData, _ := json.Marshal(map[string]interface{}{
		"player": map[string]interface{}{
			"id":       player.ID,
			"name":     player.Name,
			"level":    player.Level,
			"hp":       player.HP,
			"mp":       player.MP,
			"scene_id": player.SceneID,
			"pos_x":    player.PosX,
			"pos_y":    player.PosY,
			"gold":     player.Gold,
		},
		"online_count": s.hub.GetOnlineCount(),
	})

	msg := &WSMessage{
		Type:      MsgStateSync,
		PlayerID:  playerID,
		SceneID:   sceneID,
		Timestamp: time.Now().UnixMilli(),
		Data:      stateData,
	}
	_ = writeMessage(client, msg)
}

func (s *Server) BroadcastPlayerPosition(ctx context.Context, playerID uint, sceneID string, posX, posY int) {
	if s.hub == nil {
		return
	}

	player, err := s.repo.GetPlayerByID(ctx, playerID)
	if err != nil {
		return
	}

	posData, _ := json.Marshal(PlayerPositionData{
		PlayerID: playerID,
		Name:     player.Name,
		SceneID:  sceneID,
		PosX:     posX,
		PosY:     posY,
	})

	msg := &WSMessage{
		Type:      MsgPlayerPosition,
		PlayerID:  playerID,
		SceneID:   sceneID,
		Timestamp: time.Now().UnixMilli(),
		Data:      posData,
	}

	data, _ := json.Marshal(msg)
	s.hub.BroadcastToRoom(sceneID, data)
}

func (s *Server) BroadcastNPCState(npcID uint, code, name, sceneID, state string, posX, posY int) {
	if s.hub == nil {
		return
	}

	npcData, _ := json.Marshal(NPCStateData{
		NPCID:   npcID,
		Code:    code,
		Name:    name,
		SceneID: sceneID,
		PosX:    posX,
		PosY:    posY,
		State:   state,
	})

	msg := &WSMessage{
		Type:      MsgNPCState,
		SceneID:   sceneID,
		Timestamp: time.Now().UnixMilli(),
		Data:      npcData,
	}

	data, _ := json.Marshal(msg)
	s.hub.BroadcastToRoom(sceneID, data)
}

func (s *Server) BroadcastChatMessage(playerID uint, playerName, sceneID, channel, content string) {
	if s.hub == nil {
		return
	}

	chatData, _ := json.Marshal(ChatMessageData{
		PlayerID:   playerID,
		PlayerName: playerName,
		SceneID:    sceneID,
		Channel:    channel,
		Content:    content,
	})

	msg := &WSMessage{
		Type:      MsgChatMessage,
		PlayerID:  playerID,
		SceneID:   sceneID,
		Timestamp: time.Now().UnixMilli(),
		Data:      chatData,
	}

	data, _ := json.Marshal(msg)
	switch channel {
	case "global":
		s.hub.BroadcastToAll(data)
	default:
		s.hub.BroadcastToRoom(sceneID, data)
	}
}

func (s *Server) BroadcastCombatEvent(ctx context.Context, playerID uint, targetID uint, targetType, eventType string, damage, heal, hp int, skillCode string) {
	if s.hub == nil {
		return
	}

	player, _ := s.repo.GetPlayerByID(ctx, playerID)
	sceneID := ""
	if player != nil {
		sceneID = player.SceneID
	}

	combatData, _ := json.Marshal(CombatEventData{
		PlayerID:   playerID,
		TargetID:   targetID,
		TargetType: targetType,
		EventType:  eventType,
		Damage:     damage,
		Heal:       heal,
		HP:         hp,
		SkillCode:  skillCode,
	})

	msg := &WSMessage{
		Type:      MsgCombatEvent,
		PlayerID:  playerID,
		SceneID:   sceneID,
		Timestamp: time.Now().UnixMilli(),
		Data:      combatData,
	}

	data, _ := json.Marshal(msg)
	if sceneID != "" {
		s.hub.BroadcastToRoom(sceneID, data)
	} else {
		s.hub.BroadcastToAll(data)
	}
}

func (s *Server) BroadcastItemPickup(ctx context.Context, playerID uint, itemID uint, itemCode, itemName string, count int) {
	if s.hub == nil {
		return
	}

	player, _ := s.repo.GetPlayerByID(ctx, playerID)
	sceneID := ""
	if player != nil {
		sceneID = player.SceneID
	}

	pickupData, _ := json.Marshal(ItemPickupData{
		PlayerID: playerID,
		ItemID:   itemID,
		ItemCode: itemCode,
		ItemName: itemName,
		Count:    count,
	})

	msg := &WSMessage{
		Type:      MsgItemPickup,
		PlayerID:  playerID,
		SceneID:   sceneID,
		Timestamp: time.Now().UnixMilli(),
		Data:      pickupData,
	}

	data, _ := json.Marshal(msg)
	if sceneID != "" {
		s.hub.BroadcastToRoom(sceneID, data)
	} else {
		s.hub.BroadcastToAll(data)
	}
}

func (s *Server) BroadcastSystemMessage(message, level string, targetPlayerID uint) {
	if s.hub == nil {
		return
	}

	sysData, _ := json.Marshal(SystemMessageData{
		Message:  message,
		Level:    level,
		TargetID: targetPlayerID,
	})

	msg := &WSMessage{
		Type:      MsgSystemMessage,
		Timestamp: time.Now().UnixMilli(),
		Data:      sysData,
	}

	data, _ := json.Marshal(msg)
	if targetPlayerID > 0 {
		s.hub.SendToPlayer(targetPlayerID, data)
	} else {
		s.hub.BroadcastToAll(data)
	}
}

func (s *Server) handleGMLogin(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "username and password are required",
		})
		return
	}

	if req.Username != s.cfg.Auth.GMUsername || req.Password != s.cfg.Auth.GMPassword {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "invalid username or password",
		})
		return
	}

	token, err := GenerateJWT(s.cfg.Auth.JWTSecret, req.Username, "gm", s.cfg.Auth.TokenExpiry)
	if err != nil {
		respondInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"token": token,
		},
	})
}

func (s *Server) handleGMMe(c *gin.Context) {
	username, _ := c.Get("gm_username")
	role, _ := c.Get("gm_role")
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"username": username,
			"role":     role,
		},
	})
}

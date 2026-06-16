package network

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type MessageType string

const (
	MsgPlayerPosition MessageType = "player_position"
	MsgNPCState       MessageType = "npc_state"
	MsgChatMessage    MessageType = "chat_message"
	MsgCombatEvent    MessageType = "combat_event"
	MsgItemPickup     MessageType = "item_pickup"
	MsgSystemMessage  MessageType = "system_message"
	MsgStateSync      MessageType = "state_sync"
	MsgPing           MessageType = "ping"
	MsgPong           MessageType = "pong"
)

type WSMessage struct {
	Type      MessageType     `json:"type"`
	PlayerID  uint            `json:"player_id,omitempty"`
	SceneID   string          `json:"scene_id,omitempty"`
	Timestamp int64           `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
}

type PlayerPositionData struct {
	PlayerID uint   `json:"player_id"`
	Name     string `json:"name"`
	SceneID  string `json:"scene_id"`
	PosX     int    `json:"pos_x"`
	PosY     int    `json:"pos_y"`
}

type NPCStateData struct {
	NPCID      uint   `json:"npc_id"`
	Code       string `json:"code"`
	Name       string `json:"name"`
	SceneID    string `json:"scene_id"`
	PosX       int    `json:"pos_x"`
	PosY       int    `json:"pos_y"`
	State      string `json:"state"`
	DialogText string `json:"dialog_text,omitempty"`
}

type ChatMessageData struct {
	PlayerID  uint   `json:"player_id"`
	PlayerName string `json:"player_name"`
	SceneID   string `json:"scene_id"`
	Channel   string `json:"channel"`
	Content   string `json:"content"`
}

type CombatEventData struct {
	PlayerID   uint   `json:"player_id"`
	TargetID   uint   `json:"target_id,omitempty"`
	TargetType string `json:"target_type"`
	EventType  string `json:"event_type"`
	Damage     int    `json:"damage,omitempty"`
	Heal       int    `json:"heal,omitempty"`
	HP         int    `json:"hp,omitempty"`
	SkillCode  string `json:"skill_code,omitempty"`
}

type ItemPickupData struct {
	PlayerID uint   `json:"player_id"`
	ItemID   uint   `json:"item_id"`
	ItemCode string `json:"item_code"`
	ItemName string `json:"item_name"`
	Count    int    `json:"count"`
}

type SystemMessageData struct {
	Message  string `json:"message"`
	Level    string `json:"level"`
	TargetID uint   `json:"target_id,omitempty"`
}

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	playerID uint
	sceneID  string
	mu       sync.Mutex
}

type Hub struct {
	clients    map[*Client]bool
	rooms      map[string]map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		rooms:      make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("WebSocket client connected (player=%d, scene=%s)", client.playerID, client.sceneID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				h.removeFromRoom(client)
				close(client.send)
			}
			h.mu.Unlock()
			log.Printf("WebSocket client disconnected (player=%d)", client.playerID)

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					go func(c *Client) {
						h.unregister <- c
					}(client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) BroadcastToRoom(sceneID string, msg []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	room, ok := h.rooms[sceneID]
	if !ok {
		return
	}
	for client := range room {
		select {
		case client.send <- msg:
		default:
			go func(c *Client) {
				h.unregister <- c
			}(client)
		}
	}
}

func (h *Hub) BroadcastToAll(msg []byte) {
	select {
	case h.broadcast <- msg:
	default:
		log.Println("Broadcast channel full, dropping message")
	}
}

func (h *Hub) SendToPlayer(playerID uint, msg []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.playerID == playerID {
			select {
			case client.send <- msg:
			default:
			}
			return
		}
	}
}

func (h *Hub) addToRoom(client *Client, sceneID string) {
	if sceneID == "" {
		return
	}
	if h.rooms[sceneID] == nil {
		h.rooms[sceneID] = make(map[*Client]bool)
	}
	h.rooms[sceneID][client] = true
}

func (h *Hub) removeFromRoom(client *Client) {
	if client.sceneID == "" {
		return
	}
	if room, ok := h.rooms[client.sceneID]; ok {
		delete(room, client)
		if len(room) == 0 {
			delete(h.rooms, client.sceneID)
		}
	}
}

func (h *Hub) MoveClientToRoom(client *Client, newSceneID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.removeFromRoom(client)
	client.sceneID = newSceneID
	h.addToRoom(client, newSceneID)
}

func (h *Hub) GetRoomPlayerCount(sceneID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if room, ok := h.rooms[sceneID]; ok {
		return len(room)
	}
	return 0
}

func (h *Hub) GetOnlineCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

func (h *Hub) GetScenePlayers(sceneID string) []PlayerPositionData {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var players []PlayerPositionData
	if room, ok := h.rooms[sceneID]; ok {
		for client := range room {
			players = append(players, PlayerPositionData{
				PlayerID: client.playerID,
				SceneID:  client.sceneID,
			})
		}
	}
	return players
}

func writeMessage(client *Client, msg *WSMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	client.mu.Lock()
	defer client.mu.Unlock()
	_ = client.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return client.conn.WriteMessage(websocket.TextMessage, data)
}

func (c *Client) readPump(hub *Hub) {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}

		var wsMsg WSMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			log.Printf("Invalid WebSocket message: %v", err)
			continue
		}

		wsMsg.Timestamp = time.Now().UnixMilli()
		if wsMsg.PlayerID == 0 {
			wsMsg.PlayerID = c.playerID
		}

		hub.handleMessage(c, &wsMsg)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (h *Hub) handleMessage(client *Client, msg *WSMessage) {
	switch msg.Type {
	case MsgPlayerPosition:
		h.handlePlayerPosition(client, msg)
	case MsgChatMessage:
		h.handleChatMessage(client, msg)
	case MsgPing:
		pongData, _ := json.Marshal(map[string]string{"status": "pong"})
		pong := &WSMessage{
			Type:      MsgPong,
			Timestamp: time.Now().UnixMilli(),
			Data:      pongData,
		}
		_ = writeMessage(client, pong)
	default:
		log.Printf("Unknown message type: %s", msg.Type)
	}
}

func (h *Hub) handlePlayerPosition(client *Client, msg *WSMessage) {
	var posData PlayerPositionData
	if err := json.Unmarshal(msg.Data, &posData); err != nil {
		log.Printf("Invalid player position data: %v", err)
		return
	}

	posData.PlayerID = client.playerID
	if posData.SceneID != "" && posData.SceneID != client.sceneID {
		h.MoveClientToRoom(client, posData.SceneID)
	}

	msg.SceneID = client.sceneID
	msg.PlayerID = client.playerID
	data, _ := json.Marshal(posData)
	msg.Data = data

	broadcastData, _ := json.Marshal(msg)
	h.BroadcastToRoom(client.sceneID, broadcastData)
}

func (h *Hub) handleChatMessage(client *Client, msg *WSMessage) {
	var chatData ChatMessageData
	if err := json.Unmarshal(msg.Data, &chatData); err != nil {
		log.Printf("Invalid chat message data: %v", err)
		return
	}

	chatData.PlayerID = client.playerID
	chatData.SceneID = client.sceneID

	msg.SceneID = client.sceneID
	msg.PlayerID = client.playerID
	data, _ := json.Marshal(chatData)
	msg.Data = data

	broadcastData, _ := json.Marshal(msg)

	switch chatData.Channel {
	case "scene":
		h.BroadcastToRoom(client.sceneID, broadcastData)
	case "global":
		h.BroadcastToAll(broadcastData)
	default:
		h.BroadcastToRoom(client.sceneID, broadcastData)
	}
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
)

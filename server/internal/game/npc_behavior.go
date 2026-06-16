package game

import (
	"encoding/json"
	"fmt"
	"time"
)

// NPCBehavior defines what an NPC does autonomously
type NPCBehavior struct {
	NPCCode  string         `json:"npc_code"`
	State    string         `json:"state"` // idle, patrolling, talking, trading
	Location string         `json:"location"`
	Target   string         `json:"target"`
	Schedule []ScheduleEntry `json:"schedule"`
	Mood     string         `json:"mood"` // happy, neutral, angry, scared
	Memory   []NPCEvent     `json:"memory"`
}

// ScheduleEntry defines an NPC's scheduled action
type ScheduleEntry struct {
	Time   string `json:"time"`   // HH:MM format
	Action string `json:"action"` // open_shop, patrol, rest, etc.
	Scene  string `json:"scene"`  // target scene code
}

// NPCEvent records something that happened to an NPC
type NPCEvent struct {
	Time      string `json:"time"`
	Type      string `json:"type"` // talked, attacked, gifted, etc.
	PlayerID  uint   `json:"player_id"`
	Detail    string `json:"detail"`
}

// NPCBehaviorManager manages NPC autonomous behaviors
type NPCBehaviorManager struct{}

// NewNPCBehaviorManager creates a new behavior manager
func NewNPCBehaviorManager() *NPCBehaviorManager {
	return &NPCBehaviorManager{}
}

// UpdateBehavior runs each game tick, moves NPCs based on schedule
func (bm *NPCBehaviorManager) UpdateBehavior(behavior *NPCBehavior, currentHour int) *NPCBehavior {
	if len(behavior.Schedule) == 0 {
		behavior.State = "idle"
		return behavior
	}

	for _, entry := range behavior.Schedule {
		var hour, minute int
		fmt.Sscanf(entry.Time, "%02d:%02d", &hour, &minute)

		if currentHour == hour {
			behavior.State = bm.actionToState(entry.Action)
			if entry.Scene != "" {
				behavior.Target = entry.Scene
			}
			break
		}
	}

	return behavior
}

// ReactToPlayer changes NPC behavior when player is nearby
func (bm *NPCBehaviorManager) ReactToPlayer(behavior *NPCBehavior, playerID uint, action string) *NPCBehavior {
	// Record the event
	event := NPCEvent{
		Time:     time.Now().Format("15:04"),
		Type:     action,
		PlayerID: playerID,
	}

	switch action {
	case "talk":
		behavior.State = "talking"
		if behavior.Mood == "angry" {
			event.Detail = "NPC不情愿地交谈"
		} else {
			event.Detail = "NPC友好地交谈"
		}
	case "gift":
		if behavior.Mood == "neutral" {
			behavior.Mood = "happy"
		}
		event.Detail = "收到礼物，心情变好"
	case "attack":
		behavior.Mood = "angry"
		behavior.State = "fleeing"
		event.Detail = "受到攻击，变得愤怒"
	}

	// Keep only last 10 events
	behavior.Memory = append(behavior.Memory, event)
	if len(behavior.Memory) > 10 {
		behavior.Memory = behavior.Memory[len(behavior.Memory)-10:]
	}

	return behavior
}

// GetDialogMood affects how NPC responds based on mood/state
func (bm *NPCBehaviorManager) GetDialogMood(behavior *NPCBehavior) string {
	switch behavior.Mood {
	case "happy":
		return "npc_happy"
	case "angry":
		return "npc_angry"
	case "scared":
		return "npc_scared"
	default:
		switch behavior.State {
		case "talking":
			return "npc_friendly"
		case "trading":
			return "npc_merchant"
		case "patrolling":
			return "npc_alert"
		default:
			return "npc_neutral"
		}
	}
}

// actionToState converts schedule action to NPC state
func (bm *NPCBehaviorManager) actionToState(action string) string {
	switch action {
	case "open_shop":
		return "trading"
	case "close_shop", "go_home", "rest":
		return "idle"
	case "patrol", "go_hunt", "return_village":
		return "patrolling"
	case "stand_at_tree":
		return "idle"
	default:
		return "idle"
	}
}

// CreateDefaultBehavior creates default behavior for an NPC
func CreateDefaultBehavior(npcCode string, scheduleJSON string) *NPCBehavior {
	behavior := &NPCBehavior{
		NPCCode:  npcCode,
		State:    "idle",
		Location: "",
		Target:   "",
		Schedule: []ScheduleEntry{},
		Mood:     "neutral",
		Memory:   []NPCEvent{},
	}

	if scheduleJSON != "" && scheduleJSON != "[]" {
		json.Unmarshal([]byte(scheduleJSON), &behavior.Schedule)
	}

	return behavior
}

// SerializeNPCBehavior serializes behavior to JSON
func SerializeNPCBehavior(behavior *NPCBehavior) (string, error) {
	data, err := json.Marshal(behavior)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// DeserializeNPCBehavior deserializes behavior from JSON
func DeserializeNPCBehavior(data string) (*NPCBehavior, error) {
	var behavior NPCBehavior
	err := json.Unmarshal([]byte(data), &behavior)
	if err != nil {
		return nil, err
	}
	return &behavior, nil
}

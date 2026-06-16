package game

import (
	"encoding/json"
	"fmt"
)

// InventoryManager 背包管理器
type InventoryManager struct{}

// InventoryItem 背包物品
type InventoryItem struct {
	ItemID uint `json:"item_id"`
	Count  int  `json:"count"`
}

// Equipment 装备槽位
type Equipment struct {
	WeaponID uint `json:"weapon_id"` // 武器
	ArmorID  uint `json:"armor_id"`  // 防具
}

// EquipmentStats 装备属性加成
type EquipmentStats struct {
	Attack  int `json:"attack"`
	Defense int `json:"defense"`
	HP      int `json:"hp"`
	MP      int `json:"mp"`
}

// PlayerStats 玩家属性
type PlayerStats struct {
	BaseAttack  int `json:"base_attack"`
	BaseDefense int `json:"base_defense"`
	BaseHP      int `json:"base_hp"`
	BaseMP      int `json:"base_mp"`
	TotalAttack  int `json:"total_attack"`
	TotalDefense int `json:"total_defense"`
	TotalHP      int `json:"total_hp"`
	TotalMP      int `json:"total_mp"`
}

// NewInventoryManager 创建背包管理器
func NewInventoryManager() *InventoryManager {
	return &InventoryManager{}
}

// AddItem 添加道具到背包
func (im *InventoryManager) AddItem(itemsJSON string, itemID uint, count int) (string, error) {
	items, err := im.parseItems(itemsJSON)
	if err != nil {
		return "", err
	}

	key := fmt.Sprintf("%d", itemID)
	items[key] += count

	return im.serializeItems(items)
}

// RemoveItem 从背包移除道具
func (im *InventoryManager) RemoveItem(itemsJSON string, itemID uint, count int) (string, error) {
	items, err := im.parseItems(itemsJSON)
	if err != nil {
		return "", err
	}

	key := fmt.Sprintf("%d", itemID)
	if items[key] < count {
		return "", fmt.Errorf("道具数量不足")
	}

	items[key] -= count
	if items[key] <= 0 {
		delete(items, key)
	}

	return im.serializeItems(items)
}

// GetItems 获取背包道具列表
func (im *InventoryManager) GetItems(itemsJSON string) ([]InventoryItem, error) {
	items, err := im.parseItems(itemsJSON)
	if err != nil {
		return nil, err
	}

	var result []InventoryItem
	for key, count := range items {
		var itemID uint
		_, _ = fmt.Sscanf(key, "%d", &itemID)
		if count > 0 {
			result = append(result, InventoryItem{
				ItemID: itemID,
				Count:  count,
			})
		}
	}

	return result, nil
}

// UseItem 使用道具
func (im *InventoryManager) UseItem(itemsJSON string, itemID uint, itemEffect map[string]int) (string, map[string]int, error) {
	items, err := im.parseItems(itemsJSON)
	if err != nil {
		return "", nil, err
	}

	key := fmt.Sprintf("%d", itemID)
	if items[key] <= 0 {
		return "", nil, fmt.Errorf("道具不存在或数量为0")
	}

	// 消耗道具
	items[key]--
	if items[key] <= 0 {
		delete(items, key)
	}

	newItemsJSON, err := im.serializeItems(items)
	if err != nil {
		return "", nil, err
	}

	return newItemsJSON, itemEffect, nil
}

// EquipItem 装备道具
func (im *InventoryManager) EquipItem(equipJSON string, itemsJSON string, itemID uint, category string) (string, string, error) {
	equipment, err := im.parseEquipment(equipJSON)
	if err != nil {
		return "", "", err
	}

	items, err := im.parseItems(itemsJSON)
	if err != nil {
		return "", "", err
	}

	key := fmt.Sprintf("%d", itemID)
	if items[key] <= 0 {
		return "", "", fmt.Errorf("道具不存在")
	}

	// 根据类别装备到对应槽位
	switch category {
	case "weapon":
		// 卸下旧武器
		if equipment.WeaponID > 0 {
			oldKey := fmt.Sprintf("%d", equipment.WeaponID)
			items[oldKey]++
		}
		equipment.WeaponID = itemID
		items[key]--
	case "armor":
		// 卸下旧防具
		if equipment.ArmorID > 0 {
			oldKey := fmt.Sprintf("%d", equipment.ArmorID)
			items[oldKey]++
		}
		equipment.ArmorID = itemID
		items[key]--
	default:
		return "", "", fmt.Errorf("该道具无法装备")
	}

	if items[key] <= 0 {
		delete(items, key)
	}

	newEquipJSON, err := im.serializeEquipment(equipment)
	if err != nil {
		return "", "", err
	}

	newItemsJSON, err := im.serializeItems(items)
	if err != nil {
		return "", "", err
	}

	return newEquipJSON, newItemsJSON, nil
}

// UnequipItem 卸下装备
func (im *InventoryManager) UnequipItem(equipJSON string, itemsJSON string, slot string) (string, string, error) {
	equipment, err := im.parseEquipment(equipJSON)
	if err != nil {
		return "", "", err
	}

	items, err := im.parseItems(itemsJSON)
	if err != nil {
		return "", "", err
	}

	switch slot {
	case "weapon":
		if equipment.WeaponID == 0 {
			return "", "", fmt.Errorf("武器槽为空")
		}
		key := fmt.Sprintf("%d", equipment.WeaponID)
		items[key]++
		equipment.WeaponID = 0
	case "armor":
		if equipment.ArmorID == 0 {
			return "", "", fmt.Errorf("防具槽为空")
		}
		key := fmt.Sprintf("%d", equipment.ArmorID)
		items[key]++
		equipment.ArmorID = 0
	default:
		return "", "", fmt.Errorf("无效的装备槽位")
	}

	newEquipJSON, err := im.serializeEquipment(equipment)
	if err != nil {
		return "", "", err
	}

	newItemsJSON, err := im.serializeItems(items)
	if err != nil {
		return "", "", err
	}

	return newEquipJSON, newItemsJSON, nil
}

// GetEquipment 获取装备信息
func (im *InventoryManager) GetEquipment(equipJSON string) (*Equipment, error) {
	return im.parseEquipment(equipJSON)
}

// ItemLookupFunc 根据道具ID获取道具效果的回调函数
type ItemLookupFunc func(itemID uint) (map[string]int, error)

// EquipmentStatsFromEquip 根据装备JSON计算装备属性加成
func (im *InventoryManager) EquipmentStatsFromEquip(equipJSON string, lookup ItemLookupFunc) (EquipmentStats, error) {
	var stats EquipmentStats
	if equipJSON == "" || equipJSON == "{}" {
		return stats, nil
	}

	equip, err := im.parseEquipment(equipJSON)
	if err != nil {
		return stats, err
	}

	slotIDs := []uint{equip.WeaponID, equip.ArmorID}
	for _, itemID := range slotIDs {
		if itemID == 0 {
			continue
		}
		effect, err := lookup(itemID)
		if err != nil {
			continue
		}
		stats.Attack += effect["attack"]
		stats.Defense += effect["defense"]
		stats.HP += effect["hp"]
		stats.MP += effect["mp"]
	}

	return stats, nil
}

// CalculateStats 计算玩家属性（基础 + 装备加成）
func (im *InventoryManager) CalculateStats(baseAttack, baseDefense, baseHP, baseMP int, equipStats EquipmentStats) *PlayerStats {
	return &PlayerStats{
		BaseAttack:   baseAttack,
		BaseDefense:  baseDefense,
		BaseHP:       baseHP,
		BaseMP:       baseMP,
		TotalAttack:  baseAttack + equipStats.Attack,
		TotalDefense: baseDefense + equipStats.Defense,
		TotalHP:      baseHP + equipStats.HP,
		TotalMP:      baseMP + equipStats.MP,
	}
}

// parseItems 解析道具JSON
func (im *InventoryManager) parseItems(itemsJSON string) (map[string]int, error) {
	items := make(map[string]int)
	if itemsJSON == "" || itemsJSON == "{}" {
		return items, nil
	}
	err := json.Unmarshal([]byte(itemsJSON), &items)
	return items, err
}

// serializeItems 序列化道具JSON
func (im *InventoryManager) serializeItems(items map[string]int) (string, error) {
	data, err := json.Marshal(items)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// parseEquipment 解析装备JSON
func (im *InventoryManager) parseEquipment(equipJSON string) (*Equipment, error) {
	equipment := &Equipment{}
	if equipJSON == "" || equipJSON == "{}" {
		return equipment, nil
	}
	err := json.Unmarshal([]byte(equipJSON), equipment)
	return equipment, err
}

// serializeEquipment 序列化装备JSON
func (im *InventoryManager) serializeEquipment(equipment *Equipment) (string, error) {
	data, err := json.Marshal(equipment)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

import { describe, it, expect, beforeEach } from 'vitest';
import { InventoryManager } from '../game/systems/InventoryManager.js';

const mockGameData = {
  items: [
    { id: 1001, name: '金创药', type: 'consumable', effect: { hp: 30 }, icon: '🧪' },
    { id: 1002, name: '铁剑', type: 'weapon', effect: { attack: 5 }, icon: '⚔️' },
    { id: 1003, name: '皮甲', type: 'armor', effect: { defense: 3 }, icon: '🛡️' },
    { id: 1004, name: '草药', type: 'material', icon: '🌿' }
  ]
};

function makePlayerData(overrides = {}) {
  return {
    hp: 80,
    max_hp: 100,
    mp: 30,
    max_mp: 50,
    attack: 10,
    defense: 5,
    speed: 10,
    items: '{}',
    equipment: null,
    ...overrides
  };
}

describe('InventoryManager', () => {
  let playerData;
  let manager;

  beforeEach(() => {
    playerData = makePlayerData();
    manager = new InventoryManager(mockGameData, playerData);
  });

  describe('constructor / parseInventory', () => {
    it('initializes with empty items when playerData.items is empty', () => {
      expect(manager.items).toEqual({});
      expect(manager.equipment).toEqual({ weapon: null, armor: null, shield: null });
    });

    it('parses flat item format {item_id: count}', () => {
      playerData.items = '{"1001": 3, "1004": 5}';
      const m = new InventoryManager(mockGameData, playerData);
      expect(m.items).toEqual({ '1001': 3, '1004': 5 });
    });

    it('parses nested format {items: {...}, equipment: {...}}', () => {
      playerData.items = JSON.stringify({
        items: { '1001': 2 },
        equipment: { weapon: '1002' }
      });
      const m = new InventoryManager(mockGameData, playerData);
      expect(m.items).toEqual({ '1001': 2 });
      expect(m.equipment.weapon).toBe('1002');
    });

    it('parses server equipment format {weapon_id, armor_id}', () => {
      playerData.equipment = JSON.stringify({ weapon_id: 1002, armor_id: 1003 });
      const m = new InventoryManager(mockGameData, playerData);
      expect(m.equipment.weapon).toBe('1002');
      expect(m.equipment.armor).toBe('1003');
    });
  });

  describe('addItem', () => {
    it('adds a new item', () => {
      manager.addItem(1001, 2);
      expect(manager.items['1001']).toBe(2);
    });

    it('increments existing item count', () => {
      manager.addItem(1001, 2);
      manager.addItem(1001, 3);
      expect(manager.items['1001']).toBe(5);
    });

    it('defaults to count 1', () => {
      manager.addItem(1001);
      expect(manager.items['1001']).toBe(1);
    });

    it('persists to playerData.items', () => {
      manager.addItem(1001, 2);
      expect(JSON.parse(playerData.items)['1001']).toBe(2);
    });
  });

  describe('removeItem', () => {
    it('removes items from stack', () => {
      manager.addItem(1001, 5);
      const result = manager.removeItem(1001, 3);
      expect(result).toBe(true);
      expect(manager.items['1001']).toBe(2);
    });

    it('deletes key when count reaches 0', () => {
      manager.addItem(1001, 1);
      manager.removeItem(1001, 1);
      expect(manager.items['1001']).toBeUndefined();
    });

    it('returns false when not enough items', () => {
      manager.addItem(1001, 2);
      const result = manager.removeItem(1001, 5);
      expect(result).toBe(false);
      expect(manager.items['1001']).toBe(2);
    });

    it('returns false when item does not exist', () => {
      expect(manager.removeItem(9999, 1)).toBe(false);
    });
  });

  describe('getItemCount', () => {
    it('returns 0 for non-existent item', () => {
      expect(manager.getItemCount(9999)).toBe(0);
    });

    it('returns correct count', () => {
      manager.addItem(1001, 7);
      expect(manager.getItemCount(1001)).toBe(7);
    });
  });

  describe('getItems', () => {
    it('returns empty array when no items', () => {
      expect(manager.getItems()).toEqual([]);
    });

    it('returns items with full data', () => {
      manager.addItem(1001, 3);
      const items = manager.getItems();
      expect(items).toHaveLength(1);
      expect(items[0].name).toBe('金创药');
      expect(items[0].count).toBe(3);
    });

    it('filters out zero-count items', () => {
      manager.addItem(1001, 2);
      manager.removeItem(1001, 2);
      expect(manager.getItems()).toEqual([]);
    });
  });

  describe('getItemData', () => {
    it('returns item data for known item', () => {
      const data = manager.getItemData(1001);
      expect(data.name).toBe('金创药');
      expect(data.type).toBe('consumable');
    });

    it('returns fallback for unknown item', () => {
      const data = manager.getItemData(9999);
      expect(data.name).toBe('未知物品');
      expect(data.type).toBe('misc');
    });
  });

  describe('equipItem', () => {
    it('equips a weapon', () => {
      manager.addItem(1002, 1);
      const result = manager.equipItem(1002);
      expect(result).toBe(true);
      expect(manager.equipment.weapon).toBe('1002');
      expect(manager.getItemCount(1002)).toBe(0);
    });

    it('equips armor', () => {
      manager.addItem(1003, 1);
      const result = manager.equipItem(1003);
      expect(result).toBe(true);
      expect(manager.equipment.armor).toBe('1003');
    });

    it('swaps equipped weapon', () => {
      manager.addItem(1002, 2);
      manager.equipItem(1002);
      const newWeaponId = 1005;
      mockGameData.items.push({ id: newWeaponId, name: '钢剑', type: 'weapon', effect: { attack: 8 } });
      manager.addItem(newWeaponId, 1);
      manager.equipItem(newWeaponId);
      expect(manager.equipment.weapon).toBe(String(newWeaponId));
      expect(manager.getItemCount(1002)).toBe(1);
    });

    it('returns false for non-equippable item', () => {
      manager.addItem(1001, 1);
      expect(manager.equipItem(1001)).toBe(false);
    });

    it('returns false when item not in inventory', () => {
      expect(manager.equipItem(1002)).toBe(false);
    });
  });

  describe('unequipItem', () => {
    it('unequips item back to inventory', () => {
      manager.addItem(1002, 1);
      manager.equipItem(1002);
      const result = manager.unequipItem('weapon');
      expect(result).toBe(true);
      expect(manager.equipment.weapon).toBeNull();
      expect(manager.getItemCount(1002)).toBe(1);
    });

    it('returns false when slot is empty', () => {
      expect(manager.unequipItem('weapon')).toBe(false);
    });
  });

  describe('getStats', () => {
    it('returns base stats with no equipment', () => {
      const stats = manager.getStats();
      expect(stats.max_hp).toBe(100);
      expect(stats.attack).toBe(10);
      expect(stats.defense).toBe(5);
    });

    it('adds equipment bonuses to stats', () => {
      manager.addItem(1002, 1);
      manager.equipItem(1002);
      const stats = manager.getStats();
      expect(stats.attack).toBe(15); // 10 base + 5 from weapon
    });

    it('adds armor defense bonus', () => {
      manager.addItem(1003, 1);
      manager.equipItem(1003);
      const stats = manager.getStats();
      expect(stats.defense).toBe(8); // 5 base + 3 from armor
    });
  });

  describe('useItem', () => {
    it('uses a consumable and heals HP', () => {
      playerData.hp = 50;
      manager.addItem(1001, 2);
      const result = manager.useItem(1001);
      expect(result).not.toBeNull();
      expect(result.effect.hp).toBe(30);
      expect(playerData.hp).toBe(80);
      expect(manager.getItemCount(1001)).toBe(1);
    });

    it('does not heal above max_hp', () => {
      playerData.hp = 90;
      manager.addItem(1001, 1);
      const result = manager.useItem(1001);
      expect(result.effect.hp).toBe(10); // only 10 to cap
      expect(playerData.hp).toBe(100);
    });

    it('returns null for non-consumable', () => {
      manager.addItem(1002, 1);
      expect(manager.useItem(1002)).toBeNull();
    });

    it('returns null when item not in inventory', () => {
      expect(manager.useItem(1001)).toBeNull();
    });

    it('returns null for unknown item', () => {
      expect(manager.useItem(9999)).toBeNull();
    });
  });

  describe('getItemIcon', () => {
    it('returns item icon when present', () => {
      const item = { icon: '🗡️', type: 'weapon' };
      expect(manager.getItemIcon(item)).toBe('🗡️');
    });

    it('returns type-based icon when no custom icon', () => {
      expect(manager.getItemIcon({ type: 'weapon' })).toBe('⚔️');
      expect(manager.getItemIcon({ type: 'consumable' })).toBe('🧪');
      expect(manager.getItemIcon({ type: 'material' })).toBe('📦');
    });

    it('returns default icon for unknown type', () => {
      expect(manager.getItemIcon({ type: 'unknown' })).toBe('📦');
    });
  });
});

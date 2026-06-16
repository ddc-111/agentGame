import { describe, it, expect, beforeEach, vi } from 'vitest';
import { CombatManager } from '../game/systems/CombatManager.js';

const mockGameData = {
  items: [
    { id: 1001, name: '金创药', type: 'consumable', effect: { hp: 30 }, icon: '🧪' },
    { id: 1002, name: '铁剑', type: 'weapon', effect: { attack: 5 }, icon: '⚔️' }
  ]
};

function makePlayerData(overrides = {}) {
  return {
    hp: 100,
    max_hp: 100,
    mp: 50,
    max_mp: 50,
    attack: 10,
    defense: 5,
    speed: 10,
    gold: 50,
    exp: 0,
    level: 1,
    items: '{}',
    equipment: null,
    ...overrides
  };
}

function makeInventoryManager(playerData) {
  const items = {};
  return {
    getStats: () => ({
      hp: playerData.hp,
      max_hp: playerData.max_hp,
      mp: playerData.mp,
      max_mp: playerData.max_mp,
      attack: playerData.attack,
      defense: playerData.defense,
      speed: playerData.speed
    }),
    getItemData: (id) => mockGameData.items.find(i => i.id === id) || { id, name: '未知', type: 'misc' },
    addItem: (id, count = 1) => { items[id] = (items[id] || 0) + count; },
    useItem: (id) => {
      const item = mockGameData.items.find(i => i.id === id);
      if (!item || item.type !== 'consumable') return null;
      const healed = Math.min(30, playerData.max_hp - playerData.hp);
      playerData.hp += healed;
      return { item, effect: { hp: healed } };
    }
  };
}

describe('CombatManager', () => {
  let playerData;
  let invManager;
  let combat;

  beforeEach(() => {
    playerData = makePlayerData();
    invManager = makeInventoryManager(playerData);
    combat = new CombatManager(playerData, invManager);
  });

  describe('getEnemyData', () => {
    it('returns wolf data for "wolf" type', () => {
      const enemy = combat.getEnemyData('wolf');
      expect(enemy.name).toBe('野狼');
      expect(enemy.hp).toBe(50);
      expect(enemy.attack).toBe(12);
    });

    it('returns bandit data', () => {
      const enemy = combat.getEnemyData('bandit');
      expect(enemy.name).toBe('山贼');
      expect(enemy.hp).toBe(80);
    });

    it('returns bear data', () => {
      const enemy = combat.getEnemyData('bear');
      expect(enemy.name).toBe('黑熊');
      expect(enemy.hp).toBe(120);
    });

    it('returns tiger data', () => {
      const enemy = combat.getEnemyData('tiger');
      expect(enemy.name).toBe('猛虎');
      expect(enemy.hp).toBe(150);
    });

    it('returns ghost data', () => {
      const enemy = combat.getEnemyData('ghost');
      expect(enemy.name).toBe('厉鬼');
      expect(enemy.hp).toBe(100);
    });

    it('returns alpha_wolf data', () => {
      const enemy = combat.getEnemyData('alpha_wolf');
      expect(enemy.name).toBe('头狼');
      expect(enemy.hp).toBe(90);
    });

    it('defaults to wolf for unknown type', () => {
      const enemy = combat.getEnemyData('unknown');
      expect(enemy.name).toBe('野狼');
    });
  });

  describe('startCombat', () => {
    it('initializes combat state', () => {
      const state = combat.startCombat('wolf');
      expect(combat.inCombat).toBe(true);
      expect(combat.enemy.name).toBe('野狼');
      expect(combat.turn).toBe('player');
      expect(state.inCombat).toBe(true);
    });

    it('resets combat log', () => {
      combat.startCombat('wolf');
      expect(combat.combatLog.length).toBe(1);
      expect(combat.combatLog[0]).toContain('野狼');
    });

    it('deep copies enemy data', () => {
      combat.startCombat('wolf');
      combat.enemy.hp = 0;
      combat.startCombat('wolf');
      expect(combat.enemy.hp).toBe(50);
    });
  });

  describe('getCombatState', () => {
    it('returns full combat state', () => {
      combat.startCombat('wolf');
      const state = combat.getCombatState();
      expect(state.inCombat).toBe(true);
      expect(state.enemy.name).toBe('野狼');
      expect(state.player.hp).toBe(100);
      expect(state.turn).toBe('player');
    });
  });

  describe('calculateDamage', () => {
    it('calculates base damage correctly', () => {
      vi.spyOn(Math, 'random').mockReturnValue(0.5);
      const dmg = combat.calculateDamage(15, 5);
      expect(dmg).toBe(10);
      Math.random.mockRestore();
    });

    it('returns at least 1 damage', () => {
      const dmg = combat.calculateDamage(1, 100);
      expect(dmg).toBeGreaterThanOrEqual(1);
    });
  });

  describe('playerAttack', () => {
    it('deals damage to enemy', () => {
      vi.spyOn(Math, 'random').mockReturnValue(0.5);
      combat.startCombat('wolf');
      const result = combat.playerAttack();
      expect(result.action).toBe('attack');
      expect(result.damage).toBeGreaterThan(0);
      expect(result.enemyHp).toBeLessThan(50);
      Math.random.mockRestore();
    });

    it('switches turn to enemy', () => {
      vi.spyOn(Math, 'random').mockReturnValue(0.5);
      combat.startCombat('wolf');
      combat.playerAttack();
      expect(combat.turn).toBe('enemy');
      Math.random.mockRestore();
    });

    it('returns null when not player turn', () => {
      combat.startCombat('wolf');
      combat.turn = 'enemy';
      expect(combat.playerAttack()).toBeNull();
    });

    it('returns null when not in combat', () => {
      expect(combat.playerAttack()).toBeNull();
    });

    it('ends combat when enemy HP reaches 0', () => {
      vi.spyOn(Math, 'random').mockReturnValue(0.5);
      combat.startCombat('wolf');
      combat.enemy.hp = 1;
      const result = combat.playerAttack();
      expect(combat.inCombat).toBe(false);
      expect(combat.rewards).not.toBeNull();
      Math.random.mockRestore();
    });
  });

  describe('playerFlee', () => {
    it('succeeds when random < 0.5', () => {
      vi.spyOn(Math, 'random').mockReturnValue(0.3);
      combat.startCombat('wolf');
      const result = combat.playerFlee();
      expect(result.success).toBe(true);
      expect(combat.inCombat).toBe(false);
      Math.random.mockRestore();
    });

    it('fails when random >= 0.5', () => {
      vi.spyOn(Math, 'random').mockReturnValue(0.7);
      combat.startCombat('wolf');
      const result = combat.playerFlee();
      expect(result.success).toBe(false);
      expect(combat.inCombat).toBe(true);
      expect(combat.turn).toBe('enemy');
      Math.random.mockRestore();
    });

    it('returns null when not player turn', () => {
      combat.startCombat('wolf');
      combat.turn = 'enemy';
      expect(combat.playerFlee()).toBeNull();
    });
  });

  describe('playerUseItem', () => {
    it('uses healing item in combat', () => {
      playerData.hp = 50;
      combat.startCombat('wolf');
      invManager.addItem(1001, 1);
      const result = combat.playerUseItem(1001);
      expect(result).not.toBeNull();
      expect(result.action).toBe('useItem');
      expect(result.playerHp).toBe(80);
      expect(combat.turn).toBe('enemy');
    });

    it('returns null for non-existent item', () => {
      combat.startCombat('wolf');
      expect(combat.playerUseItem(9999)).toBeNull();
    });

    it('returns null when not player turn', () => {
      combat.startCombat('wolf');
      combat.turn = 'enemy';
      expect(combat.playerUseItem(1001)).toBeNull();
    });
  });

  describe('enemyTurn', () => {
    it('deals damage to player', () => {
      vi.spyOn(Math, 'random').mockReturnValue(0.5);
      combat.startCombat('wolf');
      combat.turn = 'enemy';
      const result = combat.enemyTurn();
      expect(result.action).toBe('enemyAttack');
      expect(result.damage).toBeGreaterThan(0);
      expect(result.playerHp).toBeLessThan(100);
      Math.random.mockRestore();
    });

    it('switches turn back to player', () => {
      vi.spyOn(Math, 'random').mockReturnValue(0.5);
      combat.startCombat('wolf');
      combat.turn = 'enemy';
      combat.enemyTurn();
      expect(combat.turn).toBe('player');
      Math.random.mockRestore();
    });

    it('returns null when not enemy turn', () => {
      combat.startCombat('wolf');
      expect(combat.enemyTurn()).toBeNull();
    });

    it('ends combat when player HP reaches 0', () => {
      vi.spyOn(Math, 'random').mockReturnValue(0.5);
      combat.startCombat('wolf');
      playerData.hp = 1;
      combat.turn = 'enemy';
      const result = combat.enemyTurn();
      expect(combat.inCombat).toBe(false);
      expect(result.victory).toBe(false);
      Math.random.mockRestore();
    });
  });

  describe('checkCombatEnd', () => {
    it('returns true and sets rewards when enemy HP <= 0', () => {
      combat.startCombat('wolf');
      combat.enemy.hp = 0;
      expect(combat.checkCombatEnd()).toBe(true);
      expect(combat.inCombat).toBe(false);
      expect(combat.rewards).not.toBeNull();
    });

    it('returns true when player HP <= 0', () => {
      combat.startCombat('wolf');
      playerData.hp = 0;
      expect(combat.checkCombatEnd()).toBe(true);
      expect(combat.inCombat).toBe(false);
    });

    it('returns false when both alive', () => {
      combat.startCombat('wolf');
      expect(combat.checkCombatEnd()).toBe(false);
    });
  });

  describe('getRewards', () => {
    it('returns exp and gold rewards', () => {
      vi.spyOn(Math, 'random').mockReturnValue(0.1);
      combat.startCombat('wolf');
      combat.enemy.hp = 0;
      const rewards = combat.getRewards();
      expect(rewards.exp).toBe(30);
      expect(rewards.gold).toBe(20);
      expect(playerData.gold).toBe(70); // 50 + 20
      Math.random.mockRestore();
    });

    it('handles level up', () => {
      playerData.exp = 95;
      playerData.level = 1;
      combat.startCombat('wolf');
      combat.enemy.hp = 0;
      const rewards = combat.getRewards();
      expect(rewards.levelUps.length).toBeGreaterThanOrEqual(1);
      expect(playerData.level).toBeGreaterThanOrEqual(2);
    });
  });

  describe('addLog', () => {
    it('adds log entries', () => {
      combat.addLog('test message');
      expect(combat.combatLog).toContain('test message');
    });

    it('trims log at 50 entries', () => {
      for (let i = 0; i < 55; i++) {
        combat.addLog(`msg ${i}`);
      }
      expect(combat.combatLog.length).toBe(50);
      expect(combat.combatLog[0]).toBe('msg 5');
    });
  });

  describe('updateCooldowns', () => {
    it('decrements cooldowns', () => {
      combat.cooldowns = { '1': 3, '2': 1 };
      combat.updateCooldowns();
      expect(combat.cooldowns['1']).toBe(2);
      expect(combat.cooldowns['2']).toBe(0);
    });

    it('does not go below 0', () => {
      combat.cooldowns = { '1': 0 };
      combat.updateCooldowns();
      expect(combat.cooldowns['1']).toBe(0);
    });
  });
});

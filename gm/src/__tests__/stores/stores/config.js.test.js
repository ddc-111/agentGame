```javascript
import { describe, it, expect, beforeEach, vi } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { useConfigStore } from './config.js';

describe('useConfigStore', () => {
  let store;

  beforeEach(() => {
    setActivePinia(createPinia());
    store = useConfigStore();
  });

  describe('初始状态', () => {
    it('应该有正确的默认配置', () => {
      expect(store.gameConfig).toBeDefined();
      expect(store.gameConfig.game).toBeDefined();
      expect(store.gameConfig.player).toBeDefined();
      expect(store.gameConfig.world).toBeDefined();
      expect(store.gameConfig.combat).toBeDefined();
      expect(store.gameConfig.economy).toBeDefined();
    });

    it('应该有正确的游戏基础配置', () => {
      expect(store.gameConfig.game.name).toBe('古风RPG');
      expect(store.gameConfig.game.version).toBe('1.0.0');
      expect(store.gameConfig.game.maxPlayers).toBe(100);
      expect(store.gameConfig.game.tickRate).toBe(20);
    });

    it('应该有正确的玩家初始配置', () => {
      expect(store.gameConfig.player.startScene).toBe('scene_001');
      expect(store.gameConfig.player.startPosition).toEqual({ x: 400, y: 300 });
      expect(store.gameConfig.player.startGold).toBe(1000);
      expect(store.gameConfig.player.startLevel).toBe(1);
      expect(store.gameConfig.player.maxLevel).toBe(100);
    });

    it('应该有正确的世界配置', () => {
      expect(store.gameConfig.world.dayNightCycle).toBe(true);
      expect(store.gameConfig.world.dayDuration).toBe(1200);
      expect(store.gameConfig.world.weatherEnabled).toBe(true);
      expect(store.gameConfig.world.weatherTypes).toEqual(['sunny', 'cloudy', 'rainy', 'snowy']);
    });

    it('应该有正确的战斗配置', () => {
      expect(store.gameConfig.combat.enabled).toBe(true);
      expect(store.gameConfig.combat.turnBased).toBe(true);
      expect(store.gameConfig.combat.maxTurns).toBe(20);
      expect(store.gameConfig.combat.criticalRate).toBe(0.1);
      expect(store.gameConfig.combat.criticalMultiplier).toBe(1.5);
    });

    it('应该有正确的经济配置', () => {
      expect(store.gameConfig.economy.inflationRate).toBe(0.01);
      expect(store.gameConfig.economy.taxRate).toBe(0.05);
      expect(store.gameConfig.economy.maxGold).toBe(999999);
    });
  });

  describe('getConfig', () => {
    it('应该返回指定配置节', () => {
      const gameSection = store.getConfig('game');
      expect(gameSection).toEqual(store.gameConfig.game);
    });

    it('对于不存在的配置节应该返回undefined', () => {
      const nonExistentSection = store.getConfig('nonExistent');
      expect(nonExistentSection).toBeUndefined();
    });
  });

  describe('updateConfig', () => {
    it('应该正确更新配置节', () => {
      const updateData = { maxPlayers: 200, name: '新游戏名称' };
      store.updateConfig('game', updateData);
      
      expect(store.gameConfig.game.maxPlayers).toBe(200);
      expect(store.gameConfig.game.name).toBe('新游戏名称');
      // 保留未更新的属性
      expect(store.gameConfig.game.version).toBe('1.0.0');
    });

    it('应该能更新多个配置节', () => {
      store.updateConfig('player', { startGold: 500 });
      store.updateConfig('world', { dayDuration: 2400 });
      
      expect(store.gameConfig.player.startGold).toBe(500);
      expect(store.gameConfig.world.dayDuration).toBe(2400);
    });

    it('对于不存在的配置节不应该修改任何内容', () => {
      const originalConfig = { ...store.gameConfig };
      store.updateConfig('nonExistent', { someKey: 'someValue' });
      
      expect(store.gameConfig).toEqual(originalConfig);
    });

    it('应该只更新提供的属性，不覆盖整个配置节', () => {
      store.updateConfig('player', { startGold: 2000 });
      
      // 确保startGold已更新
      expect(store.gameConfig.player.startGold).toBe(2000);
      // 确保其他属性保持不变
      expect(store.gameConfig.player.startLevel).toBe(1);
      expect(store.gameConfig.player.maxLevel).toBe(100);
    });
  });

  describe('exportConfig', () => {
    it('应该返回JSON字符串', () => {
      const exported = store.exportConfig();
      expect(typeof exported).toBe('string');
      expect(() => JSON.parse(exported)).not.toThrow();
    });

    it('返回的JSON应该包含所有配置节', () => {
      const exported = store.exportConfig();
      const parsed = JSON.parse(exported);
      
      expect(parsed).toHaveProperty('game');
      expect(parsed).toHaveProperty('player');
      expect(parsed).toHaveProperty('world');
      expect(parsed).toHaveProperty('combat');
      expect(parsed).toHaveProperty('economy');
    });

    it('返回的JSON应该是格式化的（包含换行和缩进）', () => {
      const exported = store.exportConfig();
      expect(exported).toContain('\n');
      expect(exported).toContain('  '); // 2个空格缩进
    });
  });

  describe('importConfig', () => {
    it('应该成功导入有效的JSON配置', () => {
      const configToImport = {
        game: { maxPlayers: 50 },
        player: { startGold: 2000 }
      };
      const jsonString = JSON.stringify(configToImport);
      
      const result = store.importConfig(jsonString);
      
      expect(result).toBe(true);
      expect(store.gameConfig.game.maxPlayers).toBe(50);
      expect(store.gameConfig.player.startGold).toBe(2000);
      // 未导入的属性应该保持原样
      expect(store.gameConfig.game.name).toBe('古风RPG');
    });

    it('对于无效的JSON应该返回false', () => {
      const invalidJson = '{invalid json}';
      const result = store.importConfig(invalidJson);
      
      expect(result).toBe(false);
      // 配置应该保持不变
      expect(store.gameConfig.game.maxPlayers).toBe(100);
    });

    it('应该合并导入的配置而不是完全替换', () => {
      const initialGold = store.gameConfig.player.startGold;
      const configToImport = {
        game: { version: '2.0.0' }
      };
      
      store.importConfig(JSON.stringify(configToImport));
      
      // 游戏版本应该更新
      expect(store.gameConfig.game.version).toBe('2.0.0');
      // 玩家金币应该保持不变
      expect(store.gameConfig.player.startGold).toBe(initialGold);
    });

    it('应该能够导入空对象而不修改配置', () => {
      const originalConfig = { ...store.gameConfig };
      const emptyConfig = '{}';
      
      store.importConfig(emptyConfig);
      
      expect(store.gameConfig).toEqual(originalConfig);
    });
  });

  describe('配置节合并行为', () => {
    it('导入配置时应该递归合并嵌套对象', () => {
      const configToImport = {
        player: {
          startPosition: { x: 500 }
        }
      };
      
      store.importConfig(JSON.stringify(configToImport));
      
      // 只有x坐标应该更新，y坐标应该保持原样
      expect(store.gameConfig.player.startPosition).toEqual({ x: 500, y: 300 });
    });

    it('updateConfig应该能够添加新属性', () => {
      store.updateConfig('game', { newProperty: 'newValue' });
      
      expect(store.gameConfig.game.newProperty).toBe('newValue');
      // 原有属性应该保持不变
      expect(store.gameConfig.game.name).toBe('古风RPG');
    });
  });

  describe('边界情况', () => {
    it('updateConfig应该处理null/undefined值', () => {
      // 这些情况不应该抛出错误
      expect(() => store.updateConfig('game', null)).not.toThrow();
      expect(() => store.updateConfig('game', undefined)).not.toThrow();
    });

    it('getConfig应该返回引用而不是副本', () => {
      const gameSection = store.getConfig('game');
      gameSection.maxPlayers = 999;
      
      // 由于返回的是引用，直接修改应该影响原始配置
      // 注意：这个测试可能会失败，取决于实现方式
      // 如果getConfig返回的是深拷贝，这个测试会失败
      expect(store.gameConfig.game.maxPlayers).toBe(999);
    });
  });
});
```
import { describe, it, expect, beforeEach, vi } from 'vitest';

vi.mock('phaser', () => ({
  default: {
    Scene: class MockScene {
      constructor(config) {
        this.sys = { settings: config };
      }
    },
    Math: {
      Between: (min, max) => min + Math.floor(Math.random() * (max - min + 1)),
      Distance: { Between: (x1, y1, x2, y2) => Math.sqrt((x2 - x1) ** 2 + (y2 - y1) ** 2) }
    },
    Input: { Keyboard: { KeyCodes: { W: 87, S: 83, A: 65, D: 68, SPACE: 32, I: 73, F5: 116, F9: 120, FORWARD_SLASH: 191, ESC: 27 } } }
  }
}));

vi.mock('../game/systems/InventoryManager.js', () => ({
  InventoryManager: vi.fn().mockImplementation(() => ({
    getStats: () => ({ hp: 100, max_hp: 100, mp: 50, max_mp: 50, attack: 10, defense: 5, speed: 10 }),
    parseInventory: vi.fn(),
    syncWithServer: vi.fn()
  }))
}));

vi.mock('../game/systems/CombatManager.js', () => ({
  CombatManager: vi.fn()
}));

vi.mock('../game/ui/InventoryUI.js', () => ({ InventoryUI: vi.fn() }));
vi.mock('../game/ui/CombatUI.js', () => ({ CombatUI: vi.fn() }));
vi.mock('../game/ui/MiniMap.js', () => {
  class MockMiniMap {
    constructor() { this.create = vi.fn(); this.update = vi.fn(); this.destroy = vi.fn(); }
  }
  return { MiniMap: MockMiniMap };
});
vi.mock('../game/ui/SaveLoadUI.js', () => ({ SaveLoadUI: vi.fn() }));

import { GameScene } from '../game/scenes/GameScene.js';

function createGameSceneMocks() {
  const tweens = { add: vi.fn() };
  const physics = {
    add: {
      sprite: vi.fn(() => ({
        setScale: vi.fn().mockReturnThis(),
        setCollideWorldBounds: vi.fn().mockReturnThis(),
        setImmovable: vi.fn().mockReturnThis(),
        setInteractive: vi.fn().mockReturnThis(),
        setData: vi.fn().mockReturnThis(),
        getData: vi.fn(),
        setPosition: vi.fn(),
        on: vi.fn()
      })),
      overlap: vi.fn()
    },
    world: { setBounds: vi.fn() }
  };
  const cameras = {
    main: {
      width: 800,
      height: 600,
      startFollow: vi.fn(),
      setZoom: vi.fn(),
      fadeOut: vi.fn(),
      fadeIn: vi.fn(),
      once: vi.fn()
    }
  };
  const add = {
    rectangle: vi.fn(() => ({ setScrollFactor: vi.fn().mockReturnThis(), setDepth: vi.fn().mockReturnThis(), setStrokeStyle: vi.fn().mockReturnThis(), setFillStyle: vi.fn().mockReturnThis(), setInteractive: vi.fn().mockReturnThis(), setOrigin: vi.fn().mockReturnThis(), on: vi.fn(), destroy: vi.fn() })),
    text: vi.fn(() => ({ setScrollFactor: vi.fn().mockReturnThis(), setDepth: vi.fn().mockReturnThis(), setOrigin: vi.fn().mockReturnThis(), setInteractive: vi.fn().mockReturnThis(), on: vi.fn(), setText: vi.fn(), setAlpha: vi.fn(), setPosition: vi.fn(), destroy: vi.fn() })),
    image: vi.fn(() => ({ setScale: vi.fn().mockReturnThis() })),
    circle: vi.fn(),
    container: vi.fn(() => ({ setDepth: vi.fn().mockReturnThis(), setScrollFactor: vi.fn().mockReturnThis(), add: vi.fn(), destroy: vi.fn() }))
  };
  const input = {
    keyboard: {
      createCursorKeys: vi.fn(() => ({})),
      addKey: vi.fn(() => ({ on: vi.fn() })),
      once: vi.fn(),
      off: vi.fn()
    }
  };
  const children = { removeAll: vi.fn() };
  const scene = { start: vi.fn() };
  const registry = { set: vi.fn(), get: vi.fn() };

  return { tweens, physics, cameras, add, input, children, scene, registry };
}

describe('GameScene', () => {
  let scene;
  let mocks;

  beforeEach(() => {
    scene = new GameScene();
    mocks = createGameSceneMocks();
    scene.tweens = mocks.tweens;
    scene.physics = mocks.physics;
    scene.cameras = mocks.cameras;
    scene.add = mocks.add;
    scene.input = mocks.input;
    scene.children = mocks.children;
    scene.scene = mocks.scene;
    scene.registry = mocks.registry;
  });

  describe('constructor', () => {
    it('sets scene key to GameScene', () => {
      expect(scene.sys.settings.key).toBe('GameScene');
    });

    it('initializes default state', () => {
      expect(scene.player).toBeNull();
      expect(scene.npcs).toEqual([]);
      expect(scene.portals).toEqual([]);
      expect(scene.gameData).toBeNull();
      expect(scene.playerData).toBeNull();
      expect(scene.currentSceneCode).toBeNull();
      expect(scene.tutorialStep).toBe(0);
      expect(scene.showingDialog).toBe(false);
      expect(scene.showingShop).toBe(false);
      expect(scene.showingTutorial).toBe(false);
      expect(scene.questLog).toEqual([]);
      expect(scene.inventory).toEqual({});
      expect(scene.visitedScenes.size).toBe(0);
      expect(scene.encounterCooldown).toBe(0);
    });

    it('initializes Maps for sprites', () => {
      expect(scene.npcSprites).toBeInstanceOf(Map);
      expect(scene.portalSprites).toBeInstanceOf(Map);
    });
  });

  describe('showTitleScreen', () => {
    it('creates title text elements', () => {
      scene.showTitleScreen();
      expect(mocks.add.text).toHaveBeenCalled();
      expect(mocks.add.rectangle).toHaveBeenCalled();
    });
  });

  describe('getNPCSpriteKey', () => {
    it('maps known NPC codes to sprite keys', () => {
      expect(scene.getNPCSpriteKey('npc_chief_chen')).toBe('npc_chief');
      expect(scene.getNPCSpriteKey('npc_merchant_li')).toBe('npc_merchant');
      expect(scene.getNPCSpriteKey('npc_tea_wang')).toBe('npc_tea');
      expect(scene.getNPCSpriteKey('npc_blacksmith_zhang')).toBe('npc_blacksmith');
      expect(scene.getNPCSpriteKey('npc_hunter_zhou')).toBe('npc_hunter');
      expect(scene.getNPCSpriteKey('npc_kid_stone')).toBe('npc_kid');
    });

    it('defaults to npc_merchant for unknown code', () => {
      expect(scene.getNPCSpriteKey('unknown_npc')).toBe('npc_merchant');
    });
  });

  describe('getShopCode', () => {
    it('maps merchant NPC to general store', () => {
      expect(scene.getShopCode('npc_merchant_li')).toBe('shop_general_store');
    });

    it('maps blacksmith NPC to blacksmith shop', () => {
      expect(scene.getShopCode('npc_blacksmith_zhang')).toBe('shop_blacksmith');
    });

    it('returns empty string for unknown NPC', () => {
      expect(scene.getShopCode('unknown')).toBe('');
    });
  });

  describe('getNPCGreeting', () => {
    it('returns neutral greeting by default', () => {
      scene.npcSprites = new Map();
      const greeting = scene.getNPCGreeting({ code: 'npc_chief_chen', id: 1 });
      expect(greeting).toContain('村长');
    });

    it('returns happy greeting when mood is happy', () => {
      const mockSprite = { getData: vi.fn(() => 'happy') };
      scene.npcSprites = new Map([[1, mockSprite]]);
      const greeting = scene.getNPCGreeting({ code: 'npc_chief_chen', id: 1 });
      expect(greeting).toContain('又见面了');
    });

    it('returns default greeting for unknown NPC', () => {
      scene.npcSprites = new Map();
      const greeting = scene.getNPCGreeting({ code: 'unknown', id: 99 });
      expect(greeting).toBe('你好，客官！');
    });

    it('returns greetings for all NPC codes', () => {
      scene.npcSprites = new Map();
      const codes = ['npc_chief_chen', 'npc_merchant_li', 'npc_tea_wang', 'npc_blacksmith_zhang', 'npc_hunter_zhou', 'npc_kid_stone'];
      codes.forEach(code => {
        const greeting = scene.getNPCGreeting({ code, id: 1 });
        expect(greeting).toBeTruthy();
        expect(greeting.length).toBeGreaterThan(0);
      });
    });
  });

  describe('hasQuestForNPC', () => {
    it('returns false when no active quests', () => {
      scene.questLog = [];
      expect(scene.hasQuestForNPC('npc_chief_chen')).toBe(false);
    });

    it('returns true when NPC has active quest objective', () => {
      scene.questLog = [{
        status: 'active',
        objectives: JSON.stringify([{ type: 'dialogue', target: 'npc_chief_chen', completed: false }])
      }];
      expect(scene.hasQuestForNPC('npc_chief_chen')).toBe(true);
    });

    it('returns false when quest objective is completed', () => {
      scene.questLog = [{
        status: 'active',
        objectives: JSON.stringify([{ type: 'dialogue', target: 'npc_chief_chen', completed: true }])
      }];
      expect(scene.hasQuestForNPC('npc_chief_chen')).toBe(false);
    });

    it('returns false for inactive quests', () => {
      scene.questLog = [{
        status: 'completed',
        objectives: JSON.stringify([{ type: 'dialogue', target: 'npc_chief_chen', completed: false }])
      }];
      expect(scene.hasQuestForNPC('npc_chief_chen')).toBe(false);
    });
  });

  describe('checkQuestComplete', () => {
    it('marks quest as completed when all objectives done', () => {
      const quest = {
        status: 'active',
        name: 'Test Quest',
        objectives: JSON.stringify([
          { type: 'visit', target: 'scene_village', completed: true },
          { type: 'dialogue', target: 'npc_chief_chen', completed: true }
        ]),
        rewards: JSON.stringify({ gold: 50, exp: 20 })
      };
      scene.questLog = [quest];
      scene.playerData = { gold: 100 };
      scene.checkQuestComplete(quest);
      expect(quest.status).toBe('completed');
    });

    it('does not mark incomplete quest as completed', () => {
      const quest = {
        status: 'active',
        name: 'Test Quest',
        objectives: JSON.stringify([
          { type: 'visit', target: 'scene_village', completed: true },
          { type: 'dialogue', target: 'npc_chief_chen', completed: false }
        ]),
        rewards: '{}'
      };
      scene.questLog = [quest];
      scene.playerData = { gold: 100 };
      scene.checkQuestComplete(quest);
      expect(quest.status).toBe('active');
    });
  });

  describe('updateVisitQuest', () => {
    it('completes visit objective for matching scene', () => {
      const quest = {
        status: 'active',
        objectives: JSON.stringify([
          { type: 'visit', target: 'scene_village', completed: false, description: 'Visit village' }
        ])
      };
      scene.questLog = [quest];
      scene.playerData = { gold: 0 };
      scene.updateVisitQuest('scene_village');
      const objectives = JSON.parse(quest.objectives);
      expect(objectives[0].completed).toBe(true);
    });

    it('ignores non-visit objectives', () => {
      const quest = {
        status: 'active',
        objectives: JSON.stringify([
          { type: 'dialogue', target: 'scene_village', completed: false }
        ])
      };
      scene.questLog = [quest];
      scene.updateVisitQuest('scene_village');
      const objectives = JSON.parse(quest.objectives);
      expect(objectives[0].completed).toBe(false);
    });

    it('ignores inactive quests', () => {
      const quest = {
        status: 'completed',
        objectives: JSON.stringify([
          { type: 'visit', target: 'scene_village', completed: false }
        ])
      };
      scene.questLog = [quest];
      scene.updateVisitQuest('scene_village');
      const objectives = JSON.parse(quest.objectives);
      expect(objectives[0].completed).toBe(false);
    });
  });

  describe('checkRandomEncounter', () => {
    it('does nothing when not in village_path', () => {
      scene.currentSceneCode = 'scene_village';
      scene.checkRandomEncounter();
      expect(scene.encounterCooldown).toBe(0);
    });

    it('does nothing when cooldown is active', () => {
      scene.currentSceneCode = 'village_path';
      scene.encounterCooldown = 100;
      scene.checkRandomEncounter();
      expect(scene.encounterCooldown).toBe(100);
    });

    it('does nothing when showing dialog', () => {
      scene.currentSceneCode = 'village_path';
      scene.showingDialog = true;
      scene.checkRandomEncounter();
      expect(scene.encounterCooldown).toBe(0);
    });
  });

  describe('loadScene', () => {
    it('clears existing sprites', () => {
      scene.gameData = { scenes: [], npcs: [] };
      scene.npcSprites.set(1, { destroy: vi.fn() });
      scene.portalSprites.set(1, { destroy: vi.fn() });
      scene.loadScene('nonexistent');
      expect(scene.npcSprites.size).toBe(0);
      expect(scene.portalSprites.size).toBe(0);
    });

    it('returns early for unknown scene', () => {
      scene.gameData = { scenes: [{ code: 'scene_a', name: 'A', description: 'Desc', width: 800, height: 600 }], npcs: [] };
      scene.loadScene('scene_b');
      expect(scene.currentSceneCode).toBeNull();
    });

    it('sets currentSceneCode for valid scene', () => {
      scene.gameData = {
        scenes: [{ code: 'scene_a', name: 'A', description: 'Desc', width: 800, height: 600 }],
        npcs: []
      };
      scene.playerData = { pos_x: 100, pos_y: 200 };
      scene.loadScene('scene_a');
      expect(scene.currentSceneCode).toBe('scene_a');
      expect(scene.visitedScenes.has('scene_a')).toBe(true);
    });
  });
});

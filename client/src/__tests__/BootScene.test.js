import { describe, it, expect, beforeEach, vi } from 'vitest';

function makeMockCanvas() {
  const ctx = {
    fillStyle: '',
    strokeStyle: '',
    lineWidth: 1,
    globalAlpha: 1,
    font: '',
    textAlign: '',
    fillRect: vi.fn(),
    strokeRect: vi.fn(),
    beginPath: vi.fn(),
    closePath: vi.fn(),
    moveTo: vi.fn(),
    lineTo: vi.fn(),
    arc: vi.fn(),
    ellipse: vi.fn(),
    fill: vi.fn(),
    stroke: vi.fn(),
    clearRect: vi.fn(),
    fillText: vi.fn()
  };
  return {
    getContext: () => ctx,
    refresh: vi.fn(),
    _ctx: ctx
  };
}

function createBootSceneMocks() {
  const mockCanvas = makeMockCanvas();
  const textures = {
    createCanvas: vi.fn(() => mockCanvas)
  };
  const cameras = { main: { width: 800, height: 600 } };
  const add = {
    graphics: vi.fn(() => ({ clear: vi.fn(), fillStyle: vi.fn(), fillRect: vi.fn(), destroy: vi.fn() })),
    text: vi.fn(() => ({ setOrigin: vi.fn() }))
  };
  const loadOn = {};
  const load = { on: vi.fn((evt, cb) => { loadOn[evt] = cb; }) };
  const registry = { set: vi.fn() };
  const scene = { start: vi.fn() };

  return { textures, cameras, add, load, loadOn, registry, scene, mockCanvas };
}

vi.mock('phaser', () => {
  return {
    default: {
      Scene: class MockScene {
        constructor(config) {
          this.sys = { settings: config };
        }
      },
      Math: { Between: (min, max) => min + Math.floor(Math.random() * (max - min + 1)) },
      Input: { Keyboard: { KeyCodes: { W: 87, S: 83, A: 65, D: 68, SPACE: 32, I: 73, F5: 116, F9: 120, FORWARD_SLASH: 191, ESC: 27 } } }
    }
  };
});

import { BootScene } from '../game/scenes/BootScene.js';

describe('BootScene', () => {
  let scene;
  let mocks;

  beforeEach(() => {
    scene = new BootScene();
    mocks = createBootSceneMocks();
    scene.textures = mocks.textures;
    scene.cameras = mocks.cameras;
    scene.add = mocks.add;
    scene.load = mocks.load;
    scene.registry = mocks.registry;
    scene.scene = mocks.scene;
  });

  describe('constructor', () => {
    it('sets scene key to BootScene', () => {
      expect(scene.sys.settings.key).toBe('BootScene');
    });
  });

  describe('createTextures', () => {
    it('creates player texture canvas', () => {
      scene.createTextures();
      expect(mocks.textures.createCanvas).toHaveBeenCalledWith('player', 32, 48);
    });

    it('creates NPC textures', () => {
      scene.createTextures();
      const npcKeys = ['npc_chief', 'npc_merchant', 'npc_tea', 'npc_blacksmith', 'npc_hunter', 'npc_kid'];
      npcKeys.forEach(key => {
        expect(mocks.textures.createCanvas).toHaveBeenCalledWith(key, 32, 48);
      });
    });

    it('creates tile textures', () => {
      scene.createTextures();
      const tileKeys = ['tile_grass', 'tile_dirt', 'tile_stone', 'tile_wood', 'tile_water'];
      tileKeys.forEach(key => {
        expect(mocks.textures.createCanvas).toHaveBeenCalledWith(key, 48, 48);
      });
    });

    it('creates portal texture', () => {
      scene.createTextures();
      expect(mocks.textures.createCanvas).toHaveBeenCalledWith('portal', 48, 48);
    });

    it('creates tree texture', () => {
      scene.createTextures();
      expect(mocks.textures.createCanvas).toHaveBeenCalledWith('tree', 48, 64);
    });

    it('creates house texture', () => {
      scene.createTextures();
      expect(mocks.textures.createCanvas).toHaveBeenCalledWith('house', 96, 80);
    });

    it('creates indicator texture', () => {
      scene.createTextures();
      expect(mocks.textures.createCanvas).toHaveBeenCalledWith('indicator', 24, 24);
    });

    it('creates enemy textures', () => {
      scene.createTextures();
      expect(mocks.textures.createCanvas).toHaveBeenCalledWith('enemy_wolf', 48, 48);
      expect(mocks.textures.createCanvas).toHaveBeenCalledWith('enemy_alpha_wolf', 48, 48);
    });

    it('refreshes all canvas textures', () => {
      scene.createTextures();
      expect(mocks.mockCanvas.refresh.mock.calls.length).toBeGreaterThan(0);
    });
  });

  describe('createNPCTexture', () => {
    it('creates NPC canvas with correct dimensions', () => {
      scene.createNPCTexture('test_npc', '#ff0000', '#fff', '#000');
      expect(mocks.textures.createCanvas).toHaveBeenCalledWith('test_npc', 32, 48);
    });

    it('calls refresh on canvas', () => {
      scene.createNPCTexture('test_npc', '#ff0000', '#fff', '#000');
      expect(mocks.mockCanvas.refresh).toHaveBeenCalled();
    });
  });

  describe('createTileTexture', () => {
    it('creates tile canvas with correct dimensions', () => {
      scene.createTileTexture('test_tile', '#abc');
      expect(mocks.textures.createCanvas).toHaveBeenCalledWith('test_tile', 48, 48);
    });

    it('draws fill and stroke rectangles', () => {
      scene.createTileTexture('test_tile', '#abc');
      const ctx = mocks.mockCanvas._ctx;
      expect(ctx.fillRect).toHaveBeenCalled();
      expect(ctx.strokeRect).toHaveBeenCalled();
    });
  });

  describe('createEnemyTexture', () => {
    it('creates enemy canvas with correct dimensions', () => {
      scene.createEnemyTexture('test_enemy', '#aaa', '#fff');
      expect(mocks.textures.createCanvas).toHaveBeenCalledWith('test_enemy', 48, 48);
    });

    it('draws body ellipse and other parts', () => {
      scene.createEnemyTexture('test_enemy', '#aaa', '#fff');
      const ctx = mocks.mockCanvas._ctx;
      expect(ctx.ellipse).toHaveBeenCalled();
      expect(ctx.arc).toHaveBeenCalled();
    });
  });

  describe('create', () => {
    it('initializes registry values', () => {
      scene.create();
      expect(mocks.registry.set).toHaveBeenCalledWith('gameData', null);
      expect(mocks.registry.set).toHaveBeenCalledWith('playerData', null);
      expect(mocks.registry.set).toHaveBeenCalledWith('currentScene', null);
    });

    it('starts GameScene', () => {
      scene.create();
      expect(mocks.scene.start).toHaveBeenCalledWith('GameScene');
    });
  });

  describe('preload', () => {
    it('creates loading bar graphics', () => {
      scene.preload();
      expect(mocks.add.graphics).toHaveBeenCalledTimes(2);
    });

    it('creates loading text', () => {
      scene.preload();
      expect(mocks.add.text).toHaveBeenCalled();
    });

    it('registers progress and complete listeners', () => {
      scene.preload();
      expect(mocks.load.on).toHaveBeenCalledWith('progress', expect.any(Function));
      expect(mocks.load.on).toHaveBeenCalledWith('complete', expect.any(Function));
    });

    it('calls createTextures', () => {
      const spy = vi.spyOn(scene, 'createTextures');
      scene.preload();
      expect(spy).toHaveBeenCalled();
    });
  });
});

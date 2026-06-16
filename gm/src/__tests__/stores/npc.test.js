import { describe, it, expect, beforeEach } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { useNPCStore } from '@/stores/npc';

describe('NPC Store', () => {
  let store;

  beforeEach(() => {
    setActivePinia(createPinia());
    store = useNPCStore();
  });

  it('should have initial npcs', () => {
    expect(store.npcs).toHaveLength(3);
    expect(store.npcs[0].id).toBe('npc_001');
    expect(store.npcs[1].id).toBe('npc_002');
    expect(store.npcs[2].id).toBe('npc_003');
  });

  it('should add a new npc', () => {
    const initialLength = store.npcs.length;
    store.addNPC({
      name: '新NPC',
      title: '测试NPC',
      description: '测试描述',
      avatar: 'test.png',
      sprite: 'test_sprite.png',
      position: { x: 100, y: 200, scene: 'scene_001' },
      agentId: 'agent_001',
      dialogues: [],
      shopId: null,
      behaviors: ['idle'],
      schedule: []
    });
    expect(store.npcs).toHaveLength(initialLength + 1);
    expect(store.npcs[store.npcs.length - 1].name).toBe('新NPC');
    expect(store.npcs[store.npcs.length - 1].id).toMatch(/^npc_\d+$/);
  });

  it('should update an existing npc', () => {
    store.updateNPC('npc_001', { name: '更新后的NPC' });
    const npc = store.getNPCById('npc_001');
    expect(npc.name).toBe('更新后的NPC');
  });

  it('should not update non-existent npc', () => {
    const initialLength = store.npcs.length;
    store.updateNPC('non_existent', { name: 'test' });
    expect(store.npcs).toHaveLength(initialLength);
  });

  it('should delete an npc', () => {
    const initialLength = store.npcs.length;
    store.deleteNPC('npc_001');
    expect(store.npcs).toHaveLength(initialLength - 1);
    expect(store.getNPCById('npc_001')).toBeUndefined();
  });

  it('should get npc by id', () => {
    const npc = store.getNPCById('npc_002');
    expect(npc).toBeDefined();
    expect(npc.name).toBe('王大娘');
  });

  it('should return undefined for non-existent npc', () => {
    const npc = store.getNPCById('non_existent');
    expect(npc).toBeUndefined();
  });

  it('should have currentNPC initially null', () => {
    expect(store.currentNPC).toBeNull();
  });
});

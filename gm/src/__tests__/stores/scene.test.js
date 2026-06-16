import { describe, it, expect, beforeEach } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { useSceneStore } from '@/stores/scene';

describe('Scene Store', () => {
  let store;

  beforeEach(() => {
    setActivePinia(createPinia());
    store = useSceneStore();
  });

  it('should have initial scenes', () => {
    expect(store.scenes).toHaveLength(2);
    expect(store.scenes[0].id).toBe('scene_001');
    expect(store.scenes[1].id).toBe('scene_002');
  });

  it('should add a new scene', () => {
    const initialLength = store.scenes.length;
    store.addScene({
      name: '新场景',
      description: '测试场景',
      background: 'test.png',
      width: 1920,
      height: 1080,
      npcs: [],
      portals: [],
      tiles: []
    });
    expect(store.scenes).toHaveLength(initialLength + 1);
    expect(store.scenes[store.scenes.length - 1].name).toBe('新场景');
    expect(store.scenes[store.scenes.length - 1].id).toMatch(/^scene_\d+$/);
  });

  it('should update an existing scene', () => {
    store.updateScene('scene_001', { name: '更新后的场景' });
    const scene = store.getSceneById('scene_001');
    expect(scene.name).toBe('更新后的场景');
  });

  it('should not update non-existent scene', () => {
    const initialLength = store.scenes.length;
    store.updateScene('non_existent', { name: 'test' });
    expect(store.scenes).toHaveLength(initialLength);
  });

  it('should delete a scene', () => {
    const initialLength = store.scenes.length;
    store.deleteScene('scene_001');
    expect(store.scenes).toHaveLength(initialLength - 1);
    expect(store.getSceneById('scene_001')).toBeUndefined();
  });

  it('should get scene by id', () => {
    const scene = store.getSceneById('scene_002');
    expect(scene).toBeDefined();
    expect(scene.name).toBe('杂货铺');
  });

  it('should return undefined for non-existent scene', () => {
    const scene = store.getSceneById('non_existent');
    expect(scene).toBeUndefined();
  });

  it('should have currentScene initially null', () => {
    expect(store.currentScene).toBeNull();
  });
});

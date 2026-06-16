import { describe, it, expect, beforeEach } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { useShopStore } from '@/stores/shop';

describe('Shop Store', () => {
  let store;

  beforeEach(() => {
    setActivePinia(createPinia());
    store = useShopStore();
  });

  it('should have initial shops', () => {
    expect(store.shops).toHaveLength(2);
    expect(store.shops[0].id).toBe('shop_001');
    expect(store.shops[1].id).toBe('shop_002');
  });

  it('should have initial items', () => {
    expect(store.items).toHaveLength(8);
    expect(store.items[0].id).toBe('item_001');
    expect(store.items[7].id).toBe('item_012');
  });

  it('should add a new shop', () => {
    const initialLength = store.shops.length;
    store.addShop({
      name: '新商店',
      type: 'general',
      description: '测试商店',
      owner: 'npc_001',
      scene: 'scene_001',
      items: [],
      openTime: '08:00',
      closeTime: '20:00',
      discount: { threshold: 3, rate: 0.9 }
    });
    expect(store.shops).toHaveLength(initialLength + 1);
    expect(store.shops[store.shops.length - 1].name).toBe('新商店');
    expect(store.shops[store.shops.length - 1].id).toMatch(/^shop_\d+$/);
  });

  it('should update an existing shop', () => {
    store.updateShop('shop_001', { name: '更新后的商店' });
    const shop = store.getShopById('shop_001');
    expect(shop.name).toBe('更新后的商店');
  });

  it('should not update non-existent shop', () => {
    const initialLength = store.shops.length;
    store.updateShop('non_existent', { name: 'test' });
    expect(store.shops).toHaveLength(initialLength);
  });

  it('should delete a shop', () => {
    const initialLength = store.shops.length;
    store.deleteShop('shop_001');
    expect(store.shops).toHaveLength(initialLength - 1);
    expect(store.getShopById('shop_001')).toBeUndefined();
  });

  it('should get shop by id', () => {
    const shop = store.getShopById('shop_002');
    expect(shop).toBeDefined();
    expect(shop.name).toBe('张记铁匠铺');
  });

  it('should return undefined for non-existent shop', () => {
    const shop = store.getShopById('non_existent');
    expect(shop).toBeUndefined();
  });

  it('should add a new item', () => {
    const initialLength = store.items.length;
    store.addItem({
      name: '新道具',
      category: 'medicine',
      description: '测试道具',
      effect: { hp: 50 },
      icon: 'test.png'
    });
    expect(store.items).toHaveLength(initialLength + 1);
    expect(store.items[store.items.length - 1].name).toBe('新道具');
    expect(store.items[store.items.length - 1].id).toMatch(/^item_\d+$/);
  });

  it('should update an existing item', () => {
    store.updateItem('item_001', { name: '更新后的道具' });
    const item = store.getItemById('item_001');
    expect(item.name).toBe('更新后的道具');
  });

  it('should not update non-existent item', () => {
    const initialLength = store.items.length;
    store.updateItem('non_existent', { name: 'test' });
    expect(store.items).toHaveLength(initialLength);
  });

  it('should delete an item', () => {
    const initialLength = store.items.length;
    store.deleteItem('item_001');
    expect(store.items).toHaveLength(initialLength - 1);
    expect(store.getItemById('item_001')).toBeUndefined();
  });

  it('should get item by id', () => {
    const item = store.getItemById('item_002');
    expect(item).toBeDefined();
    expect(item.name).toBe('灵芝');
  });

  it('should return undefined for non-existent item', () => {
    const item = store.getItemById('non_existent');
    expect(item).toBeUndefined();
  });
});

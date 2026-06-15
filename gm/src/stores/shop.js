import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useShopStore = defineStore('shop', () => {
  const shops = ref([
    {
      id: 'shop_001',
      name: '李记杂货铺',
      type: 'general',
      description: '售卖各种日常用品、药材、食材',
      owner: 'npc_001',
      scene: 'scene_002',
      items: [
        { itemId: 'item_001', price: 100, stock: 50 },
        { itemId: 'item_002', price: 200, stock: 30 },
        { itemId: 'item_003', price: 50, stock: 100 }
      ],
      openTime: '06:00',
      closeTime: '22:00',
      discount: {
        threshold: 3,
        rate: 0.9
      }
    },
    {
      id: 'shop_002',
      name: '张记铁匠铺',
      type: 'blacksmith',
      description: '打造和售卖各种兵器、工具',
      owner: 'npc_003',
      scene: 'scene_002',
      items: [
        { itemId: 'item_010', price: 500, stock: 10 },
        { itemId: 'item_011', price: 800, stock: 5 },
        { itemId: 'item_012', price: 300, stock: 20 }
      ],
      openTime: '08:00',
      closeTime: '20:00',
      discount: {
        threshold: 2,
        rate: 0.95
      }
    }
  ]);

  const items = ref([
    { id: 'item_001', name: '草药', category: 'medicine', description: '普通的草药，可恢复少量生命', effect: { hp: 20 }, icon: 'herb.png' },
    { id: 'item_002', name: '灵芝', category: 'medicine', description: '珍贵的药材，可恢复大量生命', effect: { hp: 100 }, icon: 'lingzhi.png' },
    { id: 'item_003', name: '馒头', category: 'food', description: '普通的干粮，可恢复少量体力', effect: { stamina: 10 }, icon: 'mantou.png' },
    { id: 'item_004', name: '烧酒', category: 'food', description: '烈性白酒，可驱寒保暖', effect: { cold_resist: 30 }, icon: 'wine.png' },
    { id: 'item_005', name: '麻绳', category: 'tool', description: '结实的麻绳，可用于攀爬', icon: 'rope.png' },
    { id: 'item_010', name: '铁剑', category: 'weapon', description: '普通的铁制长剑', effect: { attack: 10 }, icon: 'iron_sword.png' },
    { id: 'item_011', name: '精钢刀', category: 'weapon', description: '锻造精良的钢刀', effect: { attack: 25 }, icon: 'steel_blade.png' },
    { id: 'item_012', name: '铁甲', category: 'armor', description: '普通的铁制铠甲', effect: { defense: 15 }, icon: 'iron_armor.png' }
  ]);

  const addShop = (shop) => {
    shops.value.push({
      id: `shop_${Date.now()}`,
      ...shop
    });
  };

  const updateShop = (id, data) => {
    const index = shops.value.findIndex(s => s.id === id);
    if (index !== -1) {
      shops.value[index] = { ...shops.value[index], ...data };
    }
  };

  const deleteShop = (id) => {
    shops.value = shops.value.filter(s => s.id !== id);
  };

  const addItem = (item) => {
    items.value.push({
      id: `item_${Date.now()}`,
      ...item
    });
  };

  const updateItem = (id, data) => {
    const index = items.value.findIndex(i => i.id === id);
    if (index !== -1) {
      items.value[index] = { ...items.value[index], ...data };
    }
  };

  const deleteItem = (id) => {
    items.value = items.value.filter(i => i.id !== id);
  };

  const getItemById = (id) => {
    return items.value.find(i => i.id === id);
  };

  const getShopById = (id) => {
    return shops.value.find(s => s.id === id);
  };

  return {
    shops,
    items,
    addShop,
    updateShop,
    deleteShop,
    addItem,
    updateItem,
    deleteItem,
    getItemById,
    getShopById
  };
});

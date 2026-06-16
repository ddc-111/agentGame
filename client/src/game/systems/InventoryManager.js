export class InventoryManager {
    constructor(gameData, playerData) {
        this.gameData = gameData;
        this.playerData = playerData;
        this.items = {};
        this.equipment = { weapon: null, armor: null, shield: null };
        this.parseInventory();
    }

    parseInventory() {
        try {
            // 解析道具 - 服务器存储格式为 {item_id: count} 扁平格式
            const raw = JSON.parse(this.playerData.items || '{}');
            if (raw.items) {
                // 兼容旧的嵌套格式 {items: {...}, equipment: {...}}
                this.items = raw.items;
                this.equipment = raw.equipment || { weapon: null, armor: null, shield: null };
            } else {
                // 服务器扁平格式 {item_id: count}
                this.items = raw;
            }

            // 解析装备 - 服务器单独存储 equipment 字段
            if (this.playerData.equipment) {
                try {
                    const equip = JSON.parse(this.playerData.equipment);
                    // 服务器格式 {weapon_id: 8, armor_id: 0} -> 客户端格式 {weapon: "8", armor: null}
                    this.equipment = {
                        weapon: equip.weapon_id ? String(equip.weapon_id) : null,
                        armor: equip.armor_id ? String(equip.armor_id) : null,
                        shield: null
                    };
                } catch (e) {}
            }
        } catch (e) {
            this.items = {};
            this.equipment = { weapon: null, armor: null, shield: null };
        }
    }

    saveInventory() {
        // 保存道具为服务器扁平格式 {item_id: count}
        this.playerData.items = JSON.stringify(this.items);
    }

    getItems() {
        return Object.entries(this.items).map(([id, count]) => {
            const itemData = this.getItemData(parseInt(id));
            return {
                id: parseInt(id),
                count,
                ...itemData
            };
        }).filter(item => item.count > 0);
    }

    getItemData(itemId) {
        const items = this.gameData.items || [];
        return items.find(i => i.id === itemId) || {
            id: itemId,
            name: '未知物品',
            description: '',
            type: 'misc',
            effect: {},
            icon: '📦'
        };
    }

    addItem(itemId, count = 1) {
        itemId = String(itemId);
        this.items[itemId] = (this.items[itemId] || 0) + count;
        this.saveInventory();
        return true;
    }

    removeItem(itemId, count = 1) {
        itemId = String(itemId);
        if (!this.items[itemId] || this.items[itemId] < count) {
            return false;
        }
        this.items[itemId] -= count;
        if (this.items[itemId] <= 0) {
            delete this.items[itemId];
        }
        this.saveInventory();
        return true;
    }

    getItemCount(itemId) {
        return this.items[String(itemId)] || 0;
    }

    equipItem(itemId) {
        const itemData = this.getItemData(itemId);
        if (!itemData || !itemData.type) return false;

        let slot = null;
        if (itemData.type === 'weapon') slot = 'weapon';
        else if (itemData.type === 'armor') slot = 'armor';
        else if (itemData.type === 'shield') slot = 'shield';
        else return false;

        if (this.getItemCount(itemId) <= 0) return false;

        const currentEquipped = this.equipment[slot];
        if (currentEquipped) {
            this.addItem(currentEquipped, 1);
        }

        this.removeItem(itemId, 1);
        this.equipment[slot] = itemId;
        this.saveInventory();
        return true;
    }

    unequipItem(slot) {
        if (!this.equipment[slot]) return false;
        const itemId = this.equipment[slot];
        this.addItem(itemId, 1);
        this.equipment[slot] = null;
        this.saveInventory();
        return true;
    }

    getEquipment() {
        const result = {};
        for (const [slot, itemId] of Object.entries(this.equipment)) {
            if (itemId) {
                result[slot] = {
                    id: itemId,
                    ...this.getItemData(itemId)
                };
            } else {
                result[slot] = null;
            }
        }
        return result;
    }

    getStats() {
        const base = {
            hp: this.playerData.hp || 100,
            max_hp: this.playerData.max_hp || 100,
            mp: this.playerData.mp || 50,
            max_mp: this.playerData.max_mp || 50,
            attack: this.playerData.attack || 10,
            defense: this.playerData.defense || 5,
            speed: this.playerData.speed || 10
        };

        const equipmentBonus = { attack: 0, defense: 0, hp: 0, mp: 0, speed: 0 };
        for (const slot of ['weapon', 'armor', 'shield']) {
            const itemId = this.equipment[slot];
            if (itemId) {
                const itemData = this.getItemData(itemId);
                if (itemData.effect) {
                    let effect = itemData.effect;
                    if (typeof effect === 'string') {
                        try { effect = JSON.parse(effect); } catch (e) { effect = {}; }
                    }
                    equipmentBonus.attack += effect.attack || 0;
                    equipmentBonus.defense += effect.defense || 0;
                    equipmentBonus.hp += effect.hp || 0;
                    equipmentBonus.mp += effect.mp || 0;
                    equipmentBonus.speed += effect.speed || 0;
                }
            }
        }

        return {
            hp: base.hp,
            max_hp: base.max_hp + equipmentBonus.hp,
            mp: base.mp,
            max_mp: base.max_mp + equipmentBonus.mp,
            attack: base.attack + equipmentBonus.attack,
            defense: base.defense + equipmentBonus.defense,
            speed: base.speed + equipmentBonus.speed
        };
    }

    useItem(itemId) {
        const itemData = this.getItemData(itemId);
        if (!itemData) return null;

        if (itemData.type !== 'consumable') return null;
        if (this.getItemCount(itemId) <= 0) return null;

        let effect = itemData.effect;
        if (typeof effect === 'string') {
            try { effect = JSON.parse(effect); } catch (e) { return null; }
        }

        this.removeItem(itemId, 1);

        const result = { item: itemData, effect: {} };
        if (effect.hp) {
            const healed = Math.min(effect.hp, this.getStats().max_hp - this.playerData.hp);
            this.playerData.hp = Math.min(this.playerData.hp + effect.hp, this.getStats().max_hp);
            result.effect.hp = healed;
        }
        if (effect.mp) {
            const restored = Math.min(effect.mp, this.getStats().max_mp - this.playerData.mp);
            this.playerData.mp = Math.min(this.playerData.mp + effect.mp, this.getStats().max_mp);
            result.effect.mp = restored;
        }

        return result;
    }

    // 从服务器同步背包数据
    async syncWithServer(playerId) {
        try {
            const resp = await fetch(`/api/inventory/${playerId}`);
            const data = await resp.json();
            if (data.items) {
                // 更新本地道具数据
                this.items = {};
                data.items.forEach(item => {
                    this.items[String(item.item_id)] = item.count;
                });
            }
            if (data.equipment) {
                this.equipment = {
                    weapon: data.equipment.weapon_id ? String(data.equipment.weapon_id) : null,
                    armor: data.equipment.armor_id ? String(data.equipment.armor_id) : null,
                    shield: null
                };
            }
            // 更新玩家数据
            if (data.gold !== undefined) {
                this.playerData.gold = data.gold;
            }
            this.saveInventory();
            return true;
        } catch (e) {
            console.error('同步背包失败:', e);
            return false;
        }
    }

    getItemIcon(item) {
        const iconMap = {
            'weapon': '⚔️',
            'armor': '🛡️',
            'shield': '🛡️',
            'consumable': '🧪',
            'material': '📦',
            'misc': '📦'
        };
        return item.icon || iconMap[item.type] || '📦';
    }
}
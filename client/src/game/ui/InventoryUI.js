export class InventoryUI {
    constructor(scene, inventoryManager) {
        this.scene = scene;
        this.inventoryManager = inventoryManager;
        this.container = null;
        this.selectedItem = null;
        this.isOpen = false;
    }

    open() {
        if (this.isOpen) return;
        this.isOpen = true;
        this.createUI();
    }

    close() {
        if (!this.isOpen) return;
        this.isOpen = false;
        if (this.container) {
            this.container.destroy();
            this.container = null;
        }
        if (this.htmlOverlay) {
            this.htmlOverlay.remove();
            this.htmlOverlay = null;
        }
    }

    toggle() {
        if (this.isOpen) this.close();
        else this.open();
    }

    createUI() {
        const width = this.scene.cameras.main.width;
        const height = this.scene.cameras.main.height;

        this.container = this.scene.add.container(0, 0).setDepth(600).setScrollFactor(0);

        const overlay = this.scene.add.rectangle(width / 2, height / 2, width, height, 0x000000, 0.7);
        this.container.add(overlay);

        const panelWidth = 600;
        const panelHeight = 450;
        const panelX = width / 2;
        const panelY = height / 2;

        const bg = this.scene.add.rectangle(panelX, panelY, panelWidth, panelHeight, 0x1a1a1a, 0.95);
        this.container.add(bg);

        const border = this.scene.add.rectangle(panelX, panelY, panelWidth, panelHeight).setStrokeStyle(2, 0xd4a574);
        this.container.add(border);

        const title = this.scene.add.text(panelX, panelY - 200, '🎒 背包', {
            font: 'bold 22px Microsoft YaHei', fill: '#d4a574'
        }).setOrigin(0.5);
        this.container.add(title);

        const stats = this.inventoryManager.getStats();
        const goldText = this.scene.add.text(panelX + 200, panelY - 200, `💰 ${this.scene.playerData.gold}`, {
            font: '16px Microsoft YaHei', fill: '#f1c40f'
        }).setOrigin(0.5);
        this.container.add(goldText);

        const closeBtn = this.scene.add.text(panelX + 280, panelY - 205, '✕', {
            font: '20px Microsoft YaHei', fill: '#ff6b6b'
        }).setOrigin(0.5).setInteractive({ useHandCursor: true });
        closeBtn.on('pointerdown', () => this.close());
        this.container.add(closeBtn);

        this.createEquipmentSlots(panelX - 200, panelY - 150);
        this.createItemGrid(panelX - 120, panelY - 140);
        this.createDetailPanel(panelX + 120, panelY - 140);
        this.createStatsPanel(panelX - 200, panelY + 120);

        this.scene.input.keyboard.once('keydown-ESC', () => this.close());
    }

    createEquipmentSlots(x, y) {
        const slots = [
            { key: 'weapon', label: '武器', icon: '⚔️' },
            { key: 'armor', label: '护甲', icon: '🛡️' },
            { key: 'shield', label: '盾牌', icon: '🛡️' }
        ];

        const equipment = this.inventoryManager.getEquipment();

        const title = this.scene.add.text(x, y - 30, '装备栏', {
            font: 'bold 14px Microsoft YaHei', fill: '#d4a574'
        });
        this.container.add(title);

        slots.forEach((slot, i) => {
            const slotY = y + i * 60;
            const slotBg = this.scene.add.rectangle(x, slotY, 150, 50, 0x2a2a2a).setStrokeStyle(1, 0x555555);
            this.container.add(slotBg);

            const equipped = equipment[slot.key];
            const iconText = this.scene.add.text(x - 60, slotY, equipped ? this.inventoryManager.getItemIcon(equipped) : slot.icon, {
                font: '20px Microsoft YaHei', fill: '#fff'
            }).setOrigin(0.5);
            this.container.add(iconText);

            const nameText = this.scene.add.text(x - 30, slotY - 8, slot.label, {
                font: '11px Microsoft YaHei', fill: '#888'
            });
            this.container.add(nameText);

            const equippedName = this.scene.add.text(x - 30, slotY + 8, equipped ? equipped.name : '未装备', {
                font: '12px Microsoft YaHei', fill: equipped ? '#4a7c59' : '#666'
            });
            this.container.add(equippedName);

            if (equipped) {
                slotBg.setInteractive({ useHandCursor: true });
                slotBg.on('pointerdown', () => {
                    this.inventoryManager.unequipItem(slot.key);
                    this.refreshUI();
                });
            }
        });
    }

    createItemGrid(startX, startY) {
        const items = this.inventoryManager.getItems();
        const cols = 4;
        const cellSize = 55;
        const maxVisible = 12;

        const title = this.scene.add.text(startX, startY - 30, '物品', {
            font: 'bold 14px Microsoft YaHei', fill: '#d4a574'
        });
        this.container.add(title);

        const gridBg = this.scene.add.rectangle(startX + 100, startY + 85, 230, 180, 0x222222);
        this.container.add(gridBg);

        for (let row = 0; row < 3; row++) {
            for (let col = 0; col < cols; col++) {
                const idx = row * cols + col;
                const cellX = startX + col * cellSize;
                const cellY = startY + row * cellSize;

                const cell = this.scene.add.rectangle(cellX + 22, cellY + 22, 48, 48, 0x333333).setStrokeStyle(1, 0x555555);
                this.container.add(cell);

                if (idx < items.length) {
                    const item = items[idx];
                    const icon = this.scene.add.text(cellX + 22, cellY + 18, this.inventoryManager.getItemIcon(item), {
                        font: '20px Microsoft YaHei'
                    }).setOrigin(0.5);
                    this.container.add(icon);

                    if (item.count > 1) {
                        const countText = this.scene.add.text(cellX + 38, cellY + 35, `${item.count}`, {
                            font: '10px Microsoft YaHei', fill: '#fff', backgroundColor: '#000000aa', padding: { x: 2, y: 1 }
                        }).setOrigin(0.5);
                        this.container.add(countText);
                    }

                    cell.setInteractive({ useHandCursor: true });
                    cell.on('pointerdown', () => {
                        this.selectedItem = item;
                        this.updateDetailPanel(item);
                    });
                }
            }
        }
    }

    createDetailPanel(x, y) {
        const title = this.scene.add.text(x, y - 30, '物品详情', {
            font: 'bold 14px Microsoft YaHei', fill: '#d4a574'
        });
        this.container.add(title);

        const panelBg = this.scene.add.rectangle(x + 50, y + 85, 200, 180, 0x222222);
        this.container.add(panelBg);

        this.detailName = this.scene.add.text(x + 10, y + 10, '选择一个物品', {
            font: '14px Microsoft YaHei', fill: '#fff'
        });
        this.container.add(this.detailName);

        this.detailDesc = this.scene.add.text(x + 10, y + 35, '', {
            font: '12px Microsoft YaHei', fill: '#aaa', wordWrap: { width: 170 }
        });
        this.container.add(this.detailDesc);

        this.detailEffect = this.scene.add.text(x + 10, y + 80, '', {
            font: '12px Microsoft YaHei', fill: '#4a7c59'
        });
        this.container.add(this.detailEffect);

        this.useBtn = this.scene.add.text(x + 10, y + 130, '使用', {
            font: '14px Microsoft YaHei', fill: '#4a7c59', backgroundColor: '#2a4a3a', padding: { x: 10, y: 5 }
        }).setInteractive({ useHandCursor: true }).setVisible(false);
        this.useBtn.on('pointerdown', () => this.useSelectedItem());
        this.container.add(this.useBtn);

        this.equipBtn = this.scene.add.text(x + 70, y + 130, '装备', {
            font: '14px Microsoft YaHei', fill: '#3498db', backgroundColor: '#2a3a4a', padding: { x: 10, y: 5 }
        }).setInteractive({ useHandCursor: true }).setVisible(false);
        this.equipBtn.on('pointerdown', () => this.equipSelectedItem());
        this.container.add(this.equipBtn);

        this.dropBtn = this.scene.add.text(x + 130, y + 130, '丢弃', {
            font: '14px Microsoft YaHei', fill: '#e74c3c', backgroundColor: '#4a2a2a', padding: { x: 10, y: 5 }
        }).setInteractive({ useHandCursor: true }).setVisible(false);
        this.dropBtn.on('pointerdown', () => this.dropSelectedItem());
        this.container.add(this.dropBtn);
    }

    createStatsPanel(x, y) {
        const stats = this.inventoryManager.getStats();

        const title = this.scene.add.text(x, y - 30, '角色属性', {
            font: 'bold 14px Microsoft YaHei', fill: '#d4a574'
        });
        this.container.add(title);

        const statsText = [
            `生命: ${this.scene.playerData.hp}/${stats.max_hp}`,
            `法力: ${this.scene.playerData.mp}/${stats.max_mp}`,
            `攻击: ${stats.attack}`,
            `防御: ${stats.defense}`,
            `速度: ${stats.speed}`
        ].join('\n');

        const text = this.scene.add.text(x, y, statsText, {
            font: '12px Microsoft YaHei', fill: '#fff', lineSpacing: 8
        });
        this.container.add(text);
    }

    updateDetailPanel(item) {
        this.detailName.setText(`${this.inventoryManager.getItemIcon(item)} ${item.name}`);
        this.detailDesc.setText(item.description || '无描述');

        let effectText = '';
        if (item.effect) {
            let effect = item.effect;
            if (typeof effect === 'string') {
                try { effect = JSON.parse(effect); } catch (e) { effect = {}; }
            }
            const effects = [];
            if (effect.attack) effects.push(`攻击 +${effect.attack}`);
            if (effect.defense) effects.push(`防御 +${effect.defense}`);
            if (effect.hp) effects.push(`生命 +${effect.hp}`);
            if (effect.mp) effects.push(`法力 +${effect.mp}`);
            effectText = effects.join(', ');
        }
        this.detailEffect.setText(effectText);

        this.useBtn.setVisible(item.type === 'consumable');
        this.equipBtn.setVisible(['weapon', 'armor', 'shield'].includes(item.type));
        this.dropBtn.setVisible(true);
    }

    useSelectedItem() {
        if (!this.selectedItem) return;
        const result = this.inventoryManager.useItem(this.selectedItem.id);
        if (result) {
            this.scene.showNotification(`使用了 ${result.item.name}`);
            this.refreshUI();
        }
    }

    equipSelectedItem() {
        if (!this.selectedItem) return;
        const success = this.inventoryManager.equipItem(this.selectedItem.id);
        if (success) {
            this.scene.showNotification(`装备了 ${this.selectedItem.name}`);
            this.refreshUI();
        }
    }

    dropSelectedItem() {
        if (!this.selectedItem) return;
        this.inventoryManager.removeItem(this.selectedItem.id, 1);
        this.scene.showNotification(`丢弃了 ${this.selectedItem.name}`);
        this.selectedItem = null;
        this.refreshUI();
    }

    refreshUI() {
        if (this.container) {
            this.container.destroy();
        }
        this.createUI();
    }
}
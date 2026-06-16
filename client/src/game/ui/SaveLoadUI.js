const API_BASE = 'http://localhost:8080/api';

export class SaveLoadUI {
    constructor(scene, playerData) {
        this.scene = scene;
        this.playerData = playerData;
        this.container = null;
        this.isOpen = false;
        this.mode = 'save'; // 'save' or 'load'
    }

    open(mode = 'save') {
        if (this.isOpen) return;
        this.isOpen = true;
        this.mode = mode;
        this.createUI();
    }

    close() {
        if (!this.isOpen) return;
        this.isOpen = false;
        if (this.container) {
            this.container.destroy();
            this.container = null;
        }
    }

    toggle(mode = 'save') {
        if (this.isOpen) this.close();
        else this.open(mode);
    }

    async createUI() {
        const width = this.scene.cameras.main.width;
        const height = this.scene.cameras.main.height;

        this.container = this.scene.add.container(0, 0).setDepth(700).setScrollFactor(0);

        // Overlay
        const overlay = this.scene.add.rectangle(width / 2, height / 2, width, height, 0x000000, 0.7);
        this.container.add(overlay);

        // Panel
        const panelW = 500;
        const panelH = 480;
        const bg = this.scene.add.rectangle(width / 2, height / 2, panelW, panelH, 0x1a1a1a, 0.95);
        const border = this.scene.add.rectangle(width / 2, height / 2, panelW, panelH)
            .setStrokeStyle(2, 0xd4a574);
        this.container.add(bg);
        this.container.add(border);

        // Title
        const titleText = this.mode === 'save' ? '保存游戏' : '读取存档';
        const title = this.scene.add.text(width / 2, height / 2 - panelH / 2 + 30, titleText, {
            font: 'bold 22px Microsoft YaHei', fill: '#d4a574'
        }).setOrigin(0.5);
        this.container.add(title);

        // Close button
        const closeBtn = this.scene.add.text(width / 2 + panelW / 2 - 30, height / 2 - panelH / 2 + 15, '✕', {
            font: '20px Microsoft YaHei', fill: '#ff6b6b'
        }).setOrigin(0.5).setInteractive({ useHandCursor: true });
        closeBtn.on('pointerdown', () => this.close());
        this.container.add(closeBtn);

        // Load saves
        try {
            const resp = await fetch(`${API_BASE}/saves/${this.playerData.id}`);
            const data = await resp.json();
            const saves = data.saves || [];

            saves.forEach((save, i) => {
                const y = height / 2 - panelH / 2 + 80 + i * 40;
                if (y > height / 2 + panelH / 2 - 60) return;

                const isEmpty = save.is_empty;
                const slotBg = this.scene.add.rectangle(width / 2, y, panelW - 40, 35, isEmpty ? 0x2a2a2a : 0x333344)
                    .setInteractive({ useHandCursor: !isEmpty || this.mode === 'save' });

                const slotLabel = save.slot === 0 ? '自动' : `存档${save.slot}`;
                const infoText = isEmpty
                    ? `${slotLabel} - 空`
                    : `${slotLabel} - Lv.${save.level} ${save.name || ''} [${save.scene_id || ''}]`;

                const slotText = this.scene.add.text(width / 2 - panelW / 2 + 30, y, infoText, {
                    font: '13px Microsoft YaHei', fill: isEmpty ? '#666' : '#fff'
                });

                this.container.add(slotBg);
                this.container.add(slotText);

                if (this.mode === 'save' && save.slot > 0) {
                    const saveBtn = this.scene.add.text(width / 2 + panelW / 2 - 60, y, '保存', {
                        font: '13px Microsoft YaHei', fill: '#4a7c59'
                    }).setOrigin(0.5).setInteractive({ useHandCursor: true });
                    saveBtn.on('pointerdown', () => this.doSave(save.slot));
                    this.container.add(saveBtn);
                }

                if (this.mode === 'load' && !isEmpty) {
                    const loadBtn = this.scene.add.text(width / 2 + panelW / 2 - 60, y, '读取', {
                        font: '13px Microsoft YaHei', fill: '#4a90d9'
                    }).setOrigin(0.5).setInteractive({ useHandCursor: true });
                    loadBtn.on('pointerdown', () => this.doLoad(save.slot));
                    this.container.add(loadBtn);
                }
            });
        } catch (e) {
            const errText = this.scene.add.text(width / 2, height / 2, '加载存档列表失败', {
                font: '14px Microsoft YaHei', fill: '#ff6b6b'
            }).setOrigin(0.5);
            this.container.add(errText);
        }
    }

    async doSave(slot) {
        try {
            const resp = await fetch(`${API_BASE}/save`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    player_id: this.playerData.id,
                    slot: slot
                })
            });
            const result = await resp.json();
            if (result.error) {
                this.scene.showNotification(`保存失败: ${result.error.message || '未知错误'}`);
            } else {
                this.scene.showNotification('保存成功！');
                this.close();
            }
        } catch (e) {
            this.scene.showNotification('保存失败');
        }
    }

    async doLoad(slot) {
        try {
            // Get save list to find the save_id for this slot
            const savesResp = await fetch(`${API_BASE}/saves/${this.playerData.id}`);
            const savesData = await savesResp.json();
            const save = (savesData.saves || []).find(s => s.slot === slot);
            if (!save || save.is_empty) {
                this.scene.showNotification('存档为空');
                return;
            }

            const loadResp = await fetch(`${API_BASE}/load/${save.save_id}`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    player_id: this.playerData.id
                })
            });
            const result = await loadResp.json();
            if (result.error) {
                this.scene.showNotification(`读档失败: ${result.error.message || '未知错误'}`);
            } else {
                // Update local player data
                const loaded = result.data;
                Object.assign(this.playerData, {
                    name: loaded.name,
                    level: loaded.level,
                    exp: loaded.exp,
                    gold: loaded.gold,
                    hp: loaded.hp,
                    mp: loaded.mp,
                    attack: loaded.attack,
                    defense: loaded.defense,
                    scene_id: loaded.scene_id,
                    pos_x: loaded.pos_x,
                    pos_y: loaded.pos_y,
                    items: loaded.items,
                    equipment: loaded.equipment
                });
                this.scene.showNotification('读档成功！');
                this.close();
                // Reload the scene
                this.scene.loadScene(this.playerData.scene_id);
            }
        } catch (e) {
            this.scene.showNotification('读档失败');
        }
    }
}

import Phaser from 'phaser';
import { InventoryManager } from '../systems/InventoryManager.js';
import { CombatManager } from '../systems/CombatManager.js';
import { InventoryUI } from '../ui/InventoryUI.js';
import { CombatUI } from '../ui/CombatUI.js';
import { MiniMap } from '../ui/MiniMap.js';
import { SaveLoadUI } from '../ui/SaveLoadUI.js';

const API_BASE = 'http://localhost:8080/api';

export class GameScene extends Phaser.Scene {
    constructor() {
        super({ key: 'GameScene' });
        this.player = null;
        this.npcs = [];
        this.portals = [];
        this.cursors = null;
        this.gameData = null;
        this.playerData = null;
        this.currentSceneCode = null;
        this.npcSprites = new Map();
        this.portalSprites = new Map();
        this.tutorialStep = 0;
        this.showingDialog = false;
        this.showingShop = false;
        this.showingTutorial = false;
        this.questLog = [];
        this.inventory = {};
        this.visitedScenes = new Set();
        this.inventoryManager = null;
        this.combatManager = null;
        this.inventoryUI = null;
        this.combatUI = null;
        this.miniMap = null;
        this.saveLoadUI = null;
        this.encounterCooldown = 0;
    }

    async create() {
        // Show title screen first
        this.showTitleScreen();
    }

    showTitleScreen() {
        const width = this.cameras.main.width;
        const height = this.cameras.main.height;

        // Background
        const bg = this.add.rectangle(width/2, height/2, width, height, 0x1a1a2e);

        // Title
        this.add.text(width/2, height/3, '青石村传说', {
            font: 'bold 48px Microsoft YaHei',
            fill: '#d4a574',
            stroke: '#000',
            strokeThickness: 4
        }).setOrigin(0.5);

        this.add.text(width/2, height/3 + 60, '一个古风RPG冒险', {
            font: '20px Microsoft YaHei',
            fill: '#aaa'
        }).setOrigin(0.5);

        // Name input
        this.add.text(width/2, height/2 + 20, '输入你的名字：', {
            font: '18px Microsoft YaHei',
            fill: '#fff'
        }).setOrigin(0.5);

        // Create input element
        const input = document.createElement('input');
        input.type = 'text';
        input.value = '冒险者';
        input.style.cssText = `
            position: absolute;
            left: 50%;
            top: 55%;
            transform: translate(-50%, -50%);
            width: 200px;
            padding: 10px;
            font-size: 18px;
            text-align: center;
            background: rgba(255,255,255,0.1);
            border: 1px solid #d4a574;
            color: #fff;
            border-radius: 5px;
            outline: none;
        `;
        document.body.appendChild(input);
        this.nameInput = input;

        // Start button
        const btnBg = this.add.rectangle(width/2, height/2 + 100, 200, 50, 0x4a7c59)
            .setInteractive({ useHandCursor: true });
        const btnText = this.add.text(width/2, height/2 + 100, '开始冒险', {
            font: '22px Microsoft YaHei',
            fill: '#fff'
        }).setOrigin(0.5);

        btnBg.on('pointerover', () => btnBg.setFillStyle(0x5a8c69));
        btnBg.on('pointerout', () => btnBg.setFillStyle(0x4a7c59));
        btnBg.on('pointerdown', () => {
            const name = input.value.trim() || '冒险者';
            input.remove();
            bg.destroy();
            this.children.removeAll(true);
            this.startGame(name);
        });

        // Instructions
        this.add.text(width/2, height - 80, 'WASD/方向键移动 | 点击NPC对话 | I背包 | F5保存 | F9读档 | 空格跳过教程', {
            font: '14px Microsoft YaHei',
            fill: '#666'
        }).setOrigin(0.5);
    }

    async startGame(playerName) {
        try {
            // Load game data
            const initResp = await fetch(`${API_BASE}/game/init`);
            const initData = await initResp.json();
            this.gameData = initData;

            // Create or load player
            const playerResp = await fetch(`${API_BASE}/player/create`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    name: playerName,
                    account: 'player_' + Date.now()
                })
            });
            const playerResult = await playerResp.json();
            this.playerData = playerResult.data;

            // Parse items
            try {
                this.inventory = JSON.parse(this.playerData.items || '{}');
            } catch (e) {
                this.inventory = {};
            }

            // Initialize managers
            this.inventoryManager = new InventoryManager(this.gameData, this.playerData);
            this.combatManager = new CombatManager(this.playerData, this.inventoryManager);

            // Initialize quest log
            this.questLog = (initData.tasks || []).filter(t => t.status === 'active');

            // Load starting scene
            await this.loadScene(this.playerData.scene_id || initData.config.start_scene);

            // Setup controls
            this.cursors = this.input.keyboard.createCursorKeys();
            this.wasd = {
                up: this.input.keyboard.addKey(Phaser.Input.Keyboard.KeyCodes.W),
                down: this.input.keyboard.addKey(Phaser.Input.Keyboard.KeyCodes.S),
                left: this.input.keyboard.addKey(Phaser.Input.Keyboard.KeyCodes.A),
                right: this.input.keyboard.addKey(Phaser.Input.Keyboard.KeyCodes.D)
            };

            // Skip tutorial key
            this.input.keyboard.addKey(Phaser.Input.Keyboard.KeyCodes.SPACE).on('down', () => {
                if (this.showingTutorial) {
                    this.skipTutorial();
                }
            });

            // Inventory key
            this.input.keyboard.addKey(Phaser.Input.Keyboard.KeyCodes.I).on('down', () => {
                if (!this.showingDialog && !this.showingShop && !this.showingTutorial) {
                    this.toggleInventory();
                }
            });

            // Save key (F5)
            this.input.keyboard.addKey(Phaser.Input.Keyboard.KeyCodes.F5).on('down', () => {
                if (!this.showingDialog && !this.showingShop && !this.showingTutorial) {
                    if (!this.saveLoadUI) {
                        this.saveLoadUI = new SaveLoadUI(this, this.playerData);
                    }
                    this.saveLoadUI.toggle('save');
                }
            });

            // Load key (F9)
            this.input.keyboard.addKey(Phaser.Input.Keyboard.KeyCodes.F9).on('down', () => {
                if (!this.showingDialog && !this.showingShop && !this.showingTutorial) {
                    if (!this.saveLoadUI) {
                        this.saveLoadUI = new SaveLoadUI(this, this.playerData);
                    }
                    this.saveLoadUI.toggle('load');
                }
            });

            // Show tutorial
            this.showTutorial();

        } catch (error) {
            console.error('Failed to start game:', error);
            this.add.text(400, 300, '加载失败，请刷新重试', { font: '20px Arial', fill: '#ff0000' }).setOrigin(0.5);
        }
    }

    async loadScene(sceneCode) {
        // Clear existing sprites
        this.npcSprites.forEach(s => s.destroy());
        this.portalSprites.forEach(s => s.destroy());
        this.npcSprites.clear();
        this.portalSprites.clear();

        // Find scene data
        const scene = this.gameData.scenes.find(s => s.code === sceneCode);
        if (!scene) return;

        this.currentSceneCode = sceneCode;
        this.visitedScenes.add(sceneCode);

        // Create background
        this.createSceneBackground(scene);

        // Create portals
        if (scene.portals) {
            scene.portals.forEach(portal => {
                this.createPortal(portal);
            });
        }

        // Create NPCs
        if (scene.scene_npcs) {
            for (const sn of scene.scene_npcs) {
                const npcData = this.gameData.npcs.find(n => n.id === sn.npc_id);
                if (npcData) {
                    this.createNPC(sn, npcData);
                }
            }
        }

        // Create player
        if (!this.player) {
            const startX = this.playerData.pos_x || 200;
            const startY = this.playerData.pos_y || 450;
            this.player = this.physics.add.sprite(startX, startY, 'player');
            this.player.setCollideWorldBounds(true);
            this.player.setScale(1.5);

            // Camera follow
            this.cameras.main.startFollow(this.player, true, 0.1, 0.1);
            this.cameras.main.setZoom(1.5);
        } else {
            this.player.setPosition(this.playerData.pos_x || 200, this.playerData.pos_y || 450);
        }

        // Set world bounds
        this.physics.world.setBounds(0, 0, scene.width, scene.height);

        // Show scene name
        this.showSceneName(scene.name, scene.description);

        // Update quest for visiting scene
        this.updateVisitQuest(sceneCode);

        // Create UI
        this.createGameUI();
    }

    createSceneBackground(scene) {
        const width = scene.width;
        const height = scene.height;

        // Parse background color
        let bgColor = 0x4a7c59;
        if (scene.background && scene.background.startsWith('#')) {
            bgColor = parseInt(scene.background.replace('#', '0x'));
        }

        // Background
        this.add.rectangle(width/2, height/2, width, height, bgColor);

        // Add some decoration based on scene
        if (scene.code.includes('village_entrance')) {
            // Village entrance - add trees and path
            for (let i = 0; i < 8; i++) {
                this.add.image(Phaser.Math.Between(50, width-50), Phaser.Math.Between(50, 200), 'tree').setScale(1.2);
            }
            // Stone monument
            this.add.rectangle(800, 300, 60, 80, 0x808080);
            this.add.text(800, 300, '青石村', { font: '14px Microsoft YaHei', fill: '#fff' }).setOrigin(0.5);
        } else if (scene.code.includes('village_center')) {
            // Central square - big tree and well
            this.add.image(800, 300, 'tree').setScale(2);
            this.add.circle(800, 450, 30, 0x4a90d9);
            this.add.text(800, 450, '井', { font: '12px Microsoft YaHei', fill: '#fff' }).setOrigin(0.5);
        } else if (scene.code.includes('general_store')) {
            // Shop interior
            this.add.rectangle(600, 200, 200, 150, 0x8b7355);
            this.add.text(600, 200, '货架', { font: '14px Microsoft YaHei', fill: '#fff' }).setOrigin(0.5);
            this.add.rectangle(600, 500, 120, 60, 0x654321);
            this.add.text(600, 500, '柜台', { font: '14px Microsoft YaHei', fill: '#fff' }).setOrigin(0.5);
        } else if (scene.code.includes('blacksmith')) {
            // Forge
            this.add.rectangle(600, 250, 80, 80, 0xff4500);
            this.add.text(600, 250, '熔炉', { font: '14px Microsoft YaHei', fill: '#fff' }).setOrigin(0.5);
            this.add.rectangle(400, 400, 150, 40, 0x696969);
            this.add.text(400, 400, '铁砧', { font: '14px Microsoft YaHei', fill: '#fff' }).setOrigin(0.5);
        } else if (scene.code.includes('tea_stand')) {
            // Tea stall
            this.add.rectangle(500, 300, 100, 60, 0x8b7355);
            this.add.text(500, 300, '茶桌', { font: '14px Microsoft YaHei', fill: '#fff' }).setOrigin(0.5);
            for (let i = 0; i < 3; i++) {
                this.add.rectangle(400 + i * 80, 380, 30, 30, 0x654321);
            }
        } else if (scene.code.includes('village_path')) {
            // Path - add rocks and trees
            for (let i = 0; i < 12; i++) {
                this.add.image(Phaser.Math.Between(100, width-100), Phaser.Math.Between(100, 300), 'tree');
            }
            // Rocks
            for (let i = 0; i < 5; i++) {
                this.add.circle(Phaser.Math.Between(200, width-200), Phaser.Math.Between(400, 700), 15, 0x808080);
            }
            // Fog of war effect
            this.createFogOfWar(width, height);
        }

        // Add some grass patches
        for (let i = 0; i < 20; i++) {
            const x = Phaser.Math.Between(50, width - 50);
            const y = Phaser.Math.Between(50, height - 50);
            this.add.circle(x, y, 5, 0x3d8b37, 0.5);
        }
    }

    createFogOfWar(width, height) {
        const fogGraphics = this.add.graphics();
        fogGraphics.fillStyle(0x000000, 0.3);

        for (let i = 0; i < 30; i++) {
            const x = Phaser.Math.Between(0, width);
            const y = Phaser.Math.Between(0, height);
            const radius = Phaser.Math.Between(50, 150);
            fogGraphics.fillCircle(x, y, radius);
        }

        this.fogOfWar = fogGraphics;
    }

    createNPC(sceneNPC, npcData) {
        const spriteKey = this.getNPCSpriteKey(npcData.code);
        const npc = this.physics.add.sprite(sceneNPC.x, sceneNPC.y, spriteKey);
        npc.setScale(1.5);
        npc.setData('npcData', npcData);
        npc.setData('sceneNPC', sceneNPC);
        npc.setImmovable(true);

        // Name label
        const nameTag = this.add.text(sceneNPC.x, sceneNPC.y - 40, npcData.name, {
            font: '12px Microsoft YaHei',
            fill: '#fff',
            backgroundColor: '#000000aa',
            padding: { x: 4, y: 2 }
        }).setOrigin(0.5);
        npc.setData('nameTag', nameTag);

        // Quest indicator
        const hasQuest = this.hasQuestForNPC(npcData.code);
        if (hasQuest) {
            const indicator = this.add.image(sceneNPC.x, sceneNPC.y - 55, 'indicator');
            npc.setData('indicator', indicator);
        }

        // Make interactive
        npc.setInteractive({ useHandCursor: true });
        npc.on('pointerdown', () => {
            if (!this.showingDialog && !this.showingShop) {
                this.interactWithNPC(npcData);
            }
        });

        // Simple idle animation
        this.tweens.add({
            targets: npc,
            y: sceneNPC.y - 3,
            duration: 1500,
            yoyo: true,
            repeat: -1,
            ease: 'Sine.easeInOut'
        });

        this.npcSprites.set(npcData.id, npc);
    }

    getNPCSpriteKey(npcCode) {
        const map = {
            'npc_chief_chen': 'npc_chief',
            'npc_merchant_li': 'npc_merchant',
            'npc_tea_wang': 'npc_tea',
            'npc_blacksmith_zhang': 'npc_blacksmith',
            'npc_hunter_zhou': 'npc_hunter',
            'npc_kid_stone': 'npc_kid'
        };
        return map[npcCode] || 'npc_merchant';
    }

    createPortal(portal) {
        const pSprite = this.physics.add.sprite(portal.x, portal.y, 'portal');
        pSprite.setData('portalData', portal);
        pSprite.setImmovable(true);

        // Pulse animation
        this.tweens.add({
            targets: pSprite,
            alpha: 0.5,
            duration: 1000,
            yoyo: true,
            repeat: -1
        });

        // Label
        const targetScene = this.gameData.scenes.find(s => s.code === portal.target_scene);
        if (targetScene) {
            this.add.text(portal.x, portal.y + 30, targetScene.name, {
                font: '11px Microsoft YaHei',
                fill: '#fff',
                backgroundColor: '#000000aa',
                padding: { x: 3, y: 2 }
            }).setOrigin(0.5);
        }

        // Overlap detection
        this.physics.add.overlap(this.player, pSprite, () => {
            if (!this.showingDialog) {
                this.enterPortal(portal);
            }
        }, null, this);

        this.portalSprites.set(portal.id, pSprite);
    }

    async enterPortal(portal) {
        // Save position
        this.playerData.pos_x = portal.target_x;
        this.playerData.pos_y = portal.target_y;
        this.playerData.scene_id = portal.target_scene;

        // Update server
        fetch(`${API_BASE}/player/${this.playerData.id}/pos`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                scene_id: portal.target_scene,
                pos_x: portal.target_x,
                pos_y: portal.target_y
            })
        });

        // Smooth camera transition
        this.cameras.main.fadeOut(300, 0, 0, 0);
        await new Promise(resolve => this.cameras.main.once('camerafadeoutcomplete', resolve));

        await this.loadScene(portal.target_scene);

        this.cameras.main.fadeIn(300, 0, 0, 0);
    }

    showSceneName(name, description) {
        const width = this.cameras.main.width;
        const bg = this.add.rectangle(width/2, 50, 300, 60, 0x000000, 0.7)
            .setScrollFactor(0).setDepth(100);
        const title = this.add.text(width/2, 40, name, {
            font: 'bold 20px Microsoft YaHei',
            fill: '#d4a574'
        }).setOrigin(0.5).setScrollFactor(0).setDepth(101);
        const desc = this.add.text(width/2, 60, description.substring(0, 30) + '...', {
            font: '12px Microsoft YaHei',
            fill: '#aaa'
        }).setOrigin(0.5).setScrollFactor(0).setDepth(101);

        this.tweens.add({
            targets: [bg, title, desc],
            alpha: 0,
            delay: 2000,
            duration: 1000,
            onComplete: () => {
                bg.destroy();
                title.destroy();
                desc.destroy();
            }
        });
    }

    createGameUI() {
        // Remove old UI elements
        if (this.uiElements) {
            this.uiElements.forEach(el => el.destroy());
        }
        this.uiElements = [];

        const width = this.cameras.main.width;

        // HP Bar
        const hpBg = this.add.rectangle(width - 120, 30, 100, 16, 0x333333)
            .setScrollFactor(0).setDepth(200);
        const hpBar = this.add.rectangle(width - 120, 30, 96, 12, 0xe74c3c)
            .setScrollFactor(0).setDepth(201);
        const hpText = this.add.text(width - 120, 30, `HP: ${this.playerData.hp}`, {
            font: '10px Microsoft YaHei', fill: '#fff'
        }).setOrigin(0.5).setScrollFactor(0).setDepth(202);

        // Gold
        const goldText = this.add.text(width - 120, 50, `💰 ${this.playerData.gold}`, {
            font: '14px Microsoft YaHei', fill: '#f1c40f'
        }).setOrigin(0.5).setScrollFactor(0).setDepth(200);

        // Level
        const levelText = this.add.text(width - 120, 70, `Lv.${this.playerData.level} ${this.playerData.name}`, {
            font: '12px Microsoft YaHei', fill: '#fff'
        }).setOrigin(0.5).setScrollFactor(0).setDepth(200);

        this.uiElements.push(hpBg, hpBar, hpText, goldText, levelText);

        // Quest tracker
        this.createQuestTracker();

        // Mini-map
        if (this.miniMap) {
            this.miniMap.destroy();
        }
        this.miniMap = new MiniMap(this);
        this.miniMap.create();
    }

    createQuestTracker() {
        const activeQuests = this.questLog.filter(q => q.status === 'active').slice(0, 3);
        if (activeQuests.length === 0) return;

        const startX = 20;
        const startY = 100;

        const header = this.add.text(startX, startY - 20, '📜 当前任务', {
            font: '14px Microsoft YaHei', fill: '#d4a574'
        }).setScrollFactor(0).setDepth(200);
        this.uiElements.push(header);

        activeQuests.forEach((quest, i) => {
            const y = startY + i * 40;
            const qText = this.add.text(startX, y, `• ${quest.name}`, {
                font: '12px Microsoft YaHei', fill: '#fff'
            }).setScrollFactor(0).setDepth(200);

            // Parse objectives
            let objectives = [];
            try { objectives = JSON.parse(quest.objectives || '[]'); } catch (e) {}
            const completed = objectives.filter(o => o.completed).length;
            const total = objectives.length;
            const progress = this.add.text(startX + 10, y + 16, `进度: ${completed}/${total}`, {
                font: '10px Microsoft YaHei', fill: '#aaa'
            }).setScrollFactor(0).setDepth(200);

            this.uiElements.push(qText, progress);
        });
    }

    hasQuestForNPC(npcCode) {
        return this.questLog.some(q => {
            if (q.status !== 'active') return false;
            try {
                const objectives = JSON.parse(q.objectives || '[]');
                return objectives.some(o => o.target === npcCode && !o.completed);
            } catch (e) {
                return false;
            }
        });
    }

    updateVisitQuest(sceneCode) {
        this.questLog.forEach(quest => {
            if (quest.status !== 'active') return;
            try {
                let objectives = JSON.parse(quest.objectives || '[]');
                let changed = false;
                objectives.forEach(obj => {
                    if (obj.type === 'visit' && obj.target === sceneCode && !obj.completed) {
                        obj.completed = true;
                        changed = true;
                        this.showNotification(`✓ ${obj.description}`);
                    }
                });
                if (changed) {
                    quest.objectives = JSON.stringify(objectives);
                    this.checkQuestComplete(quest);
                }
            } catch (e) {}
        });
    }

    checkQuestComplete(quest) {
        try {
            const objectives = JSON.parse(quest.objectives || '[]');
            if (objectives.every(o => o.completed)) {
                quest.status = 'completed';
                this.showNotification(`🎉 任务完成: ${quest.name}`);

                // Parse rewards
                try {
                    const rewards = JSON.parse(quest.rewards || '{}');
                    if (rewards.gold) {
                        this.playerData.gold += rewards.gold;
                        this.showNotification(`获得 ${rewards.gold} 金币`);
                    }
                    if (rewards.exp) {
                        this.showNotification(`获得 ${rewards.exp} 经验`);
                    }
                } catch (e) {}

                // Activate next task
                if (quest.next_task) {
                    const nextQuest = this.questLog.find(q => q.code === quest.next_task);
                    if (nextQuest) {
                        nextQuest.status = 'active';
                        this.showNotification(`📜 新任务: ${nextQuest.name}`);
                    }
                }

                this.createGameUI();
            }
        } catch (e) {}
    }

    showNotification(text) {
        const width = this.cameras.main.width;
        const notif = this.add.text(width/2, 100, text, {
            font: '16px Microsoft YaHei',
            fill: '#f1c40f',
            backgroundColor: '#000000cc',
            padding: { x: 16, y: 8 }
        }).setOrigin(0.5).setScrollFactor(0).setDepth(300);

        this.tweens.add({
            targets: notif,
            y: 80,
            alpha: 0,
            delay: 1500,
            duration: 500,
            onComplete: () => notif.destroy()
        });
    }

    async interactWithNPC(npcData) {
        this.showingDialog = true;

        // Check for dialogue objectives
        this.questLog.forEach(quest => {
            if (quest.status !== 'active') return;
            try {
                let objectives = JSON.parse(quest.objectives || '[]');
                objectives.forEach(obj => {
                    if (obj.type === 'dialogue' && obj.target === npcData.code && !obj.completed) {
                        obj.completed = true;
                        quest.objectives = JSON.stringify(objectives);
                        this.showNotification(`✓ ${obj.description}`);
                        this.checkQuestComplete(quest);
                    }
                });
            } catch (e) {}
        });

        // Check if NPC has shop
        if (npcData.shop_id) {
            this.showShopDialog(npcData);
        } else {
            this.showChatDialog(npcData);
        }
    }

    showChatDialog(npcData) {
        const width = this.cameras.main.width;
        const height = this.cameras.main.height;

        // Dialog background
        const dialogBg = this.add.rectangle(width/2, height - 120, width - 100, 200, 0x000000, 0.9)
            .setScrollFactor(0).setDepth(400);

        // NPC info
        const nameText = this.add.text(80, height - 210, `${npcData.avatar || ''} ${npcData.name}`, {
            font: 'bold 18px Microsoft YaHei', fill: '#d4a574'
        }).setScrollFactor(0).setDepth(401);

        const titleText = this.add.text(80, height - 185, npcData.title || '', {
            font: '12px Microsoft YaHei', fill: '#aaa'
        }).setScrollFactor(0).setDepth(401);

        // Greeting message
        const greeting = this.getNPCGreeting(npcData);
        const msgText = this.add.text(80, height - 160, greeting, {
            font: '14px Microsoft YaHei',
            fill: '#fff',
            wordWrap: { width: width - 200 }
        }).setScrollFactor(0).setDepth(401);

        // Chat input area
        const inputBg = this.add.rectangle(width/2, height - 50, width - 140, 36, 0x333333)
            .setScrollFactor(0).setDepth(401);
        const inputHint = this.add.text(width/2, height - 50, '输入消息与NPC对话...', {
            font: '13px Microsoft YaHei', fill: '#888'
        }).setOrigin(0.5).setScrollFactor(0).setDepth(402);

        // Create HTML input
        const chatInput = document.createElement('input');
        chatInput.type = 'text';
        chatInput.placeholder = '输入消息与NPC对话...';
        chatInput.style.cssText = `
            position: fixed;
            left: 70px;
            bottom: 42px;
            width: calc(100% - 180px);
            padding: 8px 12px;
            font-size: 14px;
            background: #333;
            border: 1px solid #555;
            color: #fff;
            border-radius: 4px;
            outline: none;
            z-index: 1000;
        `;
        document.body.appendChild(chatInput);
        chatInput.focus();

        // Close button
        const closeBtn = this.add.text(width - 70, height - 210, '✕ 关闭', {
            font: '14px Microsoft YaHei', fill: '#ff6b6b'
        }).setScrollFactor(0).setDepth(401).setInteractive({ useHandCursor: true });

        const dialogElements = [dialogBg, nameText, titleText, msgText, inputBg, inputHint, closeBtn];

        closeBtn.on('pointerdown', () => {
            chatInput.remove();
            dialogElements.forEach(el => el.destroy());
            this.showingDialog = false;
        });

        // Chat history
        let chatHistory = [];

        // Handle enter key
        chatInput.addEventListener('keydown', async (e) => {
            if (e.key === 'Enter' && chatInput.value.trim()) {
                const userMsg = chatInput.value.trim();
                chatInput.value = '';

                // Add user message to history
                chatHistory.push(`你: ${userMsg}`);

                // Update display
                msgText.setText(chatHistory.slice(-5).join('\n'));

                // Send to server
                try {
                    const resp = await fetch(`${API_BASE}/npc/chat`, {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({
                            player_id: this.playerData.id,
                            npc_id: npcData.id,
                            message: userMsg
                        })
                    });
                    const result = await resp.json();

                    chatHistory.push(`${npcData.name}: ${result.reply}`);
                    msgText.setText(chatHistory.slice(-5).join('\n'));
                } catch (err) {
                    chatHistory.push(`${npcData.name}: （网络错误）`);
                    msgText.setText(chatHistory.slice(-5).join('\n'));
                }
            }
            if (e.key === 'Escape') {
                chatInput.remove();
                dialogElements.forEach(el => el.destroy());
                this.showingDialog = false;
            }
        });
    }

    getNPCGreeting(npcData) {
        const greetings = {
            'npc_chief_chen': '呵呵，年轻人，欢迎来到青石村！老朽是这里的村长。你初来乍到，有什么不懂的尽管问老朽。',
            'npc_merchant_li': '客官好！欢迎光临李记杂货铺！在下经营各种日用品和药材，物美价廉！',
            'npc_tea_wang': '哎呀，客官来了！快坐下喝杯茶歇歇脚。我跟你说，咱们青石村可是个好地方！',
            'npc_blacksmith_zhang': '嗯，客官好。俺是这儿的铁匠，姓张。俺打的兵器在方圆百里都小有名气。',
            'npc_hunter_zhou': '嗯，你好。俺是猎户老周。村外最近不太平，你出去的时候小心点。',
            'npc_kid_stone': '大哥哥好！你是新来的冒险者吗？好酷啊！你能给我看看你的武器吗？'
        };
        return greetings[npcData.code] || '你好，客官！';
    }

    showShopDialog(npcData) {
        const width = this.cameras.main.width;
        const height = this.cameras.main.height;

        // Shop background
        const shopBg = this.add.rectangle(width/2, height/2, 500, 400, 0x1a1a1a, 0.95)
            .setScrollFactor(0).setDepth(400);
        const shopBorder = this.add.rectangle(width/2, height/2, 500, 400)
            .setStrokeStyle(2, 0xd4a574).setScrollFactor(0).setDepth(401);

        // Shop title
        const title = this.add.text(width/2, height/2 - 180, `${npcData.name}的商店`, {
            font: 'bold 20px Microsoft YaHei', fill: '#d4a574'
        }).setOrigin(0.5).setScrollFactor(0).setDepth(402);

        const goldDisplay = this.add.text(width/2 + 180, height/2 - 180, `💰 ${this.playerData.gold}`, {
            font: '16px Microsoft YaHei', fill: '#f1c40f'
        }).setOrigin(0.5).setScrollFactor(0).setDepth(402);

        // Close button
        const closeBtn = this.add.text(width/2 + 230, height/2 - 185, '✕', {
            font: '20px Microsoft YaHei', fill: '#ff6b6b'
        }).setScrollFactor(0).setDepth(402).setInteractive({ useHandCursor: true });

        const shopElements = [shopBg, shopBorder, title, goldDisplay, closeBtn];

        closeBtn.on('pointerdown', () => {
            shopElements.forEach(el => el.destroy());
            itemElements.forEach(el => el.destroy());
            this.showingDialog = false;
            this.showingShop = false;
        });

        // Load shop items
        let itemElements = [];
        fetch(`${API_BASE}/game/shop/${this.getShopCode(npcData.code)}/items`)
            .then(r => r.json())
            .then(data => {
                const items = data.items || [];
                items.forEach((item, i) => {
                    const y = height/2 - 130 + i * 45;
                    if (y > height/2 + 150) return;

                    const itemBg = this.add.rectangle(width/2, y, 460, 40, 0x2a2a2a)
                        .setScrollFactor(0).setDepth(402);
                    const itemName = this.add.text(width/2 - 200, y, item.item_name || '未知', {
                        font: '14px Microsoft YaHei', fill: '#fff'
                    }).setScrollFactor(0).setDepth(403);
                    const itemDesc = this.add.text(width/2 - 100, y, item.item_description || '', {
                        font: '11px Microsoft YaHei', fill: '#aaa'
                    }).setScrollFactor(0).setDepth(403);
                    const priceText = this.add.text(width/2 + 150, y, `${item.price}文`, {
                        font: '14px Microsoft YaHei', fill: '#f1c40f'
                    }).setScrollFactor(0).setDepth(403);

                    const buyBtn = this.add.text(width/2 + 210, y, '购买', {
                        font: '13px Microsoft YaHei', fill: '#4a7c59'
                    }).setScrollFactor(0).setDepth(403).setInteractive({ useHandCursor: true });

                    buyBtn.on('pointerdown', async () => {
                        try {
                            const resp = await fetch(`${API_BASE}/shop/buy`, {
                                method: 'POST',
                                headers: { 'Content-Type': 'application/json' },
                                body: JSON.stringify({
                                    player_id: this.playerData.id,
                                    shop_code: this.getShopCode(npcData.code),
                                    item_id: item.item_id || item.id,
                                    count: 1
                                })
                            });
                            const result = await resp.json();
                            if (result.error) {
                                this.showNotification(`❌ ${result.error}`);
                            } else {
                                // 更新玩家数据
                                this.playerData.gold = result.gold;
                                this.playerData.items = result.items;
                                if (result.equipment) {
                                    this.playerData.equipment = JSON.stringify(result.equipment);
                                }
                                // 刷新背包管理器
                                if (this.inventoryManager) {
                                    this.inventoryManager.parseInventory();
                                }
                                goldDisplay.setText(`💰 ${this.playerData.gold}`);
                                this.showNotification(`✓ 购买了 ${item.item_name}`);
                            }
                        } catch (err) {
                            this.showNotification('❌ 购买失败');
                        }
                    });

                    shopElements.push(itemBg, itemName, itemDesc, priceText, buyBtn);
                });
            });
    }

    getShopCode(npcCode) {
        const map = {
            'npc_merchant_li': 'shop_general_store',
            'npc_blacksmith_zhang': 'shop_blacksmith'
        };
        return map[npcCode] || '';
    }

    showTutorial() {
        this.showingTutorial = true;
        const width = this.cameras.main.width;
        const height = this.cameras.main.height;

        const tutorials = [
            { title: '欢迎来到青石村！', text: '这是一个古风RPG冒险游戏。你将探索村庄，与NPC对话，完成任务。' },
            { title: '移动', text: '使用 WASD 键或方向键移动角色。' },
            { title: '与NPC对话', text: '点击NPC可以与其对话。NPC头上有黄色标记表示有任务。' },
            { title: '传送', text: '走到紫色传送点可以前往其他场景。' },
            { title: '任务', text: '左上角显示当前任务。完成任务可以获得奖励！' },
            { title: '准备好了吗？', text: '去找老村长聊聊吧！按空格键跳过教程。' }
        ];

        this.tutorialStep = 0;

        const overlay = this.add.rectangle(width/2, height/2, width, height, 0x000000, 0.7)
            .setScrollFactor(0).setDepth(500);
        const box = this.add.rectangle(width/2, height/2, 450, 200, 0x1a1a1a, 0.95)
            .setScrollFactor(0).setDepth(501);
        const border = this.add.rectangle(width/2, height/2, 450, 200)
            .setStrokeStyle(2, 0xd4a574).setScrollFactor(0).setDepth(502);

        const titleText = this.add.text(width/2, height/2 - 60, tutorials[0].title, {
            font: 'bold 22px Microsoft YaHei', fill: '#d4a574'
        }).setOrigin(0.5).setScrollFactor(0).setDepth(503);

        const contentText = this.add.text(width/2, height/2, tutorials[0].text, {
            font: '16px Microsoft YaHei', fill: '#fff',
            wordWrap: { width: 400 },
            align: 'center'
        }).setOrigin(0.5).setScrollFactor(0).setDepth(503);

        const stepText = this.add.text(width/2, height/2 + 60, `${this.tutorialStep + 1}/${tutorials.length}`, {
            font: '12px Microsoft YaHei', fill: '#888'
        }).setOrigin(0.5).setScrollFactor(0).setDepth(503);

        const nextBtn = this.add.text(width/2 + 150, height/2 + 80, '下一步 →', {
            font: '16px Microsoft YaHei', fill: '#4a7c59'
        }).setOrigin(0.5).setScrollFactor(0).setDepth(503).setInteractive({ useHandCursor: true });

        const skipBtn = this.add.text(width/2, height/2 + 80, '跳过教程', {
            font: '14px Microsoft YaHei', fill: '#888'
        }).setOrigin(0.5).setScrollFactor(0).setDepth(503).setInteractive({ useHandCursor: true });

        const tutorialElements = [overlay, box, border, titleText, contentText, stepText, nextBtn, skipBtn];

        const advanceTutorial = () => {
            this.tutorialStep++;
            if (this.tutorialStep >= tutorials.length) {
                tutorialElements.forEach(el => el.destroy());
                this.showingTutorial = false;
                this.showNotification('🎮 冒险开始！去找老村长对话吧！');
                return;
            }
            const t = tutorials[this.tutorialStep];
            titleText.setText(t.title);
            contentText.setText(t.text);
            stepText.setText(`${this.tutorialStep + 1}/${tutorials.length}`);
        };

        nextBtn.on('pointerdown', advanceTutorial);
        skipBtn.on('pointerdown', () => {
            tutorialElements.forEach(el => el.destroy());
            this.showingTutorial = false;
            this.showNotification('🎮 冒险开始！去找老村长对话吧！');
        });

        this.skipTutorial = () => {
            tutorialElements.forEach(el => el.destroy());
            this.showingTutorial = false;
        };
    }

    skipTutorial() {
        // Will be overridden by showTutorial
    }

    toggleInventory() {
        if (!this.inventoryUI) {
            this.inventoryUI = new InventoryUI(this, this.inventoryManager);
        }
        this.inventoryUI.toggle();
    }

    triggerCombat(enemyType) {
        if (this.combatUI && this.combatUI.isOpen) return;
        if (!this.combatUI) {
            this.combatUI = new CombatUI(this, this.combatManager, this.inventoryManager);
            this.combatUI.onCombatEnd = (result) => {
                this.createGameUI();
                if (result.victory && this.inventoryManager) {
                    // 战斗胜利后同步背包
                    this.inventoryManager.syncWithServer(this.playerData.id);
                }
                if (!result.victory && !result.fled) {
                    this.scene.restart();
                }
            };
        }
        this.combatUI.open(enemyType);
    }

    checkRandomEncounter() {
        if (this.currentSceneCode !== 'village_path') return;
        if (this.encounterCooldown > 0) return;
        if (this.showingDialog || this.showingShop || this.showingTutorial) return;
        if (this.combatUI && this.combatUI.isOpen) return;

        const encounterChance = 0.005;
        if (Math.random() < encounterChance) {
            const roll = Math.random();
            let enemyType = 'wolf';
            if (roll < 0.05) enemyType = 'ghost';
            else if (roll < 0.15) enemyType = 'tiger';
            else if (roll < 0.30) enemyType = 'bear';
            else if (roll < 0.50) enemyType = 'bandit';
            else if (roll < 0.60) enemyType = 'alpha_wolf';
            this.encounterCooldown = 180;
            this.triggerCombat(enemyType);
        }
    }

    update() {
        if (!this.player || this.showingDialog || this.showingTutorial) {
            if (this.player) this.player.setVelocity(0);
            return;
        }

        if (this.encounterCooldown > 0) {
            this.encounterCooldown--;
        }

        const speed = 200;
        let vx = 0, vy = 0;

        if (this.cursors.left.isDown || this.wasd.left.isDown) vx = -speed;
        else if (this.cursors.right.isDown || this.wasd.right.isDown) vx = speed;

        if (this.cursors.up.isDown || this.wasd.up.isDown) vy = -speed;
        else if (this.cursors.down.isDown || this.wasd.down.isDown) vy = speed;

        // Normalize diagonal movement
        if (vx !== 0 && vy !== 0) {
            vx *= 0.707;
            vy *= 0.707;
        }

        this.player.setVelocity(vx, vy);

        // Check NPC proximity for interaction hint
        this.npcSprites.forEach((npcSprite, npcId) => {
            const dist = Phaser.Math.Distance.Between(
                this.player.x, this.player.y,
                npcSprite.x, npcSprite.y
            );
            const nameTag = npcSprite.getData('nameTag');
            if (nameTag) {
                nameTag.setAlpha(dist < 100 ? 1 : 0.5);
            }
        });

        // Update mini-map
        if (this.miniMap) {
            this.miniMap.update();
        }

        // Check random encounters
        this.checkRandomEncounter();
    }
}

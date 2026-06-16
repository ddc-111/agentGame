export class CombatUI {
    constructor(scene, combatManager, inventoryManager) {
        this.scene = scene;
        this.combatManager = combatManager;
        this.inventoryManager = inventoryManager;
        this.container = null;
        this.isOpen = false;
        this.onCombatEnd = null;
        this.skills = [];
        this.onSkillUse = null;
    }

    setSkills(skills) {
        this.skills = skills || [];
    }

    open(enemyType) {
        if (this.isOpen) return;
        this.isOpen = true;
        this.combatManager.startCombat(enemyType);
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

    createUI() {
        const width = this.scene.cameras.main.width;
        const height = this.scene.cameras.main.height;

        this.container = this.scene.add.container(0, 0).setDepth(700).setScrollFactor(0);

        const overlay = this.scene.add.rectangle(width / 2, height / 2, width, height, 0x000000, 0.8);
        this.container.add(overlay);

        this.createEnemyDisplay(width / 2, 120);
        this.createPlayerStats(width / 2, height - 120);
        this.createActionButtons(width / 2, height - 200);
        this.createCombatLog(width / 2, height / 2);
    }

    createEnemyDisplay(x, y) {
        const state = this.combatManager.getCombatState();
        const enemy = state.enemy;

        const enemyBg = this.scene.add.rectangle(x, y, 300, 100, 0x2a1a1a).setStrokeStyle(2, 0xe74c3c);
        this.container.add(enemyBg);

        const enemyEmoji = this.scene.add.text(x - 100, y - 10, enemy.emoji || '🐺', {
            font: '40px Microsoft YaHei'
        }).setOrigin(0.5);
        this.container.add(enemyEmoji);

        const enemyName = this.scene.add.text(x, y - 30, enemy.name, {
            font: 'bold 18px Microsoft YaHei', fill: '#e74c3c'
        }).setOrigin(0.5);
        this.container.add(enemyName);

        const hpBarBg = this.scene.add.rectangle(x + 20, y + 10, 200, 20, 0x333333);
        this.container.add(hpBarBg);

        const hpPercent = enemy.hp / enemy.max_hp;
        const hpBar = this.scene.add.rectangle(x + 20 - (200 * (1 - hpPercent) / 2), y + 10, 200 * hpPercent, 16, 0xe74c3c);
        this.container.add(hpBar);

        this.enemyHpText = this.scene.add.text(x + 20, y + 10, `${enemy.hp}/${enemy.max_hp}`, {
            font: '12px Microsoft YaHei', fill: '#fff'
        }).setOrigin(0.5);
        this.container.add(this.enemyHpText);

        this.enemyHpBar = hpBar;
    }

    createPlayerStats(x, y) {
        const state = this.combatManager.getCombatState();
        const player = state.player;

        const playerBg = this.scene.add.rectangle(x, y, 400, 80, 0x1a1a2a).setStrokeStyle(2, 0x3498db);
        this.container.add(playerBg);

        const hpLabel = this.scene.add.text(x - 180, y - 25, 'HP', {
            font: '12px Microsoft YaHei', fill: '#e74c3c'
        });
        this.container.add(hpLabel);

        const hpBarBg = this.scene.add.rectangle(x, y - 20, 250, 16, 0x333333);
        this.container.add(hpBarBg);

        const hpPercent = player.hp / player.max_hp;
        const hpBar = this.scene.add.rectangle(x - (250 * (1 - hpPercent) / 2), y - 20, 250 * hpPercent, 12, 0xe74c3c);
        this.container.add(hpBar);

        this.playerHpText = this.scene.add.text(x, y - 20, `${player.hp}/${player.max_hp}`, {
            font: '10px Microsoft YaHei', fill: '#fff'
        }).setOrigin(0.5);
        this.container.add(this.playerHpText);

        this.playerHpBar = hpBar;

        const mpLabel = this.scene.add.text(x - 180, y + 5, 'MP', {
            font: '12px Microsoft YaHei', fill: '#3498db'
        });
        this.container.add(mpLabel);

        const mpBarBg = this.scene.add.rectangle(x, y + 10, 250, 16, 0x333333);
        this.container.add(mpBarBg);

        const mpPercent = player.mp / player.max_mp;
        const mpBar = this.scene.add.rectangle(x - (250 * (1 - mpPercent) / 2), y + 10, 250 * mpPercent, 12, 0x3498db);
        this.container.add(mpBar);

        this.playerMpText = this.scene.add.text(x, y + 10, `${player.mp}/${player.max_mp}`, {
            font: '10px Microsoft YaHei', fill: '#fff'
        }).setOrigin(0.5);
        this.container.add(this.playerMpText);

        this.playerMpBar = mpBar;
    }

    createActionButtons(x, y) {
        const buttons = [
            { text: '⚔️ 攻击', action: 'attack', color: 0xe74c3c },
            { text: '✨ 技能', action: 'skill', color: 0x9b59b6 },
            { text: '🧪 物品', action: 'item', color: 0x2ecc71 },
            { text: '🏃 逃跑', action: 'flee', color: 0xf39c12 }
        ];

        buttons.forEach((btn, i) => {
            const bx = x - 180 + i * 120;
            const btnBg = this.scene.add.rectangle(bx, y, 100, 40, btn.color).setStrokeStyle(1, 0xffffff);
            const btnText = this.scene.add.text(bx, y, btn.text, {
                font: '14px Microsoft YaHei', fill: '#fff'
            }).setOrigin(0.5);

            btnBg.setInteractive({ useHandCursor: true });
            btnBg.on('pointerdown', () => this.handleAction(btn.action));
            btnBg.on('pointerover', () => btnBg.setAlpha(0.8));
            btnBg.on('pointerout', () => btnBg.setAlpha(1));

            this.container.add(btnBg);
            this.container.add(btnText);
        });
    }

    createCombatLog(x, y) {
        const logBg = this.scene.add.rectangle(x, y, 400, 150, 0x1a1a1a, 0.9).setStrokeStyle(1, 0x555555);
        this.container.add(logBg);

        const logTitle = this.scene.add.text(x - 180, y - 65, '战斗记录', {
            font: 'bold 12px Microsoft YaHei', fill: '#d4a574'
        });
        this.container.add(logTitle);

        this.logText = this.scene.add.text(x - 180, y - 45, '', {
            font: '11px Microsoft YaHei', fill: '#ccc', wordWrap: { width: 360 }, lineSpacing: 4
        });
        this.container.add(this.logText);

        this.updateCombatLog();
    }

    updateCombatLog() {
        const state = this.combatManager.getCombatState();
        const logs = state.combatLog.slice(-8).join('\n');
        this.logText.setText(logs);
    }

    handleAction(action) {
        if (this.combatManager.turn !== 'player') return;

        let result = null;

        switch (action) {
            case 'attack':
                result = this.combatManager.playerAttack();
                if (result) {
                    this.showDamageAnimation(result.damage, 'enemy');
                    this.updateEnemyDisplay();
                }
                break;

            case 'skill':
                this.showSkillSelection();
                return;

            case 'item':
                this.showItemSelection();
                return;

            case 'flee':
                result = this.combatManager.playerFlee();
                if (result && result.success) {
                    this.endCombat(false, true);
                    return;
                }
                break;
        }

        this.updateCombatLog();

        if (this.combatManager.checkCombatEnd()) {
            this.endCombat(true, false);
            return;
        }

        if (this.combatManager.turn === 'enemy') {
            this.scene.time.delayedCall(800, () => {
                this.enemyTurnAction();
            });
        }
    }

    enemyTurnAction() {
        const result = this.combatManager.enemyTurn();
        if (result) {
            this.showDamageAnimation(result.damage, 'player');
            this.updatePlayerDisplay();
            this.updateCombatLog();

            if (result.combatEnded) {
                this.endCombat(false, false);
            }
        }
    }

    showDamageAnimation(damage, target) {
        const x = target === 'enemy' ? this.scene.cameras.main.width / 2 : this.scene.cameras.main.width / 2;
        const y = target === 'enemy' ? 120 : this.scene.cameras.main.height - 120;

        const damageText = this.scene.add.text(x + Phaser.Math.Between(-30, 30), y, `-${damage}`, {
            font: 'bold 28px Microsoft YaHei', fill: '#ff0000', stroke: '#000', strokeThickness: 3
        }).setOrigin(0.5).setDepth(800).setScrollFactor(0);

        this.scene.tweens.add({
            targets: damageText,
            y: y - 60,
            alpha: 0,
            duration: 1000,
            ease: 'Power2',
            onComplete: () => damageText.destroy()
        });

        this.createAttackParticles(x, y, target === 'enemy' ? 0xe74c3c : 0x3498db);
        this.scene.cameras.main.shake(100, 0.01);
    }

    createAttackParticles(x, y, color) {
        for (let i = 0; i < 10; i++) {
            const particle = this.scene.add.circle(
                x + Phaser.Math.Between(-20, 20),
                y + Phaser.Math.Between(-20, 20),
                Phaser.Math.Between(2, 5),
                color
            ).setDepth(750).setScrollFactor(0);

            this.scene.tweens.add({
                targets: particle,
                x: particle.x + Phaser.Math.Between(-50, 50),
                y: particle.y + Phaser.Math.Between(-50, 50),
                alpha: 0,
                scale: 0,
                duration: Phaser.Math.Between(300, 600),
                ease: 'Power2',
                onComplete: () => particle.destroy()
            });
        }
    }

    updateEnemyDisplay() {
        const state = this.combatManager.getCombatState();
        const enemy = state.enemy;
        const hpPercent = Math.max(0, enemy.hp / enemy.max_hp);

        if (this.enemyHpBar) {
            this.enemyHpBar.setSize(200 * hpPercent, 16);
            this.enemyHpBar.x = this.scene.cameras.main.width / 2 - (200 * (1 - hpPercent) / 2);
        }
        if (this.enemyHpText) {
            this.enemyHpText.setText(`${Math.max(0, enemy.hp)}/${enemy.max_hp}`);
        }
    }

    updatePlayerDisplay() {
        const state = this.combatManager.getCombatState();
        const player = state.player;
        const x = this.scene.cameras.main.width / 2;
        const y = this.scene.cameras.main.height - 120;

        const hpPercent = Math.max(0, player.hp / player.max_hp);
        if (this.playerHpBar) {
            this.playerHpBar.setSize(250 * hpPercent, 12);
            this.playerHpBar.x = x - (250 * (1 - hpPercent) / 2);
        }
        if (this.playerHpText) {
            this.playerHpText.setText(`${Math.max(0, player.hp)}/${player.max_hp}`);
        }

        const mpPercent = Math.max(0, player.mp / player.max_mp);
        if (this.playerMpBar) {
            this.playerMpBar.setSize(250 * mpPercent, 12);
            this.playerMpBar.x = x - (250 * (1 - mpPercent) / 2);
        }
        if (this.playerMpText) {
            this.playerMpText.setText(`${Math.max(0, player.mp)}/${player.max_mp}`);
        }
    }

    showItemSelection() {
        const items = this.inventoryManager.getItems().filter(i => i.type === 'consumable');
        if (items.length === 0) {
            this.scene.showNotification('没有可用的物品！');
            return;
        }

        const width = this.scene.cameras.main.width;
        const height = this.scene.cameras.main.height;

        const itemContainer = this.scene.add.container(0, 0).setDepth(750).setScrollFactor(0);
        const overlay = this.scene.add.rectangle(width / 2, height / 2, width, height, 0x000000, 0.5);
        itemContainer.add(overlay);

        const panelBg = this.scene.add.rectangle(width / 2, height / 2, 300, 250, 0x1a1a1a, 0.95).setStrokeStyle(2, 0x2ecc71);
        itemContainer.add(panelBg);

        const title = this.scene.add.text(width / 2, height / 2 - 100, '选择物品', {
            font: 'bold 16px Microsoft YaHei', fill: '#2ecc71'
        }).setOrigin(0.5);
        itemContainer.add(title);

        items.slice(0, 5).forEach((item, i) => {
            const iy = height / 2 - 60 + i * 40;
            const itemBg = this.scene.add.rectangle(width / 2, iy, 260, 35, 0x2a2a2a).setStrokeStyle(1, 0x555555);
            itemContainer.add(itemBg);

            const itemText = this.scene.add.text(width / 2, iy, `${this.inventoryManager.getItemIcon(item)} ${item.name} x${item.count}`, {
                font: '14px Microsoft YaHei', fill: '#fff'
            }).setOrigin(0.5);
            itemContainer.add(itemText);

            itemBg.setInteractive({ useHandCursor: true });
            itemBg.on('pointerdown', () => {
                const result = this.combatManager.playerUseItem(item.id);
                if (result) {
                    this.updatePlayerDisplay();
                    this.updateCombatLog();
                    itemContainer.destroy();

                    if (this.combatManager.turn === 'enemy') {
                        this.scene.time.delayedCall(800, () => {
                            this.enemyTurnAction();
                        });
                    }
                }
            });
        });

        const closeBtn = this.scene.add.text(width / 2, height / 2 + 100, '取消', {
            font: '14px Microsoft YaHei', fill: '#ff6b6b'
        }).setOrigin(0.5).setInteractive({ useHandCursor: true });
        closeBtn.on('pointerdown', () => itemContainer.destroy());
        itemContainer.add(closeBtn);
    }

    showSkillSelection() {
        const availableSkills = this.skills.filter(s => s.level <= (this.combatManager.playerData.level || 1));
        if (availableSkills.length === 0) {
            this.scene.showNotification('没有可用的技能！');
            return;
        }

        const width = this.scene.cameras.main.width;
        const height = this.scene.cameras.main.height;

        const skillContainer = this.scene.add.container(0, 0).setDepth(750).setScrollFactor(0);
        const overlay = this.scene.add.rectangle(width / 2, height / 2, width, height, 0x000000, 0.5);
        skillContainer.add(overlay);

        const panelBg = this.scene.add.rectangle(width / 2, height / 2, 320, 280, 0x1a1a2a, 0.95).setStrokeStyle(2, 0x9b59b6);
        skillContainer.add(panelBg);

        const title = this.scene.add.text(width / 2, height / 2 - 120, '选择技能', {
            font: 'bold 16px Microsoft YaHei', fill: '#9b59b6'
        }).setOrigin(0.5);
        skillContainer.add(title);

        const playerMp = this.combatManager.playerData.mp || 0;

        availableSkills.slice(0, 5).forEach((skill, i) => {
            const iy = height / 2 - 70 + i * 45;
            const canAfford = playerMp >= skill.mp_cost;
            const bgColor = canAfford ? 0x2a2a4a : 0x3a2a2a;

            const skillBg = this.scene.add.rectangle(width / 2, iy, 280, 40, bgColor)
                .setStrokeStyle(1, canAfford ? 0x9b59b6 : 0x555555);
            skillContainer.add(skillBg);

            const icon = this.getSkillIcon(skill.type);
            const skillText = this.scene.add.text(width / 2 - 120, iy - 8, `${icon} ${skill.name}`, {
                font: '14px Microsoft YaHei', fill: canAfford ? '#fff' : '#888'
            });
            skillContainer.add(skillText);

            const costText = this.scene.add.text(width / 2 + 100, iy - 8, `${skill.mp_cost}MP`, {
                font: '12px Microsoft YaHei', fill: canAfford ? '#3498db' : '#e74c3c'
            });
            skillContainer.add(costText);

            const descText = this.scene.add.text(width / 2 - 120, iy + 8, skill.description, {
                font: '9px Microsoft YaHei', fill: '#888',
                wordWrap: { width: 240 }
            });
            skillContainer.add(descText);

            if (canAfford) {
                skillBg.setInteractive({ useHandCursor: true });
                skillBg.on('pointerdown', () => {
                    skillContainer.destroy();
                    if (this.onSkillUse) {
                        this.onSkillUse(skill);
                    }
                });
                skillBg.on('pointerover', () => skillBg.setAlpha(0.8));
                skillBg.on('pointerout', () => skillBg.setAlpha(1));
            }
        });

        const closeBtn = this.scene.add.text(width / 2, height / 2 + 120, '取消', {
            font: '14px Microsoft YaHei', fill: '#ff6b6b'
        }).setOrigin(0.5).setInteractive({ useHandCursor: true });
        closeBtn.on('pointerdown', () => skillContainer.destroy());
        skillContainer.add(closeBtn);
    }

    getSkillIcon(type) {
        switch (type) {
            case 'attack': return '⚔️';
            case 'heal': return '💚';
            case 'buff': return '🛡️';
            case 'debuff': return '☠️';
            default: return '✨';
        }
    }

    endCombat(victory, fled) {
        if (fled) {
            this.scene.showNotification('成功逃脱！');
            this.close();
            if (this.onCombatEnd) this.onCombatEnd({ victory: false, fled: true });
            return;
        }

        if (victory) {
            const rewards = this.combatManager.rewards;
            this.showVictoryScreen(rewards);
        } else {
            this.showDefeatScreen();
        }
    }

    showVictoryScreen(rewards) {
        const width = this.scene.cameras.main.width;
        const height = this.scene.cameras.main.height;

        const victoryContainer = this.scene.add.container(0, 0).setDepth(800).setScrollFactor(0);
        const overlay = this.scene.add.rectangle(width / 2, height / 2, width, height, 0x000000, 0.8);
        victoryContainer.add(overlay);

        const panelBg = this.scene.add.rectangle(width / 2, height / 2, 380, 350, 0x1a2a1a, 0.95).setStrokeStyle(2, 0xf1c40f);
        victoryContainer.add(panelBg);

        const victoryText = this.scene.add.text(width / 2, height / 2 - 140, '🎉 战斗胜利！', {
            font: 'bold 24px Microsoft YaHei', fill: '#f1c40f'
        }).setOrigin(0.5);
        victoryContainer.add(victoryText);

        const state = this.combatManager.getCombatState();
        const enemyName = state.enemy ? state.enemy.name : '敌人';
        const enemyLabel = this.scene.add.text(width / 2, height / 2 - 105, `击败了 ${enemyName}`, {
            font: '14px Microsoft YaHei', fill: '#e74c3c'
        }).setOrigin(0.5);
        victoryContainer.add(enemyLabel);

        let rewardY = height / 2 - 65;
        const lineH = 28;

        const expLine = this.scene.add.text(width / 2, rewardY, `✨ 经验: +${rewards.exp}`, {
            font: '16px Microsoft YaHei', fill: '#9b59b6'
        }).setOrigin(0.5);
        victoryContainer.add(expLine);
        rewardY += lineH;

        const goldLine = this.scene.add.text(width / 2, rewardY, `💰 金币: +${rewards.gold}`, {
            font: '16px Microsoft YaHei', fill: '#f1c40f'
        }).setOrigin(0.5);
        victoryContainer.add(goldLine);
        rewardY += lineH;

        if (rewards.items && rewards.items.length > 0) {
            const itemsLabel = this.scene.add.text(width / 2, rewardY, '🎒 获得物品:', {
                font: '14px Microsoft YaHei', fill: '#2ecc71'
            }).setOrigin(0.5);
            victoryContainer.add(itemsLabel);
            rewardY += 22;

            rewards.items.forEach(item => {
                const itemLine = this.scene.add.text(width / 2, rewardY, `  ${item.name} x${item.count}`, {
                    font: '13px Microsoft YaHei', fill: '#ccc'
                }).setOrigin(0.5);
                victoryContainer.add(itemLine);
                rewardY += 20;
            });
        }

        // Check for level up
        const prevLevel = this.scene.playerData.level || 1;
        const newExp = (this.scene.playerData.exp || 0) + rewards.exp;
        const expNeeded = prevLevel * 100;
        if (newExp >= expNeeded) {
            const levelUpText = this.scene.add.text(width / 2, rewardY + 10, '🎊 等级提升！', {
                font: 'bold 18px Microsoft YaHei', fill: '#f39c12'
            }).setOrigin(0.5);
            victoryContainer.add(levelUpText);

            this.scene.tweens.add({
                targets: levelUpText,
                scaleX: 1.2,
                scaleY: 1.2,
                duration: 500,
                yoyo: true,
                repeat: 2,
                ease: 'Sine.easeInOut'
            });
        }

        const continueBtn = this.scene.add.rectangle(width / 2, height / 2 + 140, 160, 45, 0x4a7c59)
            .setStrokeStyle(2, 0x6aaa79).setScrollFactor(0).setDepth(801);
        const continueText = this.scene.add.text(width / 2, height / 2 + 140, '▶ 继续', {
            font: 'bold 18px Microsoft YaHei', fill: '#fff'
        }).setOrigin(0.5).setScrollFactor(0).setDepth(802);

        continueBtn.setInteractive({ useHandCursor: true });
        continueBtn.on('pointerover', () => continueBtn.setFillStyle(0x5a8c69));
        continueBtn.on('pointerout', () => continueBtn.setFillStyle(0x4a7c59));
        continueBtn.on('pointerdown', () => {
            victoryContainer.destroy();
            this.close();
            if (this.onCombatEnd) this.onCombatEnd({ victory: true, rewards });
        });
        victoryContainer.add(continueBtn);
        victoryContainer.add(continueText);
    }

    showDefeatScreen() {
        const width = this.scene.cameras.main.width;
        const height = this.scene.cameras.main.height;

        const state = this.combatManager.getCombatState();
        const enemyName = state.enemy ? state.enemy.name : '敌人';

        const defeatContainer = this.scene.add.container(0, 0).setDepth(800).setScrollFactor(0);
        const overlay = this.scene.add.rectangle(width / 2, height / 2, width, height, 0x000000, 0.9);
        defeatContainer.add(overlay);

        const panelBg = this.scene.add.rectangle(width / 2, height / 2, 350, 220, 0x2a1a1a, 0.95).setStrokeStyle(2, 0xe74c3c);
        defeatContainer.add(panelBg);

        const defeatText = this.scene.add.text(width / 2, height / 2 - 60, '💀 战斗失败', {
            font: 'bold 24px Microsoft YaHei', fill: '#e74c3c'
        }).setOrigin(0.5);
        defeatContainer.add(defeatText);

        const messageText = this.scene.add.text(width / 2, height / 2 - 20, `你被 ${enemyName} 击败了...`, {
            font: '16px Microsoft YaHei', fill: '#aaa'
        }).setOrigin(0.5);
        defeatContainer.add(messageText);

        const respawnBtn = this.scene.add.rectangle(width / 2, height / 2 + 60, 140, 42, 0x4a2a2a)
            .setStrokeStyle(2, 0xe74c3c).setScrollFactor(0).setDepth(801);
        const respawnText = this.scene.add.text(width / 2, height / 2 + 60, '复活', {
            font: 'bold 18px Microsoft YaHei', fill: '#e74c3c'
        }).setOrigin(0.5).setScrollFactor(0).setDepth(802);

        respawnBtn.setInteractive({ useHandCursor: true });
        respawnBtn.on('pointerdown', () => {
            this.scene.playerData.hp = Math.floor(this.inventoryManager.getStats().max_hp * 0.5);
            defeatContainer.destroy();
            this.close();
            if (this.onCombatEnd) this.onCombatEnd({ victory: false, fled: false });
        });
        defeatContainer.add(respawnBtn);
        defeatContainer.add(respawnText);
    }
}
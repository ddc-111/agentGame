export class SkillBar {
    constructor(scene, playerData, combatManager) {
        this.scene = scene;
        this.playerData = playerData;
        this.combatManager = combatManager;
        this.container = null;
        this.skills = [];
        this.cooldowns = {};
        this.onSkillUse = null;
    }

    setSkills(skills) {
        this.skills = skills || [];
    }

    create(x, y, width) {
        if (this.container) {
            this.container.destroy();
        }

        this.container = this.scene.add.container(x, y).setDepth(710).setScrollFactor(0);

        const barWidth = Math.min(width - 40, 500);
        const bg = this.scene.add.rectangle(0, 0, barWidth, 50, 0x1a1a2a, 0.9)
            .setStrokeStyle(1, 0x555555);
        this.container.add(bg);

        const title = this.scene.add.text(-barWidth/2 + 10, -18, '技能', {
            font: '10px Microsoft YaHei', fill: '#888'
        });
        this.container.add(title);

        this.skillButtons = [];
        const startX = -barWidth/2 + 30;
        const spacing = 60;

        this.skills.slice(0, 5).forEach((skill, i) => {
            const bx = startX + i * spacing;
            this.createSkillButton(bx, 5, skill, i + 1);
        });
    }

    createSkillButton(x, y, skill, keyNum) {
        const btnSize = 40;
        const isOnCooldown = this.cooldowns[skill.id] > 0;
        const canAfford = this.playerData.mp >= skill.mp_cost;

        let btnColor = 0x2a2a4a;
        if (isOnCooldown) btnColor = 0x333333;
        else if (!canAfford) btnColor = 0x4a2a2a;

        const btn = this.scene.add.rectangle(x, y, btnSize, btnSize, btnColor)
            .setStrokeStyle(1, isOnCooldown ? 0x555555 : 0x9b59b6);
        this.container.add(btn);

        const icon = this.getSkillIcon(skill.type);
        const iconText = this.scene.add.text(x, y - 5, icon, {
            font: '16px Microsoft YaHei'
        }).setOrigin(0.5);
        this.container.add(iconText);

        const keyText = this.scene.add.text(x - btnSize/2 + 4, y - btnSize/2 + 2, `${keyNum}`, {
            font: '9px Microsoft YaHei', fill: '#aaa'
        });
        this.container.add(keyText);

        const costText = this.scene.add.text(x, y + 12, `${skill.mp_cost}MP`, {
            font: '8px Microsoft YaHei', fill: canAfford ? '#3498db' : '#e74c3c'
        }).setOrigin(0.5);
        this.container.add(costText);

        if (isOnCooldown) {
            const cdText = this.scene.add.text(x, y, `${this.cooldowns[skill.id]}`, {
                font: 'bold 16px Microsoft YaHei', fill: '#ff0000'
            }).setOrigin(0.5);
            this.container.add(cdText);
        }

        btn.setInteractive({ useHandCursor: !isOnCooldown && canAfford });
        btn.on('pointerdown', () => {
            if (!isOnCooldown && canAfford) {
                this.useSkill(skill);
            }
        });
        btn.on('pointerover', () => {
            if (!isOnCooldown && canAfford) btn.setAlpha(0.8);
            this.showTooltip(skill, x, y - 40);
        });
        btn.on('pointerout', () => {
            btn.setAlpha(1);
            this.hideTooltip();
        });

        this.skillButtons.push({ btn, skill, iconText, costText });
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

    showTooltip(skill, x, y) {
        this.hideTooltip();
        const tipWidth = 200;
        const tipHeight = 80;
        const tipX = x;
        const tipY = y - tipHeight/2;

        this.tooltip = this.scene.add.container(tipX, tipY).setDepth(720).setScrollFactor(0);

        const bg = this.scene.add.rectangle(0, 0, tipWidth, tipHeight, 0x1a1a1a, 0.95)
            .setStrokeStyle(1, 0x9b59b6);
        this.tooltip.add(bg);

        const nameText = this.scene.add.text(0, -25, skill.name, {
            font: 'bold 12px Microsoft YaHei', fill: '#9b59b6'
        }).setOrigin(0.5);
        this.tooltip.add(nameText);

        const descText = this.scene.add.text(0, 0, skill.description, {
            font: '10px Microsoft YaHei', fill: '#ccc',
            wordWrap: { width: tipWidth - 20 }
        }).setOrigin(0.5);
        this.tooltip.add(descText);

        const statsText = this.scene.add.text(0, 25, `消耗: ${skill.mp_cost} MP | 冷却: ${skill.cooldown}回合`, {
            font: '9px Microsoft YaHei', fill: '#888'
        }).setOrigin(0.5);
        this.tooltip.add(statsText);
    }

    hideTooltip() {
        if (this.tooltip) {
            this.tooltip.destroy();
            this.tooltip = null;
        }
    }

    useSkill(skill) {
        if (this.onSkillUse) {
            this.onSkillUse(skill);
        }
        if (skill.cooldown > 0) {
            this.cooldowns[skill.id] = skill.cooldown;
        }
    }

    updateCooldowns() {
        for (const id in this.cooldowns) {
            if (this.cooldowns[id] > 0) {
                this.cooldowns[id]--;
            }
        }
    }

    update() {
        if (!this.container || !this.container.active) return;
        this.container.destroy();
        this.create(this.container.x, this.container.y, 500);
    }

    destroy() {
        this.hideTooltip();
        if (this.container) {
            this.container.destroy();
            this.container = null;
        }
    }
}

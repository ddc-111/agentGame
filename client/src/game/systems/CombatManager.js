export class CombatManager {
    constructor(playerData, inventoryManager) {
        this.playerData = playerData;
        this.inventoryManager = inventoryManager;
        this.inCombat = false;
        this.enemy = null;
        this.combatLog = [];
        this.turn = 'player';
        this.rewards = null;
        this.skills = [];
        this.cooldowns = {};
    }

    setSkills(skills) {
        this.skills = skills || [];
    }

    getEnemyData(enemyType) {
        const enemies = {
            'wolf': {
                name: '野狼',
                hp: 30,
                max_hp: 30,
                attack: 8,
                defense: 3,
                speed: 12,
                exp: 15,
                gold: 5,
                drops: [
                    { item_id: 1001, chance: 0.3, count: 1 }
                ],
                emoji: '🐺'
            },
            'alpha_wolf': {
                name: '头狼',
                hp: 60,
                max_hp: 60,
                attack: 15,
                defense: 8,
                speed: 10,
                exp: 40,
                gold: 20,
                drops: [
                    { item_id: 1002, chance: 0.5, count: 1 },
                    { item_id: 1001, chance: 0.8, count: 2 }
                ],
                emoji: '🐺'
            }
        };
        return enemies[enemyType] || enemies['wolf'];
    }

    startCombat(enemyType) {
        this.inCombat = true;
        this.enemy = { ...this.getEnemyData(enemyType) };
        this.combatLog = [];
        this.turn = 'player';
        this.rewards = null;
        this.cooldowns = {};
        this.addLog(`遭遇了 ${this.enemy.name}！`);
        return this.getCombatState();
    }

    addLog(text) {
        this.combatLog.push(text);
        if (this.combatLog.length > 50) {
            this.combatLog.shift();
        }
    }

    getCombatState() {
        const stats = this.inventoryManager.getStats();
        return {
            inCombat: this.inCombat,
            enemy: { ...this.enemy },
            player: {
                hp: this.playerData.hp,
                max_hp: stats.max_hp,
                mp: this.playerData.mp,
                max_mp: stats.max_mp,
                attack: stats.attack,
                defense: stats.defense
            },
            combatLog: [...this.combatLog],
            turn: this.turn,
            rewards: this.rewards
        };
    }

    calculateDamage(attackerAtk, defenderDef) {
        const baseDmg = Math.max(1, attackerAtk - defenderDef);
        const variance = Math.floor(baseDmg * 0.2);
        return baseDmg + Math.floor(Math.random() * variance * 2) - variance;
    }

    playerAttack() {
        if (this.turn !== 'player' || !this.inCombat) return null;

        const stats = this.inventoryManager.getStats();
        const damage = this.calculateDamage(stats.attack, this.enemy.defense);
        this.enemy.hp = Math.max(0, this.enemy.hp - damage);
        this.addLog(`你攻击了 ${this.enemy.name}，造成 ${damage} 点伤害！`);

        const result = {
            action: 'attack',
            damage,
            target: 'enemy',
            enemyHp: this.enemy.hp,
            enemyMaxHp: this.enemy.max_hp
        };

        if (this.checkCombatEnd()) {
            return result;
        }

        this.turn = 'enemy';
        return result;
    }

    playerUseItem(itemId) {
        if (this.turn !== 'player' || !this.inCombat) return null;

        const useResult = this.inventoryManager.useItem(itemId);
        if (!useResult) return null;

        this.addLog(`你使用了 ${useResult.item.name}！`);
        if (useResult.effect.hp) {
            this.addLog(`恢复了 ${useResult.effect.hp} 点生命！`);
        }
        if (useResult.effect.mp) {
            this.addLog(`恢复了 ${useResult.effect.mp} 点法力！`);
        }

        this.turn = 'enemy';
        return {
            action: 'useItem',
            item: useResult.item,
            effect: useResult.effect,
            playerHp: this.playerData.hp,
            playerMp: this.playerData.mp
        };
    }

    playerFlee() {
        if (this.turn !== 'player' || !this.inCombat) return null;

        const fleeChance = 0.5;
        const success = Math.random() < fleeChance;

        if (success) {
            this.addLog('你成功逃离了战斗！');
            this.inCombat = false;
            return { action: 'flee', success: true };
        } else {
            this.addLog('逃跑失败！');
            this.turn = 'enemy';
            return { action: 'flee', success: false };
        }
    }

    playerUseSkill(skill) {
        if (this.turn !== 'player' || !this.inCombat) return null;

        if (this.playerData.mp < skill.mp_cost) {
            return null;
        }

        if (this.cooldowns[skill.id] > 0) {
            return null;
        }

        this.playerData.mp -= skill.mp_cost;

        const stats = this.inventoryManager.getStats();
        let result = { action: 'skill', skill, damage: 0, heal: 0 };

        switch (skill.type) {
            case 'attack': {
                const baseDmg = skill.damage + stats.attack - this.enemy.defense;
                const damage = Math.max(1, baseDmg + Math.floor(Math.random() * 6) - 3);
                this.enemy.hp = Math.max(0, this.enemy.hp - damage);
                result.damage = damage;
                this.addLog(`使用【${skill.name}】对 ${this.enemy.name} 造成 ${damage} 点伤害！`);
                break;
            }
            case 'heal': {
                const healAmount = skill.heal;
                const maxHp = stats.max_hp;
                this.playerData.hp = Math.min(maxHp, this.playerData.hp + healAmount);
                result.heal = healAmount;
                this.addLog(`使用【${skill.name}】恢复了 ${healAmount} 点生命！`);
                break;
            }
            case 'buff': {
                this.addLog(`使用【${skill.name}】！增益效果生效！`);
                result.effect = skill.effect;
                break;
            }
            case 'debuff': {
                const baseDmg = skill.damage + stats.attack - this.enemy.defense;
                const damage = Math.max(1, baseDmg + Math.floor(Math.random() * 6) - 3);
                this.enemy.hp = Math.max(0, this.enemy.hp - damage);
                result.damage = damage;
                this.addLog(`对 ${this.enemy.name} 使用【${skill.name}】，造成 ${damage} 点伤害并附加减益效果！`);
                break;
            }
        }

        if (skill.cooldown > 0) {
            this.cooldowns[skill.id] = skill.cooldown;
        }

        if (this.checkCombatEnd()) {
            return result;
        }

        this.turn = 'enemy';
        return result;
    }

    updateCooldowns() {
        for (const id in this.cooldowns) {
            if (this.cooldowns[id] > 0) {
                this.cooldowns[id]--;
            }
        }
    }

    enemyTurn() {
        if (this.turn !== 'enemy' || !this.inCombat) return null;

        const stats = this.inventoryManager.getStats();
        const damage = this.calculateDamage(this.enemy.attack, stats.defense);
        this.playerData.hp = Math.max(0, this.playerData.hp - damage);
        this.addLog(`${this.enemy.name} 攻击了你，造成 ${damage} 点伤害！`);

        const result = {
            action: 'enemyAttack',
            damage,
            target: 'player',
            playerHp: this.playerData.hp,
            playerMaxHp: stats.max_hp
        };

        this.turn = 'player';
        this.updateCooldowns();

        if (this.checkCombatEnd()) {
            result.combatEnded = true;
            result.victory = false;
        }

        return result;
    }

    checkCombatEnd() {
        if (this.enemy.hp <= 0) {
            this.inCombat = false;
            this.rewards = this.getRewards();
            this.addLog(`击败了 ${this.enemy.name}！`);
            this.addLog(`获得 ${this.rewards.exp} 经验，${this.rewards.gold} 金币！`);
            if (this.rewards.items.length > 0) {
                this.rewards.items.forEach(item => {
                    this.addLog(`获得了 ${item.name} x${item.count}！`);
                });
            }
            return true;
        }

        if (this.playerData.hp <= 0) {
            this.inCombat = false;
            this.addLog('你被击败了...');
            return true;
        }

        return false;
    }

    getRewards() {
        if (!this.enemy) return null;

        const rewards = {
            exp: this.enemy.exp,
            gold: this.enemy.gold,
            items: []
        };

        if (this.enemy.drops) {
            this.enemy.drops.forEach(drop => {
                if (Math.random() < drop.chance) {
                    const itemData = this.inventoryManager.getItemData(drop.item_id);
                    rewards.items.push({
                        id: drop.item_id,
                        name: itemData.name,
                        count: drop.count
                    });
                    this.inventoryManager.addItem(drop.item_id, drop.count);
                }
            });
        }

        this.playerData.gold += rewards.gold;
        this.playerData.exp = (this.playerData.exp || 0) + rewards.exp;

        return rewards;
    }

    applyRewards() {
        if (!this.rewards) return;
        this.playerData.gold += this.rewards.gold;
        this.playerData.exp = (this.playerData.exp || 0) + this.rewards.exp;
    }
}
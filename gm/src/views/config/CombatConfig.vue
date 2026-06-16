<template>
  <div class="combat-config">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>战斗配置</span>
          <el-button type="primary" @click="handleSave">保存配置</el-button>
        </div>
      </template>

      <el-tabs v-model="activeTab">
        <el-tab-pane label="敌人配置" name="enemies">
          <div class="section-header">
            <span>敌人属性配置</span>
            <el-button type="primary" size="small" @click="addEnemy">
              <el-icon><Plus /></el-icon> 添加敌人
            </el-button>
          </div>
          <el-table :data="enemies" style="width: 100%">
            <el-table-column prop="name" label="名称" width="120" />
            <el-table-column prop="hp" label="生命" width="80" />
            <el-table-column prop="attack" label="攻击" width="80" />
            <el-table-column prop="defense" label="防御" width="80" />
            <el-table-column prop="speed" label="速度" width="80" />
            <el-table-column label="奖励" width="150">
              <template #default="{ row }">
                <span>经验:{{ row.rewardExp }} 金币:{{ row.rewardGold }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="description" label="描述" />
            <el-table-column label="操作" width="150" fixed="right">
              <template #default="{ row, $index }">
                <el-button type="primary" text @click="editEnemy($index)">编辑</el-button>
                <el-button type="danger" text @click="removeEnemy($index)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="战斗公式" name="formulas">
          <el-form :model="formulas" label-width="160px">
            <el-divider content-position="left">伤害计算</el-divider>
            <el-form-item label="物理伤害公式">
              <el-input v-model="formulas.physicalDamage" />
              <span class="form-tip">变量: ATK(攻击), DEF(防御), CRIT(暴击), RAND(随机0-1)</span>
            </el-form-item>
            <el-form-item label="魔法伤害公式">
              <el-input v-model="formulas.magicDamage" />
              <span class="form-tip">变量: MATK(魔攻), MDEF(魔防), SKILL(技能加成)</span>
            </el-form-item>
            <el-form-item label="治疗公式">
              <el-input v-model="formulas.healAmount" />
              <span class="form-tip">变量: HEAL(治疗力), SKILL(技能加成)</span>
            </el-form-item>

            <el-divider content-position="left">概率计算</el-divider>
            <el-form-item label="暴击率公式">
              <el-input v-model="formulas.criticalRate" />
            </el-form-item>
            <el-form-item label="暴击倍率">
              <el-input-number v-model="formulas.criticalMultiplier" :min="1" :max="5" :step="0.1" />
            </el-form-item>
            <el-form-item label="闪避率公式">
              <el-input v-model="formulas.dodgeRate" />
            </el-form-item>
            <el-form-item label="逃跑成功率">
              <el-input v-model="formulas.fleeRate" />
            </el-form-item>

            <el-divider content-position="left">经验计算</el-divider>
            <el-form-item label="经验公式">
              <el-input v-model="formulas.expFormula" />
              <span class="form-tip">变量: BASE(基础经验), LEVEL(敌人等级), DIFF(等级差)</span>
            </el-form-item>
            <el-form-item label="升级经验公式">
              <el-input v-model="formulas.levelUpExp" />
              <span class="form-tip">变量: LEVEL(当前等级)</span>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="战斗模拟" name="simulate">
          <el-row :gutter="20">
            <el-col :span="12">
              <el-card>
                <template #header>玩家配置</template>
                <el-form :model="simPlayer" label-width="80px">
                  <el-form-item label="等级">
                    <el-input-number v-model="simPlayer.level" :min="1" :max="100" />
                  </el-form-item>
                  <el-form-item label="生命">
                    <el-input-number v-model="simPlayer.hp" :min="1" :max="9999" />
                  </el-form-item>
                  <el-form-item label="攻击">
                    <el-input-number v-model="simPlayer.attack" :min="1" :max="999" />
                  </el-form-item>
                  <el-form-item label="防御">
                    <el-input-number v-model="simPlayer.defense" :min="0" :max="999" />
                  </el-form-item>
                  <el-form-item label="速度">
                    <el-input-number v-model="simPlayer.speed" :min="1" :max="999" />
                  </el-form-item>
                </el-form>
              </el-card>
            </el-col>
            <el-col :span="12">
              <el-card>
                <template #header>敌人选择</template>
                <el-select v-model="selectedEnemy" placeholder="选择敌人" style="width: 100%">
                  <el-option
                    v-for="(enemy, index) in enemies"
                    :key="index"
                    :label="enemy.name"
                    :value="index"
                  />
                </el-select>
                <div v-if="selectedEnemy !== null" class="enemy-preview">
                  <p>HP: {{ enemies[selectedEnemy].hp }}</p>
                  <p>攻击: {{ enemies[selectedEnemy].attack }}</p>
                  <p>防御: {{ enemies[selectedEnemy].defense }}</p>
                </div>
              </el-card>
            </el-col>
          </el-row>

          <div style="margin-top: 20px; text-align: center">
            <el-button type="primary" size="large" @click="runSimulation" :loading="simulating">
              开始模拟
            </el-button>
          </div>

          <el-card v-if="simLog.length > 0" style="margin-top: 20px">
            <template #header>
              <div class="card-header">
                <span>战斗日志</span>
                <el-tag :type="simResult === 'win' ? 'success' : 'danger'">
                  {{ simResult === 'win' ? '胜利' : '失败' }}
                </el-tag>
              </div>
            </template>
            <div class="sim-log">
              <div v-for="(log, index) in simLog" :key="index" class="log-entry" :class="log.type">
                <span class="log-turn">回合{{ log.turn }}</span>
                <span class="log-text">{{ log.text }}</span>
                <span class="log-hp">玩家HP: {{ log.playerHP }} | 敌人HP: {{ log.enemyHP }}</span>
              </div>
            </div>
          </el-card>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <el-dialog v-model="enemyDialogVisible" :title="editingEnemyIndex === -1 ? '添加敌人' : '编辑敌人'" width="500px">
      <el-form :model="enemyForm" label-width="80px">
        <el-form-item label="名称">
          <el-input v-model="enemyForm.name" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="enemyForm.description" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item label="生命">
          <el-input-number v-model="enemyForm.hp" :min="1" :max="99999" />
        </el-form-item>
        <el-form-item label="攻击">
          <el-input-number v-model="enemyForm.attack" :min="1" :max="9999" />
        </el-form-item>
        <el-form-item label="防御">
          <el-input-number v-model="enemyForm.defense" :min="0" :max="9999" />
        </el-form-item>
        <el-form-item label="速度">
          <el-input-number v-model="enemyForm.speed" :min="1" :max="9999" />
        </el-form-item>
        <el-form-item label="经验奖励">
          <el-input-number v-model="enemyForm.rewardExp" :min="0" :max="99999" />
        </el-form-item>
        <el-form-item label="金币奖励">
          <el-input-number v-model="enemyForm.rewardGold" :min="0" :max="99999" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="enemyDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveEnemy">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue';
import { ElMessage } from 'element-plus';

const activeTab = ref('enemies');

const enemies = ref([
  { name: '野狼', description: '村外常见的野兽', hp: 80, attack: 12, defense: 3, speed: 15, rewardExp: 50, rewardGold: 30 },
  { name: '头狼', description: '狼群的首领，体型巨大', hp: 200, attack: 25, defense: 8, speed: 20, rewardExp: 150, rewardGold: 100 },
  { name: '野兔', description: '温顺的小动物', hp: 20, attack: 3, defense: 1, speed: 30, rewardExp: 10, rewardGold: 5 },
  { name: '山贼', description: '在山路上抢劫的强盗', hp: 150, attack: 20, defense: 10, speed: 12, rewardExp: 100, rewardGold: 80 },
  { name: '毒蛇', description: '带有剧毒的蛇', hp: 40, attack: 15, defense: 2, speed: 25, rewardExp: 30, rewardGold: 15 }
]);

const formulas = reactive({
  physicalDamage: 'max(1, ATK * (1 + CRIT * 0.5) - DEF * 0.5) * (0.9 + RAND * 0.2)',
  magicDamage: 'MATK * SKILL * (1 - MDEF * 0.003)',
  healAmount: 'HEAL * SKILL',
  criticalRate: 'min(0.5, 0.05 + AGI * 0.002)',
  criticalMultiplier: 1.5,
  dodgeRate: 'min(0.3, AGI * 0.001)',
  fleeRate: '0.6 + (PLAYER_SPD - ENEMY_SPD) * 0.01',
  expFormula: 'BASE * (1 + DIFF * 0.1)',
  levelUpExp: '100 * LEVEL * (1 + LEVEL * 0.1)'
});

const simPlayer = reactive({
  level: 5,
  hp: 100,
  attack: 20,
  defense: 5,
  speed: 15
});

const selectedEnemy = ref(0);
const simulating = ref(false);
const simLog = ref([]);
const simResult = ref('');

const enemyDialogVisible = ref(false);
const editingEnemyIndex = ref(-1);
const enemyForm = reactive({
  name: '',
  description: '',
  hp: 100,
  attack: 10,
  defense: 5,
  speed: 10,
  rewardExp: 50,
  rewardGold: 30
});

const addEnemy = () => {
  editingEnemyIndex.value = -1;
  Object.assign(enemyForm, { name: '', description: '', hp: 100, attack: 10, defense: 5, speed: 10, rewardExp: 50, rewardGold: 30 });
  enemyDialogVisible.value = true;
};

const editEnemy = (index) => {
  editingEnemyIndex.value = index;
  Object.assign(enemyForm, enemies.value[index]);
  enemyDialogVisible.value = true;
};

const removeEnemy = (index) => {
  enemies.value.splice(index, 1);
};

const saveEnemy = () => {
  if (editingEnemyIndex.value === -1) {
    enemies.value.push({ ...enemyForm });
  } else {
    enemies.value[editingEnemyIndex.value] = { ...enemyForm };
  }
  enemyDialogVisible.value = false;
  ElMessage.success('保存成功');
};

const runSimulation = () => {
  simulating.value = true;
  simLog.value = [];

  let playerHP = simPlayer.hp;
  const playerATK = simPlayer.attack;
  const playerDEF = simPlayer.defense;
  const playerSPD = simPlayer.speed;

  const enemy = enemies.value[selectedEnemy.value];
  let enemyHP = enemy.hp;
  const enemyATK = enemy.attack;
  const enemyDEF = enemy.defense;
  const enemySPD = enemy.speed;

  let turn = 1;
  const maxTurns = 20;

  while (playerHP > 0 && enemyHP > 0 && turn <= maxTurns) {
    const playerFirst = playerSPD >= enemySPD;

    if (playerFirst) {
      const dmg = Math.max(1, playerATK - enemyDEF * 0.5);
      const actualDmg = Math.round(dmg * (0.9 + Math.random() * 0.2));
      enemyHP = Math.max(0, enemyHP - actualDmg);
      simLog.value.push({ turn, type: 'player', text: `玩家攻击造成${actualDmg}点伤害`, playerHP, enemyHP });

      if (enemyHP <= 0) break;

      const eDmg = Math.max(1, enemyATK - playerDEF * 0.5);
      const actualEDmg = Math.round(eDmg * (0.9 + Math.random() * 0.2));
      playerHP = Math.max(0, playerHP - actualEDmg);
      simLog.value.push({ turn, type: 'enemy', text: `敌人攻击造成${actualEDmg}点伤害`, playerHP, enemyHP });
    } else {
      const eDmg = Math.max(1, enemyATK - playerDEF * 0.5);
      const actualEDmg = Math.round(eDmg * (0.9 + Math.random() * 0.2));
      playerHP = Math.max(0, playerHP - actualEDmg);
      simLog.value.push({ turn, type: 'enemy', text: `敌人攻击造成${actualEDmg}点伤害`, playerHP, enemyHP });

      if (playerHP <= 0) break;

      const dmg = Math.max(1, playerATK - enemyDEF * 0.5);
      const actualDmg = Math.round(dmg * (0.9 + Math.random() * 0.2));
      enemyHP = Math.max(0, enemyHP - actualDmg);
      simLog.value.push({ turn, type: 'player', text: `玩家攻击造成${actualDmg}点伤害`, playerHP, enemyHP });
    }

    turn++;
  }

  simResult.value = enemyHP <= 0 ? 'win' : 'lose';
  simulating.value = false;
};

const handleSave = () => {
  ElMessage.success('战斗配置已保存');
};
</script>

<style scoped>
.combat-config {
  max-width: 1200px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.form-tip {
  margin-left: 12px;
  color: #909399;
  font-size: 12px;
}

.enemy-preview {
  margin-top: 16px;
  padding: 12px;
  background: #f5f7fa;
  border-radius: 4px;
}

.sim-log {
  max-height: 400px;
  overflow-y: auto;
}

.log-entry {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  border-bottom: 1px solid #ebeef5;
}

.log-entry.player {
  background: #f0f9ff;
}

.log-entry.enemy {
  background: #fef0f0;
}

.log-turn {
  width: 60px;
  color: #909399;
  font-size: 12px;
}

.log-text {
  flex: 1;
}

.log-hp {
  width: 200px;
  text-align: right;
  color: #606266;
  font-size: 12px;
}
</style>

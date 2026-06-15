<template>
  <div class="game-config">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>游戏配置</span>
          <el-button type="primary" @click="handleSave">保存配置</el-button>
        </div>
      </template>

      <el-tabs v-model="activeTab">
        <el-tab-pane label="基础配置" name="basic">
          <el-form :model="config.game" label-width="120px">
            <el-form-item label="游戏名称">
              <el-input v-model="config.game.name" />
            </el-form-item>
            <el-form-item label="版本号">
              <el-input v-model="config.game.version" />
            </el-form-item>
            <el-form-item label="游戏描述">
              <el-input v-model="config.game.description" type="textarea" :rows="3" />
            </el-form-item>
            <el-form-item label="最大玩家数">
              <el-input-number v-model="config.game.maxPlayers" :min="1" :max="10000" />
            </el-form-item>
            <el-form-item label="Tick频率">
              <el-input-number v-model="config.game.tickRate" :min="1" :max="60" />
              <span class="form-tip">每秒更新次数</span>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="玩家配置" name="player">
          <el-form :model="config.player" label-width="120px">
            <el-form-item label="初始场景">
              <el-select v-model="config.player.startScene">
                <el-option
                  v-for="scene in sceneStore.scenes"
                  :key="scene.id"
                  :label="scene.name"
                  :value="scene.id"
                />
              </el-select>
            </el-form-item>
            <el-form-item label="初始金币">
              <el-input-number v-model="config.player.startGold" :min="0" :max="999999" :step="100" />
            </el-form-item>
            <el-form-item label="初始等级">
              <el-input-number v-model="config.player.startLevel" :min="1" :max="100" />
            </el-form-item>
            <el-form-item label="最大等级">
              <el-input-number v-model="config.player.maxLevel" :min="1" :max="999" />
            </el-form-item>
            <el-form-item label="基础生命">
              <el-input-number v-model="config.player.baseHP" :min="1" :max="99999" :step="10" />
            </el-form-item>
            <el-form-item label="基础法力">
              <el-input-number v-model="config.player.baseMP" :min="0" :max="99999" :step="10" />
            </el-form-item>
            <el-form-item label="基础攻击">
              <el-input-number v-model="config.player.baseAttack" :min="1" :max="9999" />
            </el-form-item>
            <el-form-item label="基础防御">
              <el-input-number v-model="config.player.baseDefense" :min="0" :max="9999" />
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="世界配置" name="world">
          <el-form :model="config.world" label-width="120px">
            <el-form-item label="昼夜循环">
              <el-switch v-model="config.world.dayNightCycle" />
            </el-form-item>
            <el-form-item v-if="config.world.dayNightCycle" label="昼夜时长">
              <el-input-number v-model="config.world.dayDuration" :min="60" :max="7200" :step="60" />
              <span class="form-tip">秒</span>
            </el-form-item>
            <el-form-item label="天气系统">
              <el-switch v-model="config.world.weatherEnabled" />
            </el-form-item>
            <el-form-item v-if="config.world.weatherEnabled" label="天气类型">
              <el-select v-model="config.world.weatherTypes" multiple>
                <el-option label="晴天" value="sunny" />
                <el-option label="多云" value="cloudy" />
                <el-option label="雨天" value="rainy" />
                <el-option label="雪天" value="snowy" />
                <el-option label="雾天" value="foggy" />
              </el-select>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="战斗配置" name="combat">
          <el-form :model="config.combat" label-width="120px">
            <el-form-item label="启用战斗">
              <el-switch v-model="config.combat.enabled" />
            </el-form-item>
            <el-form-item v-if="config.combat.enabled" label="回合制">
              <el-switch v-model="config.combat.turnBased" />
            </el-form-item>
            <el-form-item v-if="config.combat.enabled" label="最大回合数">
              <el-input-number v-model="config.combat.maxTurns" :min="1" :max="100" />
            </el-form-item>
            <el-form-item v-if="config.combat.enabled" label="暴击率">
              <el-slider v-model="config.combat.criticalRate" :min="0" :max="1" :step="0.01" show-input />
            </el-form-item>
            <el-form-item v-if="config.combat.enabled" label="暴击倍率">
              <el-input-number v-model="config.combat.criticalMultiplier" :min="1" :max="5" :step="0.1" />
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="经济配置" name="economy">
          <el-form :model="config.economy" label-width="120px">
            <el-form-item label="通胀率">
              <el-slider v-model="config.economy.inflationRate" :min="0" :max="0.1" :step="0.001" show-input />
            </el-form-item>
            <el-form-item label="税率">
              <el-slider v-model="config.economy.taxRate" :min="0" :max="0.5" :step="0.01" show-input />
            </el-form-item>
            <el-form-item label="最大金币">
              <el-input-number v-model="config.economy.maxGold" :min="1000" :max="999999999" :step="1000" />
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useConfigStore, useSceneStore } from '@/stores';
import { ElMessage } from 'element-plus';

const configStore = useConfigStore();
const sceneStore = useSceneStore();

const activeTab = ref('basic');

const config = computed(() => configStore.gameConfig);

const handleSave = () => {
  ElMessage.success('配置已保存');
};
</script>

<style scoped>
.game-config {
  max-width: 1000px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.form-tip {
  margin-left: 12px;
  color: #909399;
  font-size: 12px;
}
</style>

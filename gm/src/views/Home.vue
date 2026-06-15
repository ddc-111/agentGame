<template>
  <div>
    <h1>GM管理端</h1>
    <el-row :gutter="20">
      <el-col :span="8">
        <el-card>
          <template #header>在线玩家</template>
          <div class="stat">0</div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card>
          <template #header>NPC数量</template>
          <div class="stat">0</div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card>
          <template #header>服务器状态</template>
          <div class="stat">正常</div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
</script>

<style scoped>
.stat {
  font-size: 24px;
  text-align: center;
  padding: 20px;
}
</style>

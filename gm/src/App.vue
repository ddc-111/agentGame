<template>
  <el-container>
    <el-aside width="200px">
      <el-menu :router="true">
        <el-menu-item index="/">
          <span>首页</span>
        </el-menu-item>
        <el-menu-item index="/npc">
          <span>NPC管理</span>
        </el-menu-item>
        <el-menu-item index="/player">
          <span>玩家管理</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-main>
      <router-view />
    </el-main>
  </el-container>
</template>

<script setup>
</script>

<style>
body {
  margin: 0;
  padding: 0;
}
</style>

<template>
  <div>
    <h2>玩家管理</h2>
    <el-table :data="playerList" style="width: 100%;">
      <el-table-column prop="id" label="ID" />
      <el-table-column prop="name" label="名称" />
      <el-table-column prop="level" label="等级" />
      <el-table-column prop="status" label="状态" />
      <el-table-column label="操作">
        <template #default>
          <el-button size="small">查看</el-button>
          <el-button size="small" type="warning">封禁</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup>
import { ref } from 'vue';

const playerList = ref([]);
</script>

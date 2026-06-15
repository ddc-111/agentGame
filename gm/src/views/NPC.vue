<template>
  <div>
    <h2>NPC管理</h2>
    <el-button type="primary" @click="addNPC">添加NPC</el-button>
    <el-table :data="npcList" style="width: 100%; margin-top: 20px;">
      <el-table-column prop="id" label="ID" />
      <el-table-column prop="name" label="名称" />
      <el-table-column prop="persona" label="人设" />
      <el-table-column prop="location" label="位置" />
      <el-table-column label="操作">
        <template #default>
          <el-button size="small">编辑</el-button>
          <el-button size="small" type="danger">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup>
import { ref } from 'vue';

const npcList = ref([]);

const addNPC = () => {
  console.log('Add NPC');
};
</script>

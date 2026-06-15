<template>
  <div class="tileset-manager">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>图块管理</span>
          <el-button type="primary" @click="handleUpload">
            <el-icon><Upload /></el-icon>
            上传图块
          </el-button>
        </div>
      </template>

      <el-tabs v-model="activeCategory">
        <el-tab-pane label="地形" name="terrain">
          <div class="tile-grid">
            <div v-for="tile in terrainTiles" :key="tile.id" class="tile-item" @click="handleSelect(tile)">
              <div class="tile-preview" :style="{ backgroundColor: tile.color }">
                <img v-if="tile.image" :src="tile.image" :alt="tile.name" />
              </div>
              <div class="tile-name">{{ tile.name }}</div>
            </div>
          </div>
        </el-tab-pane>
        <el-tab-pane label="建筑" name="building">
          <div class="tile-grid">
            <div v-for="tile in buildingTiles" :key="tile.id" class="tile-item" @click="handleSelect(tile)">
              <div class="tile-preview" :style="{ backgroundColor: tile.color }">
                <img v-if="tile.image" :src="tile.image" :alt="tile.name" />
              </div>
              <div class="tile-name">{{ tile.name }}</div>
            </div>
          </div>
        </el-tab-pane>
        <el-tab-pane label="装饰" name="decoration">
          <div class="tile-grid">
            <div v-for="tile in decorationTiles" :key="tile.id" class="tile-item" @click="handleSelect(tile)">
              <div class="tile-preview" :style="{ backgroundColor: tile.color }">
                <img v-if="tile.image" :src="tile.image" :alt="tile.name" />
              </div>
              <div class="tile-name">{{ tile.name }}</div>
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';

const activeCategory = ref('terrain');

const tiles = ref([
  { id: 'tile_001', name: '草地', category: 'terrain', color: '#90EE90', image: '' },
  { id: 'tile_002', name: '土地', category: 'terrain', color: '#D2B48C', image: '' },
  { id: 'tile_003', name: '石板路', category: 'terrain', color: '#808080', image: '' },
  { id: 'tile_004', name: '水面', category: 'terrain', color: '#87CEEB', image: '' },
  { id: 'tile_005', name: '木屋', category: 'building', color: '#8B4513', image: '' },
  { id: 'tile_006', name: '石屋', category: 'building', color: '#A9A9A9', image: '' },
  { id: 'tile_007', name: '围墙', category: 'building', color: '#696969', image: '' },
  { id: 'tile_008', name: '树木', category: 'decoration', color: '#228B22', image: '' },
  { id: 'tile_009', name: '花草', category: 'decoration', color: '#FF69B4', image: '' },
  { id: 'tile_010', name: '石头', category: 'decoration', color: '#A9A9A9', image: '' }
]);

const terrainTiles = computed(() => tiles.value.filter(t => t.category === 'terrain'));
const buildingTiles = computed(() => tiles.value.filter(t => t.category === 'building'));
const decorationTiles = computed(() => tiles.value.filter(t => t.category === 'decoration'));

const handleUpload = () => {
  console.log('Upload tile');
};

const handleSelect = (tile) => {
  console.log('Selected:', tile);
};
</script>

<style scoped>
.tileset-manager {
  max-width: 1200px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.tile-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
  gap: 16px;
  padding: 16px 0;
}

.tile-item {
  cursor: pointer;
  text-align: center;
  padding: 8px;
  border: 2px solid transparent;
  border-radius: 8px;
  transition: all 0.2s;
}

.tile-item:hover {
  border-color: #409eff;
  background-color: #ecf5ff;
}

.tile-preview {
  width: 80px;
  height: 80px;
  margin: 0 auto 8px;
  border-radius: 4px;
  overflow: hidden;
}

.tile-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.tile-name {
  font-size: 12px;
  color: #606266;
}
</style>

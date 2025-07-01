<template>
  <div class="blocks">
    <div class="info">
      {{ blocksData.length }} 个方块
    </div>
    <div class="blocks-container">
      <block class="block" v-for="block in sortedBlocks" :key="block.file_name" :block="block"></block>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue';
import { loadBlocksData, type BlockInfo } from './mcolor';
import Block from './Block.vue';

const props = defineProps<{
  currentColor: { r: number; g: number; b: number };
}>();

const blocksData = ref<BlockInfo[]>([]);
const sortedBlocks = ref<BlockInfo[]>([]);


onMounted(async () => {
  blocksData.value = await loadBlocksData();
  sortedBlocks.value = [...blocksData.value];
});

// 计算颜色相似度（欧几里得距离）
function calculateColorDistance(color1: { r: number; g: number; b: number },
  color2: { r: number; g: number; b: number }) {
  const dr = color1.r - color2.r;
  const dg = color1.g - color2.g;
  const db = color1.b - color2.b;
  return Math.sqrt(dr * dr + dg * dg + db * db);
}

watch(() => props.currentColor, () => {

  sortedBlocks.value = [...blocksData.value].sort((a, b) => {
    const distanceA = calculateColorDistance(
      { r: a.avg_r, g: a.avg_g, b: a.avg_b },
      props.currentColor
    );
    const distanceB = calculateColorDistance(
      { r: b.avg_r, g: b.avg_g, b: b.avg_b },
      props.currentColor
    );
    return distanceA - distanceB;
  });

}, { deep: true });
</script>

<style scoped>
.info{
  margin-bottom: 10px;
  text-align: center;
}
.blocks {
  width: 100%;
  height: 100%;
}

.blocks-container {
  padding-top: 10px;
  padding-bottom: 10px;
  width: 100%;
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: center;
  overflow-y: scroll;
  gap: 15px;
}
</style>

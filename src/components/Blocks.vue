<template>
  <div class="blocks">
    <div class="panel">
      <div class="filter-controls">
        <div class="filter-item">
          <span class="filter-label">材质总数: {{ filteredBlocks.length }}</span>
        </div>

        <div class="filter-item">
          <input id="filterFull" v-model="filterFull" class="filter-checkbox" type="checkbox" />
          <label for="filterFull" class="filter-label">只显示完整材质</label>
        </div>

        <div class="filter-item">
          <label for="typeFilter" class="filter-label">材质类型:</label>
          <select id="typeFilter" v-model="filterType" class="filter-select">
            <option value="all">全部类型</option>
            <option value="side">侧面材质</option>
            <option value="top">顶部材质</option>
            <option value="bottom">底部材质</option>
            <option value="null">其他</option>
          </select>
        </div>
      </div>
    </div>
    <div class="blocks-container">
      <block class="block" v-for="block in sortedBlocks" :key="block.file_name" :block="block"></block>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';
import { loadBlocksData, useThrottledUpdate, type BlockInfo, rgbToLab, type Lab } from './mcolor';
import Block from './Block.vue';

const props = defineProps<{
  currentColor: { r: number; g: number; b: number };
}>();

const blocksData = ref<BlockInfo[]>([]);
const sortedBlocks = ref<BlockInfo[]>([]);
const filterType = ref<string>('all'); // 添加筛选状态变量
const filterFull = ref<boolean>(false);

// 计算过滤后的方块列表
const filteredBlocks = ref<BlockInfo[]>([])


onMounted(async () => {
  blocksData.value = await loadBlocksData();
  filteredBlocks.value = [...blocksData.value];
  sortedBlocks.value = [...blocksData.value];
});

// CIEDE2000 算法实现
function deltaE2000(lab1: Lab,
  lab2: Lab): number {
  const kL = 1;
  const kC = 1;
  const kH = 1;

  const C1 = Math.sqrt(lab1.a * lab1.a + lab1.b * lab1.b);
  const C2 = Math.sqrt(lab2.a * lab2.a + lab2.b * lab2.b);

  const aC = (C1 + C2) / 2;
  const aC7 = Math.pow(aC, 7);
  const G = 0.5 * (1 - Math.sqrt(aC7 / (aC7 + 6103515625))); // 6103515625 = 25^7

  const a1p = (1 + G) * lab1.a;
  const a2p = (1 + G) * lab2.a;

  const C1p = Math.sqrt(a1p * a1p + lab1.b * lab1.b);
  const C2p = Math.sqrt(a2p * a2p + lab2.b * lab2.b);

  const h1p = Math.atan2(lab1.b, a1p) * (180 / Math.PI);
  const h2p = Math.atan2(lab2.b, a2p) * (180 / Math.PI);

  const dLp = lab2.L - lab1.L;
  const dCp = C2p - C1p;

  let dh = 0;
  if (C1p * C2p !== 0) {
    dh = h2p - h1p;
    if (dh > 180) dh -= 360;
    else if (dh < -180) dh += 360;
  }

  const dHp = 2 * Math.sqrt(C1p * C2p) * Math.sin(dh * Math.PI / 360);

  const aL = (lab1.L + lab2.L) / 2;
  const aCp = (C1p + C2p) / 2;

  let aHp = (h1p + h2p) / 2;
  if (Math.abs(h1p - h2p) > 180) {
    aHp += 180;
  }

  const T = 1
    - 0.17 * Math.cos((aHp - 30) * Math.PI / 180)
    + 0.24 * Math.cos(2 * aHp * Math.PI / 180)
    + 0.32 * Math.cos((3 * aHp + 6) * Math.PI / 180)
    - 0.20 * Math.cos((4 * aHp - 63) * Math.PI / 180);

  const dTheta = 30 * Math.exp(-Math.pow((aHp - 275) / 25, 2));
  const aCp7 = Math.pow(aCp, 7);
  const RC = 2 * Math.sqrt(aCp7 / (aCp7 + 6103515625));

  const SL = 1 + (0.015 * Math.pow(aL - 50, 2)) / Math.sqrt(20 + Math.pow(aL - 50, 2));
  const SC = 1 + 0.045 * aCp;
  const SH = 1 + 0.015 * aCp * T;

  const RT = -Math.sin(2 * dTheta * Math.PI / 180) * RC;

  return Math.sqrt(
    Math.pow(dLp / (kL * SL), 2) +
    Math.pow(dCp / (kC * SC), 2) +
    Math.pow(dHp / (kH * SH), 2) +
    RT * (dCp / (kC * SC)) * (dHp / (kH * SH))
  );
}

// 使用预计算的LAB值计算颜色距离
function calculateColorDistance(
  lab1: Lab,
  lab2: Lab
): number {
  return deltaE2000(lab1, lab2);
}

watch(() => [filterFull.value, filterType.value], () => {
  filteredBlocks.value = [...blocksData.value].filter((block) => {
    if (filterFull.value && !block.full) {
      return false;
    }
    if (filterType.value !== 'all' && block.type !== filterType.value) {
      return false;
    }
    return true;
  })
})

function updateSortedBlocksAndTitle() {
  // 计算当前颜色的LAB值（仅一次）
  const currentLab = rgbToLab(props.currentColor.r, props.currentColor.g, props.currentColor.b);

  sortedBlocks.value = [...filteredBlocks.value].sort((a, b) => {
    const distanceA = calculateColorDistance(
      a.lab,
      currentLab
    );
    const distanceB = calculateColorDistance(
      b.lab,
      currentLab
    );
    return distanceA - distanceB;
  });
  if (sortedBlocks.value.length > 0) {
    document.head.getElementsByTagName('link')[0].href = `/mcolor/${sortedBlocks.value[0].file_path}`
    document.title = `MColor | ${sortedBlocks.value[0].file_name}`
  }
}
const { triggerUpdate } = useThrottledUpdate(updateSortedBlocksAndTitle, 100);

watch(() => [props.currentColor, filteredBlocks.value], () => {
  triggerUpdate();
}, { deep: true });
</script>

<style scoped>
.panel {
  margin-top: 10px;
  margin-left: 5%;
  margin-right: 5%;
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

.filter-controls {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
}

.filter-controls {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: center;
  gap: 15px;
  padding: 10px;
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(10px);
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  margin-bottom: 15px;
}

.filter-item {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
}

.filter-label {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}

.filter-select {
  padding: 6px 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  background: white;
  font-size: 13px;
  transition: all 0.2s;
  min-width: 120px;
}

.filter-select:focus {
  outline: none;
  border-color: #4d90fe;
  box-shadow: 0 0 0 2px rgba(77, 144, 254, 0.2);
}

.filter-checkbox {
  width: 16px;
  height: 16px;
  accent-color: #4d90fe;
  cursor: pointer;
}

/* 响应式调整 */
@media (max-width: 768px) {
  .filter-controls {
    justify-self: center;
    align-items: center;
    gap: 10px;
  }
}
</style>

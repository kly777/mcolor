<template>
  <div class="blocks">
    <div class="panel">
      <div class="filter-controls">
        <span>共 {{ filteredBlocks.length }} 个方块材质</span>

        <input v-model="filterFull" class="filter-select" type="checkbox" />
        <span for="filterFull">只显示完整材质</span>

        <select v-model="filterType" class="filter-select">
          <option value="all">全部类型</option>
          <option value="side">侧面材质</option>
          <option value="top">顶部材质</option>
          <option value="bottom">底部材质</option>
          <option value="null">其他</option>
        </select>
      </div>
    </div>
    <div class="blocks-container">
      <block class="block" v-for="block in sortedBlocks" :key="block.file_name" :block="block"></block>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';
import { loadBlocksData, type BlockInfo } from './mcolor';
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

// 计算颜色相似度（欧几里得距离）
// function calculateColorDistance(
//   color1: { r: number; g: number; b: number },
//   color2: { r: number; g: number; b: number }
// ) {
//   const dr = color1.r - color2.r;
//   const dg = color1.g - color2.g;
//   const db = color1.b - color2.b;
//   return Math.sqrt(dr * dr + dg * dg + db * db);
// }
// 在 Blocks.vue 中添加以下函数
function rgbToLab(r: number, g: number, b: number): { L: number; a: number; b: number } {
  // RGB 转 XYZ
  let r_ = r / 255;
  let g_ = g / 255;
  let b_ = b / 255;

  r_ = r_ > 0.04045 ? Math.pow((r_ + 0.055) / 1.055, 2.4) : r_ / 12.92;
  g_ = g_ > 0.04045 ? Math.pow((g_ + 0.055) / 1.055, 2.4) : g_ / 12.92;
  b_ = b_ > 0.04045 ? Math.pow((b_ + 0.055) / 1.055, 2.4) : b_ / 12.92;

  r_ *= 100;
  g_ *= 100;
  b_ *= 100;

  const X = r_ * 0.4124564 + g_ * 0.3575761 + b_ * 0.1804375;
  const Y = r_ * 0.2126729 + g_ * 0.7151522 + b_ * 0.0721750;
  const Z = r_ * 0.0193339 + g_ * 0.1191920 + b_ * 0.9503041;

  // XYZ 转 Lab
  const Xn = 95.047;
  const Yn = 100.000;
  const Zn = 108.883;

  let x = X / Xn;
  let y = Y / Yn;
  let z = Z / Zn;

  x = x > 0.008856 ? Math.pow(x, 1 / 3) : (7.787 * x) + 16 / 116;
  y = y > 0.008856 ? Math.pow(y, 1 / 3) : (7.787 * y) + 16 / 116;
  z = z > 0.008856 ? Math.pow(z, 1 / 3) : (7.787 * z) + 16 / 116;

  const L = (116 * y) - 16;
  const a = 500 * (x - y);
  const b__ = 200 * (y - z);

  return { L, a, b: b__ };
}

// CIEDE2000 算法实现
function deltaE2000(lab1: { L: number; a: number; b: number },
  lab2: { L: number; a: number; b: number }): number {
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

// 修改后的颜色距离计算函数
function calculateColorDistance(
  color1: { r: number; g: number; b: number },
  color2: { r: number; g: number; b: number }
): number {
  const lab1 = rgbToLab(color1.r, color1.g, color1.b);
  const lab2 = rgbToLab(color2.r, color2.g, color2.b);
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

watch(() => [props.currentColor, filteredBlocks.value], () => {

  sortedBlocks.value = [...filteredBlocks.value].sort((a, b) => {
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
  if (sortedBlocks.value.length > 0) {
    document.head.getElementsByTagName('link')[0].href = `/mcolor/${sortedBlocks.value[0].file_path}`
    document.title = `MColor | ${sortedBlocks.value[0].file_name}`
  }

}, { deep: true });


</script>

<style scoped>
.panel {
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
</style>

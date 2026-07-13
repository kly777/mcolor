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
    <div ref="scrollContainerRef" class="blocks-container" @scroll="onScroll">
      <div :style="spacerStyle">
        <block
          v-for="item in visibleItems"
          :key="item.block.file_name"
          :block="item.block"
          :style="item.style"
          class="block"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch, onBeforeUnmount } from 'vue';
import { loadBlocksData, useThrottledUpdate, type BlockInfo, rgbToLab, type Lab } from './mcolor';
import Block from './Block.vue';

const BLOCK_SIZE = 120;
const GAP = 15;
const CELL = BLOCK_SIZE + GAP; // 135

const props = defineProps<{
  currentColor: { r: number; g: number; b: number };
  dragging: boolean;
}>();

const blocksData = ref<BlockInfo[]>([]);
const sortedBlocks = ref<BlockInfo[]>([]);
const filterType = ref<string>('all');
const filterFull = ref<boolean>(false);

const filteredBlocks = ref<BlockInfo[]>([]);

// 虚拟滚动状态
const scrollContainerRef = ref<HTMLElement | null>(null);
const scrollTop = ref(0);
const containerWidth = ref(800);
let resizeObserver: ResizeObserver | null = null;

onMounted(async () => {
  blocksData.value = await loadBlocksData();
  filteredBlocks.value = [...blocksData.value];
  sortedBlocks.value = [...blocksData.value];

  // 监听容器宽度变化
  if (scrollContainerRef.value) {
    containerWidth.value = scrollContainerRef.value.clientWidth;
    resizeObserver = new ResizeObserver((entries) => {
      for (const entry of entries) {
        containerWidth.value = entry.contentBoxSize?.[0]?.inlineSize ?? entry.contentRect.width;
      }
    });
    resizeObserver.observe(scrollContainerRef.value);
  }
});

onBeforeUnmount(() => {
  resizeObserver?.disconnect();
});

const columns = computed(() => Math.max(1, Math.floor((containerWidth.value + GAP) / CELL)));

const totalRows = computed(() => Math.ceil(sortedBlocks.value.length / columns.value));

// 撑开滚动高度的 spacer
const spacerStyle = computed(() => ({
  height: `${Math.max(1, totalRows.value * CELL)}px`,
  position: 'relative' as const,
}));

// 水平居中偏移
const leftOffset = computed(() => {
  const totalGridWidth = columns.value * CELL - GAP;
  return Math.max(0, Math.floor((containerWidth.value - totalGridWidth) / 2));
});

// 可视范围
const visibleItems = computed(() => {
  const col = columns.value;
  if (col === 0 || sortedBlocks.value.length === 0) return [];

  const rowHeight = CELL;
  const containerHeight = scrollContainerRef.value?.clientHeight ?? 800;
  const startRow = Math.max(0, Math.floor(scrollTop.value / rowHeight));
  const visibleRows = Math.ceil(containerHeight / rowHeight) + 3; // 上下各多 1.5 行缓冲
  const endRow = Math.min(totalRows.value, startRow + visibleRows);

  const start = startRow * col;
  const end = Math.min(sortedBlocks.value.length, endRow * col);

  const items = [];
  for (let i = start; i < end; i++) {
    const row = Math.floor(i / col);
    const column = i % col;
    items.push({
      block: sortedBlocks.value[i],
      style: {
        position: 'absolute' as const,
        left: `${column * CELL + leftOffset.value}px`,
        top: `${row * CELL}px`,
        width: `${BLOCK_SIZE}px`,
        height: `${BLOCK_SIZE}px`,
      },
    });
  }
  return items;
});

function onScroll(e: Event) {
  scrollTop.value = (e.target as HTMLElement).scrollTop;
}

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
  const G = 0.5 * (1 - Math.sqrt(aC7 / (aC7 + 6103515625)));

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

// RGB 空间欧氏距离（纯加减乘除，零 trig，用于拖拽时实时排序）
function rgbEuclideanDistance(
  rgb1: { r: number; g: number; b: number },
  rgb2: { r: number; g: number; b: number }
): number {
  const dr = rgb1.r - rgb2.r;
  const dg = rgb1.g - rgb2.g;
  const db = rgb1.b - rgb2.b;
  return dr * dr + dg * dg + db * db;
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
  // 拖拽时用 RGB 欧氏距离（极轻量），否则用 CIEDE2000（精确）
  const useAccurate = !props.dragging;

  const withDistances = filteredBlocks.value.map(block => ({
    block,
    distance: useAccurate
      ? calculateColorDistance(block.lab, rgbToLab(props.currentColor.r, props.currentColor.g, props.currentColor.b))
      : rgbEuclideanDistance(block.rgb, props.currentColor)
  }));
  withDistances.sort((a, b) => a.distance - b.distance);
  sortedBlocks.value = withDistances.map(item => item.block);

  if (useAccurate && sortedBlocks.value.length > 0) {
    document.head.getElementsByTagName('link')[0].href = `/mcolor/${sortedBlocks.value[0].file_path}`
    document.title = `MColor | ${sortedBlocks.value[0].file_name}`
  }
}
const { triggerUpdate } = useThrottledUpdate(updateSortedBlocksAndTitle, 100);

watch(() => [props.currentColor, filteredBlocks.value], () => {
  triggerUpdate();
}, { deep: true });

// 拖拽结束时执行一次最终排序（绕过节流）
watch(() => props.dragging, (newVal, oldVal) => {
  if (oldVal === true && newVal === false) {
    updateSortedBlocksAndTitle();
  }
});
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
  display: flex;
  flex-direction: column;
}

.blocks-container {
  flex: 1;
  overflow-y: auto;
  position: relative;
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

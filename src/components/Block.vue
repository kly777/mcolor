<template>
  <div class="block" :style="{ backgroundColor: `rgb(${block.avg_r}, ${block.avg_g}, ${block.avg_b})` }"
    :title="block.file_name">
    <!-- 可选：显示颜色值 -->
    <div class="image-container">
      <img :src="block.file_path" :alt="block.file_name" class="block-image" />
    </div>
    <div class="hex-value">
      {{ rgbToHex(block.avg_r, block.avg_g, block.avg_b) }}
    </div>
  </div>
</template>

<script lang="ts" setup>
import { defineProps, computed } from 'vue';
import type { BlockInfo } from './mcolor';

const props = defineProps<{
  block: BlockInfo;
}>();

// 可选：是否显示十六进制颜色值


// RGB转十六进制工具函数
const rgbToHex = (r: number, g: number, b: number) => {
  return `#${[r, g, b]
    .map(x => Math.round(x).toString(16).padStart(2, '0'))
    .join('')}`;
};
</script>

<style scoped>
.block {
  width: 120px;
  height: 120px;
  position: relative;
  border: 1px solid rgba(0, 0, 0, 0.1);
  border-radius: 8px;
  box-sizing: border-box;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s, box-shadow 0.2s;

}

.block:hover {
  transform: translateY(-5px);
}

.image-container {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 10px;
  box-sizing: border-box;

}

.hex-value {
  position: absolute;
  bottom: 0;
  left: 0;
  font-size: 11px;
  font-weight: 500;
  white-space: nowrap;

  padding: 4px 8px;
  background-color: rgba(255, 255, 255, 0.8);
  border-radius: 0 0 0 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  z-index: 10;
}

.block-image {
  width: 50%;
  height: 50%;
  display: block;
  object-fit: cover;
  border: 1px solid rgba(118, 118, 118, 0.1);
}
</style>
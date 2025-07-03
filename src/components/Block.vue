<template>
  <div class="block" :style="{ backgroundColor: `rgb(${block.rgb.r}, ${block.rgb.g}, ${block.rgb.b})` }"
    :title="block.file_name">
    <div class="block-name" @click="copyFileName">
      {{ block.file_name }}
    </div>
    <div class="image-container">
      <img :src="block.file_path" :alt="block.file_name" class="block-image" />
    </div>
    <div class="hex-value">
      {{ rgbToHex(block.rgb) }}
    </div>
    <div v-if="block.type !== 'null'" class="block-type">
      {{ block.type }}
    </div>
    <div v-if="showCopiedHint" class="copied-hint">已复制</div>
  </div>
</template>

<script lang="ts" setup>
import {  computed, ref } from 'vue';
import type { BlockInfo, RGB } from './mcolor';

const props = defineProps<{
  block: BlockInfo;
}>();

// RGB转十六进制工具函数
const rgbToHex = (rgb: RGB) => {
  return `#${[rgb.r, rgb.g, rgb.b]
    .map(x => Math.round(x).toString(16).padStart(2, '0'))
    .join('')}`;
};
const showCopiedHint = ref(false) // 添加提示状态
function copyFileName() {
  navigator.clipboard.writeText(props.block.file_name)
    .then(() => {
      showCopiedHint.value = true // 显示提示
      setTimeout(() => {
        showCopiedHint.value = false // 2秒后隐藏提示
      }, 2000)
    })
    .catch(err => {
      console.error('复制失败:', err);
    });
}
</script>

<style scoped>
.copied-hint {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background-color: rgba(0, 0, 0, 0.7);
  color: white;
  padding: 5px 10px;
  border-radius: 4px;
  z-index: 20;
  font-size: 12px;
  animation: fadeOut 2s forwards;
}

@keyframes fadeOut {
  0% {
    opacity: 1;
  }

  70% {
    opacity: 1;
  }

  100% {
    opacity: 0;
  }
}

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

.block-name {
  box-sizing: border-box;
  position: absolute;
  top: 0;
  left: 0;
  font-size: 11px;
  font-weight: 500;

  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;

  padding: 4px 8px;
  background-color: rgba(255, 255, 255, 0.8);
  border-radius: 8px 0 0 0;
  max-width: 85%;
}

.block-name:hover {
  text-overflow: unset;
  white-space: unset;
  text-decoration: underline;
  cursor: pointer;
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

.block-type {
  position: absolute;
  bottom: 0;
  right: 0;
  font-size: 11px;
  font-weight: 500;
  white-space: nowrap;
  padding: 4px 8px;
  background-color: rgba(255, 255, 255, 0.8);
  border-radius: 0 0 8px 0;
}
</style>
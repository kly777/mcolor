<template>
  <div class="color-picker">
    <!-- 添加外层容器保持宽高比 -->
    <div class="picker-container">
      <div class="color-controls">
        <div ref="colorPane" class="color-pane" @mousedown="startPaneDrag" @touchstart="startPaneDrag">
          <div class="pointer" :style="pointerStyle"></div>
        </div>
        <div ref="valueSlider" class="value-slider" @mousedown="startValueDrag" @touchstart="startValueDrag">
          <div class="value-pointer" :style="valuePointerStyle"></div>
        </div>
      </div>
    </div>

    <div class="color-output">
      <div>HSV: {{ hsv.h }}°, {{ hsv.s }}%, {{ hsv.v }}%</div>
      <div>RGB: {{ rgb.r }}, {{ rgb.g }}, {{ rgb.b }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">

import { ref, computed, onMounted, watch, nextTick } from 'vue';



const props = defineProps({
  initRgbColor: {
    type: Object,
    required: true,
    default: () => ({ r: 255, g: 0, b: 0 }) // 默认RGB值（红色）
  }
});

defineExpose({
  update
});
function update(newRgb: { r: number; g: number; b: number }) {
  const hsvColor = rgbToHsv(newRgb.r, newRgb.g, newRgb.b);
  position.value.x = hsvColor.h / 360;
  position.value.y = 1 - (hsvColor.s / 100);
  valuePosition.value = hsvColor.v / 100;
}


// 响应式数据
const colorPane = ref<HTMLElement | null>(null);
const valueSlider = ref<HTMLElement | null>(null);
const paneDragging = ref(false);
const valueDragging = ref(false);
const position = ref({ x: 0.5, y: 0.5 });
const valuePosition = ref(1); // 初始为100%（顶部）

// HSV颜色值（色相H:0-360, 饱和度S:0-100, 明度V:0-100）
const hsv = computed(() => ({
  h: Math.min(Math.round(position.value.x * 360), 359),
  s: Math.round((1 - position.value.y) * 100),
  v: Math.round(valuePosition.value * 100) // 由滑块控制
}));

// 从props初始化
onMounted(() => {
  // 将RGB转换为HSV用于内部位置计算
  const { r, g, b } = props.initRgbColor;
  const hsvColor = rgbToHsv(r, g, b);

  position.value.x = hsvColor.h / 360;
  position.value.y = 1 - (hsvColor.s / 100);
  valuePosition.value = hsvColor.v / 100;
});
// RGB转HSV函数
function rgbToHsv(r: number, g: number, b: number) {
  r /= 255;
  g /= 255;
  b /= 255;

  const max = Math.max(r, g, b);
  const min = Math.min(r, g, b);
  let h = 0, s = 0, v = max;

  const d = max - min;
  s = max === 0 ? 0 : d / max;

  if (max === min) {
    h = 0; // 无色相
  } else {
    switch (max) {
      case r: h = (g - b) / d + (g < b ? 6 : 0); break;
      case g: h = (b - r) / d + 2; break;
      case b: h = (r - g) / d + 4; break;
    }
    h /= 6;
  }

  return {
    h: Math.round(h * 360),
    s: Math.round(s * 100),
    v: Math.round(v * 100)
  };
}

// RGB颜色值
// RGB颜色值
const rgb = computed(() => {
  // 处理色相边界：360° 等同于 0°
  const normalizedH = hsv.value.h % 360;
  const h = normalizedH / 360;

  const s = hsv.value.s / 100;
  const v = hsv.value.v / 100; // 明度值（0-1）

  let r = 0, g = 0, b = 0;
  if (s === 0) {
    r = g = b = v * 255; // 灰度，乘以明度
  } else {
    const hue = h * 6;
    const i = Math.floor(hue);
    const f = hue - i;
    const p = v * 255 * (1 - s);
    const q = v * 255 * (1 - s * f);
    const t = v * 255 * (1 - s * (1 - f));

    switch (i) {
      case 0: r = v * 255; g = t; b = p; break;
      case 1: r = q; g = v * 255; b = p; break;
      case 2: r = p; g = v * 255; b = t; break;
      case 3: r = p; g = q; b = v * 255; break;
      case 4: r = t; g = p; b = v * 255; break;
      case 5: r = v * 255; g = p; b = q; break; // 修改了 default 为 case 5
      default: r = v * 255; g = p; b = q; // 保留 default 以防万一
    }
  }

  return {
    r: Math.round(r),
    g: Math.round(g),
    b: Math.round(b)
  };
});

const emit = defineEmits(['update:rgbColor']);

// watch(() => props.rgbColor, (newValue) => {
//   // 将RGB转换为HSV用于内部位置计算
//   const { r, g, b } = newValue;
//   const hsvColor = rgbToHsv(r, g, b);

//   position.value.x = hsvColor.h / 360;
//   position.value.y = 1 - (hsvColor.s / 100);
//   valuePosition.value = hsvColor.v / 100;
// }, { immediate: true });
watch([position, valuePosition], () => {
  // 防抖处理，避免频繁触发
  nextTick(() => {
    emit('update:rgbColor', rgb.value);
  });
}, { deep: true });
// 指针位置样式
const pointerStyle = computed(() => ({
  left: `${position.value.x * 100}%`,
  top: `${position.value.y * 100}%`,
  backgroundColor: `rgb(${rgb.value.r}, ${rgb.value.g}, ${rgb.value.b})`
}));
// 明度滑块指针样式
const valuePointerStyle = computed(() => ({
  top: `${valuePosition.value * 100}%`,
  backgroundColor: `rgb(${Math.round(rgb.value.r * valuePosition.value)}, ${Math.round(rgb.value.g * valuePosition.value)}, ${Math.round(rgb.value.b * valuePosition.value)})`
}));
// 开始拖动颜色面板
const startPaneDrag = (e: MouseEvent | TouchEvent) => {
  paneDragging.value = true;
  updatePanePosition(e);
  window.addEventListener('mousemove', updatePanePosition);
  window.addEventListener('touchmove', updatePanePosition);
  window.addEventListener('mouseup', stopPaneDrag);
  window.addEventListener('touchend', stopPaneDrag);
};

// 更新颜色面板指针位置
const updatePanePosition = (e: MouseEvent | TouchEvent) => {
  if (!paneDragging.value || !colorPane.value) return;

  const rect = colorPane.value.getBoundingClientRect();
  let clientX, clientY;

  if (e instanceof MouseEvent) {
    clientX = e.clientX;
    clientY = e.clientY;
  } else if (e instanceof TouchEvent && e.touches.length > 0) {
    clientX = e.touches[0].clientX;
    clientY = e.touches[0].clientY;
  } else {
    return;
  }

  // 计算相对位置（0-1）
  position.value.x = Math.max(0, Math.min(1, (clientX - rect.left) / rect.width));
  position.value.y = Math.max(0, Math.min(1, (clientY - rect.top) / rect.height));
};

// 停止拖动
const stopPaneDrag = () => {
  paneDragging.value = false;
  window.removeEventListener('mousemove', updatePanePosition);
  window.removeEventListener('touchmove', updatePanePosition);
  window.removeEventListener('mouseup', stopPaneDrag);
  window.removeEventListener('touchend', stopPaneDrag);
};


// 开始拖动明度滑块
const startValueDrag = (e: MouseEvent | TouchEvent) => {
  valueDragging.value = true;
  updateValuePosition(e);
  window.addEventListener('mousemove', updateValuePosition);
  window.addEventListener('touchmove', updateValuePosition);
  window.addEventListener('mouseup', stopValueDrag);
  window.addEventListener('touchend', stopValueDrag);
};

// 更新明度滑块位置
const updateValuePosition = (e: MouseEvent | TouchEvent) => {
  if (!valueDragging.value || !valueSlider.value) return;

  const rect = valueSlider.value.getBoundingClientRect();
  let clientY;

  if (e instanceof MouseEvent) {
    clientY = e.clientY;
  } else if (e instanceof TouchEvent && e.touches.length > 0) {
    clientY = e.touches[0].clientY;
  } else {
    return;
  }

  // 计算相对位置（0-1），顶部为100%（1），底部为0%（0）
  valuePosition.value = Math.max(0, Math.min(1, (clientY - rect.top) / rect.height));
};

// 停止拖动明度滑块
const stopValueDrag = () => {
  valueDragging.value = false;
  window.removeEventListener('mousemove', updateValuePosition);
  window.removeEventListener('touchmove', updateValuePosition);
  window.removeEventListener('mouseup', stopValueDrag);
  window.removeEventListener('touchend', stopValueDrag);
};

// 设置色板背景
// onMounted(() => {
//   if (colorPane.value) {
//     colorPane.value.style.background = `
//       linear-gradient(to top, #000, rgba(0,0,0,0)),
//     `;
//   }
// });
</script>

<style scoped>
.picker-container {
  aspect-ratio: 1.5;
  /* 宽高比为 3:2 */
  max-height: 400px;
  width: 100%;
}

.color-picker {
  display: flex;
  flex-direction: column;
  width: 100%;
  gap: 20px;
}

.color-controls {
  display: flex;
  width: 100%;
  height: 100%;
  gap: 15px;
}

.color-pane {
  box-sizing: border-box;
  position: relative;
  width: 100%;
  height: 100%;
  background:
    linear-gradient(to bottom, transparent, rgba(255, 255, 255, 1) 100%),
    linear-gradient(to right, rgb(255, 0, 0) 0%, rgb(255, 255, 0) 17%, rgb(0, 255, 0) 34%, rgb(0, 255, 255) 50%, rgb(0, 0, 255) 67%, rgb(255, 0, 255) 84%, rgb(255, 0, 0) 100%);
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  touch-action: none;
  border: #000 2px solid;
}

.value-slider {
  box-sizing: border-box;
  position: relative;
  width: 25px;
  height: 100%;
  margin-left: 10px;
  background: linear-gradient(to top,
      #fff 0%,
      rgba(255, 255, 255, 0) 50%,
      #000 100%);
  border-radius: 4px;
  overflow: hidden;
  cursor: pointer;
  touch-action: none;
  border: #000 2px solid;
}

.pointer {
  position: absolute;
  width: 20px;
  height: 20px;
  border: 2px solid white;
  border-radius: 50%;
  transform: translate(-50%, -50%);
  box-shadow: 0 0 3px rgba(0, 0, 0, 0.5);
  pointer-events: none;
}

.value-pointer {
  box-sizing: border-box;
  position: absolute;
  left: 0;
  width: 100%;
  height: min(5%, 10px);
  border: 2px solid white;
  border-radius: 2px;
  transform: translateY(-50%);
  box-shadow: 0 0 3px rgba(0, 0, 0, 0.5);
  pointer-events: none;
}

.color-output {
  margin-top: 15px;
  font-family: monospace;
  font-size: 14px;
  line-height: 1.5;
}
</style>
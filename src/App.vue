<template>
  <div class="main" :style="mainStyle">
    <color-picker ref="colorPicker" v-model:rgbColor="rgbColor" :initRgbColor="initRgbColor" />
    <div>当前颜色: RGB({{ rgbColor.r }}, {{ rgbColor.g }}, {{ rgbColor.b }})</div>
    <input type="text" v-model="rgbColor.r" @input="updateColor" />
    <input type="text" v-model="rgbColor.g" @input="updateColor" />
    <input type="text" v-model="rgbColor.b" @input="updateColor" />
  </div>
</template>


<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import ColorPicker from './components/ColorPicker.vue';
const colorPicker = ref<InstanceType<typeof ColorPicker> | null>(null);

const initRgbColor = { r: 0, g: 70, b: 140 };
const rgbColor = ref({ r: 180, g: 100, b: 100 });

const mainStyle = computed(() => ({
  backgroundColor: `rgb(${Number(rgbColor.value.r)}, ${Number(rgbColor.value.g)}, ${Number(rgbColor.value.b)})`
}))

function updateColor() {
  if (rgbColor.value.r > 255) {
    rgbColor.value.r = 255;
  } else if (rgbColor.value.r < 0) {
    rgbColor.value.r = 0;
  }
  if (rgbColor.value.g > 255) {
    rgbColor.value.g = 255;
  } else if (rgbColor.value.g < 0) {
    rgbColor.value.g = 0;
  }
  if (rgbColor.value.b > 255) {
    rgbColor.value.b = 255;
  } else if (rgbColor.value.b < 0) {
    rgbColor.value.b = 0;
  }
  colorPicker.value?.update(rgbColor.value);
}
</script>
<style scoped>
.main {
  height: 100vh;
  width: 100vw;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

button {
  margin: 10px;
  padding: 8px 16px;
  background-color: #f0f0f0;
  border: 1px solid #ccc;
  border-radius: 4px;
  cursor: pointer;
}
</style>

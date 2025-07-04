<template>
  <div class="main" :style="mainStyle" @click.self="randomColor">
    <ToGithub to="https://github.com/kly777/mcolor" />
    <div class="color-panel" @click.self="randomColor">
      <div class="color-controls">
        <color-picker ref="colorPicker" v-model:rgbColor="rgbColor" :initRgbColor="initRgbColor" />

        <div class="rgb-display">
          <div class="rgb-inputs">
            <div class="input-group">
              <label>R:</label>
              <input type="text" v-model.number="rgbColor.r" @change="updateColor" min="0" max="255" />
            </div>
            <div class="input-group">
              <label>G:</label>
              <input type="text" v-model.number="rgbColor.g" @change="updateColor" min="0" max="255" />
            </div>
            <div class="input-group">
              <label>B:</label>
              <input type="text" v-model.number="rgbColor.b" @change="updateColor" min="0" max="255" />
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="blocks-panel">
      <blocks :current-color="reasonableColor" />
    </div>
  </div>
</template>


<script setup lang="ts">

import { computed, ref, watch } from 'vue';
import ToGithub from './components/ToGithub.vue';
import ColorPicker from './components/ColorPicker.vue';
import Blocks from './components/Blocks.vue';
const colorPicker = ref<InstanceType<typeof ColorPicker> | null>(null);


const initRgbColor = (() => {
  return {
    r: Math.floor(Math.random() * 256),
    g: Math.floor(Math.random() * 256),
    b: Math.floor(Math.random() * 256)
  };
})();
const rgbColor = ref({ r: 180, g: 100, b: 100 });
const reasonableColor = computed(() => {
  return {
    r: Math.max(0, Math.min(rgbColor.value.r, 255)),
    g: Math.max(0, Math.min(rgbColor.value.g, 255)),
    b: Math.max(0, Math.min(rgbColor.value.b, 255))
  }
})

const mainStyle = computed(() => ({
  backgroundColor: `rgb(${Number(rgbColor.value.r)}, ${Number(rgbColor.value.g)}, ${Number(rgbColor.value.b)})`
}))

function updateColor() {
  colorPicker.value?.update(reasonableColor.value);
}
function randomColor() {
  rgbColor.value = {
    r: Math.floor(Math.random() * 256),
    g: Math.floor(Math.random() * 256),
    b: Math.floor(Math.random() * 256)
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
  flex-direction: row;
  align-items: stretch;
  padding: 20px;
  gap: 20px;
  background-color: #f5f5f5;
  transition: background 0.4s ease-in-out;
}

.color-panel {
  flex: 1;
  min-width: 450px;
  max-width: 100%;
  max-height: 100%;

  border-radius: 12px;
  padding: 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.blocks-panel {
  padding: 20px;
  flex: 1;
  background: white;
  border-radius: 12px;

  overflow: scroll;
  display: flex;
  flex-direction: column;

  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(10px);
  border-radius: 12px;
  padding: 15px;
  box-shadow: 0 4px 30px rgba(0, 0, 0, 0.1);
}

.color-controls {
  display: flex;
  flex-direction: column;
  gap: 20px;
  width: 50%;
  height: auto;

  background: rgba(255, 255, 255, 0.3);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(10px);
  border-radius: 12px;
  padding: 15px;
  box-shadow: 0 4px 30px rgba(0, 0, 0, 0.1);
}

@media (max-width: 768px) {
  .color-controls {
    width: 90%;
    height: auto;
  }

  .main {
    flex-direction: column;
  }
}

.rgb-display {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.current-color {
  font-size: 16px;
  font-weight: 500;
  color: #333;
  text-align: center;
  padding: 8px;
  background: #f8f9fa;
  border-radius: 6px;
}

.rgb-inputs {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.input-group {
  display: flex;
  align-items: center;
  gap: 10px;
}

.input-group label {
  width: 30px;
  font-weight: 500;
  color: #555;
}

.input-group input {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
  transition: border-color 0.2s;
  width: 100%;
  min-width: none;
}

.input-group input:focus {
  outline: none;
  border-color: #4d90fe;
  box-shadow: 0 0 0 2px rgba(77, 144, 254, 0.2);
}
</style>

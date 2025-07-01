<template>
  <div>
    {{ c }}:
  </div>
  <input type="text" v-model="value" />
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue'

const props = defineProps<
  {
    c: string;
    value: number
  }>()
</script>
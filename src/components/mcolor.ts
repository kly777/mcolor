import { onUnmounted, ref } from "vue";

export async function loadBlocksData() {
  var blocks = [];
  try {
    const response = await fetch('/mcolor/data.json'); // 绝对路径
    console.log(response);
    const jsonData = await response.json();
    console.log(jsonData); // 访问具体字段
    blocks = jsonData;
  } catch (error) {
    console.error('加载JSON失败', error);
  }

  return blocks as BlockInfo[];
};

export type BlockInfo = {
  file_path: string;
  file_name: string,
  type: "null" | "side" | "top" | "bottom",
  full: boolean,
  avg_r: number,
  avg_g: number,
  avg_b: number,
  var_r: number,
  var_g: number,
  var_b: number
}

export function useThrottledUpdate(updateFn: () => void, interval = 100) {
  const throttleTimer = ref<number | null>(null);
  const lastUpdateTime = ref<number>(0);
  const pendingUpdate = ref(false);

  const executeUpdate = () => {
    updateFn();
    lastUpdateTime.value = Date.now();
    pendingUpdate.value = false;
  };

  const triggerUpdate = () => {
    const now = Date.now();
    const elapsed = now - lastUpdateTime.value;

    if (elapsed >= interval) {
      if (throttleTimer.value) {
        clearTimeout(throttleTimer.value);
        throttleTimer.value = null;
      }
      executeUpdate();
    } else {
      pendingUpdate.value = true;
      if (!throttleTimer.value) {
        throttleTimer.value = setTimeout(() => {
          if (pendingUpdate.value) {
            executeUpdate();
          }
          throttleTimer.value = null;
        }, interval - elapsed);
      }
    }
  };

  onUnmounted(() => {
    if (throttleTimer.value) {
      clearTimeout(throttleTimer.value);
    }
  });

  return {
    triggerUpdate
  };
}
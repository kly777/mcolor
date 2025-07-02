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
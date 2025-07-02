package main

import (
	"encoding/json"
	"image"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// 结果数据结构
type ImageStats struct {
	FilePath string  `json:"file_path"`
	FileName string  `json:"file_name"`
	Type     string  `json:"type"` // top/button/side/null
	Full     bool    `json:"full"`
	AvgR     float64 `json:"avg_r"` // 红色通道平均值
	AvgG     float64 `json:"avg_g"` // 绿色通道平均值
	AvgB     float64 `json:"avg_b"` // 蓝色通道平均值
	VarR     float64 `json:"var_r"` // 红色通道方差
	VarG     float64 `json:"var_g"` // 绿色通道方差
	VarB     float64 `json:"var_b"` // 蓝色通道方差
}

func main() {
	files, err := filepath.Glob("block/*.png")
	if err != nil {
		log.Fatal("文件遍历失败:", err)
	}

	results := []ImageStats{}

	for _, file := range files {
		// 2. 打开并解码PNG图片
		f, err := os.Open(file)
		if err != nil {
			log.Printf("跳过无法打开的文件: %s, 错误: %v", file, err)
			continue
		}
		defer f.Close()

		img, _, err := image.Decode(f)
		if err != nil {
			log.Printf("跳过解码失败的文件: %s, 错误: %v", file, err)
			continue
		}

		// 3. 确定图片类型
		imgType := "null"
		baseName := strings.TrimSuffix(file, ".png")
		switch {
		case strings.HasSuffix(baseName, "_top"):
			imgType = "top"
		case strings.HasSuffix(baseName, "_bottom"):
			imgType = "bottom"
		case strings.HasSuffix(baseName, "_side"):
			imgType = "side"
		}

		//  计算统计信息

		stats := calculateImageStats(img, file, imgType)
		results = append(results, stats)
	}

	// 5. 生成JSON输出
	outputJSON(results)
}
func calculateImageStats(img image.Image, file, imgType string) ImageStats {
	filename := strings.Split(filepath.Base(file), ".")[0]
	filename = strings.ReplaceAll(filename, "_", " ")
	bounds := img.Bounds()

	// 初始化累加器和有效像素计数器
	var sumR, sumG, sumB float64
	var sumSqR, sumSqG, sumSqB float64
	totalPixels := bounds.Dx() * bounds.Dy()
	validPixels := 0 // 记录非完全透明像素的数量

	// 遍历所有像素
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			// 转换到0-255范围
			a8 := float64(a >> 8)
			if a8 == 0 { // 完全透明的像素跳过计算
				continue
			}

			r8 := float64(r >> 8)
			g8 := float64(g >> 8)
			b8 := float64(b >> 8)

			sumR += r8
			sumG += g8
			sumB += b8

			sumSqR += r8 * r8
			sumSqG += g8 * g8
			sumSqB += b8 * b8

			validPixels++ // 增加有效像素计数
		}
	}
	if validPixels == 0 {
		// 返回空值或默认值
		return ImageStats{FilePath: file, FileName: filename, Type: imgType}
	}

	// 计算平均值（使用有效像素数）
	avgR := sumR / float64(validPixels)
	avgG := sumG / float64(validPixels)
	avgB := sumB / float64(validPixels)

	// 计算方差（使用有效像素数）
	varR := (sumSqR/float64(validPixels) - avgR*avgR)
	varG := (sumSqG/float64(validPixels) - avgG*avgG)
	varB := (sumSqB/float64(validPixels) - avgB*avgB)

	return ImageStats{
		FilePath: file,
		FileName: filename,
		Type:     imgType,
		Full:     totalPixels == validPixels,
		AvgR:     avgR,
		AvgG:     avgG,
		AvgB:     avgB,
		VarR:     varR,
		VarG:     varG,
		VarB:     varB,
	}
}
func outputJSON(results []ImageStats) {
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Fatal("JSON编码失败:", err)
	}

	// 写入结果文件
	if err := os.WriteFile("results.json", jsonData, 0644); err != nil {
		log.Fatal("文件写入失败:", err)
	}
	log.Println("结果已保存到 results.json")
}

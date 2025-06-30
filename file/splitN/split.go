package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	srcPath := flag.String("src", "", "源ZIP文件路径")
	count := flag.Int("n", 2, "切分文件数量")
	flag.Parse()

	if *srcPath == "" {
		fmt.Println("请指定源文件路径 -src=xxx.zip")
		return
	}

	// 读取整个ZIP文件
	zipData, err := os.ReadFile(*srcPath)
	if err != nil {
		panic(err)
	}

	// 计算每个分片大小
	totalSize := int64(len(zipData))
	partSize := totalSize / int64(*count)
	remainder := totalSize % int64(*count)

	// 创建输出目录
	outDir := filepath.Join(filepath.Dir(*srcPath), "part")
	os.MkdirAll(outDir, 0755)

	// 写入分片文件
	for i := 0; i < *count; i++ {
		start := int64(i) * partSize
		end := start + partSize
		if i == *count-1 {
			end += remainder
		}
		partData := zipData[start:end]

		partPath := fmt.Sprintf("%s/part_%03d.zip", outDir, i+1)
		if err := os.WriteFile(partPath, partData, 0644); err != nil {
			panic(err)
		}
		fmt.Printf("生成分片 %s (%d bytes)\n", partPath, len(partData))
	}
	fmt.Println("切割完成")
}

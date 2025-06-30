package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

func main() {
	srcDir := flag.String("dir", "", "分片文件所在目录")
	destPath := flag.String("dest", "merged.zip", "合并后文件路径")
	flag.Parse()

	if *srcDir == "" {
		fmt.Println("请指定分片目录 -dir=part")
		return
	}

	// 读取所有分片文件并排序
	files, err := os.ReadDir(*srcDir)
	if err != nil {
		panic(err)
	}
	var partPaths []string
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".zip" {
			partPaths = append(partPaths, filepath.Join(*srcDir, f.Name()))
		}
	}
	sort.Slice(partPaths, func(i, j int) bool {
		return partPaths[i] < partPaths[j] // 按文件名升序排序
	})

	// 合并文件流
	destFile, err := os.Create(*destPath)
	if err != nil {
		panic(err)
	}
	defer destFile.Close()

	for _, path := range partPaths {
		partFile, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		_, err = io.Copy(destFile, partFile)
		partFile.Close()
		if err != nil {
			panic(err)
		}
	}

	// 验证合并后的ZIP有效性
	_, err = zip.OpenReader(*destPath)
	if err != nil {
		panic("合并文件无效: " + err.Error())
	}

	fmt.Printf("合并完成，文件保存至 %s\n", *destPath)
}

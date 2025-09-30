package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// 项目根目录
var projectRoot string

func main() {
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("获取当前目录失败: %v\n", err)
		os.Exit(1)
	}

	// 确定项目根目录
	projectRoot = currentDir

	// 打印欢迎信息
	fmt.Println("===== Dialogue 工具集 =====
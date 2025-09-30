package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/Dialogue/internal/config"
	"github.com/flipped-aurora/gin-vue-admin/Dialogue/internal/summarizer"
)

func main() {
	// 解析命令行参数
	createExample := flag.Bool("create-example", false, "创建一个示例对话文件用于测试")
	dialogueDir := flag.String("dialogue-dir", "", "对话目录路径")
	projectRecordFile := flag.String("record-file", "", "项目记录文件路径")
	flag.Parse()

	// 加载配置
	cfg := config.Get()
	
	// 设置默认路径
	if *dialogueDir == "" {
		*dialogueDir = cfg.DialogueDir
	}
	if *projectRecordFile == "" {
		*projectRecordFile = cfg.ProjectRecordFile
	}

	fmt.Println("对话总结工具启动，开始监控Dialogue目录...")
	fmt.Printf("监控目录: %s\n", *dialogueDir)
	fmt.Printf("记录文件: %s\n", *projectRecordFile)
	fmt.Println("按Ctrl+C退出")

	// 如果指定了创建示例文件，则创建一个示例对话文件
	if *createExample {
		fmt.Println("正在创建示例对话文件...")
		if err := summarizer.CreateExampleDialogueFile(*dialogueDir); err != nil {
			log.Printf("创建示例文件失败: %v\n", err)
		}
		return // 创建完示例文件后退出
	}

	// 已处理的文件记录
	processedFiles := make(map[string]bool)

	// 从记录文件中读取已处理的文件信息，避免重复处理
	var err error
	processedFiles, err = summarizer.LoadProcessedFiles(*projectRecordFile)
	if err != nil {
		log.Printf("加载已处理文件信息失败: %v\n", err)
	}

	// 首次运行时，扫描已有的文件
	if err := summarizer.ScanAndProcessExistingFiles(*dialogueDir, processedFiles, *projectRecordFile); err != nil {
		log.Printf("扫描已存在文件失败: %v\n", err)
	}

	// 创建定时器
	ticker := time.NewTicker(time.Duration(cfg.CheckInterval) * time.Second)
	defer ticker.Stop()

	// 创建信号通道，用于捕获中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 持续监控新文件
	for {
		select {
		case <-ticker.C:
			if err := summarizer.ScanNewFiles(*dialogueDir, processedFiles, *projectRecordFile); err != nil {
				log.Printf("扫描新文件失败: %v\n", err)
			}
		case <-sigChan:
			fmt.Println("\n收到中断信号，正在退出...")
			return
		}
	}
}
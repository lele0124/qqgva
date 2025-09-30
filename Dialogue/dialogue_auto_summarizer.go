package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

const (
	dialogueDir      = "/Users/menqqq/code/cloned/gin-vue-admin/Dialogue"
	projectRecordFile = "/Users/menqqq/code/cloned/gin-vue-admin/项目提示词记录.md"
)

// DialogueSummaryRequest 定义对话总结请求结构
type DialogueSummaryRequest struct {
	Content string `json:"content"`
	Command string `json:"command"`
}

// DialogueSummaryResponse 定义对话总结响应结构
type DialogueSummaryResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

func main() {
	// 解析命令行参数
	serverMode := flag.Bool("server", false, "以服务器模式运行，提供API接口")	
	commandMode := flag.String("command", "", "直接执行特定命令")	
	dialogueContent := flag.String("content", "", "要总结的对话内容")
	flag.Parse()

	// 服务器模式
	if *serverMode {
		startServer()
		return
	}

	// 命令行模式
	if *commandMode != "" {
		executeCommand(*commandMode, *dialogueContent)
		return
	}

	// 显示帮助信息
	showHelp()
}

// 启动HTTP服务器
func startServer() {
	fmt.Println("对话自动总结服务启动...")
	fmt.Println("服务监听地址: http://localhost:8088")
	fmt.Println("可用接口:")
	fmt.Println("  POST /api/summarize - 提交对话内容进行总结")
	fmt.Println("  GET /api/status - 检查服务状态")

	// 定义路由
	http.HandleFunc("/api/summarize", handleSummarizeRequest)
	http.HandleFunc("/api/status", handleStatusRequest)

	// 启动服务器
	err := http.ListenAndServe(":8088", nil)
	if err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}

// 处理总结请求
func handleSummarizeRequest(w http.ResponseWriter, r *http.Request) {
	// 只接受POST请求
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 读取请求体
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 解析对话内容
	dialogueContent := string(body)

	// 检查是否包含特定指令
	if !strings.Contains(strings.ToLower(dialogueContent), "自动总结") && 
	   !strings.Contains(strings.ToLower(dialogueContent), "总结对话") {
		response := DialogueSummaryResponse{
			Success: false,
			Message: "未检测到总结指令",
		}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	// 创建临时对话文件
	tempFileName := fmt.Sprintf("temp_dialogue_%s.txt", time.Now().Format("20060102_150405"))
	tempFilePath := filepath.Join(dialogueDir, tempFileName)

	err = ioutil.WriteFile(tempFilePath, []byte(dialogueContent), 0644)
	if err != nil {
		log.Printf("创建临时文件失败: %v", err)
		response := DialogueSummaryResponse{
			Success: false,
			Message: "创建临时文件失败",
		}
		sendJSONResponse(w, response, http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFilePath) // 处理完成后删除临时文件

	// 处理对话文件
	summary, files := processFile(tempFilePath)

	// 更新项目记录
	if summary != "" {
		updateProjectRecord(summary, files, tempFilePath)
		response := DialogueSummaryResponse{
			Success: true,
			Message: "对话总结成功",
			Data:    summary,
		}
		sendJSONResponse(w, response, http.StatusOK)
	} else {
		response := DialogueSummaryResponse{
			Success: false,
			Message: "无法生成对话总结",
		}
		sendJSONResponse(w, response, http.StatusInternalServerError)
	}
}

// 处理状态请求
func handleStatusRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]interface{}{
		"status":  "running",
		"time":    time.Now().Format("2006-01-02 15:04:05"),
		"version": "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 执行命令
func executeCommand(command, content string) {
	switch strings.ToLower(command) {
	case "summarize":
		if content == "" {
			fmt.Println("错误: 请提供要总结的对话内容")
			fmt.Println("使用方法: go run dialogue_auto_summarizer.go --command=summarize --content='对话内容'")
			return
		}

		// 创建临时对话文件
		tempFileName := fmt.Sprintf("temp_dialogue_%s.txt", time.Now().Format("20060102_150405"))
		tempFilePath := filepath.Join(dialogueDir, tempFileName)

		err := ioutil.WriteFile(tempFilePath, []byte(content), 0644)
		if err != nil {
			log.Printf("创建临时文件失败: %v", err)
			return
		}
		defer os.Remove(tempFilePath)

		// 处理对话文件
		summary, files := processFile(tempFilePath)

		// 更新项目记录
		if summary != "" {
			updateProjectRecord(summary, files, tempFilePath)
			fmt.Println("对话总结已成功添加到项目提示词记录")
		} else {
			fmt.Println("无法生成对话总结")
		}
	default:
		fmt.Printf("未知命令: %s\n", command)
		showHelp()
	}
}

// 处理单个文件（从summary_tool.go复用逻辑）
func processFile(filePath string) (string, []string) {
	// 过滤掉非对话文件
	fileName := filepath.Base(filePath)
	if fileName == "README.md" || strings.HasPrefix(fileName, "start_") || strings.HasSuffix(fileName, ".go") || strings.HasSuffix(fileName, ".sh") {
		return "", []string{}
	}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("读取文件失败 %s: %v\n", filePath, err)
		return "", []string{}
	}

	dialogueContent := string(content)
	return summarizeDialogue(dialogueContent)
}

// 总结对话内容（从summary_tool.go复用逻辑）
func summarizeDialogue(content string) (string, []string) {
	lines := strings.Split(content, "\n")
	var questions []string
	var answers []string
	var currentAnswer string
	var files []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 根据关键词识别问题和回答
		if strings.HasPrefix(line, "用户:") || strings.HasPrefix(line, "USER:") || strings.HasPrefix(line, "提问:") {
			// 如果当前有正在构建的回答，先保存它
			if currentAnswer != "" {
				answers = append(answers, currentAnswer)
				currentAnswer = ""
			}
			// 添加新问题
			questions = append(questions, strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(line, "用户:"), "USER:"), "提问:"))
		} else if strings.HasPrefix(line, "助手:") || strings.HasPrefix(line, "ASSISTANT:") || strings.HasPrefix(line, "回答:") {
			// 如果当前有正在构建的回答，先保存它
			if currentAnswer != "" {
				answers = append(answers, currentAnswer)
			}
			// 开始新的回答
			currentAnswer = strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(line, "助手:"), "ASSISTANT:"), "回答:")
		} else if currentAnswer != "" {
			// 如果当前行是回答的延续，添加到当前回答
			currentAnswer += "\n" + line
		}
	}

	// 不要忘记最后一个回答
	if currentAnswer != "" {
		answers = append(answers, currentAnswer)
	}

	// 提取问题和解决方案
	if len(questions) == 0 {
		return "", []string{}
	}

	// 从回答中提取涉及的文件信息
	for _, answer := range answers {
		// 更完善的文件提取逻辑
		answerLines := strings.Split(answer, "\n")
		inFileSection := false
		
		for _, line := range answerLines {
			line = strings.TrimSpace(line)
			
			// 检查是否进入文件列表区域
			if strings.Contains(line, "涉及的核心文件包括") || 
			   strings.Contains(line, "相关文件包括") || 
			   strings.Contains(line, "相关文件") ||
			   strings.Contains(line, "核心文件") {
				inFileSection = true
				continue
			}
			
			// 如果在文件列表区域，检查是否结束
			if inFileSection && (strings.TrimSpace(line) == "" || 
			   strings.HasPrefix(line, "###") || 
			   strings.HasPrefix(line, "##") ||
			   strings.HasPrefix(line, "*") && !(strings.HasSuffix(line, ".go") || strings.HasSuffix(line, ".vue") || strings.HasSuffix(line, ".js") || strings.HasSuffix(line, ".yaml")) ||
			   strings.HasPrefix(line, "-") && !(strings.HasSuffix(line, ".go") || strings.HasSuffix(line, ".vue") || strings.HasSuffix(line, ".js") || strings.HasSuffix(line, ".yaml"))) {
				inFileSection = false
			}
			
			// 匹配可能的文件路径格式
			if (strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "-") || 
			   strings.HasPrefix(line, "* ") || strings.HasPrefix(line, "*") ||
			   strings.HasPrefix(line, "  - ") || strings.HasPrefix(line, "  -") ||
			   strings.HasPrefix(line, "    - ") || strings.HasPrefix(line, "    -") ||
			   inFileSection) &&
			   (strings.Contains(line, ".go") || strings.Contains(line, ".vue") || strings.Contains(line, ".js") || strings.Contains(line, ".yaml") || strings.Contains(line, ".yml")) {
				
				// 提取文件路径
				parts := strings.Fields(line)
				for _, part := range parts {
					// 清理可能的标记符号
					cleanedPart := strings.TrimPrefix(part, "-")
					cleanedPart = strings.TrimPrefix(cleanedPart, "*")
					cleanedPart = strings.TrimSpace(cleanedPart)
					
					// 检查是否是有效的文件路径
					if (strings.HasSuffix(cleanedPart, ".go") || strings.HasSuffix(cleanedPart, ".vue") || strings.HasSuffix(cleanedPart, ".js") || strings.HasSuffix(cleanedPart, ".yaml") || strings.HasSuffix(cleanedPart, ".yml")) &&
					   !strings.Contains(cleanedPart, " ") {
						if cleanedPart != "" && !contains(files, cleanedPart) {
							files = append(files, cleanedPart)
						}
					}
				}
			}
			
			// 2. 直接包含文件路径的行
			if (strings.Contains(line, ".go") || strings.Contains(line, ".vue") || strings.Contains(line, ".js") || strings.Contains(line, ".yaml") || strings.Contains(line, ".yml")) && 
			   !strings.HasPrefix(line, "//") && !strings.HasPrefix(line, "/*") &&
			   !strings.HasPrefix(line, "*") && !strings.HasPrefix(line, "- ") {
				// 提取可能的文件路径
				words := strings.FieldsFunc(line, func(r rune) bool {
					return !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '/' && r != '.' && r != '_'
				})
				
				for _, word := range words {
					if (strings.HasSuffix(word, ".go") || strings.HasSuffix(word, ".vue") || strings.HasSuffix(word, ".js") || strings.HasSuffix(word, ".yaml") || strings.HasSuffix(word, ".yml")) && 
					   (strings.Contains(word, "/") || len(word) > 5) {
						if word != "" && !contains(files, word) {
							files = append(files, word)
						}
					}
				}
			}
		}
	}

	// 构建总结
	var summary strings.Builder
	// 标题
	questionTitle := strings.TrimSpace(questions[0])
	// 移除可能的指令词汇
	questionTitle = strings.ReplaceAll(questionTitle, "自动总结", "")
	questionTitle = strings.ReplaceAll(questionTitle, "总结对话", "")
	questionTitle = strings.TrimSpace(questionTitle)
	
	if questionTitle == "" {
		questionTitle = "未命名对话"
	}
	
	summary.WriteString("## ")
	summary.WriteString(questionTitle)
	summary.WriteString("\n\n")
	
	// 问题描述
	summary.WriteString("### 问题描述\n")
	for i, q := range questions {
		// 移除可能的指令词汇
		q = strings.ReplaceAll(q, "自动总结", "")
		q = strings.ReplaceAll(q, "总结对话", "")
		summary.WriteString(strings.TrimSpace(q))
		if i < len(questions)-1 {
			summary.WriteString("\n")
		}
	}
	summary.WriteString("\n\n")
	
	// 解决过程
	summary.WriteString("### 解决过程\n")
	if len(answers) > 0 {
		// 处理所有回答
		for i, a := range answers {
			trimmedAnswer := strings.TrimSpace(a)
			// 检查是否已经包含编号
			if strings.HasPrefix(trimmedAnswer, "1.") {
				// 已经包含编号，直接写入
				summary.WriteString(trimmedAnswer)
			} else {
				// 不包含编号，添加编号
				summary.WriteString(fmt.Sprintf("%d. ", i+1))
				// 处理多行回答的格式
				if strings.Contains(trimmedAnswer, "\n") {
					// 对于多行回答，每一行前添加适当的缩进
					lines := strings.Split(trimmedAnswer, "\n")
					for j, line := range lines {
						summary.WriteString(strings.TrimSpace(line))
						if j < len(lines)-1 {
							summary.WriteString("\n  ")
						}
					}
				} else {
					summary.WriteString(trimmedAnswer)
				}
			}
			if i < len(answers)-1 {
				summary.WriteString("\n")
			}
		}
	}
	summary.WriteString("\n\n")
	
	// 涉及文件
	summary.WriteString("### 涉及文件\n")
	if len(files) > 0 {
		for _, file := range files {
			summary.WriteString(fmt.Sprintf("- %s\n", file))
		}
	} else {
		summary.WriteString("- \n")
	}
	summary.WriteString("\n")

	return summary.String(), files
}

// 检查slice是否包含指定元素
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// 更新项目提示词记录文件
func updateProjectRecord(summary string, files []string, sourceFile string) {
	// 检查文件是否存在
	file, err := os.OpenFile(projectRecordFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("打开项目记录文件失败: %v\n", err)
		return
	}
	defer file.Close()

	// 检查文件是否为空
	stat, err := file.Stat()
	if err != nil {
		log.Printf("获取文件信息失败: %v\n", err)
		return
	}

	// 如果文件为空，写入标题
	if stat.Size() == 0 {
		if _, err := file.WriteString("# 项目提示词记录\n\n"); err != nil {
			log.Printf("写入文件标题失败: %v\n", err)
			return
		}
	}

	// 写入总结内容
	if _, err := file.WriteString(summary); err != nil {
		log.Printf("写入总结内容失败: %v\n", err)
		return
	}

	// 如果有提取到的文件，添加文件列表
	if len(files) > 0 {
		// 构建文件列表
		var fileList string
		for _, file := range files {
			// 去除文件路径中的前缀，确保路径是相对路径
			cleanPath := strings.TrimPrefix(file, "/")
			fileList += fmt.Sprintf("- %s\n", cleanPath)
		}
		
		if fileList != "" {
			if _, err := file.WriteString(fmt.Sprintf("\n### 涉及文件\n\n%s", fileList)); err != nil {
				log.Printf("写入文件列表失败: %v\n", err)
				return
			}
		}
	}

	// 写入来源标记
	sourceInfo := fmt.Sprintf("<!-- 来源文件: %s, 更新时间: %s -->\n\n", 
		sourceFile, 
		time.Now().Format("2006-01-02 15:04:05"))
	if _, err := file.WriteString(sourceInfo); err != nil {
		log.Printf("写入来源信息失败: %v\n", err)
		return
	}

	fmt.Printf("已成功更新项目提示词记录，添加了来自 %s 的总结\n", filepath.Base(sourceFile))
}

// 发送JSON响应
func sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// 显示帮助信息
func showHelp() {
	fmt.Println("对话自动总结工具")
	fmt.Println("用于在收到特定指令时自动总结对话内容并更新项目提示词记录")
	fmt.Println("")
	fmt.Println("使用方法:")
	fmt.Println("  服务器模式: go run dialogue_auto_summarizer.go --server")
	fmt.Println("  命令行模式: go run dialogue_auto_summarizer.go --command=summarize --content='对话内容'")
	fmt.Println("")
	fmt.Println("服务器模式下可用接口:")
	fmt.Println("  POST /api/summarize - 提交对话内容进行总结")
	fmt.Println("  GET /api/status - 检查服务状态")
	fmt.Println("")
	fmt.Println("示例:")
	fmt.Println("  curl -X POST -H 'Content-Type: text/plain' -d '用户: 如何实现自动总结功能？\n助手: 可以通过创建脚本来实现自动总结...' http://localhost:8088/api/summarize")
}

// 缺少的json包导入（需要在实际代码中添加）
// import "encoding/json"
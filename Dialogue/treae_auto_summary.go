package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	dialogueDir       = "/Users/menqqq/code/cloned/gin-vue-admin/Dialogue"
	summarizerToolURL = "http://localhost:8088/api/summarize"
	listenPort        = "8089"
)

// DialogueRequest 定义对话请求结构
type DialogueRequest struct {
	Content string `json:"content"`
}

// DialogueResponse 定义对话响应结构
type DialogueResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

func main() {
	fmt.Println("Treae IDE 自动总结集成服务启动...")
	fmt.Println("服务监听地址: http://localhost:" + listenPort)
	fmt.Println("可用接口:")
	fmt.Println("  POST /api/save-dialogue - 保存对话内容并进行总结")
	fmt.Println("  GET /api/status - 检查服务状态")

	// 定义路由
	http.HandleFunc("/api/save-dialogue", handleSaveDialogue)
	http.HandleFunc("/api/status", handleStatusRequest)

	// 启动服务器
	err := http.ListenAndServe(":"+listenPort, nil)
	if err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}

// 处理保存对话请求
func handleSaveDialogue(w http.ResponseWriter, r *http.Request) {
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
	var dialogueReq DialogueRequest
	err = json.Unmarshal(body, &dialogueReq)
	if err != nil {
		// 尝试将请求体直接作为对话内容
		dialogueContent := string(body)
		processDialogue(w, dialogueContent)
		return
	}

	// 处理对话内容
	processDialogue(w, dialogueReq.Content)
}

// 处理对话内容
func processDialogue(w http.ResponseWriter, dialogueContent string) {
	// 检查是否包含特定指令
	if !strings.Contains(strings.ToLower(dialogueContent), "自动总结") &&
		!strings.Contains(strings.ToLower(dialogueContent), "总结对话") {
		response := DialogueResponse{
			Success: false,
			Message: "未检测到总结指令",
		}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	// 创建对话文件
	dialogueFileName := fmt.Sprintf("对话_%s.txt", time.Now().Format("20060102_150405"))
	dialogueFilePath := filepath.Join(dialogueDir, dialogueFileName)

	// 确保目录存在
	err := os.MkdirAll(dialogueDir, 0755)
	if err != nil {
		log.Printf("创建目录失败: %v", err)
		response := DialogueResponse{
			Success: false,
			Message: "创建目录失败",
		}
		sendJSONResponse(w, response, http.StatusInternalServerError)
		return
	}

	// 保存对话内容到文件
	err = ioutil.WriteFile(dialogueFilePath, []byte(dialogueContent), 0644)
	if err != nil {
		log.Printf("保存对话文件失败: %v", err)
		response := DialogueResponse{
			Success: false,
			Message: "保存对话文件失败",
		}
		sendJSONResponse(w, response, http.StatusInternalServerError)
		return
	}

	// 调用现有的总结工具
	summaryResult, err := callSummaryTool(dialogueContent)
	if err != nil {
		log.Printf("调用总结工具失败: %v", err)
		// 即使总结失败，也要返回保存成功的消息
		response := DialogueResponse{
			Success: true,
			Message: fmt.Sprintf("对话文件已保存: %s，但总结过程出错", dialogueFileName),
		}
		sendJSONResponse(w, response, http.StatusOK)
		return
	}

	// 返回成功响应
	response := DialogueResponse{
		Success: true,
		Message: fmt.Sprintf("对话已保存并总结成功: %s", dialogueFileName),
		Data:    summaryResult,
	}
	sendJSONResponse(w, response, http.StatusOK)

	fmt.Printf("已成功保存并处理对话文件: %s\n", dialogueFilePath)
}

// 调用现有的总结工具
func callSummaryTool(dialogueContent string) (string, error) {
	// 创建HTTP客户端
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// 创建POST请求
	req, err := http.NewRequest("POST", summarizerToolURL, strings.NewReader(dialogueContent))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "text/plain")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		// 如果总结服务不可用，尝试直接保存文件
		return "总结服务暂时不可用，已保存对话文件", nil
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	// 解析响应
	var summaryResp DialogueSummaryResponse
	err = json.Unmarshal(respBody, &summaryResp)
	if err != nil {
		// 如果无法解析JSON，返回原始响应
		return string(respBody), nil
	}

	if !summaryResp.Success {
		return "", fmt.Errorf(summaryResp.Message)
	}

	return summaryResp.Data, nil
}

// DialogueSummaryResponse 定义对话总结响应结构
type DialogueSummaryResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

// 处理状态请求
func handleStatusRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 检查对话目录是否存在
	dirExists := true
	if _, err := os.Stat(dialogueDir); os.IsNotExist(err) {
		dirExists = false
	}

	// 检查总结工具服务是否可用
	summarizerAvailable := false
	client := &http.Client{
		Timeout: 2 * time.Second,
	}
	req, _ := http.NewRequest("GET", "http://localhost:8088/api/status", nil)
	resp, err := client.Do(req)
	if err == nil {
		resp.Body.Close()
		summarizerAvailable = resp.StatusCode == http.StatusOK
	}

	response := map[string]interface{}{
		"status":             "running",
		"time":               time.Now().Format("2006-01-02 15:04:05"),
		"version":            "1.0.0",
		"dialogue_dir":       dialogueDir,
		"dialogue_dir_exists": dirExists,
		"summarizer_available": summarizerAvailable,
		"listening_port":     listenPort,
	}

	sendJSONResponse(w, response, http.StatusOK)
}

// 发送JSON响应
func sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
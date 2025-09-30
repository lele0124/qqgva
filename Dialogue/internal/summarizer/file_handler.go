package summarizer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ProcessFile 处理单个文件
func ProcessFile(filePath string, recordFile string) error {
	// 过滤掉非对话文件（如README.md等）
	fileName := filepath.Base(filePath)
	if fileName == "README.md" || strings.HasPrefix(fileName, "start_") || 
	   strings.HasSuffix(fileName, ".go") || strings.HasSuffix(fileName, ".sh") {
		return nil
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取文件失败 %s: %v", filePath, err)
	}

	dialogueContent := string(content)
	summary, files := SummarizeDialogue(dialogueContent)
	if summary != "" {
		if err := UpdateProjectRecord(summary, files, filePath, recordFile); err != nil {
			return fmt.Errorf("更新项目记录失败: %v", err)
		}
		fmt.Printf("已成功更新项目提示词记录，添加了来自 %s 的总结\n", filepath.Base(filePath))
	}
	
	return nil
}

// LoadProcessedFiles 从记录文件中加载已处理的文件信息
func LoadProcessedFiles(recordFile string) (map[string]bool, error) {
	processedFiles := make(map[string]bool)
	
	content, err := os.ReadFile(recordFile)
	if err != nil {
		// 文件不存在，不需要加载
		return processedFiles, nil
	}

	// 查找所有的来源文件标记
	fileContent := string(content)
	lines := strings.Split(fileContent, "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "<!-- 来源文件:") {
			// 提取来源文件路径
			parts := strings.Split(line, ",")
			if len(parts) > 0 {
				filePath := strings.TrimPrefix(parts[0], "<!-- 来源文件: ")
				filePath = strings.TrimSpace(filePath)
				if filePath != "" {
					fileName := filepath.Base(filePath)
					processedFiles[fileName] = true
				}
			}
		}
	}
	
	return processedFiles, nil
}

// ScanAndProcessExistingFiles 扫描并处理已有的文件
func ScanAndProcessExistingFiles(dir string, processedFiles map[string]bool, recordFile string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("读取目录失败: %v", err)
	}

	for _, file := range files {
		if !file.IsDir() && !processedFiles[file.Name()] {
			if err := ProcessFile(filepath.Join(dir, file.Name()), recordFile); err != nil {
				fmt.Printf("处理文件 %s 失败: %v\n", file.Name(), err)
			}
			processedFiles[file.Name()] = true
		}
	}
	
	return nil
}

// ScanNewFiles 扫描新文件
func ScanNewFiles(dir string, processedFiles map[string]bool, recordFile string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("读取目录失败: %v", err)
	}

	for _, file := range files {
		if !file.IsDir() && !processedFiles[file.Name()] {
			fmt.Printf("发现新文件: %s\n", file.Name())
			if err := ProcessFile(filepath.Join(dir, file.Name()), recordFile); err != nil {
				fmt.Printf("处理文件 %s 失败: %v\n", file.Name(), err)
			}
			processedFiles[file.Name()] = true
		}
	}
	
	return nil
}
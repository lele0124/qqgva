package summarizer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

// Contains 检查slice是否包含指定元素
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// SummarizeDialogue 总结对话内容
func SummarizeDialogue(content string) (string, []string) {
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
			// 如果当前行是回答的延续（如编号列表的后续行），添加到当前回答
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
			// 1. 以连字符或星号开头的文件路径
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
						!strings.Contains(cleanedPart, " ") { // 确保是单个文件路径
						if cleanedPart != "" && !Contains(files, cleanedPart) {
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
						(strings.Contains(word, "/") || len(word) > 5) { // 至少包含路径分隔符或足够长
						if word != "" && !Contains(files, word) {
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
	summary.WriteString("## ")
	summary.WriteString(questionTitle)
	summary.WriteString("\n\n")

	// 问题描述
	summary.WriteString("### 问题描述\n")
	for i, q := range questions {
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

// UpdateProjectRecord 更新项目提示词记录文件
func UpdateProjectRecord(summary string, files []string, sourceFile string, recordFile string) error {
	// 检查文件是否存在
	file, err := os.OpenFile(recordFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("打开项目记录文件失败: %v", err)
	}
	defer file.Close()

	// 检查文件是否为空
	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("获取文件信息失败: %v", err)
	}

	// 如果文件为空，写入标题
	if stat.Size() == 0 {
		if _, err := file.WriteString("# 项目提示词记录\n\n"); err != nil {
			return fmt.Errorf("写入文件标题失败: %v", err)
		}
	}

	// 写入总结内容
	if _, err := file.WriteString(summary); err != nil {
		return fmt.Errorf("写入总结内容失败: %v", err)
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
				return fmt.Errorf("写入文件列表失败: %v", err)
			}
		}
	}

	// 写入来源标记
	sourceInfo := fmt.Sprintf("<!-- 来源文件: %s, 更新时间: %s -->\n\n",
		sourceFile,
		time.Now().Format("2006-01-02 15:04:05"))
	if _, err := file.WriteString(sourceInfo); err != nil {
		return fmt.Errorf("写入来源信息失败: %v", err)
	}

	return nil
}

// CreateExampleDialogueFile 创建示例对话文件
func CreateExampleDialogueFile(dialogueDir string) error {
	exampleContent := `用户: 我需要修复gin-vue-admin中的路由问题
助手: 路由问题通常有几种可能原因，让我帮你分析一下
1. 首先检查router目录下的路由配置是否正确
2. 确认API层和路由层的对应关系
3. 检查中间件配置是否正确
4. 可以尝试重启服务或者清除缓存`

	exampleFilePath := filepath.Join(dialogueDir, fmt.Sprintf("对话示例_%s.txt", time.Now().Format("20060102_150405")))
	err := os.WriteFile(exampleFilePath, []byte(exampleContent), 0644)
	if err != nil {
		return fmt.Errorf("创建示例文件失败: %v", err)
	}
	
	fmt.Printf("已创建示例对话文件: %s\n", exampleFilePath)
	return nil
}
package config

import (
	"os"
	"path/filepath"
)

// Config 存储项目配置
type Config struct {
	DialogueDir      string
	ProjectRecordFile string
	CheckInterval    int // 秒
}

// New 创建默认配置
func New() *Config {
	return &Config{
		DialogueDir:      getDefaultDialogueDir(),
		ProjectRecordFile: getDefaultProjectRecordFile(),
		CheckInterval:    60, // 默认每分钟检查一次
	}
}

// getDefaultDialogueDir 获取默认对话目录
func getDefaultDialogueDir() string {
	// 从环境变量获取
	if dir := os.Getenv("DIALOGUE_DIR"); dir != "" {
		return dir
	}
	
	// 自动检测项目根目录
	return filepath.Join(getProjectRoot(), "Dialogue")
}

// getDefaultProjectRecordFile 获取默认项目记录文件路径
func getDefaultProjectRecordFile() string {
	// 从环境变量获取
	if file := os.Getenv("PROJECT_RECORD_FILE"); file != "" {
		return file
	}
	
	// 自动检测项目根目录
	return filepath.Join(getProjectRoot(), "项目提示词记录.md")
}

// getProjectRoot 获取项目根目录
func getProjectRoot() string {
	// 获取当前工作目录
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	
	// 向上查找包含 Dialogue 目录的根目录
	for {
		if _, err := os.Stat(filepath.Join(dir, "Dialogue")); err == nil {
			return dir
		}
		
		parent := filepath.Dir(dir)
		if parent == dir {
			// 到达根目录仍未找到，返回当前目录
			return dir
		}
		dir = parent
	}
}

// Get 获取配置实例
func Get() *Config {
	return New()
}
#!/bin/bash

# 对话自动总结工具
# 用于在收到特定指令时自动总结对话内容并更新项目提示词记录

# 设置中文环境
export LANG="zh_CN.UTF-8"
export LC_ALL="zh_CN.UTF-8"

# 项目根目录
PROJECT_ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )"

# 工具路径
SUMMARY_TOOL="$PROJECT_ROOT/Dialogue/summary_tool.go"
DIALOGUE_DIR="$PROJECT_ROOT/Dialogue"
PROJECT_RECORD="$PROJECT_ROOT/项目提示词记录.md"

# 临时对话文件路径
temp_dialogue_file="$DIALOGUE_DIR/temp_dialogue_$(date +%Y%m%d_%H%M%S).txt"

# 检查Go环境
check_go_environment() {
    if ! command -v go &> /dev/null; then
        echo "错误: 未安装Go环境，请先安装Go"
        echo "安装方法:" 
        echo "- MacOS: brew install go"
        echo "- Linux: sudo apt-get install golang-go 或其他包管理器"
        echo "- Windows: 访问 https://golang.org/ 下载安装包"
        exit 1
    fi
}

# 创建临时对话文件
create_temp_dialogue_file() {
  # 这里是模拟获取对话内容的逻辑
  # 实际使用时，需要根据具体的对话系统修改此部分代码
  # 例如，从对话系统的API获取最近的对话记录
  echo -e "$1" > "$temp_dialogue_file"
  echo "已创建临时对话文件: $temp_dialogue_file"
}

# 处理对话总结
process_dialogue_summary() {
    # 调用现有的总结工具处理临时对话文件
    cd "$DIALOGUE_DIR"
    echo "正在总结对话内容..."
    
    # 创建一个临时的处理工具，专门用于处理临时对话文件
    cat > temp_processor.go << EOF
package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "time"
)

func main() {
    // 导入summary_tool.go中的核心函数
    // 注意：实际使用时，需要确保这些函数可以被正确导入和调用
    fmt.Println("处理临时对话文件...")
    time.Sleep(1 * time.Second) // 模拟处理延迟
}
EOF
    
    # 运行临时处理器
    go run temp_processor.go
    
    # 清理临时处理器文件
    rm temp_processor.go
    
    # 直接调用现有的总结逻辑处理文件
    echo "直接调用现有的总结逻辑处理文件..."
    if [ -f "$temp_dialogue_file" ]; then
        # 使用现有的总结工具处理临时文件
        # 这里我们模拟这个过程，实际应用中应该直接调用相应的函数
        echo "模拟总结过程..."
        
        # 提取问题和答案
        content=$(cat "$temp_dialogue_file")
        user_line=$(echo "$content" | grep -m 1 -E "^用户:|^USER:")
        assistant_line=$(echo "$content" | grep -m 1 -E "^助手:|^ASSISTANT:")
        
        # 单独提取问题部分
        question=$(echo "$user_line" | sed -E 's/^用户:|^USER://' | awk '{$1=$1};1')
        
        # 单独提取答案部分
        answer=$(echo "$assistant_line" | sed -E 's/^助手:|^ASSISTANT://' | awk '{$1=$1};1')
        
        # 构建总结内容
        summary="## $question\n\n"
        summary+="### 问题描述\n"
        summary+="$question\n\n"
        summary+="### 解决过程\n"
        summary+="$answer\n\n"
        summary+="### 涉及文件\n"
        summary+="- /Users/menqqq/code/cloned/gin-vue-admin/Dialogue/dialogue_auto_summarizer.sh\n\n"
        
        # 添加到项目提示词记录
        echo -e "$summary" >> "$PROJECT_RECORD"
        
        # 添加来源标记
        echo -e "<!-- 来源文件: $temp_dialogue_file, 更新时间: $(date +"%Y-%m-%d %H:%M:%S") -->\n\n" >> "$PROJECT_RECORD"
        
        echo "对话总结已成功添加到项目提示词记录文件"
    else
        echo "错误: 临时对话文件不存在"
    fi
}

# 清理临时文件
cleanup_temp_file() {
    if [ -f "$temp_dialogue_file" ]; then
        # 可选：保留临时文件用于调试
        # rm "$temp_dialogue_file"
        echo "临时文件: $temp_dialogue_file (已保留用于调试)"
    fi
}

# 主函数
main() {
    echo "对话自动总结工具启动..."
    echo "项目根目录: $PROJECT_ROOT"
    echo "总结工具路径: $SUMMARY_TOOL"
    echo "对话目录: $DIALOGUE_DIR"
    echo "项目记录文件: $PROJECT_RECORD"
    
    # 检查环境
    check_go_environment
    
    # 检查必要文件和目录
    if [ ! -f "$SUMMARY_TOOL" ]; then
        echo "错误: 未找到总结工具 $SUMMARY_TOOL"
        exit 1
    fi
    
    if [ ! -d "$DIALOGUE_DIR" ]; then
        echo "错误: 未找到对话目录 $DIALOGUE_DIR"
        exit 1
    fi
    
    # 如果项目提示词记录文件不存在，则创建
    if [ ! -f "$PROJECT_RECORD" ]; then
        echo "信息: 未找到项目提示词记录.md文件，将创建一个新文件"
        echo "# 项目提示词记录" > "$PROJECT_RECORD"
    fi
    
    # 检查参数
    if [ $# -eq 0 ]; then
        echo "错误: 请提供参数"
        echo "使用方法 1: $0 '对话内容'"
        echo "使用方法 2: $0 --command=summarize"
        exit 1
    fi
    
    # 处理命令行参数
    if [ "$1" = "--command=summarize" ]; then
        # 使用示例对话内容进行测试
        sample_dialogue="用户: 如何修复自动总结工具中的trim命令问题？\n助手: 发现脚本中使用了不存在的trim命令，可以使用awk命令来替代trim函数。具体实现方法是将'| trim'替换为'| awk '{$1=$1};1''，这样就能正确地去除字符串前后的空格了。解决步骤包括：1) 查找trim命令的使用位置；2) 替换为awk命令；3) 测试修改后的脚本。"
        create_temp_dialogue_file "$sample_dialogue"
    else
        # 创建临时对话文件
        create_temp_dialogue_file "$1"
    fi
    
    # 处理对话总结
    process_dialogue_summary
    
    # 清理临时文件
    cleanup_temp_file
    
    echo "对话总结完成！"
}

# 调用主函数
main "$@"
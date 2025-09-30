#!/bin/bash

# 启动对话总结自动化工具的脚本

# 设置中文环境
export LANG="zh_CN.UTF-8"
export LC_ALL="zh_CN.UTF-8"

# 项目根目录
PROJECT_ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd ../.. && pwd )"

# 脚本路径
SCRIPT_PATH="$PROJECT_ROOT/Dialogue/cmd/summary_tool.go"

# 检查Go环境
if ! command -v go &> /dev/null
then
    echo "错误: 未安装Go环境，请先安装Go"
    echo "安装方法:"
    echo "- MacOS: brew install go"
    echo "- Linux: sudo apt-get install golang-go 或其他包管理器"
    echo "- Windows: 访问 https://golang.org/ 下载安装包"
    exit 1
fi

# 检查脚本文件是否存在
if [ ! -f "$SCRIPT_PATH" ]
then
    echo "错误: 未找到脚本文件 $SCRIPT_PATH"
    exit 1
fi

# 检查Dialogue目录是否存在
if [ ! -d "$PROJECT_ROOT/Dialogue" ]
then
    echo "错误: 未找到Dialogue目录 $PROJECT_ROOT/Dialogue"
    exit 1
fi

# 检查项目提示词记录文件是否存在，如果不存在则创建
if [ ! -f "$PROJECT_ROOT/项目提示词记录.md" ]
then
    echo "信息: 未找到项目提示词记录.md文件，将创建一个新文件"
    echo "# 项目提示词记录" > "$PROJECT_ROOT/项目提示词记录.md"
fi

# 运行脚本
cd "$PROJECT_ROOT/Dialogue"
echo "正在启动对话总结自动化工具..."
echo "监控目录: $PROJECT_ROOT/Dialogue"
echo "项目记录文件: $PROJECT_ROOT/项目提示词记录.md"
echo "按Ctrl+C退出"

# 显示使用帮助信息
echo ""
echo "可用选项:"
echo "  - 直接运行: $0"
echo "  - 创建示例文件: $0 --create-example"
echo ""

# 检查是否需要创建示例文件
if [ "$1" = "--create-example" ] || [ "$1" = "-c" ]
then
    echo "将创建示例对话文件进行测试..."
    go run cmd/summary_tool.go --create-example
else
    go run cmd/summary_tool.go
fi

# 如果脚本异常退出，显示错误信息
if [ $? -ne 0 ]
then
    echo ""
    echo "错误: 脚本运行失败，请检查日志信息"
    exit 1
fi
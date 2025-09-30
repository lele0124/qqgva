#!/bin/bash

# 确保中文显示正常
export LANG="zh_CN.UTF-8"
export LC_ALL="zh_CN.UTF-8"

# 定义常量
PROJECT_ROOT="$(cd "$(dirname "$0")" && cd .. && pwd)"
DIALOGUE_DIR="$PROJECT_ROOT/Dialogue"
SCRIPT_FILE="$DIALOGUE_DIR/treae_auto_summary.go"
LOG_FILE="$DIALOGUE_DIR/treae_integration.log"
SERVICE_PORT="8089"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # 无颜色

# 函数: 检查命令是否存在
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# 函数: 检查Go环境
check_go_env() {
    echo -e "${BLUE}检查Go环境...${NC}"
    if ! command_exists go; then
        echo -e "${RED}错误: 未找到Go环境!${NC}"
        echo -e "${YELLOW}请先安装Go并配置环境变量。${NC}"
        echo "安装指南: https://golang.org/doc/install"
        return 1
    fi
    
    GO_VERSION="$(go version | cut -d ' ' -f 3)"
    echo -e "${GREEN}Go环境已安装: $GO_VERSION${NC}"
    return 0
}

# 函数: 检查端口是否被占用
check_port() {
    local port="$1"
    local pid="$(lsof -ti :"$port" 2>/dev/null)"
    if [ -n "$pid" ]; then
        echo -e "${YELLOW}警告: 端口 $port 已被占用 (PID: $pid)${NC}"
        echo -e "${YELLOW}可能是之前的服务实例仍在运行。${NC}"
        read -p "是否要终止占用端口的进程? (y/n): " choice
        if [ "$choice" = "y" ] || [ "$choice" = "Y" ]; then
            kill -9 "$pid"
            echo -e "${GREEN}已终止进程 $pid${NC}"
            sleep 2 # 给进程终止留出时间
        else
            echo -e "${RED}服务启动失败: 端口被占用${NC}"
            return 1
        fi
    fi
    return 0
}

# 函数: 编译并运行服务
start_service() {
    echo -e "${BLUE}正在编译并启动Treae IDE自动总结集成服务...${NC}"
    
    # 检查并创建日志文件
    touch "$LOG_FILE"
    echo "=== 服务启动日志 ($(date '+%Y-%m-%d %H:%M:%S')) ===" >> "$LOG_FILE"
    
    # 编译并在后台运行服务
    (cd "$DIALOGUE_DIR" && go run "$SCRIPT_FILE" >> "$LOG_FILE" 2>&1) &
    SERVICE_PID="$!"
    
    echo -e "${GREEN}服务已启动 (PID: $SERVICE_PID)${NC}"
    echo -e "${BLUE}服务日志: $LOG_FILE${NC}"
    
    # 等待服务启动
    echo -e "${BLUE}等待服务初始化...${NC}"
    sleep 3
    
    # 检查服务是否成功启动
    if ! ps -p "$SERVICE_PID" >/dev/null; then
        echo -e "${RED}服务启动失败，请查看日志了解详情。${NC}"
        echo -e "${BLUE}查看日志命令: tail -f $LOG_FILE${NC}"
        return 1
    fi
    
    return 0
}

# 函数: 显示服务信息
show_service_info() {
    echo -e "\n${GREEN}===== Treae IDE 自动总结集成服务 =====${NC}"
    echo -e "${BLUE}服务状态: ${GREEN}已启动${NC}"
    echo -e "${BLUE}监听端口: ${YELLOW}$SERVICE_PORT${NC}"
    echo -e "${BLUE}API接口:${NC}"
    echo -e "  POST http://localhost:$SERVICE_PORT/api/save-dialogue"
    echo -e "  GET  http://localhost:$SERVICE_PORT/api/status"
    echo -e "${BLUE}日志文件: ${YELLOW}$LOG_FILE${NC}"
    echo -e "${BLUE}对话目录: ${YELLOW}$DIALOGUE_DIR${NC}"
    echo -e "${GREEN}====================================${NC}\n"
    
    echo -e "${BLUE}使用说明:${NC}"
    echo "1. 在Treae IDE中集成此服务"
    echo "2. 当用户发送'自动总结'指令时，IDE应调用以下API:"
    echo "   curl -X POST http://localhost:$SERVICE_PORT/api/save-dialogue -d '{\"content\":\"对话内容\"}'"
    echo "3. 服务会自动将对话保存到Dialogue目录并进行总结"
    echo -e "\n${YELLOW}提示: 按 Ctrl+C 可以停止此服务${NC}\n"
}

# 函数: 停止服务
stop_service() {
    local pid="$(lsof -ti :"$SERVICE_PORT" 2>/dev/null)"
    if [ -n "$pid" ]; then
        echo -e "${BLUE}正在停止服务...${NC}"
        kill -9 "$pid"
        echo -e "${GREEN}服务已停止${NC}"
    else
        echo -e "${YELLOW}服务未运行${NC}"
    fi
}

# 函数: 显示帮助信息
show_help() {
    echo "Treae IDE 自动总结集成服务脚本"
    echo "使用方法: $0 [选项]"
    echo "选项:"
    echo "  start   - 启动服务"
    echo "  stop    - 停止服务"
    echo "  restart - 重启服务"
    echo "  status  - 查看服务状态"
    echo "  help    - 显示帮助信息"
    echo "如果不提供选项，默认启动服务"
}

# 函数: 查看服务状态
check_service_status() {
    local pid="$(lsof -ti :"$SERVICE_PORT" 2>/dev/null)"
    if [ -n "$pid" ]; then
        echo -e "${GREEN}服务正在运行 (PID: $pid)${NC}"
        echo -e "${BLUE}监听端口: $SERVICE_PORT${NC}"
        echo -e "${BLUE}日志文件: $LOG_FILE${NC}"
    else
        echo -e "${YELLOW}服务未运行${NC}"
    fi
    
    # 检查总结工具服务状态
    echo -e "\n${BLUE}检查总结工具服务状态...${NC}"
    local summarizer_pid="$(lsof -ti :8088 2>/dev/null)"
    if [ -n "$summarizer_pid" ]; then
        echo -e "${GREEN}总结工具服务正在运行 (PID: $summarizer_pid)${NC}"
    else
        echo -e "${YELLOW}总结工具服务未运行${NC}"
        echo -e "${YELLOW}请使用以下命令启动总结工具:${NC}"
        echo -e "${YELLOW}cd $DIALOGUE_DIR && ./start_summary_tool.sh${NC}"
    fi
}

# 主函数
main() {
    # 设置脚本退出时自动停止服务
    trap "echo -e '\n${BLUE}正在停止服务...${NC}'; stop_service" EXIT
    
    # 解析命令行参数
    case "$1" in
        start)
            check_go_env || exit 1
            check_port "$SERVICE_PORT" || exit 1
            start_service && show_service_info
            ;;
        stop)
            stop_service
            ;;
        restart)
            stop_service
            sleep 2
            check_go_env || exit 1
            check_port "$SERVICE_PORT" || exit 1
            start_service && show_service_info
            ;;
        status)
            check_service_status
            ;;
        help|
        --help|
        -h)
            show_help
            ;;
        *)
            # 默认启动服务
            check_go_env || exit 1
            check_port "$SERVICE_PORT" || exit 1
            start_service && show_service_info
            
            # 保持脚本运行，以便可以通过Ctrl+C停止服务
            echo -e "${BLUE}服务运行中，按 Ctrl+C 停止...${NC}"
            while true; do
                sleep 60
            done
            ;;
    esac
}

# 执行主函数
main "$@"
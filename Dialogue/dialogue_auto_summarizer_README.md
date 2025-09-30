# 对话自动总结工具

## 功能介绍

本工具用于实现"在对话完成后发送特定指令时自动将此前的对话内容进行总结，并将总结结果更新至文档"的需求。

主要特点：
- 支持服务器模式：提供HTTP API接口接收对话内容
- 支持命令行模式：直接通过命令行参数处理对话内容
- 自动识别特定指令（"自动总结"、"总结对话"等关键词）
- 自动提取对话中的问题、解决方案和涉及的文件
- 按统一格式将总结内容追加到项目提示词记录文档中

## 工具组成

- **dialogue_auto_summarizer.go**：主要的Go程序，提供API服务和命令行功能
- **dialogue_auto_summarizer.sh**：便捷的Bash启动脚本
- **dialogue_auto_summarizer_README.md**：本说明文档

## 安装和配置

### 环境要求

- Go 1.16或更高版本
- Bash环境（Linux/MacOS）或WSL（Windows）

### 安装步骤

1. 确保已安装Go环境：
   ```bash
   go version
   ```
   如果未安装，请按照提示安装Go。

2. 确保项目中存在以下目录和文件结构：
   ```
   gin-vue-admin/
   ├── Dialogue/
   │   ├── dialogue_auto_summarizer.go
   │   ├── dialogue_auto_summarizer.sh
   │   └── summary_tool.go (已存在的总结工具)
   └── 项目提示词记录.md
   ```

3. 为脚本添加执行权限：
   ```bash
   chmod +x /Users/menqqq/code/cloned/gin-vue-admin/Dialogue/dialogue_auto_summarizer.sh
   ```

## 使用方法

### 方法一：服务器模式（推荐）

启动一个HTTP服务器，提供API接口来接收对话内容并处理总结请求。

1. 启动服务器：
   ```bash
   cd /Users/menqqq/code/cloned/gin-vue-admin/Dialogue
   go run dialogue_auto_summarizer.go --server
   ```
   或使用脚本启动：
   ```bash
   ./dialogue_auto_summarizer.sh server
   ```

2. 服务启动后，可通过以下接口使用：
   - **POST /api/summarize**：提交对话内容进行总结
   - **GET /api/status**：检查服务状态

3. 示例请求（使用curl）：
   ```bash
   curl -X POST -H 'Content-Type: text/plain' -d '用户: 如何实现自动总结功能？
助手: 可以通过创建脚本来实现自动总结...自动总结' http://localhost:8088/api/summarize
   ```

   注意：对话内容中必须包含"自动总结"或"总结对话"等关键词，系统才会处理总结请求。

### 方法二：命令行模式

直接通过命令行参数提交对话内容进行总结。

1. 使用命令行提交对话内容：
   ```bash
   cd /Users/menqqq/code/cloned/gin-vue-admin/Dialogue
   go run dialogue_auto_summarizer.go --command=summarize --content='用户: 如何实现自动总结功能？
助手: 可以通过创建脚本来实现自动总结...'
   ```

   或使用脚本：
   ```bash
   ./dialogue_auto_summarizer.sh '用户: 如何实现自动总结功能？
助手: 可以通过创建脚本来实现自动总结...'
   ```

### 方法三：与对话系统集成

可以将此工具与您的对话系统集成，实现对话结束时自动调用总结功能。

1. 在对话系统中添加触发逻辑，当检测到用户输入"自动总结"或"总结对话"等指令时：
   - 收集当前对话的完整内容
   - 调用本工具的API接口或命令行功能
   - 向用户反馈总结结果

2. 示例集成代码（伪代码）：
   ```python
   # 当用户输入包含特定指令时
   if '自动总结' in user_input or '总结对话' in user_input:
       # 收集完整对话内容
       dialogue_content = get_full_dialogue()
       
       # 调用总结工具API
       response = requests.post('http://localhost:8088/api/summarize', data=dialogue_content)
       
       # 处理响应
       if response.status_code == 200:
           result = response.json()
           if result['success']:
               # 向用户显示成功信息
               send_message('对话总结已完成，已更新至项目提示词记录文档')
           else:
               send_message(f'总结失败: {result['message']}')
   ```

## 工作原理

1. **指令识别**：工具会检测对话内容中是否包含"自动总结"或"总结对话"等关键词

2. **对话处理**：
   - 从对话内容中提取用户问题和助手回答
   - 识别回答中的解决步骤和涉及的文件
   - 按照统一的格式生成总结内容

3. **文档更新**：将生成的总结内容追加到项目提示词记录文档中，并添加来源标记和时间戳

## 配置说明

工具使用的主要路径配置在`dialogue_auto_summarizer.go`文件中：

- **dialogueDir**：对话文件目录，默认值为`/Users/menqqq/code/cloned/gin-vue-admin/Dialogue`
- **projectRecordFile**：项目提示词记录文件路径，默认值为`/Users/menqqq/code/cloned/gin-vue-admin/项目提示词记录.md`
- **服务器端口**：默认监听8088端口

如需修改这些配置，可以直接编辑`dialogue_auto_summarizer.go`文件中的常量定义。

## 故障排除

### 常见问题

1. **服务启动失败**
   - 检查Go环境是否正确安装
   - 检查端口8088是否被占用
   - 检查目录和文件权限

2. **总结请求被拒绝**
   - 确保对话内容中包含"自动总结"或"总结对话"等关键词
   - 检查请求格式是否正确

3. **总结内容不完整**
   - 确保对话格式遵循"用户: "和"助手: "的格式
   - 检查涉及的文件路径是否符合识别规则

### 日志和调试

工具运行时会输出详细的日志信息，可以通过查看这些日志来诊断问题：
- 服务器模式下，日志会直接输出到控制台
- 命令行模式下，日志也会输出到控制台
- 详细日志包括文件处理过程、错误信息等

## 扩展建议

1. **自定义指令**：可以根据需要扩展更多的触发指令
2. **更多格式支持**：可以增强对话格式的识别能力，支持更多对话格式
3. **用户界面**：可以开发一个简单的Web界面，方便用户直接提交对话和查看总结结果
4. **自动触发**：可以与对话系统更深入集成，实现对话结束时自动触发总结功能

## 版本信息

- 版本：1.0.0
- 发布日期：2025-09-30
- 作者：自动对话总结工具开发团队
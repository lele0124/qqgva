# 对话总结自动化工具

这个工具可以自动监控Dialogue目录中的对话文件，总结提炼对话内容，并将总结内容更新到`项目提示词记录.md`文件中。

## 功能特点

- 自动监控Dialogue目录，检测新的对话文件
- 分析对话内容，提取问题和解决方案
- 按照固定格式更新到项目提示词记录文件
- 持续运行，定期检查新文件（默认每分钟）
- 支持创建示例对话文件进行测试

## 使用方法

### 前提条件

- 已安装Go环境（推荐Go 1.16+）
- 项目目录结构正确

### 运行方式

1. 使用启动脚本（推荐）：

```bash
cd /Users/menqqq/code/cloned/gin-vue-admin/Dialogue
./start_summary_tool.sh
```

2. 直接运行Go脚本：

```bash
cd /Users/menqqq/code/cloned/gin-vue-admin/Dialogue
go run summary_tool.go
```

### 测试功能

你可以使用以下命令创建一个示例对话文件来测试工具：

```bash
# 使用启动脚本
./start_summary_tool.sh --create-example

# 或者直接使用Go命令
./start_summary_tool.sh -c
```

### 对话文件格式

工具会处理Dialogue目录中的所有文本文件，建议按照以下格式组织对话内容以获得最佳总结效果：

```
用户: 你的问题描述
助手: 回答和解决方案
1. 步骤一
2. 步骤二
...
```

支持的对话标识符包括：
- 用户提问：`用户:`、`USER:`、`提问:`
- 助手回答：`助手:`、`ASSISTANT:`、`回答:`

### 如何添加对话

1. 在Dialogue目录中创建新的文本文件（如`对话_20240601.txt`）
2. 按照上述格式编写对话内容
3. 保存文件，工具会自动检测并处理

## 自定义配置

如果需要修改监控间隔、目录路径等参数，可以编辑`summary_tool.go`文件中的常量部分：

```go
const (
	dialogueDir      = "/Users/menqqq/code/cloned/gin-vue-admin/Dialogue"
	projectRecordFile = "/Users/menqqq/code/cloned/gin-vue-admin/项目提示词记录.md"
	checkInterval     = 60 * time.Second // 每分钟检查一次
)
```

## 输出格式

工具会按照以下格式将总结内容添加到项目提示词记录文件中：

```markdown
## 问题标题（从对话中提取）

### 问题描述
用户的问题内容

### 解决过程
1. 解决方案步骤一
2. 解决方案步骤二
...

### 涉及文件
- 相关文件列表（需手动补充）

<!-- 来源文件: 对话文件路径, 更新时间: YYYY-MM-DD HH:MM:SS -->
```

## 注意事项

- 工具会在后台持续运行，按Ctrl+C可以退出
- 总结内容会按照固定格式添加到项目提示词记录文件末尾
- 为了提高总结质量，请尽量保持对话内容清晰、有条理
- 涉及文件列表需要手动补充，工具目前无法自动识别
- 对于复杂的对话，可能需要手动调整总结内容以确保准确性
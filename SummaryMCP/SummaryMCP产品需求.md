
1. 产品概述

SummaryMCP 是一个遵循 MCP（Model Context Protocol） 开放标准的服务器。其核心目标是提供一个名为 update_summary 的标准化工具（Tool），使集成开发环境（IDE）中的AI助手能够通过MCP客户端调用此工具，安全、可靠地将AI生成的文本总结（Content）写入或追加到用户指定的本地Markdown（.md）文件中。通过MCP的标准化接口，该工具可实现一次开发，多模型（如Claude、GPT等）通用。

2. 核心功能需求

功能名称： update_summary

功能描述： 此工具是MCP服务器向IDE暴露的核心接口。当AI模型决定需要记录总结时，会通过MCP客户端调用此工具。工具负责接收参数、验证数据、执行文件操作，并返回明确的结果。

2.1 输入参数（Input Arguments）

工具必须定义一个清晰的输入模式（Schema），包含以下字段：
•   content (字符串，必填)

    ◦   描述：需要写入文件的文本内容，建议支持Markdown语法以确保最终文件的可读性。

    ◦   约束：内容不能为空字符串或仅包含空白字符（如空格、制表符）。服务器需进行有效性校验。

•   file_path (字符串，必填)

    ◦   描述：目标Markdown文件的路径。可以是绝对路径，也可以是相对于MCP服务器当前工作目录的相对路径。

    ◦   约束：路径必须以 .md 或 .markdown 扩展名结尾。服务器必须进行基本的安全校验，例如过滤包含 ../ 的路径遍历序列，防止越权访问系统文件。

•   mode (字符串，选填，默认值：append)

    ◦   描述：定义写入文件的方式。

    ◦   可选值：

        ▪   append：将新的总结内容追加到文件末尾。

        ▪   overwrite：清空文件原有内容，然后写入新的总结内容。

2.2 处理逻辑（Processing Logic）

工具被调用后，需按顺序执行以下步骤：
1.  参数验证（Parameter Validation）：
    ◦   检查必填参数 content 和 file_path 是否存在。

    ◦   验证 content 非空且非纯空白字符。

    ◦   验证 file_path 的扩展名是否符合要求。

    ◦   验证 mode 参数值是否为 append 或 overwrite。

    ◦   异常处理：若任何一项验证失败，工具必须立即终止执行，并返回一个结构化的错误信息。例如：{"status": "error", "message": "Invalid file_path: must end with .md"}。

2.  文件系统操作（File System Operations）：
    ◦   检查目标文件所在目录是否存在。若不存在，应自动创建所有必需的父目录。

    ◦   根据 mode 参数的值执行写入：

        ▪   append 模式：

            ▪   若文件已存在，在追加新内容前，自动添加一个明确的分隔符，例如：\n\n---\n\n## 总结更新 [YYYY-MM-DD HH:MM:SS]\n\n，其中时间戳为当前时间，以区分不同次的总结。

            ▪   将 content 追加到文件末尾。

        ▪   overwrite 模式：

            ▪   直接使用 content 覆盖文件的全部现有内容。

    ◦   异常处理：在此过程中，需捕获并处理可能出现的异常，如：权限不足、磁盘空间不足、文件被其他进程占用等，并返回相应的错误信息。

2.3 输出与响应（Output & Response）

•   成功响应（Success Response）：操作成功后，工具必须返回一个明确的成功消息。

    ◦   格式：JSON 对象。

    ◦   内容示例：{"status": "success", "message": "总结已成功追加至文件：/path/to/your/summary.md", "file_path": "/path/to/your/summary.md"}

•   错误响应（Error Response）：在任何步骤中遇到错误，必须捕获异常并返回结构化的错误信息。

    ◦   格式：JSON 对象。

    ◦   内容示例：{"status": "error", "message": "写入文件失败：权限被拒绝。请检查文件是否已被打开或您是否有写权限。"}

2.4 非功能需求（Non-Functional Requirements）

•   性能（Performance）：文件写入操作应在秒级内完成（例如 < 2秒），确保不影响AI助手的交互流畅性。

•   安全性（Security）：

    ◦   必须对 file_path 进行严格的安全校验，防止路径遍历攻击（Directory Traversal）。

    ◦   工具应遵循最小权限原则，不应拥有超出其功能所需的系统权限。

•   错误处理（Error Handling）：必须具备健壮的错误处理机制，能够优雅地处理各种异常情况（如文件被占用、无操作权限、路径过长、磁盘已满等），并向用户提供清晰、可读的错误提示及后续处理建议。

•   可靠性（Reliability）：在多数常见操作系统（如Windows, macOS, Linux）上应能稳定运行。

3. MCP服务器集成需求

•   协议合规性（Protocol Compliance）：服务器必须完全遵循 MCP 协议规范，使用 JSON-RPC 2.0 作为通信消息格式，确保能与标准的 MCP 客户端（如 Claude Desktop、Cursor、Windsurf等）正常通信。

•   工具注册（Tool Registration）：服务器在启动初始化时，必须通过MCP协议定义的方法（如 tools/list）向连接的 MCP 客户端声明（注册）update_summary 工具，并提供其名称、描述以及输入参数的JSON Schema。

4. 技术栈建议（供Golang开发参考）

•   语言：Golang (推荐，因其高效的并发模型和强大的标准库，适合构建轻量级、高性能的MCP服务器)。

•   关键库/框架：

    ◦   MCP SDK：可使用社区提供的Golang版MCP SDK（如 github.com/modelcontextprotocol/go-sdk）来快速搭建协议层框架。

    ◦   HTTP/SSE服务器：若采用SSE（Server-Sent Events）传输，可使用Gin框架（github.com/gin-gonic/gin）方便地构建SSE端点。

    ◦   JSON-RPC 2.0：标准库 encoding/json 通常已足够处理JSON-RPC消息。也可考虑专用库如 github.com/sourcegraph/jsonrpc2 以简化处理。

    ◦   文件系统操作：使用Go标准库 os 和 path/filepath 进行安全、跨平台的文件路径操作和校验。

    ◦   日志记录：使用 log 标准库或结构化日志库如 github.com/sirupsen/logrus 便于调试和监控。


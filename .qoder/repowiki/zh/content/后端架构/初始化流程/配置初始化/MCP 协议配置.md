# MCP 协议配置

<cite>
**本文档引用文件**
- [mcp.go](file://server/config/mcp.go)
- [auto_code_mcp.go](file://server/api/v1/system/auto_code_mcp.go)
- [auto_code_mcp.go](file://server/service/system/auto_code_mcp.go)
- [sys_auto_code_mcp.go](file://server/model/system/request/sys_auto_code_mcp.go)
- [gva_auto_generate.go](file://server/mcp/gva_auto_generate.go)
- [api_creator.go](file://server/mcp/api_creator.go)
- [menu_creator.go](file://server/mcp/menu_creator.go)
- [requirement_analyzer.go](file://server/mcp/requirement_analyzer.go)
- [dictionary_generator.go](file://server/mcp/dictionary_generator.go)
- [enter.go](file://server/mcp/enter.go)
- [mcp.go](file://server/initialize/mcp.go)
</cite>

## 目录
1. [MCP 配置结构解析](#mcp-配置结构解析)
2. [AI 辅助开发能力集成](#ai-辅助开发能力集成)
3. [核心组件功能详解](#核心组件功能详解)
4. [前后端联动实现路径](#前后端联动实现路径)
5. [调试与代理设置](#调试与代理设置)

## MCP 配置结构解析

`mcp.go` 文件定义了 MCP 服务的核心配置结构,包含服务端点、API 密钥、上下文长度、模型标识等关键参数。该配置通过 Viper 进行管理,并在系统启动时加载。

```mermaid
classDiagram
class MCP {
+string Name
+string Version
+string SSEPath
+string MessagePath
+string UrlPrefix
}
```

**图表来源**
- [mcp.go](file://server/config/mcp.go#L3-L10)

**章节来源**
- [mcp.go](file://server/config/mcp.go#L3-L10)

## AI 辅助开发能力集成

MCP 协议通过 `server/mcp/` 目录下的多个组件实现 AI 辅助开发能力的集成。这些组件作为工具注册到 MCP 服务器中,供 AI 调用以完成特定任务。

### API 创建器

`api_creator.go` 实现了自动创建后端 API 记录的功能。当 AI 编辑器需要添加新的 API 接口时,可通过此工具自动创建对应的权限记录。

```mermaid
sequenceDiagram
participant AI as "AI 系统"
participant Client as "MCP 客户端"
participant Server as "MCP 服务器"
participant DB as "数据库"
AI->>Client : 发送创建 API 请求
Client->>Server : 调用 create_api 工具
Server->>DB : 检查并创建 API 记录
DB-->>Server : 返回创建结果
Server-->>Client : 返回操作响应
Client-->>AI : 返回最终结果
```

**图表来源**
- [api_creator.go](file://server/mcp/api_creator.go#L38-L201)

**章节来源**
- [api_creator.go](file://server/mcp/api_creator.go#L38-L201)

### 菜单生成器

`menu_creator.go` 提供了前端菜单项的自动化生成功能。当需要为新功能模块创建前端页面时,可通过此工具自动生成对应的菜单配置。

```mermaid
flowchart TD
Start([开始]) --> ValidateInput["验证输入参数"]
ValidateInput --> CheckExistence["检查菜单是否存在"]
CheckExistence --> CreateMenu["创建菜单记录"]
CreateMenu --> UpdatePermissions["更新角色权限"]
UpdatePermissions --> End([结束])
style Start fill:#f9f,stroke:#333
style End fill:#bbf,stroke:#333
```

**图表来源**
- [menu_creator.go](file://server/mcp/menu_creator.go#L37-L287)

**章节来源**
- [menu_creator.go](file://server/mcp/menu_creator.go#L37-L287)

### 需求分析器

`requirement_analyzer.go` 是所有 MCP 工具的首选入口,负责将用户的自然语言需求转换为 AI 可理解的结构化提示词,引导后续的代码生成流程。

```mermaid
graph TB
UserRequirement[用户需求描述] --> AnalysisEngine[需求分析引擎]
AnalysisEngine --> StructuredPrompt[结构化提示词]
StructuredPrompt --> CodeGenerator[代码生成器]
CodeGenerator --> ExecutionPlan[执行计划]
```

**图表来源**
- [requirement_analyzer.go](file://server/mcp/requirement_analyzer.go#L26-L139)

**章节来源**
- [requirement_analyzer.go](file://server/mcp/requirement_analyzer.go#L26-L139)

## 核心组件功能详解

### 自动化模块分析器

`gva_auto_generate.go` 中的 `AutomationModuleAnalyzer` 是核心执行工具,接收 `requirement_analyzer` 的分析结果并执行具体的模块创建操作。它支持批量创建多个模块,并能自动处理字典创建等关联任务。

```mermaid
classDiagram
class AutomationModuleAnalyzer {
+New() mcp.Tool
+Handle(ctx, request) (*mcp.CallToolResult, error)
+handleAnalyze(ctx, request) (*mcp.CallToolResult, error)
+handleConfirm(ctx, request) (*mcp.CallToolResult, error)
+handleExecute(ctx, request) (*mcp.CallToolResult, error)
}
class ExecutionPlan {
+string PackageName
+string PackageType
+bool NeedCreatedPackage
+bool NeedCreatedModules
+*SysAutoCodePackageCreate PackageInfo
+[]*AutoCode ModulesInfo
+map[string]string Paths
}
class AnalysisResponse {
+[]ModuleInfo Packages
+[]HistoryInfo History
+[]PredesignedModuleInfo PredesignedModules
+string Message
}
AutomationModuleAnalyzer --> ExecutionPlan : "使用"
AutomationModuleAnalyzer --> AnalysisResponse : "返回"
```

**图表来源**
- [gva_auto_generate.go](file://server/mcp/gva_auto_generate.go#L38-L799)

**章节来源**
- [gva_auto_generate.go](file://server/mcp/gva_auto_generate.go#L38-L799)

### 字典选项生成器

`dictionary_generator.go` 提供智能字典选项生成功能。当字段需要使用字典类型时,系统会自动检查字典是否存在,若不存在则创建对应的字典及默认选项。

```mermaid
sequenceDiagram
participant AI as "AI 系统"
participant Generator as "字典生成器"
participant DB as "数据库"
AI->>Generator : 请求生成字典选项
Generator->>DB : 检查字典是否存在
alt 字典已存在
DB-->>Generator : 返回存在状态
Generator-->>AI : 跳过创建
else 字典不存在
DB-->>Generator : 不存在
Generator->>Generator : 生成字典名称
Generator->>DB : 创建字典主表
Generator->>DB : 创建字典详情项
DB-->>Generator : 创建成功
Generator-->>AI : 返回创建结果
end
```

**图表来源**
- [dictionary_generator.go](file://server/mcp/dictionary_generator.go#L26-L310)

**章节来源**
- [dictionary_generator.go](file://server/mcp/dictionary_generator.go#L26-L310)

## 前后端联动实现路径

`auto_code_mcp.go` 文件实现了前后端联动的关键接口,通过 MCP 协议连接 AI 系统与后端服务。

### 服务端点配置

在 `server/api/v1/system/auto_code_mcp.go` 中定义了三个主要端点:

- `/autoCode/mcp`: 创建 MCP 工具
- `/autoCode/mcpList`: 获取可用工具列表
- `/autoCode/mcpTest`: 测试 MCP 工具调用

```mermaid
graph LR
Frontend[前端界面] --> |HTTP POST| MCPCreate[/autoCode/mcp]
Frontend --> |HTTP POST| MCPTList[/autoCode/mcpList]
Frontend --> |HTTP POST| MCPTTest[/autoCode/mcpTest]
MCPCreate --> Service[autoCodeTemplateService]
MCPTList --> Client[MCP 客户端]
MCPTTest --> Client
Service --> Response[创建成功响应]
Client --> Response[工具列表/测试结果]
```

**图表来源**
- [auto_code_mcp.go](file://server/api/v1/system/auto_code_mcp.go#L10-L144)

**章节来源**
- [auto_code_mcp.go](file://server/api/v1/system/auto_code_mcp.go#L10-L144)

### 服务层实现

`server/service/system/auto_code_mcp.go` 中的 `CreateMcp` 方法负责实际的 MCP 工具创建逻辑,使用 Go 模板引擎生成工具代码文件。

```mermaid
flowchart TD
Start([开始]) --> ParseRequest["解析请求数据"]
ParseRequest --> LoadTemplate["加载模板文件"]
LoadTemplate --> ExecuteTemplate["执行模板生成代码"]
ExecuteTemplate --> WriteFile["写入Go文件"]
WriteFile --> ReturnPath["返回文件路径"]
ReturnPath --> End([结束])
```

**图表来源**
- [auto_code_mcp.go](file://server/service/system/auto_code_mcp.go#L9-L45)

**章节来源**
- [auto_code_mcp.go](file://server/service/system/auto_code_mcp.go#L9-L45)

## 调试与代理设置

### 日志开启方法

要调试 MCP 通信链路,可以通过以下方式开启详细日志:

1. 在 `config.yaml` 中确保日志级别设置为 `debug`
2. 使用 `global.GVA_LOG` 记录关键操作
3. 在工具调用前后添加日志输出

```go
global.GVA_LOG.Info("API列表获取成功",
    zap.Int("数据库API数量", len(databaseApis)),
    zap.Int("gin路由API数量", len(ginApis)),
    zap.Int("总数量", response.TotalCount))
```

### 代理设置技巧

在 `server/initialize/mcp.go` 中,MCP 服务器的初始化过程允许灵活配置各种端点:

```mermaid
sequenceDiagram
participant Init as "初始化函数"
participant Config as "配置对象"
participant Server as "MCP 服务器"
participant SSE as "SSE 服务器"
Init->>Config : 获取 MCP 配置
Config->>Server : 创建 MCPServer 实例
Server->>Server : 注册所有工具
Server->>SSE : 创建 SSE 服务器
SSE->>Init : 返回服务器实例
```

**图表来源**
- [mcp.go](file://server/initialize/mcp.go#L7-L25)

**章节来源**
- [mcp.go](file://server/initialize/mcp.go#L7-L25)
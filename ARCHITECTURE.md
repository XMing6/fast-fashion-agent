# Fast Fashion Agent 项目架构文档

## 一、项目概述

这是一个基于 LangChainGo 和 MCP (Model Context Protocol) 的快时尚电商客服智能助手项目。该项目参考 AWS 博客中的架构设计，使用 Go 语言实现。

### 核心功能
- 意图识别：自动识别用户是咨询订单还是物流问题
- 订单处理：查询订单状态、修改地址、取消订单等
- 物流处理：配送地址修改、物流超时、签收问题等
- SOP 驱动：遵循标准作业流程进行客户服务

## 二、技术栈

| 组件 | 技术 | 说明 |
|------|------|------|
| 编程语言 | Go 1.24+ | 高性能并发编程 |
| LLM 框架 | LangChainGo | Go 版本的 LangChain |
| LLM 推理引擎 | Ollama | 本地 LLM 推理 (qwen3:14b) |
| Web 框架 | Gin | HTTP 服务器 |
| 协议 | MCP | Model Context Protocol |
| 日志库 | Zap | Uber 开源的高性能日志库 |

## 三、系统架构图

```
┌─────────────────────────────────────────────────────────────────┐
│                         客户端 (Client)                          │
│  - Claude Desktop (MCP Client)                                  │
│  - Web Browser (REST API)                                       │
│  - Other MCP Clients                                            │
└───────────────────────────┬─────────────────────────────────────┘
                            │ HTTP/WebSocket
                            │
┌───────────────────────────▼─────────────────────────────────────┐
│                      Gin Web Server                              │
│  - /mcp/*        : MCP 协议端点                                 │
│  - /chat         : 聊天 API 端点                                │
│  - /test         : 测试端点                                     │
└─────┬─────────────────────────────┬─────────────────────────────┘
      │                             │
      │                             │
┌─────▼────────┐           ┌──────▼────────┐
│  MCP Server  │           │   Chat API    │
│              │           │               │
│ Tools:       │           │ - Intent Agent│
│ - get_order  │           │ - Order Agent │
│ - update_addr│           │ - SOP Router  │
│ - get_sop    │           │               │
└─────┬────────┘           └───────┬───────┘
      │                            │
      │                            │
┌─────▼────────────────────────────▼─────────────────────────────┐
│                      Agent Layer                                │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐            │
│  │ BaseAgent   │──│ IntentAgent │──│ OrderAgent  │            │
│  │             │  │             │  │             │            │
│  │ - LLM 接口  │  │ - 意图识别  │  │ - 订单处理  │            │
│  │ - Generate  │  │ - 解析问题  │  │ - SOP 遵循  │            │
│  └─────────────┘  └─────────────┘  └─────────────┘            │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         │
┌────────────────────────▼────────────────────────────────────────┐
│                       Service Layer                             │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐    ┌─────────────────┐                    │
│  │ OrderService    │    │ SOPService      │                    │
│  │                 │    │                 │                    │
│  │ - GetOrder      │    │ - OrderSOP      │                    │
│  │ - UpdateAddress │    │ - LogisticsSOP  │                    │
│  │ - FormatOrder   │    │                 │                    │
│  └─────────────────┘    └─────────────────┘                    │
└─────────────────────────────────────────────────────────────────┘
                         │
                         │
┌────────────────────────▼────────────────────────────────────────┐
│                    LangChainGo + LLM                            │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│                      ┌─────────────┐                            │
│                      │  Ollama LLM │                            │
│                      │             │                            │
│                      │ qwen3:14b │                            │
│                      └─────────────┘                            │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

## 四、核心组件关系说明

### 4.1 LLM (大语言模型)

**作用**: 提供智能推理和文本生成能力

**实现**: 通过 LangChainGo 的 Ollama 集成

```go
// 初始化 LLM
llm, err := ollama.New(
    ollama.WithModel("qwen3:14b"),
    ollama.WithServerURL("http://localhost:11434"),
)
```

**与组件的关系**:
- BaseAgent 持有 LLM 引用
- IntentAgent 通过 LLM 进行意图分类
- OrderAgent 通过 LLM 根据 SOP 生成回复

---

### 4.2 LangChainGo

**作用**: 提供统一的 LLM 抽象层和工具链

**主要功能**:
- LLM 统一接口 (llms.Model)
- Prompt 模板引擎 (prompts)
- 工具调用 (Tools)

**使用示例**:

```go
import "github.com/tmc/langchaingo/llms/ollama"
import "github.com/tmc/langchaingo/prompts"

// 1. 初始化 LLM
llm, _ := ollama.New(ollama.WithModel("qwen3:14b"))

// 2. 创建 Prompt 模板
tmpl := prompts.NewPromptTemplate(
    "你是{{.role}}，当前问题：{{.question}}",
    []string{"role", "question"},
)

// 3. 格式化 Prompt
prompt, _ := prompts.FormatPrompt(tmpl, map[string]any{
    "role": "客服",
    "question": "我的订单在哪",
})

// 4. 调用 LLM
response, _ := llm.Call(ctx, prompt)
```

**与其他组件关系**:
- Agent 层依赖 LangChainGo 的 LLM 接口
- 使用 Prompt 模板构建结构化输入
- 未来可扩展使用 Chains、Memory 等功能

---

### 4.3 MCP (Model Context Protocol)

**作用**: 标准化的 AI 工具调用协议，让 LLM 能调用外部工具

**MCP Server 提供**:

| 工具名称 | 功能 | 参数 |
|----------|------|------|
| `get_order_info` | 查询订单详情 | `order_id` (必填) |
| `update_order_address` | 更新收货地址 | `order_id`, `new_address` |
| `get_sop_tree` | 获取 SOP 决策树 | `sop_type` (order/logistics) |

**MCP 工作流程**:

```
Client (Claude Desktop) → MCP Server → Tool → Service → 数据库
                              ↓
                          LLM 决定调用哪个工具
```

**MCP 与其他组件关系**:

```
┌─────────────────────────────────────────────────────────────┐
│  MCP Server                                                   │
│  ├─ 注册 Tools (get_order, update_addr, get_sop)            │
│  ├─ 接收 MCP 客户端请求                                        │
│  └─ 调用 Service 层执行实际操作                               │
└─────────────────────────────────────────────────────────────┘
         ↓                                          ↓
┌──────────────────┐                    ┌──────────────────┐
│ MCP Client       │                    │ Service Layer    │
│ (Claude Desktop)│                    │ (Order, SOP)     │
└──────────────────┘                    └──────────────────┘
```

---

### 4.4 Agent (智能体)

**架构层次**:

```
BaseAgent (基础)
    ↓ 继承
IntentAgent (意图识别)
    ↓ 组合
OrderAgent (订单处理)
```

#### BaseAgent
- 提供统一的 LLM 调用接口
- 所有 Agent 的基类

#### IntentAgent
- 功能：识别用户对话意图
- 输入：用户问题 + 对话历史
- 输出：ORDER / LOGISTICS / UNKNOWN

```go
intent, err := intentAgent.Recognize(ctx, "我的订单在哪？", "")
// 返回: "ORDER"
```

#### OrderAgent
- 功能：处理订单相关问题
- 输入：用户问题 + 对话历史 + 订单信息 + SOP
- 输出：客服回复

```go
response, err := orderAgent.Handle(ctx,
    "什么时候发货？",
    "客户: 我的订单在哪",
    "订单号: 123, 状态: 处理中")
```

---

### 4.5 订单 (Order)

**数据结构**:

```go
type Order struct {
    OrderID      string   // 订单号
    CustomerName string   // 客户姓名
    Items        []string // 商品列表
    Address      string   // 收货地址
    Status       string   // 订单状态
    TotalAmount  float64  // 总金额
    CreateDate   string   // 下单日期
}
```

**OrderService 功能**:

| 方法 | 功能 |
|------|------|
| `GetOrder(id)` | 查询单个订单 |
| `UpdateAddress(id, addr)` | 更新收货地址 |
| `GetAllOrders()` | 获取所有订单 |
| `FormatOrder(order)` | 格式化订单信息 |
| `LoadOrdersFromFile(path)` | 从文件加载订单 |
| `SaveOrdersToFile(path)` | 保存订单到文件 |

**与 Agent 关系**:

```
OrderAgent → OrderService → Orders (内存/文件)
     ↓            ↓
   LLM       CRUD 操作
```

---

### 4.6 SOP (标准作业流程)

**作用**: 确保客服回复符合公司规范和业务流程

**OrderSOP**: 订单问题决策树

```
1. 订单状态查询
   ├─ 我的订单在哪？ → 查询订单状态
   └─ 订单什么时候发货？ → 返回预计发货时间

2. 订单修改
   ├─ 可以修改/删除订单吗？ → 检查订单状态
   ├─ 我想给订单添加商品 → 检查订单状态
   └─ 修改收货地址 → 检查是否已发货

3. 订单取消
   └─ 如何取消订单？ → 检查订单状态

4. 支付问题
   └─ 支付失败怎么办？ → 检查支付状态
```

**LogisticsSOP**: 物流问题决策树

```
1. 包裹超时
   ├─ 超过预计配送时间 >7天无物流更新 → 提供重发或退款
   └─ 超过预计配送时间 3-7天 → 联系物流公司

2. 配送地址
   ├─ 修改配送地址 → 订单未发货时可更新
   └─ 地址填写错误 → 联系派送员

3. 包裹显示已签收但未收到
   └─ 首次客户，地址正确 → 提供重发或退款

4. 物流跟踪
   └─ 查询物流进度 → 提供物流单号和跟踪链接
```

**SOP 与 Agent 关系**:

```
用户提问 → IntentAgent (识别意图) → OrderAgent
                                                ↓
                            将 SOP 作为 Prompt 的一部分
                                                ↓
                                LLM 根据 SOP 生成回复
```

**SOP Prompt 示例**:

```go
prompt := fmt.Sprintf(`你是一个客服AI，严格按照以下决策树回复:

%s

订单信息: %s
客户问题: %s

请根据决策树流程回复客户。`, service.OrderSOP, orderInfo, question)
```

## 五、数据流示例

### 示例 1: 查询订单状态

```
用户: "查询订单123的状态"
    ↓
Client → POST /chat
    ↓
Gin Server → Chat API
    ↓
IntentAgent.Recognize("查询订单123的状态")
    ↓
BaseAgent.LLM.Call("识别意图...")
    ↓
返回: "ORDER"
    ↓
OrderAgent.Handle("查询订单123的状态", "", "订单123信息")
    ↓
BaseAgent.LLM.Call("根据SOP生成回复...")
    ↓
返回: "您的订单123当前状态为处理中..."
    ↓
Response: {"intent": "ORDER", "response": "..."}
```

### 示例 2: MCP 工具调用

```
Claude Desktop → MCP Server
    ↓
Request: {"method": "tools/call", "params": {"name": "get_order_info", "arguments": {"order_id": "123"}}}
    ↓
MCP Server → OrderService.GetOrder("123")
    ↓
返回: {订单信息}
    ↓
MCP Response: {"result": {"order": {...}}}
    ↓
Claude Desktop 使用订单信息继续对话
```

## 六、目录结构

```
fast-fashion-agent/
├── main.go                  # 程序入口
├── go.mod                   # Go 模块定义
├── go.sum                   # 依赖锁定文件
├── test.http                # REST API 测试文件
├── ARCHITECTURE.md          # 本架构文档
└── internal/
    ├── agent/               # Agent 层
    │   ├── base.go          # 基础 Agent
    │   ├── intent.go        # 意图识别 Agent
    │   └── order.go         # 订单处理 Agent
    ├── mcp/                 # MCP Server
    │   └── server.go        # MCP 协议实现
    ├── service/             # Service 层
    │   ├── order.go         # 订单服务
    │   └── sop.go           # SOP 定义
    ├── logger/              # 日志模块
    │   └── logger.go        # Zap 日志封装
    └── middleware/          # 中间件
        └── logger.go        # Gin 日志中间件
```

## 七、环境配置

### 必需环境
- Go 1.24+
- Ollama (本地部署 qwen3:14b 模型)

### 环境变量
```bash
OLLAMA_BASE_URL=http://localhost:11434  # Ollama 服务地址
OLLAMA_MODEL=qwen3:14b                # 使用的模型
SERVER_PORT=8080                        # 服务端口
```

### 安装依赖
```bash
go mod download
```

## 八、运行和测试

### 启动服务
```bash
go run main.go
```

### 测试 API
```bash
# 意图识别测试
curl http://localhost:8080/test

# 聊天对话测试
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{"message":"查询订单123的状态","history":""}'
```

### MCP 配置 (Claude Desktop)

在 Claude Desktop 配置文件中添加:

```json
{
  "mcpServers": {
    "fast-fashion-agent": {
      "command": "node",
      "args": ["-e", "require('http').createServer((req, res) => req.pipe(require('http').request('http://localhost:8080/mcp' + req.url, (r) => r.pipe(res))).listen(3000))"]
    }
  }
}
```

## 九、扩展方向

1. **数据库集成**: 从内存存储迁移到 PostgreSQL/MongoDB
2. **更多 Agent**: 添加 LogisticsAgent、PaymentAgent 等
3. **Memory 机制**: 使用 LangChainGo 的 Memory 实现多轮对话记忆
4. **RAG 集成**: 添加知识库检索增强生成
5. **流式输出**: 实现 SSE/WebSocket 流式响应
6. **链式调用**: 使用 LangChainGo Chains 构建复杂工作流

## 十、参考资料

- [AWS 博客：快时尚电商代理设计思路与应用实践 (第二部分)](https://aws.amazon.com/cn/blogs/china/fast-fashion-e-commerce-agent-design-ideas-and-application-practice-part-two/)
- [LangChainGo 文档](https://github.com/tmc/langchaingo)
- [MCP 规范](https://modelcontextprotocol.io/)
- [Ollama 文档](https://ollama.ai/)

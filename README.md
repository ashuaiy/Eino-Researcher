# Eino-Researcher

基于 Go + Eino 的 Deep Research 多智能体检索增强系统。

[![Go](https://img.shields.io/badge/Go-1.23-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Gin](https://img.shields.io/badge/Gin-1.10-008ECF)](https://gin-gonic.com/)
[![Eino](https://img.shields.io/badge/Eino-planned-orange)](https://github.com/cloudwego/eino)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-pgvector-4169E1?logo=postgresql&logoColor=white)](https://github.com/pgvector/pgvector)
[![Redis](https://img.shields.io/badge/Redis-7-DC382D?logo=redis&logoColor=white)](https://redis.io/)

> 当前状态：V1 RAG 非流式闭环已经实现。系统可以上传 Markdown/txt 文档、生成 Embedding、写入 PostgreSQL + pgvector、执行 TopK 余弦检索，并调用可配置的 OpenAI-compatible 模型生成带来源的回答。Eino Agent、Redis 队列和 SSE 流式任务事件仍属于后续阶段。

## 项目介绍

普通 RAG 系统通常只完成“用户提问 -> 检索文档 -> 生成回答”的单轮问答流程，难以处理需要任务拆解、多轮检索、多来源融合和结构化输出的复杂研究问题。

Eino-Researcher 的目标不是实现一个简单聊天机器人，而是构建一个具备以下能力的 AI 应用后端：

- Agent 编排与多阶段研究工作流
- 本地知识库 RAG 检索
- 外部搜索工具调用
- 多子问题并发检索
- 带来源引用的 Markdown 研究报告
- Redis 异步任务与 Worker
- SSE 中间状态和报告流式输出
- PostgreSQL + pgvector 向量存储
- Docker Compose 一键部署

项目同时用于展示 Go 后端工程能力，包括强类型接口、模块边界、并发控制、异步任务、可观测性和容器化部署。

## 当前完成情况

项目当前处于 **第 2 阶段：V1 RAG 闭环**。

| 模块 | 状态 | 当前实现 |
| --- | --- | --- |
| Go 工程结构 | 已完成 | 已按项目说明建立 `cmd`、`internal`、`migrations`、`docs` 和 `examples` |
| Gin HTTP 服务 | 已完成 | 服务入口、路由、恢复中间件和请求日志已建立 |
| 健康检查 | 已完成 | `GET /health` 返回 `{"status":"ok"}`，包含自动化测试 |
| 环境配置 | 已完成 | 支持应用、PostgreSQL、Redis、LLM 和 Embedding 环境变量 |
| Docker Compose | 已完成 | 包含 API、PostgreSQL + pgvector、Redis 和健康检查 |
| 数据库迁移 | 已完成 | 已建立 `documents`、`chunks`、`research_tasks`、`evidences` |
| 文档上传 API | 已完成 | 校验并读取 `.md` / `.txt`，切分、生成 Embedding，并在事务中写入 document 与 chunks |
| RAG API | 已完成 | 生成 query embedding，执行 TopK 检索，调用模型并返回回答与引用来源 |
| LLM / Embedding | 已完成 | 真实调用可配置的 OpenAI-compatible `/chat/completions` 与 `/embeddings` |
| pgvector 检索 | 已完成 | 使用 HNSW cosine 索引、`<=>` 距离和 `1 - distance` 相关度分数 |
| Agent Orchestrator | 基础占位 | 可创建和查询内存任务，Planner 等 Agent 尚未执行 |
| Planner / Retriever / Reader / Writer / Evaluator | 接口占位 | 接口和 Noop 实现已建立 |
| Redis 队列与 Worker | 接口占位 | Queue 和 Worker 边界已建立，尚未连接 Redis |
| SSE | 基础占位 | 可返回示例 `step` 和 `done` 事件，尚未连接真实任务事件 |
| Web Search | 接口占位 | 已预留工具接口，尚未接入 SearXNG |
| 指标与调用统计 | 未开始 | Prometheus、Token 用量和工具耗时统计尚未实现 |

按里程碑估算：

- **第 1 阶段项目骨架：约 95%**
- **第 2 阶段 RAG 模块：约 80%**
- **第 3 阶段 Agent 工作流：约 10%**
- **第 4 阶段异步任务与流式输出：约 10%**
- **完整项目规划整体：约 45%**

这里的百分比表示工程里程碑完成度，不代表生产可用性。当前版本已经具备真实的非流式知识库问答能力，但尚未包含文档管理、异步摄取、质量评测、Deep Research Agent 或生产级权限隔离。

## 功能规划

### V1：RAG 知识库问答

- 上传 Markdown / txt 文档
- 文档解析与 chunk 切分
- 调用 Embedding 模型生成向量
- 将向量和原文写入 pgvector
- TopK 语义检索
- 将检索结果作为上下文调用 LLM
- 返回回答和引用来源
- 支持流式回答（后续）

### V2：Deep Research Agent

| Agent | 职责 |
| --- | --- |
| Planner Agent | 将复杂研究问题拆分为多个具体、可独立检索的子问题 |
| Retriever Agent | 并发调用本地知识库和可选 Web Search，生成 evidence cards |
| Reader Agent | 清洗、去重、摘要并提取与子问题相关的证据 |
| Writer Agent | 根据研究计划和证据生成结构化 Markdown 报告 |
| Evaluator Agent | 检查子问题覆盖率、报告结构和关键结论的引用支持 |

### V3：工程增强

- Redis 异步任务队列
- 多 Worker 并发执行
- 研究任务状态查询
- SSE 推送规划、检索、写作和报告增量
- Prometheus 指标
- LLM 调用日志、Token 用量和工具耗时统计

## 系统架构

```text
Client / Web UI / API Caller
        |
        v
Gin API Server
        |
        v
Research Service
        |
        +---------------------+
        |                     |
        v                     v
RAG Service              Agent Orchestrator
        |                     |
        v                     v
Vector Store             Planner Agent
PostgreSQL + pgvector         |
                              v
                         Retriever Agent
                              |
                              v
                         Reader Agent
                              |
                              v
                         Writer Agent
                              |
                              v
                         Evaluator Agent
        |
        v
Streaming Response / Markdown Report
```

## Deep Research 工作流

```text
用户提交复杂研究问题
        |
        v
创建 research_task
        |
        v
Planner Agent 生成研究计划和子问题
        |
        v
Retriever Agent 并发检索本地知识库和外部搜索
        |
        v
Reader Agent 整理 evidence cards
        |
        v
Writer Agent 生成带引用的 Markdown 报告
        |
        v
Evaluator Agent 检查完整性与引用支持
        |
        v
保存 report 和 evidences
        |
        v
通过查询 API / SSE 返回任务过程与最终报告
```

并发检索将遵循以下原则：

- 设置最大 goroutine 并发数
- 使用 `context` 设置整体超时和取消
- 单个子问题失败不终止整体任务
- 在最终报告中明确标记证据不足的部分

## 技术栈

| 类型 | 选型 | 说明 |
| --- | --- | --- |
| 后端语言 | Go 1.23 | 强类型、并发友好、部署简单 |
| Web 框架 | Gin | 提供 HTTP API 和中间件 |
| Agent 主框架 | Eino | 计划用于 ChatModel、Tool、Retriever、Agent 和 Workflow |
| 数据库 | PostgreSQL 16 | 保存文档、任务、报告和证据 |
| 向量存储 | pgvector | 与 PostgreSQL 集成，降低部署复杂度 |
| 缓存 / 队列 | Redis 7 | 计划用于异步任务和事件分发 |
| 模型协议 | OpenAI-compatible API | 便于接入 OpenAI、DeepSeek、Doubao、Qwen 和 Ollama |
| 部署 | Docker Compose | 本地开发和演示环境 |

注意：Eino 目前是规划中的主框架，尚未加入 `go.mod`，将在 LLM 和 Agent 阶段正式接入。

## 参考项目

| 类型 | 项目 | 地址 |
| --- | --- | --- |
| 主框架 | Eino | [cloudwego/eino](https://github.com/cloudwego/eino) |
| Eino 扩展 | Eino Ext | [cloudwego/eino-ext](https://github.com/cloudwego/eino-ext) |
| Deep Research 参考 | DeerFlow | [bytedance/deer-flow](https://github.com/bytedance/deer-flow) |
| Go Agent 框架参考 | tRPC-Agent-Go | [trpc-group/trpc-agent-go](https://github.com/trpc-group/trpc-agent-go) |
| Go Agent 框架参考 | Google ADK-Go | [google/adk-go](https://github.com/google/adk-go) |
| Go LLM 编排参考 | LangChainGo | [tmc/langchaingo](https://github.com/tmc/langchaingo) |
| 向量数据库客户端 | Qdrant Go Client | [qdrant/go-client](https://github.com/qdrant/go-client) |
| pgvector Go 支持 | pgvector-go | [pgvector/pgvector-go](https://github.com/pgvector/pgvector-go) |
| 搜索工具 | SearXNG | [searxng/searxng](https://github.com/searxng/searxng) |

## 目录结构

```text
eino-researcher/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   ├── middleware.go
│   │   ├── router.go
│   │   └── router_test.go
│   ├── config/
│   ├── llm/
│   ├── rag/
│   ├── agent/
│   ├── tools/
│   ├── store/
│   ├── queue/
│   ├── model/
│   └── utils/
├── migrations/
│   └── 001_init.sql
├── docs/
│   ├── api.md
│   ├── architecture.md
│   └── demo.md
├── examples/
│   └── sample_docs/
├── docker-compose.yml
├── Dockerfile
├── .env.example
├── go.mod
└── README.md
```

## 数据库设计

当前迁移包含 [`migrations/001_init.sql`](migrations/001_init.sql) 和 [`migrations/002_hnsw_index.sql`](migrations/002_hnsw_index.sql)。

| 表 | 用途 |
| --- | --- |
| `documents` | 文档标题、来源、文件类型和创建时间 |
| `chunks` | 文档切片、Token 数、1536 维向量和文档关联 |
| `research_tasks` | 研究问题、状态、计划、报告和错误信息 |
| `evidences` | 子问题对应的来源、内容、相关度分数和任务关联 |

当前 `embedding` 字段固定为 `VECTOR(1536)`，应用启动时也会强制校验 `EMBEDDING_DIM=1536`。可以更换 OpenAI-compatible Embedding 服务和模型，但新模型必须返回 1536 维；切换模型后必须清理并重建全部文档索引，因为不同模型的向量不可直接比较。

## API

| 方法 | 路径 | 当前状态 |
| --- | --- | --- |
| `GET` | `/health` | 可用 |
| `POST` | `/api/v1/documents` | 可用 |
| `POST` | `/api/v1/rag/query` | 可用（非流式） |
| `POST` | `/api/v1/research/tasks` | 基础占位 |
| `GET` | `/api/v1/research/tasks/{task_id}` | 基础占位 |
| `GET` | `/api/v1/research/tasks/{task_id}/report` | 基础占位 |
| `GET` | `/api/v1/research/tasks/{task_id}/stream` | SSE 示例事件 |

### 健康检查

```bash
curl http://localhost:8080/health
```

```json
{
  "status": "ok"
}
```

### 上传文档

```bash
curl -X POST http://localhost:8080/api/v1/documents \
  -F "title=Eino Overview" \
  -F "file=@examples/sample_docs/eino_overview.md"
```

接口只接受 UTF-8 编码的 `.md` 和 `.txt` 文件，默认最大 2 MiB。原始文件不会写入磁盘；文档元数据、切片内容和向量会在同一 PostgreSQL 事务中写入，Embedding 或数据库任一步失败都会回滚。

成功响应：

```json
{
  "document_id": "7a18d53f-...",
  "title": "Eino Overview",
  "status": "indexed"
}
```

### RAG 问答

```bash
curl -X POST http://localhost:8080/api/v1/rag/query \
  -H "Content-Type: application/json" \
  -d '{"question":"Eino 框架适合开发 Agent 系统吗？","top_k":5,"stream":false}'
```

`top_k` 默认值为 5，最大值为 20。当前 `stream` 字段会被接受，但仍返回普通 JSON；流式回答将在后续 SSE 阶段实现。

响应示例：

```json
{
  "answer": "Eino 是面向 Go 的 LLM 应用开发框架 [1]。",
  "sources": [
    {
      "id": "chunk-uuid",
      "document_id": "document-uuid",
      "source_type": "local_doc",
      "title": "Eino Overview",
      "content": "Eino is a Go framework for building LLM applications.",
      "source": "eino_overview.md",
      "score": 0.91
    }
  ]
}
```

### 创建研究任务

```bash
curl -X POST http://localhost:8080/api/v1/research/tasks \
  -H "Content-Type: application/json" \
  -d '{"question":"Go 语言开发 Deep Research Agent 的优势和局限是什么？","use_web_search":false,"max_sub_questions":5}'
```

任务当前保存在进程内存中，应用重启后会丢失。

### 查询任务与报告

```bash
curl http://localhost:8080/api/v1/research/tasks/{task_id}
curl http://localhost:8080/api/v1/research/tasks/{task_id}/report
```

### SSE 事件

```bash
curl -N http://localhost:8080/api/v1/research/tasks/{task_id}/stream
```

当前只发送示例 `step` 和 `done` 事件。

更多接口说明见 [`docs/api.md`](docs/api.md)。

## 快速启动

### 使用 Docker Compose

Linux / macOS：

```bash
cp .env.example .env
docker compose up --build
```

Windows PowerShell：

```powershell
Copy-Item .env.example .env
docker compose up --build
```

服务启动后访问：

```text
http://localhost:8080/health
```

### 本地运行

本地需要 Go 1.23，并需要单独启动 PostgreSQL 和 Redis：

```bash
go mod download
go run ./cmd/server
```

运行测试：

```bash
go test ./...
```

## 环境变量

```dotenv
APP_ENV=development
APP_PORT=8080
MAX_UPLOAD_BYTES=2097152
MODEL_REQUEST_TIMEOUT_SECONDS=60

POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=eino_researcher
POSTGRES_SSLMODE=disable

REDIS_ADDR=localhost:6379
REDIS_PORT=6379

LLM_BASE_URL=https://api.openai.com/v1
LLM_API_KEY=your_api_key
LLM_MODEL=gpt-4o-mini

EMBEDDING_BASE_URL=https://api.openai.com/v1
EMBEDDING_API_KEY=your_api_key
EMBEDDING_MODEL=text-embedding-3-small
EMBEDDING_DIM=1536
```

完整示例见 [`.env.example`](.env.example)。

### 自定义模型与本地 Ollama

LLM 与 Embedding 可以分别配置为任何实现 OpenAI-compatible API 的服务。API Key 为空时应用不会发送 `Authorization` Header，便于连接本地服务：

```dotenv
LLM_BASE_URL=http://localhost:11434/v1
LLM_API_KEY=
LLM_MODEL=qwen2.5:7b

EMBEDDING_BASE_URL=http://localhost:11434/v1
EMBEDDING_API_KEY=
EMBEDDING_MODEL=your-1536-dimension-embedding-model
EMBEDDING_DIM=1536
```

注意：并非所有 Ollama Embedding 模型都输出 1536 维。当前数据库结构固定为 1536 维，选择模型前需要确认其输出维度。

### 隐私说明

- API Key 仅从环境变量读取，不会写入数据库或日志。
- 应用不会记录模型请求体、模型响应体或上传文档内容。
- 使用远程模型服务时，文档切片会发送给 Embedding 服务，检索片段和问题会发送给 LLM 服务。
- 对隐私敏感的数据应使用受信任的自托管 OpenAI-compatible 服务。

## RAG 与普通问答的区别

普通问答直接让模型根据已有参数知识生成答案，内容可能缺乏可追踪来源。RAG 会先从知识库检索相关文档片段，再将检索结果作为上下文交给模型生成回答。

本项目以 RAG 作为 Deep Research 的证据底座，使研究报告能够：

- 基于本地文档和外部搜索结果
- 保留证据与原始来源
- 为关键结论提供引用
- 明确信息不足或证据冲突
- 降低无来源生成和幻觉风险

## 为什么使用 Go + Eino

Go 适合构建需要高并发检索、异步任务和稳定部署的 Agent 后端：

- goroutine 和 `context` 适合并发子任务、超时和取消控制
- 强类型接口有利于约束 Agent、Tool、Retriever 和 Repository 边界
- 单二进制部署简单，适合容器化和多 Worker 扩展
- 工程工具链统一，便于测试、观测和长期维护

Eino 是项目计划采用的核心 LLM 应用框架。它提供符合 Go 语言习惯的模型、工具、检索器、Agent 和 Workflow 抽象，并可通过 Eino Ext 接入更多模型与组件。

## 开发路线图

### 第 1 阶段：项目骨架

- [x] 创建 Go 项目和推荐目录结构
- [x] 使用 Gin 搭建 HTTP 服务
- [x] 接入环境变量配置和基础日志
- [x] 创建 PostgreSQL + pgvector + Redis Docker Compose
- [x] 创建数据库迁移
- [x] 实现 `/health`
- [x] 建立基础 API、服务接口和 TODO
- [x] 补充 API handler 和核心组件单元测试

### 第 2 阶段：RAG 模块

- [x] 将文档元数据和切片内容持久化到 PostgreSQL
- [x] 支持 Markdown / txt 文档读取
- [x] 实现 chunk 切分策略和测试
- [x] 接入 OpenAI-compatible Embedding
- [x] 实现 pgvector 向量入库
- [x] 实现 TopK 语义检索
- [x] 接入 OpenAI-compatible Chat Completions
- [x] 实现带引用的 RAG 回答
- [ ] 实现回答流式输出
- [ ] 增加文档删除、更新和重建索引
- [ ] 增加检索质量评测集

### 第 3 阶段：Agent 工作流

- [ ] 实现 Planner Agent
- [ ] 实现 Retriever Agent
- [ ] 实现 Reader Agent
- [ ] 实现 Writer Agent
- [ ] 实现 Evaluator Agent
- [ ] 实现受限并发检索和超时控制
- [ ] 完成研究计划、证据和报告持久化

### 第 4 阶段：异步任务与流式输出

- [ ] 实现 Redis 任务队列
- [ ] 实现 Worker 执行研究任务
- [ ] 持久化任务状态
- [ ] 实现真实 SSE 任务事件
- [ ] 支持报告增量输出
- [ ] 添加 Prometheus 指标
- [ ] 添加 Token 用量和工具耗时统计

## 当前限制

- 尚未通过 Eino 组件调用模型，目前使用直接的 OpenAI-compatible HTTP 客户端
- 原始上传文件不会保存，仅保存文档元数据和切片内容
- 文档摄取为同步流程，大文件和高并发场景需要迁移到异步队列
- Embedding 维度固定为 1536，切换模型必须重建索引
- Research Task 使用内存 Map，进程重启后数据丢失
- Agent 均为 Noop / placeholder 实现
- Redis 尚未接入应用运行路径
- SSE 尚未发送真实任务进度或报告增量
- 尚未实现用户、知识库和文档级权限隔离

## 测试

单元测试：

```bash
go test ./...
```

PostgreSQL + pgvector 集成测试复用 Docker Compose：

```bash
docker compose up -d postgres
go test -tags=integration ./internal/api -count=1
```

可以通过 `TEST_DATABASE_URL` 指定单独的测试数据库。集成测试会创建随机 schema，并在结束后清理。

## License

项目暂未添加开源许可证。正式对外发布前需要补充 License 文件。

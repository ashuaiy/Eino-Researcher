# API

Base URL：`http://localhost:8080`

## 健康检查

```bash
curl http://localhost:8080/health
```

```json
{"status":"ok"}
```

## 上传并索引文档

```bash
curl -X POST http://localhost:8080/api/v1/documents \
  -F "title=Eino Overview" \
  -F "file=@examples/sample_docs/eino_overview.md"
```

约束：

- 仅支持 UTF-8 `.md` / `.txt`。
- 默认最大 2 MiB，由 `MAX_UPLOAD_BYTES` 调整。
- 上传、切片、Embedding 和数据库写入同步完成。
- document 和 chunks 在同一事务中写入。

成功响应：

```json
{
  "document_id": "uuid",
  "title": "Eino Overview",
  "status": "indexed"
}
```

常见状态码：

- `400`：缺少文件、类型不支持、空内容或非法 UTF-8。
- `413`：文件超过大小限制。
- `502`：Embedding 模型服务失败。
- `500`：数据库写入失败。

## RAG 问答

```bash
curl -X POST http://localhost:8080/api/v1/rag/query \
  -H "Content-Type: application/json" \
  -d '{"question":"Eino 适合开发 Agent 系统吗？","top_k":5,"stream":false}'
```

`top_k` 默认 5，最大 20。`stream` 当前保留但仍返回普通 JSON。

```json
{
  "answer": "Eino 是面向 Go 的 LLM 应用开发框架 [1]。",
  "sources": [
    {
      "id": "chunk-uuid",
      "document_id": "document-uuid",
      "source_type": "local_doc",
      "title": "Eino Overview",
      "content": "retrieved chunk content",
      "source": "eino_overview.md",
      "score": 0.91
    }
  ]
}
```

回答模型被要求只使用检索上下文，并使用 `[1]`、`[2]` 标记引用。模型服务失败返回 `502`，数据库或检索失败返回 `500`。

## 创建研究任务

以下接口仍为后续 Agent 阶段的占位实现：

```bash
curl -X POST http://localhost:8080/api/v1/research/tasks \
  -H "Content-Type: application/json" \
  -d '{"question":"Go 语言开发 Deep Research Agent 的优势和局限是什么？","use_web_search":false,"max_sub_questions":5}'
```

## 查询任务状态

```bash
curl http://localhost:8080/api/v1/research/tasks/{task_id}
```

## 获取任务报告

```bash
curl http://localhost:8080/api/v1/research/tasks/{task_id}/report
```

## 任务 SSE

```bash
curl -N http://localhost:8080/api/v1/research/tasks/{task_id}/stream
```

当前只返回示例事件，尚未接入真实任务执行过程。

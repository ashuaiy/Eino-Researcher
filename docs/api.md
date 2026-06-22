# API

Base URL: `http://localhost:8080`

## Health

```bash
curl http://localhost:8080/health
```

Response:

```json
{"status":"ok"}
```

## Upload Document

```bash
curl -X POST http://localhost:8080/api/v1/documents \
  -F "title=Eino Overview" \
  -F "file=@examples/sample_docs/eino_overview.md"
```

Current behavior: stores document metadata in memory and returns `accepted`.

## RAG Query

```bash
curl -X POST http://localhost:8080/api/v1/rag/query \
  -H "Content-Type: application/json" \
  -d '{"question":"Eino 适合开发 Agent 系统吗？","top_k":5,"stream":false}'
```

Current behavior: returns placeholder content until Eino, embeddings, and pgvector search are wired.

## Research Task

```bash
curl -X POST http://localhost:8080/api/v1/research/tasks \
  -H "Content-Type: application/json" \
  -d '{"question":"Go 语言开发 Deep Research Agent 的优势和局限是什么？","use_web_search":false,"max_sub_questions":5}'
```

## Task Status

```bash
curl http://localhost:8080/api/v1/research/tasks/{task_id}
```

## Task Report

```bash
curl http://localhost:8080/api/v1/research/tasks/{task_id}/report
```

## Task Stream

```bash
curl -N http://localhost:8080/api/v1/research/tasks/{task_id}/stream
```

# Eino-Researcher

Eino-Researcher is a Go + Eino Deep Research backend skeleton. The goal is to evolve from a basic RAG knowledge-base QA service into a multi-agent research system that can plan sub-questions, retrieve evidence, synthesize reports, evaluate coverage, and stream progress.

This repository currently implements the initial engineering skeleton requested in the project document: Gin API service, configuration, Docker Compose, pgvector schema, base API contracts, service interfaces, and clear TODO boundaries.

## Architecture

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
```

## Tech Stack

- Go
- Gin
- Eino, planned for LLM, tool, retriever, agent, and workflow orchestration
- PostgreSQL + pgvector
- Redis, planned for async task queue and worker execution
- Docker Compose
- OpenAI-compatible LLM and embedding APIs

## Current Features

- `GET /health`
- `POST /api/v1/documents` upload contract with metadata placeholder
- `POST /api/v1/rag/query` RAG query contract with placeholder generation
- `POST /api/v1/research/tasks` task creation placeholder
- `GET /api/v1/research/tasks/{task_id}` task status placeholder
- `GET /api/v1/research/tasks/{task_id}/report` report placeholder
- `GET /api/v1/research/tasks/{task_id}/stream` SSE placeholder
- PostgreSQL migration for `documents`, `chunks`, `research_tasks`, and `evidences`

## Quick Start

Copy the environment file:

```bash
cp .env.example .env
```

Start dependencies and API:

```bash
docker compose up --build
```

Check health:

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{"status":"ok"}
```

Run locally without Docker:

```bash
go mod tidy
go run ./cmd/server
```

## API Examples

Upload a document:

```bash
curl -X POST http://localhost:8080/api/v1/documents \
  -F "title=Eino Overview" \
  -F "file=@examples/sample_docs/eino_overview.md"
```

Ask a RAG question:

```bash
curl -X POST http://localhost:8080/api/v1/rag/query \
  -H "Content-Type: application/json" \
  -d '{"question":"Eino 框架适合开发 Agent 系统吗？","top_k":5,"stream":false}'
```

Create a research task:

```bash
curl -X POST http://localhost:8080/api/v1/research/tasks \
  -H "Content-Type: application/json" \
  -d '{"question":"Go 语言开发 Deep Research Agent 的优势和局限是什么？","use_web_search":false,"max_sub_questions":5}'
```

Stream task events:

```bash
curl -N http://localhost:8080/api/v1/research/tasks/{task_id}/stream
```

## Deep Research Workflow

The intended workflow is:

1. Create a `research_task`.
2. Planner Agent decomposes the question.
3. Retriever Agent searches local knowledge base and optional web sources.
4. Reader Agent cleans and summarizes evidence cards.
5. Writer Agent generates a Markdown report.
6. Evaluator Agent checks coverage and citation support.
7. Report and evidences are persisted.
8. API and SSE stream expose task progress and final report.

## RAG vs Normal QA

Normal QA asks the model to answer directly. RAG first retrieves relevant source chunks, then asks the model to answer with retrieved context. This project uses RAG as the foundation for Deep Research so final reports can include traceable evidence and source references.

## Why Go + Eino

Go gives the backend strong typing, simple concurrency, good deployment ergonomics, and predictable production behavior. Eino is planned as the main LLM application framework because it brings Go-native abstractions for models, tools, retrievers, agents, and workflows while leaving room for production-grade API, queue, and storage layers.

## TODO

- Wire Eino ChatModel and embedding implementations.
- Implement Markdown/txt file persistence and parsing.
- Implement chunking, embedding, and pgvector indexing.
- Implement TopK vector retrieval.
- Implement RAG answer generation with citations.
- Implement Planner, Retriever, Reader, Writer, and Evaluator Agents.
- Add Redis queue and worker execution.
- Replace SSE placeholder with real task events and report deltas.
- Add Prometheus metrics and token/tool timing logs.
- Add demo screenshots or example output after the first full RAG flow works.

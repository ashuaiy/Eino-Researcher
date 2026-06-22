# Architecture

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

## Current Implementation

- `cmd/server`: application entrypoint.
- `internal/api`: Gin router, middleware, and handlers.
- `internal/config`: environment-based configuration.
- `internal/llm`: real OpenAI-compatible chat and embedding HTTP clients.
- `internal/rag`: fixed-size chunking, transactional indexing, pgvector retrieval, and grounded generation.
- `internal/agent`: Deep Research Agent workflow interfaces.
- `internal/tools`: local retrieval, web search, and document reader tool boundaries.
- `internal/store`: pgx-based PostgreSQL repositories and transactional document/chunk persistence.
- `internal/queue`: Redis queue and worker boundaries.
- `migrations`: PostgreSQL + pgvector schema.

## Next Integration Points

1. Add document update/delete and index rebuild workflows.
2. Add retrieval evaluation, hybrid search, and reranking.
3. Replace direct model clients with Eino components where orchestration benefits from it.
4. Add Redis-backed queue and worker lifecycle.
5. Back `/stream` with real task events and report deltas.

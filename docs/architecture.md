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

## Current Skeleton

- `cmd/server`: application entrypoint.
- `internal/api`: Gin router, middleware, and handlers.
- `internal/config`: environment-based configuration.
- `internal/llm`: OpenAI-compatible chat and embedding interfaces.
- `internal/rag`: chunking, retrieval, indexing, and generation interfaces.
- `internal/agent`: Deep Research Agent workflow interfaces.
- `internal/tools`: local retrieval, web search, and document reader tool boundaries.
- `internal/store`: PostgreSQL repository boundaries plus in-memory placeholders.
- `internal/queue`: Redis queue and worker boundaries.
- `migrations`: PostgreSQL + pgvector schema.

## Next Integration Points

1. Replace LLM stubs with Eino ChatModel and embedding components.
2. Replace in-memory repositories with PostgreSQL repositories.
3. Implement pgvector TopK search.
4. Add Redis-backed queue and worker lifecycle.
5. Back `/stream` with real task events and report deltas.

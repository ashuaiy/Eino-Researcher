# Demo

This file records the intended demo flow for the first runnable milestone.

1. Start PostgreSQL, Redis, and API with Docker Compose.
2. Check `GET /health`.
3. Upload a Markdown document.
4. Ask a RAG question.
5. Create a Deep Research task.
6. Query task status.
7. Open the SSE stream.

Current skeleton responses are intentionally placeholder responses. They keep the API contract stable while RAG, Eino Agent orchestration, Redis queueing, and SSE streaming are implemented incrementally.

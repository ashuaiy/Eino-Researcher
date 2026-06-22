package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"eino-researcher/internal/agent"
	"eino-researcher/internal/config"
	"eino-researcher/internal/llm"
	"eino-researcher/internal/rag"
	"eino-researcher/internal/store"
)

func TestHealthReturnsOK(t *testing.T) {
	documents := store.NewInMemoryDocumentRepository()
	ragService := rag.NewService(
		documents,
		llm.NewOpenAICompatibleEmbedder(config.EmbeddingConfig{}),
		llm.NewOpenAICompatibleClient(config.LLMConfig{}),
	)

	router := NewRouter(Dependencies{
		Config:       config.Load(),
		Documents:    documents,
		RAG:          ragService,
		Orchestrator: agent.NewOrchestrator(agent.OrchestratorConfig{RAG: ragService}),
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if rec.Body.String() != `{"status":"ok"}` {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

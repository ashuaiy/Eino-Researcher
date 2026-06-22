package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"eino-researcher/internal/agent"
	"eino-researcher/internal/api"
	"eino-researcher/internal/config"
	"eino-researcher/internal/llm"
	"eino-researcher/internal/rag"
	"eino-researcher/internal/store"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}

func run() error {
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := store.NewPostgresPool(ctx, cfg.Postgres)
	if err != nil {
		return err
	}
	defer pool.Close()

	llmClient := llm.NewOpenAICompatibleClient(cfg.LLM)
	embedder := llm.NewOpenAICompatibleEmbedder(cfg.Embedding)
	documentRepo := store.NewPostgresDocumentRepository(pool)
	chunkRepo := store.NewPostgresChunkRepository(pool, cfg.Embedding.Dim)
	documentChunkStore := store.NewPostgresDocumentChunkStore(pool, cfg.Embedding.Dim)

	indexer := rag.NewPgvectorIndexer(
		rag.NewFixedSizeChunker(1000, 100),
		embedder,
		documentChunkStore,
		cfg.Embedding.Dim,
	)
	retriever := rag.NewPgvectorRetriever(embedder, chunkRepo, cfg.Embedding.Dim)
	generator := rag.NewLLMGenerator(llmClient)
	ragService := rag.NewService(retriever, generator)
	researchOrchestrator := agent.NewOrchestrator(agent.OrchestratorConfig{
		RAG: ragService,
	})

	router := api.NewRouter(api.Dependencies{
		Config:       cfg,
		Documents:    documentRepo,
		Indexer:      indexer,
		RAG:          ragService,
		Orchestrator: researchOrchestrator,
	})

	addr := ":" + cfg.App.Port
	log.Printf("starting Eino-Researcher API on %s", addr)
	server := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}
	errCh := make(chan error, 1)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return server.Shutdown(shutdownCtx)
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}

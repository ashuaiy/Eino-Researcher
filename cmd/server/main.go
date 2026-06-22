package main

import (
	"log"

	"eino-researcher/internal/agent"
	"eino-researcher/internal/api"
	"eino-researcher/internal/config"
	"eino-researcher/internal/llm"
	"eino-researcher/internal/rag"
	"eino-researcher/internal/store"
)

func main() {
	cfg := config.Load()

	llmClient := llm.NewOpenAICompatibleClient(cfg.LLM)
	embedder := llm.NewOpenAICompatibleEmbedder(cfg.Embedding)
	documentRepo := store.NewInMemoryDocumentRepository()

	ragService := rag.NewService(documentRepo, embedder, llmClient)
	researchOrchestrator := agent.NewOrchestrator(agent.OrchestratorConfig{
		RAG: ragService,
	})

	router := api.NewRouter(api.Dependencies{
		Config:       cfg,
		Documents:    documentRepo,
		RAG:          ragService,
		Orchestrator: researchOrchestrator,
	})

	addr := ":" + cfg.App.Port
	log.Printf("starting Eino-Researcher API on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}

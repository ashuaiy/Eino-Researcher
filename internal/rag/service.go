package rag

import (
	"context"

	"eino-researcher/internal/llm"
	"eino-researcher/internal/model"
	"eino-researcher/internal/store"
)

type Service interface {
	Query(ctx context.Context, req QueryRequest) (QueryResponse, error)
}

type QueryRequest struct {
	Question string `json:"question"`
	TopK     int    `json:"top_k"`
	Stream   bool   `json:"stream"`
}

type QueryResponse struct {
	Answer  string           `json:"answer"`
	Sources []model.Evidence `json:"sources"`
}

type BasicService struct {
	documents store.DocumentRepository
	embedder  llm.Embedder
	generator Generator
	retriever Retriever
}

func NewService(documents store.DocumentRepository, embedder llm.Embedder, chat llm.ChatClient) *BasicService {
	return &BasicService{
		documents: documents,
		embedder:  embedder,
		generator: NewLLMGenerator(chat),
		retriever: EmptyRetriever{},
	}
}

func (s *BasicService) Query(ctx context.Context, req QueryRequest) (QueryResponse, error) {
	evidences, err := s.retriever.Search(ctx, req.Question, req.TopK)
	if err != nil {
		return QueryResponse{}, err
	}

	answer, err := s.generator.Generate(ctx, req.Question, evidences)
	if err != nil {
		return QueryResponse{}, err
	}

	return QueryResponse{
		Answer:  answer,
		Sources: evidences,
	}, nil
}

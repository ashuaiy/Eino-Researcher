package rag

import (
	"context"

	"eino-researcher/internal/model"
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
	generator Generator
	retriever Retriever
}

func NewService(retriever Retriever, generator Generator) *BasicService {
	return &BasicService{
		generator: generator,
		retriever: retriever,
	}
}

func (s *BasicService) Query(ctx context.Context, req QueryRequest) (QueryResponse, error) {
	topK := req.TopK
	if topK <= 0 {
		topK = 5
	}
	if topK > 20 {
		topK = 20
	}

	evidences, err := s.retriever.Search(ctx, req.Question, topK)
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

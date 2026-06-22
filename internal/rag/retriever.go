package rag

import (
	"context"

	"eino-researcher/internal/model"
)

type Retriever interface {
	Search(ctx context.Context, question string, topK int) ([]model.Evidence, error)
}

type EmptyRetriever struct{}

func (r EmptyRetriever) Search(ctx context.Context, question string, topK int) ([]model.Evidence, error) {
	// TODO: implement pgvector TopK semantic search.
	return []model.Evidence{}, nil
}

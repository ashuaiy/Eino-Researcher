package rag

import (
	"context"
	"fmt"

	"eino-researcher/internal/llm"
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

type ChunkSearcher interface {
	Search(ctx context.Context, embedding []float32, topK int) ([]model.Evidence, error)
}

type PgvectorRetriever struct {
	embedder llm.Embedder
	chunks   ChunkSearcher
	dim      int
}

func NewPgvectorRetriever(embedder llm.Embedder, chunks ChunkSearcher, dim int) *PgvectorRetriever {
	return &PgvectorRetriever{embedder: embedder, chunks: chunks, dim: dim}
}

func (r *PgvectorRetriever) Search(ctx context.Context, question string, topK int) ([]model.Evidence, error) {
	vector, err := r.embedder.Embed(ctx, question)
	if err != nil {
		return nil, err
	}
	if r.dim > 0 && len(vector) != r.dim {
		return nil, fmt.Errorf("query embedding dimension mismatch: expected %d, got %d", r.dim, len(vector))
	}
	return r.chunks.Search(ctx, vector, topK)
}

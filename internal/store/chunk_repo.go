package store

import (
	"context"

	"eino-researcher/internal/model"
)

type ChunkRepository interface {
	CreateMany(ctx context.Context, chunks []model.Chunk) error
	Search(ctx context.Context, embedding []float32, topK int) ([]model.Evidence, error)
}

type NoopChunkRepository struct{}

func (r NoopChunkRepository) CreateMany(ctx context.Context, chunks []model.Chunk) error {
	// TODO: insert chunk rows with pgvector embeddings.
	return nil
}

func (r NoopChunkRepository) Search(ctx context.Context, embedding []float32, topK int) ([]model.Evidence, error) {
	// TODO: implement vector similarity search with pgvector.
	return []model.Evidence{}, nil
}

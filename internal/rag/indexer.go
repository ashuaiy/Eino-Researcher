package rag

import (
	"context"

	"eino-researcher/internal/model"
)

type Indexer interface {
	Index(ctx context.Context, doc model.Document, content string) error
}

type NoopIndexer struct{}

func (i NoopIndexer) Index(ctx context.Context, doc model.Document, content string) error {
	// TODO: split content, generate embeddings, and write chunks to pgvector.
	return nil
}

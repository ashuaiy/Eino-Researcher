package rag

import (
	"context"
	"fmt"
	"time"

	"eino-researcher/internal/llm"
	"eino-researcher/internal/model"
	"eino-researcher/internal/utils"
)

type Indexer interface {
	Index(ctx context.Context, doc model.Document, content string) error
}

type NoopIndexer struct{}

func (i NoopIndexer) Index(ctx context.Context, doc model.Document, content string) error {
	// TODO: split content, generate embeddings, and write chunks to pgvector.
	return nil
}

type DocumentChunkStore interface {
	Store(ctx context.Context, doc model.Document, chunks []model.Chunk) error
}

type PgvectorIndexer struct {
	chunker  Chunker
	embedder llm.Embedder
	store    DocumentChunkStore
	dim      int
}

func NewPgvectorIndexer(chunker Chunker, embedder llm.Embedder, store DocumentChunkStore, dim int) *PgvectorIndexer {
	return &PgvectorIndexer{
		chunker:  chunker,
		embedder: embedder,
		store:    store,
		dim:      dim,
	}
}

func (i *PgvectorIndexer) Index(ctx context.Context, doc model.Document, content string) error {
	parts := i.chunker.Split(content)
	if len(parts) == 0 {
		return fmt.Errorf("document produced no chunks")
	}

	now := time.Now().UTC()
	chunks := make([]model.Chunk, 0, len(parts))
	for index, part := range parts {
		vector, err := i.embedder.Embed(ctx, part)
		if err != nil {
			return fmt.Errorf("embed chunk %d: %w", index, err)
		}
		if i.dim > 0 && len(vector) != i.dim {
			return fmt.Errorf(
				"chunk %d embedding dimension mismatch: expected %d, got %d",
				index,
				i.dim,
				len(vector),
			)
		}
		chunks = append(chunks, model.Chunk{
			ID:         utils.NewID(),
			DocumentID: doc.ID,
			Content:    part,
			Index:      index,
			TokenCount: 0,
			Embedding:  vector,
			CreatedAt:  now,
		})
	}

	return i.store.Store(ctx, doc, chunks)
}

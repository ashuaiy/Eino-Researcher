package store

import (
	"context"
	"fmt"

	"eino-researcher/internal/model"
)

type PostgresDocumentChunkStore struct {
	db  Beginner
	dim int
}

func NewPostgresDocumentChunkStore(db Beginner, dim int) *PostgresDocumentChunkStore {
	return &PostgresDocumentChunkStore{db: db, dim: dim}
}

func (s *PostgresDocumentChunkStore) Store(ctx context.Context, doc model.Document, chunks []model.Chunk) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin document transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	documents := NewPostgresDocumentRepository(tx)
	if err := documents.Create(ctx, doc); err != nil {
		return err
	}
	chunkRepo := NewPostgresChunkRepository(tx, s.dim)
	if err := chunkRepo.CreateMany(ctx, chunks); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit document transaction: %w", err)
	}
	return nil
}

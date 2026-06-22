package store

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v4"

	"eino-researcher/internal/model"
)

func TestPostgresDocumentChunkStoreUsesSingleTransaction(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("new mock: %v", err)
	}
	defer mock.Close()

	now := time.Now().UTC().Truncate(time.Microsecond)
	doc := model.Document{
		ID:        "00000000-0000-4000-8000-000000000001",
		Title:     "Title",
		Source:    "notes.md",
		FileType:  ".md",
		CreatedAt: now,
	}
	chunk := model.Chunk{
		ID:         "00000000-0000-4000-8000-000000000002",
		DocumentID: doc.ID,
		Content:    "content",
		Index:      0,
		Embedding:  []float32{1, 0},
		CreatedAt:  now,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO documents").
		WithArgs(doc.ID, doc.Title, doc.Source, doc.FileType, doc.CreatedAt).
		WillReturnResult(pgconn.NewCommandTag("INSERT 0 1"))
	mock.ExpectExec("INSERT INTO chunks").
		WithArgs(
			chunk.ID,
			chunk.DocumentID,
			chunk.Content,
			chunk.Index,
			chunk.TokenCount,
			pgxmock.AnyArg(),
			chunk.CreatedAt,
		).
		WillReturnResult(pgconn.NewCommandTag("INSERT 0 1"))
	mock.ExpectCommit()

	store := NewPostgresDocumentChunkStore(mock, 2)
	if err := store.Store(context.Background(), doc, []model.Chunk{chunk}); err != nil {
		t.Fatalf("store: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestPostgresDocumentChunkStoreRollsBackWhenChunkInsertFails(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("new mock: %v", err)
	}
	defer mock.Close()

	now := time.Now().UTC().Truncate(time.Microsecond)
	doc := model.Document{
		ID:        "00000000-0000-4000-8000-000000000001",
		Title:     "Title",
		Source:    "notes.md",
		FileType:  ".md",
		CreatedAt: now,
	}
	chunk := model.Chunk{
		ID:         "00000000-0000-4000-8000-000000000002",
		DocumentID: doc.ID,
		Content:    "content",
		Index:      0,
		Embedding:  []float32{1, 0},
		CreatedAt:  now,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO documents").
		WithArgs(doc.ID, doc.Title, doc.Source, doc.FileType, doc.CreatedAt).
		WillReturnResult(pgconn.NewCommandTag("INSERT 0 1"))
	mock.ExpectExec("INSERT INTO chunks").
		WithArgs(
			chunk.ID,
			chunk.DocumentID,
			chunk.Content,
			chunk.Index,
			chunk.TokenCount,
			pgxmock.AnyArg(),
			chunk.CreatedAt,
		).
		WillReturnError(errors.New("insert failed"))
	mock.ExpectRollback()

	store := NewPostgresDocumentChunkStore(mock, 2)
	if err := store.Store(context.Background(), doc, []model.Chunk{chunk}); err == nil {
		t.Fatal("expected store failure")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

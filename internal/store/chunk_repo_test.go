package store

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v4"

	"eino-researcher/internal/model"
)

func TestPostgresChunkRepositoryCreatesAndSearchesChunks(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("new mock: %v", err)
	}
	defer mock.Close()

	chunk := model.Chunk{
		ID:         "00000000-0000-4000-8000-000000000002",
		DocumentID: "00000000-0000-4000-8000-000000000001",
		Content:    "content",
		Index:      0,
		TokenCount: 0,
		Embedding:  []float32{1, 0},
		CreatedAt:  time.Now().UTC().Truncate(time.Microsecond),
	}
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

	repo := NewPostgresChunkRepository(mock, 2)
	if err := repo.CreateMany(context.Background(), []model.Chunk{chunk}); err != nil {
		t.Fatalf("create chunks: %v", err)
	}

	mock.ExpectQuery("SELECT\\s+c.id::text").
		WithArgs(pgxmock.AnyArg(), 3).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "document_id", "title", "source", "content", "score"}).
				AddRow(chunk.ID, chunk.DocumentID, "Title", "notes.md", chunk.Content, 0.91),
		)
	results, err := repo.Search(context.Background(), []float32{1, 0}, 3)
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected one result, got %#v", results)
	}
	got := results[0]
	if got.ID != chunk.ID || got.DocumentID != chunk.DocumentID || got.SourceType != "local_doc" || got.Score != 0.91 {
		t.Fatalf("unexpected evidence: %#v", got)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestPostgresChunkRepositoryRejectsWrongDimension(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("new mock: %v", err)
	}
	defer mock.Close()

	repo := NewPostgresChunkRepository(mock, 2)
	if err := repo.CreateMany(context.Background(), []model.Chunk{{Embedding: []float32{1}}}); err == nil {
		t.Fatal("expected dimension error")
	}
	if _, err := repo.Search(context.Background(), []float32{1}, 3); err == nil {
		t.Fatal("expected query dimension error")
	}
}

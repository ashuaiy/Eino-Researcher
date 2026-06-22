package store

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v4"

	"eino-researcher/internal/model"
)

func TestPostgresDocumentRepositoryCreateAndGet(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("new mock: %v", err)
	}
	defer mock.Close()

	doc := model.Document{
		ID:        "00000000-0000-4000-8000-000000000001",
		Title:     "Title",
		Source:    "notes.md",
		FileType:  ".md",
		CreatedAt: time.Now().UTC().Truncate(time.Microsecond),
	}
	mock.ExpectExec("INSERT INTO documents").
		WithArgs(doc.ID, doc.Title, doc.Source, doc.FileType, doc.CreatedAt).
		WillReturnResult(pgconn.NewCommandTag("INSERT 0 1"))

	repo := NewPostgresDocumentRepository(mock)
	if err := repo.Create(context.Background(), doc); err != nil {
		t.Fatalf("create: %v", err)
	}

	mock.ExpectQuery("SELECT id::text, title").
		WithArgs(doc.ID).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "title", "source", "file_type", "created_at"}).
				AddRow(doc.ID, doc.Title, doc.Source, doc.FileType, doc.CreatedAt),
		)
	got, err := repo.Get(context.Background(), doc.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got != doc {
		t.Fatalf("expected %#v, got %#v", doc, got)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

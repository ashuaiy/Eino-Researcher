package store

import (
	"context"
	"fmt"
	"sync"

	"eino-researcher/internal/model"
)

type DocumentRepository interface {
	Create(ctx context.Context, doc model.Document) error
	Get(ctx context.Context, id string) (model.Document, error)
	List(ctx context.Context) ([]model.Document, error)
}

type PostgresDocumentRepository struct {
	db DBTX
}

func NewPostgresDocumentRepository(db DBTX) *PostgresDocumentRepository {
	return &PostgresDocumentRepository{db: db}
}

func (r *PostgresDocumentRepository) Create(ctx context.Context, doc model.Document) error {
	_, err := r.db.Exec(
		ctx,
		`INSERT INTO documents (id, title, source, file_type, created_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		doc.ID,
		doc.Title,
		doc.Source,
		doc.FileType,
		doc.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert document: %w", err)
	}
	return nil
}

func (r *PostgresDocumentRepository) Get(ctx context.Context, id string) (model.Document, error) {
	var doc model.Document
	err := r.db.QueryRow(
		ctx,
		`SELECT id::text, title, COALESCE(source, ''), COALESCE(file_type, ''), created_at
		 FROM documents
		 WHERE id = $1`,
		id,
	).Scan(&doc.ID, &doc.Title, &doc.Source, &doc.FileType, &doc.CreatedAt)
	if err != nil {
		return model.Document{}, fmt.Errorf("get document %s: %w", id, err)
	}
	return doc, nil
}

func (r *PostgresDocumentRepository) List(ctx context.Context) ([]model.Document, error) {
	rows, err := r.db.Query(
		ctx,
		`SELECT id::text, title, COALESCE(source, ''), COALESCE(file_type, ''), created_at
		 FROM documents
		 ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list documents: %w", err)
	}
	defer rows.Close()

	var documents []model.Document
	for rows.Next() {
		var doc model.Document
		if err := rows.Scan(&doc.ID, &doc.Title, &doc.Source, &doc.FileType, &doc.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan document: %w", err)
		}
		documents = append(documents, doc)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate documents: %w", err)
	}
	return documents, nil
}

type InMemoryDocumentRepository struct {
	mu        sync.RWMutex
	documents map[string]model.Document
}

func NewInMemoryDocumentRepository() *InMemoryDocumentRepository {
	return &InMemoryDocumentRepository{documents: make(map[string]model.Document)}
}

func (r *InMemoryDocumentRepository) Create(ctx context.Context, doc model.Document) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.documents[doc.ID] = doc
	return nil
}

func (r *InMemoryDocumentRepository) Get(ctx context.Context, id string) (model.Document, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	doc, ok := r.documents[id]
	if !ok {
		return model.Document{}, fmt.Errorf("document %s not found", id)
	}
	return doc, nil
}

func (r *InMemoryDocumentRepository) List(ctx context.Context) ([]model.Document, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	docs := make([]model.Document, 0, len(r.documents))
	for _, doc := range r.documents {
		docs = append(docs, doc)
	}
	return docs, nil
}

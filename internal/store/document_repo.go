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

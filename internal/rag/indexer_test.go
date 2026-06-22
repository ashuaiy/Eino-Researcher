package rag

import (
	"context"
	"errors"
	"testing"

	"eino-researcher/internal/model"
)

type sequenceEmbedder struct {
	vectors [][]float32
	errAt   int
	calls   int
}

func (e *sequenceEmbedder) Embed(ctx context.Context, input string) ([]float32, error) {
	e.calls++
	if e.errAt > 0 && e.calls == e.errAt {
		return nil, errors.New("embed failed")
	}
	return e.vectors[e.calls-1], nil
}

type capturingDocumentChunkStore struct {
	doc    model.Document
	chunks []model.Chunk
	calls  int
}

func (s *capturingDocumentChunkStore) Store(ctx context.Context, doc model.Document, chunks []model.Chunk) error {
	s.calls++
	s.doc = doc
	s.chunks = chunks
	return nil
}

func TestPgvectorIndexerEmbedsAllChunksBeforeAtomicStore(t *testing.T) {
	embedder := &sequenceEmbedder{
		vectors: [][]float32{{1, 0}, {0, 1}},
	}
	store := &capturingDocumentChunkStore{}
	indexer := NewPgvectorIndexer(NewFixedSizeChunker(4, 0), embedder, store, 2)
	doc := model.NewDocument("title", "notes.md", ".md")

	if err := indexer.Index(context.Background(), doc, "甲乙丙丁戊己"); err != nil {
		t.Fatalf("index: %v", err)
	}
	if store.calls != 1 || store.doc.ID != doc.ID || len(store.chunks) != 2 {
		t.Fatalf("unexpected store call: calls=%d doc=%#v chunks=%#v", store.calls, store.doc, store.chunks)
	}
	for i, chunk := range store.chunks {
		if chunk.DocumentID != doc.ID || chunk.Index != i || len(chunk.Embedding) != 2 || chunk.ID == "" {
			t.Fatalf("unexpected chunk %d: %#v", i, chunk)
		}
	}
}

func TestPgvectorIndexerDoesNotStorePartialEmbeddings(t *testing.T) {
	embedder := &sequenceEmbedder{
		vectors: [][]float32{{1, 0}, {0, 1}},
		errAt:   2,
	}
	store := &capturingDocumentChunkStore{}
	indexer := NewPgvectorIndexer(NewFixedSizeChunker(4, 0), embedder, store, 2)

	err := indexer.Index(context.Background(), model.NewDocument("title", "notes.md", ".md"), "甲乙丙丁戊己")
	if err == nil {
		t.Fatal("expected embedding failure")
	}
	if store.calls != 0 {
		t.Fatalf("expected no store call, got %d", store.calls)
	}
}

func TestPgvectorIndexerRejectsWrongEmbeddingDimension(t *testing.T) {
	embedder := &sequenceEmbedder{vectors: [][]float32{{1}}}
	store := &capturingDocumentChunkStore{}
	indexer := NewPgvectorIndexer(NewFixedSizeChunker(100, 0), embedder, store, 2)

	err := indexer.Index(context.Background(), model.NewDocument("title", "notes.md", ".md"), "content")
	if err == nil {
		t.Fatal("expected dimension error")
	}
	if store.calls != 0 {
		t.Fatalf("expected no store call, got %d", store.calls)
	}
}

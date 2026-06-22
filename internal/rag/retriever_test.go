package rag

import (
	"context"
	"reflect"
	"testing"

	"eino-researcher/internal/model"
)

type staticEmbedder struct {
	input  string
	vector []float32
	err    error
}

func (e *staticEmbedder) Embed(ctx context.Context, input string) ([]float32, error) {
	e.input = input
	return e.vector, e.err
}

type capturingChunkSearcher struct {
	vector []float32
	topK   int
	out    []model.Evidence
}

func (s *capturingChunkSearcher) Search(ctx context.Context, embedding []float32, topK int) ([]model.Evidence, error) {
	s.vector = embedding
	s.topK = topK
	return s.out, nil
}

func TestPgvectorRetrieverEmbedsQuestionAndSearchesChunks(t *testing.T) {
	embedder := &staticEmbedder{vector: []float32{0.1, 0.2}}
	searcher := &capturingChunkSearcher{out: []model.Evidence{{ID: "chunk-1"}}}
	retriever := NewPgvectorRetriever(embedder, searcher, 2)

	got, err := retriever.Search(context.Background(), "question", 3)
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if embedder.input != "question" {
		t.Fatalf("unexpected embedding input: %q", embedder.input)
	}
	if !reflect.DeepEqual(searcher.vector, embedder.vector) || searcher.topK != 3 {
		t.Fatalf("unexpected search call: vector=%v topK=%d", searcher.vector, searcher.topK)
	}
	if len(got) != 1 || got[0].ID != "chunk-1" {
		t.Fatalf("unexpected results: %#v", got)
	}
}

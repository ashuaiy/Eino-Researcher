package rag

import (
	"context"
	"testing"

	"eino-researcher/internal/model"
)

type capturingRetriever struct {
	topK int
	out  []model.Evidence
}

func (r *capturingRetriever) Search(ctx context.Context, question string, topK int) ([]model.Evidence, error) {
	r.topK = topK
	return r.out, nil
}

type staticGenerator struct {
	answer string
}

func (g staticGenerator) Generate(ctx context.Context, question string, evidences []model.Evidence) (string, error) {
	return g.answer, nil
}

func TestServiceClampsTopKAndReturnsSources(t *testing.T) {
	sources := []model.Evidence{{ID: "chunk-1", DocumentID: "doc-1"}}
	retriever := &capturingRetriever{out: sources}
	service := NewService(retriever, staticGenerator{answer: "answer"})

	resp, err := service.Query(context.Background(), QueryRequest{
		Question: "question",
		TopK:     100,
	})
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if retriever.topK != 20 {
		t.Fatalf("expected topK 20, got %d", retriever.topK)
	}
	if len(resp.Sources) != 1 || resp.Sources[0].DocumentID != "doc-1" {
		t.Fatalf("unexpected sources: %#v", resp.Sources)
	}
}

func TestServiceUsesDefaultTopK(t *testing.T) {
	retriever := &capturingRetriever{}
	service := NewService(retriever, staticGenerator{})

	if _, err := service.Query(context.Background(), QueryRequest{Question: "question"}); err != nil {
		t.Fatalf("query: %v", err)
	}
	if retriever.topK != 5 {
		t.Fatalf("expected default topK 5, got %d", retriever.topK)
	}
}

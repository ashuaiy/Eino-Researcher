package rag

import (
	"context"
	"strings"
	"testing"

	"eino-researcher/internal/llm"
	"eino-researcher/internal/model"
)

type capturingChatClient struct {
	req llm.ChatRequest
}

func (c *capturingChatClient) Generate(ctx context.Context, req llm.ChatRequest) (llm.ChatResponse, error) {
	c.req = req
	return llm.ChatResponse{Content: "answer [1]"}, nil
}

func TestLLMGeneratorBuildsGroundedPromptWithNumberedSources(t *testing.T) {
	client := &capturingChatClient{}
	generator := NewLLMGenerator(client)
	evidences := []model.Evidence{{
		ID:         "chunk-1",
		DocumentID: "doc-1",
		Title:      "Eino Overview",
		Source:     "eino.md",
		Content:    "Eino is a Go framework.",
		Score:      0.91,
	}}

	answer, err := generator.Generate(context.Background(), "What is Eino?", evidences)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if answer != "answer [1]" {
		t.Fatalf("unexpected answer: %q", answer)
	}
	for _, expected := range []string{
		"untrusted reference data",
		"only",
		"[1]",
		"Eino Overview",
		"eino.md",
		"Eino is a Go framework.",
		"What is Eino?",
	} {
		if !strings.Contains(client.req.Prompt, expected) && !strings.Contains(client.req.SystemPrompt, expected) {
			t.Fatalf("prompt missing %q: system=%q prompt=%q", expected, client.req.SystemPrompt, client.req.Prompt)
		}
	}
}

func TestLLMGeneratorExplainsInsufficientEvidenceWhenNoSources(t *testing.T) {
	client := &capturingChatClient{}
	generator := NewLLMGenerator(client)

	if _, err := generator.Generate(context.Background(), "unknown", nil); err != nil {
		t.Fatalf("generate: %v", err)
	}
	if !strings.Contains(client.req.Prompt, "insufficient") {
		t.Fatalf("expected insufficient-evidence instruction, got %q", client.req.Prompt)
	}
}

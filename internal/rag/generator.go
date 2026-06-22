package rag

import (
	"context"

	"eino-researcher/internal/llm"
	"eino-researcher/internal/model"
)

type Generator interface {
	Generate(ctx context.Context, question string, evidences []model.Evidence) (string, error)
}

type LLMGenerator struct {
	client llm.ChatClient
}

func NewLLMGenerator(client llm.ChatClient) LLMGenerator {
	return LLMGenerator{client: client}
}

func (g LLMGenerator) Generate(ctx context.Context, question string, evidences []model.Evidence) (string, error) {
	resp, err := g.client.Generate(ctx, llm.ChatRequest{
		SystemPrompt: "You are a RAG answer generator. Use only the supplied context.",
		Prompt:       question,
	})
	if err != nil {
		return "", err
	}
	return resp.Content, nil
}

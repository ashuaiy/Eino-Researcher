package llm

import (
	"context"
	"fmt"

	"eino-researcher/internal/config"
)

type ChatClient interface {
	Generate(ctx context.Context, req ChatRequest) (ChatResponse, error)
}

type OpenAICompatibleClient struct {
	cfg config.LLMConfig
}

func NewOpenAICompatibleClient(cfg config.LLMConfig) *OpenAICompatibleClient {
	return &OpenAICompatibleClient{cfg: cfg}
}

func (c *OpenAICompatibleClient) Generate(ctx context.Context, req ChatRequest) (ChatResponse, error) {
	// TODO: wire Eino ChatModel with OpenAI-compatible endpoint.
	if req.Prompt == "" {
		return ChatResponse{}, fmt.Errorf("prompt is required")
	}
	return ChatResponse{
		Content: "TODO: connect Eino ChatModel and generate an answer.",
		Model:   c.cfg.Model,
	}, nil
}

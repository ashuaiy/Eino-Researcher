package llm

import (
	"context"
	"fmt"

	"eino-researcher/internal/config"
)

type Embedder interface {
	Embed(ctx context.Context, input string) ([]float32, error)
}

type OpenAICompatibleEmbedder struct {
	cfg config.EmbeddingConfig
}

func NewOpenAICompatibleEmbedder(cfg config.EmbeddingConfig) *OpenAICompatibleEmbedder {
	return &OpenAICompatibleEmbedder{cfg: cfg}
}

func (e *OpenAICompatibleEmbedder) Embed(ctx context.Context, input string) ([]float32, error) {
	// TODO: wire Eino embedding model and call the configured provider.
	if input == "" {
		return nil, fmt.Errorf("embedding input is required")
	}
	dim := e.cfg.Dim
	if dim <= 0 {
		dim = 1536
	}
	return make([]float32, dim), nil
}

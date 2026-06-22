package tools

import (
	"context"

	"eino-researcher/internal/model"
)

type WebSearchTool interface {
	Search(ctx context.Context, query string, topK int) ([]model.Evidence, error)
}

type NoopWebSearchTool struct{}

func (t NoopWebSearchTool) Search(ctx context.Context, query string, topK int) ([]model.Evidence, error) {
	// TODO: integrate SearXNG or another search provider.
	return []model.Evidence{}, nil
}

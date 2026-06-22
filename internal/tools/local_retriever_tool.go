package tools

import (
	"context"

	"eino-researcher/internal/model"
	"eino-researcher/internal/rag"
)

type LocalRetrieverTool struct {
	retriever rag.Retriever
}

func NewLocalRetrieverTool(retriever rag.Retriever) LocalRetrieverTool {
	return LocalRetrieverTool{retriever: retriever}
}

func (t LocalRetrieverTool) Search(ctx context.Context, question string, topK int) ([]model.Evidence, error) {
	// TODO: expose this as an Eino tool for Retriever Agent.
	return t.retriever.Search(ctx, question, topK)
}

package agent

import (
	"context"

	"eino-researcher/internal/model"
)

type RetrieveResult struct {
	SubQuestion string
	Evidences   []model.Evidence
	Err         error
}

type RetrieverAgent interface {
	Retrieve(ctx context.Context, subQuestion string) RetrieveResult
}

type NoopRetrieverAgent struct{}

func (a NoopRetrieverAgent) Retrieve(ctx context.Context, subQuestion string) RetrieveResult {
	// TODO: combine local retriever tool and optional web search tool.
	return RetrieveResult{SubQuestion: subQuestion, Evidences: []model.Evidence{}}
}

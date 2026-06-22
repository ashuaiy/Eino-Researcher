package agent

import (
	"context"

	"eino-researcher/internal/model"
)

type Writer interface {
	Write(ctx context.Context, question string, plan ResearchPlan, evidences []model.Evidence) (string, error)
}

type NoopWriter struct{}

func (w NoopWriter) Write(ctx context.Context, question string, plan ResearchPlan, evidences []model.Evidence) (string, error) {
	// TODO: generate Markdown research report with citations.
	return "# Research Report\n\nTODO: Writer Agent is not implemented yet.\n", nil
}

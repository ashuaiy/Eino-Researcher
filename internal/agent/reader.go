package agent

import (
	"context"

	"eino-researcher/internal/model"
)

type Reader interface {
	Read(ctx context.Context, subQuestion string, evidences []model.Evidence) ([]model.Evidence, error)
}

type NoopReader struct{}

func (r NoopReader) Read(ctx context.Context, subQuestion string, evidences []model.Evidence) ([]model.Evidence, error) {
	// TODO: clean, deduplicate, and summarize evidence cards.
	return evidences, nil
}

package store

import (
	"context"

	"eino-researcher/internal/model"
)

type TaskRepository interface {
	Create(ctx context.Context, task model.ResearchTask) error
	Get(ctx context.Context, id string) (model.ResearchTask, error)
	Update(ctx context.Context, task model.ResearchTask) error
}

type NoopTaskRepository struct{}

func (r NoopTaskRepository) Create(ctx context.Context, task model.ResearchTask) error {
	// TODO: persist research task to PostgreSQL.
	return nil
}

func (r NoopTaskRepository) Get(ctx context.Context, id string) (model.ResearchTask, error) {
	// TODO: load research task from PostgreSQL.
	return model.ResearchTask{}, nil
}

func (r NoopTaskRepository) Update(ctx context.Context, task model.ResearchTask) error {
	// TODO: update task status, plan, report, and error fields.
	return nil
}

package queue

import "context"

type ResearchJob struct {
	TaskID string `json:"task_id"`
}

type Queue interface {
	Enqueue(ctx context.Context, job ResearchJob) error
	Dequeue(ctx context.Context) (ResearchJob, error)
}

type NoopQueue struct{}

func (q NoopQueue) Enqueue(ctx context.Context, job ResearchJob) error {
	// TODO: push research jobs to Redis.
	return nil
}

func (q NoopQueue) Dequeue(ctx context.Context) (ResearchJob, error) {
	// TODO: block/pop research jobs from Redis.
	return ResearchJob{}, nil
}

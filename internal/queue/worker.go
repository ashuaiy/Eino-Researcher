package queue

import "context"

type Worker struct {
	queue Queue
}

func NewWorker(queue Queue) Worker {
	return Worker{queue: queue}
}

func (w Worker) Run(ctx context.Context) error {
	// TODO: consume Redis jobs and run Agent Orchestrator.
	<-ctx.Done()
	return ctx.Err()
}

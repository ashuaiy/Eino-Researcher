package agent

import (
	"context"
	"fmt"
	"sync"
	"time"

	"eino-researcher/internal/model"
	"eino-researcher/internal/rag"
	"eino-researcher/internal/utils"
)

type Orchestrator interface {
	CreateTask(ctx context.Context, req CreateTaskRequest) (model.ResearchTask, error)
	GetTask(ctx context.Context, id string) (model.ResearchTask, error)
	GetReport(ctx context.Context, id string) (ReportResponse, error)
}

type OrchestratorConfig struct {
	RAG rag.Service
}

type CreateTaskRequest struct {
	Question        string
	UseWebSearch    bool
	MaxSubQuestions int
}

type ReportResponse struct {
	TaskID    string           `json:"task_id"`
	Status    string           `json:"status"`
	Report    string           `json:"report"`
	Evidences []model.Evidence `json:"evidences"`
}

type BasicOrchestrator struct {
	cfg   OrchestratorConfig
	mu    sync.RWMutex
	tasks map[string]model.ResearchTask
}

func NewOrchestrator(cfg OrchestratorConfig) *BasicOrchestrator {
	return &BasicOrchestrator{
		cfg:   cfg,
		tasks: make(map[string]model.ResearchTask),
	}
}

func (o *BasicOrchestrator) CreateTask(ctx context.Context, req CreateTaskRequest) (model.ResearchTask, error) {
	now := time.Now().UTC()
	task := model.ResearchTask{
		ID:          utils.NewID(),
		Question:    req.Question,
		Status:      model.TaskStatusPending,
		CurrentStep: "queued",
		Plan: map[string]any{
			"research_question": req.Question,
			"sub_questions":     []string{},
			"todo":              "Planner Agent not implemented yet",
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	o.mu.Lock()
	defer o.mu.Unlock()
	o.tasks[task.ID] = task
	return task, nil
}

func (o *BasicOrchestrator) GetTask(ctx context.Context, id string) (model.ResearchTask, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	task, ok := o.tasks[id]
	if !ok {
		return model.ResearchTask{}, fmt.Errorf("task %s not found", id)
	}
	return task, nil
}

func (o *BasicOrchestrator) GetReport(ctx context.Context, id string) (ReportResponse, error) {
	task, err := o.GetTask(ctx, id)
	if err != nil {
		return ReportResponse{}, err
	}

	return ReportResponse{
		TaskID:    task.ID,
		Status:    task.Status,
		Report:    task.Report,
		Evidences: []model.Evidence{},
	}, nil
}

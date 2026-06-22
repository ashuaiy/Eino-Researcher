package agent

import "context"

type ResearchPlan struct {
	ResearchQuestion string   `json:"research_question"`
	SubQuestions     []string `json:"sub_questions"`
}

type Planner interface {
	Plan(ctx context.Context, question string, maxSubQuestions int) (ResearchPlan, error)
}

type NoopPlanner struct{}

func (p NoopPlanner) Plan(ctx context.Context, question string, maxSubQuestions int) (ResearchPlan, error) {
	// TODO: implement Planner Agent using Eino graph/workflow primitives.
	return ResearchPlan{
		ResearchQuestion: question,
		SubQuestions:     []string{},
	}, nil
}

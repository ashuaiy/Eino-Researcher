package agent

import "context"

type EvaluationResult struct {
	Passed              bool     `json:"passed"`
	MissingSubQuestions []string `json:"missing_sub_questions"`
	UnsupportedClaims   []string `json:"unsupported_claims"`
	Suggestions         []string `json:"suggestions"`
}

type Evaluator interface {
	Evaluate(ctx context.Context, question string, plan ResearchPlan, report string) (EvaluationResult, error)
}

type NoopEvaluator struct{}

func (e NoopEvaluator) Evaluate(ctx context.Context, question string, plan ResearchPlan, report string) (EvaluationResult, error) {
	// TODO: validate report coverage and citation support.
	return EvaluationResult{Passed: false, Suggestions: []string{"Evaluator Agent is not implemented yet."}}, nil
}

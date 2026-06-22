package model

import "time"

type Evidence struct {
	ID          string    `json:"id,omitempty"`
	DocumentID  string    `json:"document_id,omitempty"`
	TaskID      string    `json:"task_id,omitempty"`
	SubQuestion string    `json:"sub_question"`
	SourceType  string    `json:"source_type"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Source      string    `json:"source"`
	Score       float64   `json:"score"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}

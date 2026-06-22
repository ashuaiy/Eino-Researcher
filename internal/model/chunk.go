package model

import "time"

type Chunk struct {
	ID         string    `json:"id"`
	DocumentID string    `json:"document_id"`
	Content    string    `json:"content"`
	Index      int       `json:"chunk_index"`
	TokenCount int       `json:"token_count"`
	Embedding  []float32 `json:"embedding,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

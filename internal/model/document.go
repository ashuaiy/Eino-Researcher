package model

import (
	"time"

	"eino-researcher/internal/utils"
)

type Document struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Source    string    `json:"source"`
	FileType  string    `json:"file_type"`
	CreatedAt time.Time `json:"created_at"`
}

func NewDocument(title, source, fileType string) Document {
	return Document{
		ID:        utils.NewID(),
		Title:     title,
		Source:    source,
		FileType:  fileType,
		CreatedAt: time.Now().UTC(),
	}
}

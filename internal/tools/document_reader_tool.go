package tools

import (
	"context"
	"fmt"
	"path/filepath"
)

type DocumentReaderTool interface {
	Read(ctx context.Context, path string) (string, error)
}

type MarkdownTextReader struct{}

func (r MarkdownTextReader) Read(ctx context.Context, path string) (string, error) {
	// TODO: read Markdown and txt content from persisted uploads.
	ext := filepath.Ext(path)
	if ext != ".md" && ext != ".txt" {
		return "", fmt.Errorf("unsupported document type: %s", ext)
	}
	return "", nil
}

package handlers

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"eino-researcher/internal/llm"
	"eino-researcher/internal/model"
)

type capturingIndexer struct {
	doc     model.Document
	content string
	err     error
}

func (i *capturingIndexer) Index(ctx context.Context, doc model.Document, content string) error {
	i.doc = doc
	i.content = content
	return i.err
}

func newUploadRequest(t *testing.T, filename string, content []byte) *http.Request {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := part.Write(content); err != nil {
		t.Fatalf("write form file: %v", err)
	}
	_ = writer.WriteField("title", "Test Document")
	if err := writer.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/documents", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func serveUpload(t *testing.T, indexer *capturingIndexer, maxBytes int64, req *http.Request) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/v1/documents", NewDocumentHandler(indexer, maxBytes).Upload)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func TestDocumentUploadIndexesMarkdown(t *testing.T) {
	indexer := &capturingIndexer{}
	req := newUploadRequest(t, "notes.MD", []byte("# Hello\ncontent"))
	rec := serveUpload(t, indexer, 1024, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
	if indexer.content != "# Hello\ncontent" || indexer.doc.FileType != ".md" {
		t.Fatalf("unexpected indexed document: %#v content=%q", indexer.doc, indexer.content)
	}
	if !strings.Contains(rec.Body.String(), `"status":"indexed"`) {
		t.Fatalf("unexpected response: %s", rec.Body.String())
	}
}

func TestDocumentUploadValidation(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		content  []byte
		maxBytes int64
		status   int
	}{
		{name: "unsupported extension", filename: "notes.pdf", content: []byte("x"), maxBytes: 10, status: http.StatusBadRequest},
		{name: "empty content", filename: "notes.txt", content: []byte(" \n"), maxBytes: 10, status: http.StatusBadRequest},
		{name: "invalid utf8", filename: "notes.txt", content: []byte{0xff}, maxBytes: 10, status: http.StatusBadRequest},
		{name: "too large", filename: "notes.txt", content: []byte("12345"), maxBytes: 4, status: http.StatusRequestEntityTooLarge},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := serveUpload(t, &capturingIndexer{}, tt.maxBytes, newUploadRequest(t, tt.filename, tt.content))
			if rec.Code != tt.status {
				t.Fatalf("expected %d, got %d: %s", tt.status, rec.Code, rec.Body.String())
			}
		})
	}
}

func TestDocumentUploadMapsProviderAndDatabaseErrors(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		status int
	}{
		{name: "provider", err: llm.ErrProvider, status: http.StatusBadGateway},
		{name: "database", err: errors.New("database unavailable"), status: http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := serveUpload(t, &capturingIndexer{err: tt.err}, 1024, newUploadRequest(t, "notes.txt", []byte("content")))
			if rec.Code != tt.status {
				t.Fatalf("expected %d, got %d: %s", tt.status, rec.Code, rec.Body.String())
			}
		})
	}
}

package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"eino-researcher/internal/llm"
	"eino-researcher/internal/rag"
)

type failingRAGService struct {
	err error
}

func (s failingRAGService) Query(ctx context.Context, req rag.QueryRequest) (rag.QueryResponse, error) {
	return rag.QueryResponse{}, s.err
}

func TestRAGHandlerMapsProviderErrorsToBadGateway(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/v1/rag/query", NewRAGHandler(failingRAGService{err: llm.ErrProvider}).Query)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/rag/query",
		strings.NewReader(`{"question":"question","top_k":5}`),
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadGateway {
		t.Fatalf("expected 502, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestRAGHandlerRejectsWhitespaceQuestion(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/v1/rag/query", NewRAGHandler(failingRAGService{}).Query)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/rag/query",
		strings.NewReader(`{"question":"   ","top_k":5}`),
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", rec.Code, rec.Body.String())
	}
}

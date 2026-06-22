//go:build integration

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxvec "github.com/pgvector/pgvector-go/pgx"

	"eino-researcher/internal/agent"
	"eino-researcher/internal/config"
	"eino-researcher/internal/llm"
	"eino-researcher/internal/rag"
	"eino-researcher/internal/store"
)

func TestRAGUploadAndQueryWithPgvector(t *testing.T) {
	ctx := context.Background()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/eino_researcher?sslmode=disable"
	}

	var admin *pgxpool.Pool
	var err error
	if os.Getenv("TEST_DATABASE_URL") != "" {
		admin, err = pgxpool.New(ctx, dsn)
	} else {
		admin, err = store.NewPostgresPool(ctx, config.PostgresConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "postgres",
			Password: "postgres",
			Database: "eino_researcher",
			SSLMode:  "disable",
		})
	}
	if err != nil {
		t.Fatalf("connect integration database: %v", err)
	}
	defer admin.Close()
	if _, err := admin.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS vector"); err != nil {
		t.Fatalf("create vector extension: %v", err)
	}
	if _, err := admin.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS pgcrypto"); err != nil {
		t.Fatalf("create pgcrypto extension: %v", err)
	}

	schema := fmt.Sprintf("rag_test_%d", time.Now().UnixNano())
	quotedSchema := pgx.Identifier{schema}.Sanitize()
	if _, err := admin.Exec(ctx, "CREATE SCHEMA "+quotedSchema); err != nil {
		t.Fatalf("create test schema: %v", err)
	}
	defer func() {
		_, _ = admin.Exec(ctx, "DROP SCHEMA "+quotedSchema+" CASCADE")
	}()

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatalf("parse test database URL: %v", err)
	}
	poolConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		if err := pgxvec.RegisterTypes(ctx, conn); err != nil {
			return err
		}
		_, err := conn.Exec(ctx, "SET search_path TO "+quotedSchema+", public")
		return err
	}
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		t.Fatalf("create test pool: %v", err)
	}
	defer pool.Close()

	applyMigrations(t, ctx, pool)

	var chatPrompt string
	modelServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/embeddings":
			vector := make([]float32, 1536)
			vector[0] = 1
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": []any{map[string]any{"embedding": vector}},
			})
		case "/v1/chat/completions":
			var body struct {
				Messages []llm.Message `json:"messages"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode chat request: %v", err)
			}
			chatPrompt = body.Messages[len(body.Messages)-1].Content
			_ = json.NewEncoder(w).Encode(map[string]any{
				"model": "test-chat",
				"choices": []any{
					map[string]any{"message": map[string]any{"content": "Eino is a Go framework [1]."}},
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer modelServer.Close()

	modelTimeout := 5 * time.Second
	embedder := llm.NewOpenAICompatibleEmbedder(config.EmbeddingConfig{
		BaseURL: modelServer.URL + "/v1",
		Model:   "test-embedding",
		Dim:     1536,
		Timeout: modelTimeout,
	})
	chatClient := llm.NewOpenAICompatibleClient(config.LLMConfig{
		BaseURL: modelServer.URL + "/v1",
		Model:   "test-chat",
		Timeout: modelTimeout,
	})
	documents := store.NewPostgresDocumentRepository(pool)
	chunks := store.NewPostgresChunkRepository(pool, 1536)
	documentStore := store.NewPostgresDocumentChunkStore(pool, 1536)
	indexer := rag.NewPgvectorIndexer(rag.NewFixedSizeChunker(1000, 100), embedder, documentStore, 1536)
	retriever := rag.NewPgvectorRetriever(embedder, chunks, 1536)
	ragService := rag.NewService(retriever, rag.NewLLMGenerator(chatClient))
	orchestrator := agent.NewOrchestrator(agent.OrchestratorConfig{RAG: ragService})
	router := NewRouter(Dependencies{
		Config: config.Config{
			App:       config.AppConfig{MaxUploadBytes: 2 * 1024 * 1024},
			Embedding: config.EmbeddingConfig{Dim: 1536},
		},
		Documents:    documents,
		Indexer:      indexer,
		RAG:          ragService,
		Orchestrator: orchestrator,
	})

	upload := makeMultipartUpload(t, "eino.md", "# Eino\nEino is a Go framework for LLM applications.")
	uploadRec := httptest.NewRecorder()
	router.ServeHTTP(uploadRec, upload)
	if uploadRec.Code != http.StatusCreated {
		t.Fatalf("upload status %d: %s", uploadRec.Code, uploadRec.Body.String())
	}

	var documentCount, chunkCount int
	if err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM documents").Scan(&documentCount); err != nil {
		t.Fatalf("count documents: %v", err)
	}
	if err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM chunks").Scan(&chunkCount); err != nil {
		t.Fatalf("count chunks: %v", err)
	}
	if documentCount != 1 || chunkCount != 1 {
		t.Fatalf("expected one document and chunk, got documents=%d chunks=%d", documentCount, chunkCount)
	}

	queryReq := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/rag/query",
		strings.NewReader(`{"question":"What is Eino?","top_k":5,"stream":false}`),
	)
	queryReq.Header.Set("Content-Type", "application/json")
	queryRec := httptest.NewRecorder()
	router.ServeHTTP(queryRec, queryReq)
	if queryRec.Code != http.StatusOK {
		t.Fatalf("query status %d: %s", queryRec.Code, queryRec.Body.String())
	}

	var response rag.QueryResponse
	if err := json.Unmarshal(queryRec.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode query response: %v", err)
	}
	if len(response.Sources) != 1 {
		t.Fatalf("expected one source, got %#v", response.Sources)
	}
	source := response.Sources[0]
	if source.DocumentID == "" || source.Source != "eino.md" || source.Score <= 0 {
		t.Fatalf("unexpected source: %#v", source)
	}
	if !strings.Contains(chatPrompt, "Eino is a Go framework") {
		t.Fatalf("chat prompt did not receive retrieved context: %q", chatPrompt)
	}
}

func applyMigrations(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate integration test file")
	}
	migrationDir := filepath.Join(filepath.Dir(currentFile), "..", "..", "migrations")
	for _, name := range []string{"001_init.sql", "002_hnsw_index.sql"} {
		sqlBytes, err := os.ReadFile(filepath.Join(migrationDir, name))
		if err != nil {
			t.Fatalf("read migration %s: %v", name, err)
		}
		if _, err := pool.Exec(ctx, string(sqlBytes)); err != nil {
			t.Fatalf("apply migration %s: %v", name, err)
		}
	}
}

func makeMultipartUpload(t *testing.T, filename, content string) *http.Request {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatalf("create multipart file: %v", err)
	}
	if _, err := part.Write([]byte(content)); err != nil {
		t.Fatalf("write multipart file: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/documents", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

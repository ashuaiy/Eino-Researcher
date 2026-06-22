package llm

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"eino-researcher/internal/config"
)

func TestOpenAICompatibleEmbedderSendsCompatibleRequest(t *testing.T) {
	var gotAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/embeddings" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		gotAuth = r.Header.Get("Authorization")

		var body struct {
			Model      string `json:"model"`
			Input      string `json:"input"`
			Dimensions int    `json:"dimensions"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if body.Model != "embedding-model" || body.Input != "hello" || body.Dimensions != 3 {
			t.Fatalf("unexpected request: %#v", body)
		}

		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": []any{map[string]any{"embedding": []float32{0.1, 0.2, 0.3}}},
		})
	}))
	defer server.Close()

	embedder := NewOpenAICompatibleEmbedder(config.EmbeddingConfig{
		BaseURL: server.URL + "/v1/",
		APIKey:  "secret-key",
		Model:   "embedding-model",
		Dim:     3,
		Timeout: time.Second,
	})
	vector, err := embedder.Embed(context.Background(), "hello")
	if err != nil {
		t.Fatalf("embed: %v", err)
	}
	if len(vector) != 3 {
		t.Fatalf("expected 3 dimensions, got %d", len(vector))
	}
	if gotAuth != "Bearer secret-key" {
		t.Fatalf("unexpected authorization header: %q", gotAuth)
	}
}

func TestOpenAICompatibleEmbedderOmitsAuthorizationWithoutAPIKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "" {
			t.Fatalf("expected no authorization header, got %q", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": []any{map[string]any{"embedding": []float32{1}}},
		})
	}))
	defer server.Close()

	embedder := NewOpenAICompatibleEmbedder(config.EmbeddingConfig{
		BaseURL: server.URL,
		Model:   "local-model",
		Dim:     1,
		Timeout: time.Second,
	})
	if _, err := embedder.Embed(context.Background(), "private input"); err != nil {
		t.Fatalf("embed: %v", err)
	}
}

func TestOpenAICompatibleEmbedderReturnsSafeProviderErrors(t *testing.T) {
	const secret = "top-secret"
	const input = "private document text"

	tests := []struct {
		name    string
		handler http.HandlerFunc
	}{
		{
			name: "non-2xx",
			handler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, secret+" "+input, http.StatusBadGateway)
			},
		},
		{
			name: "invalid json",
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("{"))
			},
		},
		{
			name: "empty data",
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(`{"data":[]}`))
			},
		},
		{
			name: "wrong dimensions",
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(`{"data":[{"embedding":[1]}]}`))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			embedder := NewOpenAICompatibleEmbedder(config.EmbeddingConfig{
				BaseURL: server.URL,
				APIKey:  secret,
				Model:   "model",
				Dim:     2,
				Timeout: time.Second,
			})
			_, err := embedder.Embed(context.Background(), input)
			if err == nil {
				t.Fatal("expected error")
			}
			if !errors.Is(err, ErrProvider) {
				t.Fatalf("expected provider error, got %v", err)
			}
			if strings.Contains(err.Error(), secret) || strings.Contains(err.Error(), input) {
				t.Fatalf("error leaked sensitive content: %v", err)
			}
		})
	}
}

package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"eino-researcher/internal/config"
)

func TestOpenAICompatibleClientCallsChatCompletions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/chat/completions" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer chat-secret" {
			t.Fatalf("unexpected authorization: %q", got)
		}

		var body struct {
			Model    string    `json:"model"`
			Messages []Message `json:"messages"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if body.Model != "chat-model" || len(body.Messages) != 2 {
			t.Fatalf("unexpected request: %#v", body)
		}

		_ = json.NewEncoder(w).Encode(map[string]any{
			"model": "chat-model",
			"choices": []any{
				map[string]any{"message": map[string]any{"content": "grounded answer [1]"}},
			},
		})
	}))
	defer server.Close()

	client := NewOpenAICompatibleClient(config.LLMConfig{
		BaseURL: server.URL + "/api/",
		APIKey:  "chat-secret",
		Model:   "chat-model",
		Timeout: time.Second,
	})
	resp, err := client.Generate(context.Background(), ChatRequest{
		SystemPrompt: "system",
		Prompt:       "question",
	})
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if resp.Content != "grounded answer [1]" || resp.Model != "chat-model" {
		t.Fatalf("unexpected response: %#v", resp)
	}
}

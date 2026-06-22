package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"eino-researcher/internal/config"
)

type Embedder interface {
	Embed(ctx context.Context, input string) ([]float32, error)
}

type OpenAICompatibleEmbedder struct {
	cfg    config.EmbeddingConfig
	client *http.Client
}

func NewOpenAICompatibleEmbedder(cfg config.EmbeddingConfig) *OpenAICompatibleEmbedder {
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 60 * time.Second
	}
	return &OpenAICompatibleEmbedder{
		cfg:    cfg,
		client: &http.Client{Timeout: timeout},
	}
}

func (e *OpenAICompatibleEmbedder) Embed(ctx context.Context, input string) ([]float32, error) {
	if strings.TrimSpace(input) == "" {
		return nil, fmt.Errorf("embedding input is required")
	}
	dim := e.cfg.Dim
	if dim <= 0 {
		dim = 1536
	}

	payload := struct {
		Model      string `json:"model"`
		Input      string `json:"input"`
		Dimensions int    `json:"dimensions"`
	}{
		Model:      e.cfg.Model,
		Input:      input,
		Dimensions: dim,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("%w: encode embedding request", ErrProvider)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		modelEndpoint(e.cfg.BaseURL, "embeddings"),
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: create embedding request", ErrProvider)
	}
	req.Header.Set("Content-Type", "application/json")
	if e.cfg.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+e.cfg.APIKey)
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: embedding request failed", ErrProvider)
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("%w: embedding provider returned status %d", ErrProvider, resp.StatusCode)
	}

	var decoded struct {
		Data []struct {
			Embedding []float32 `json:"embedding"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return nil, fmt.Errorf("%w: decode embedding response", ErrProvider)
	}
	if len(decoded.Data) == 0 {
		return nil, fmt.Errorf("%w: embedding response contained no data", ErrProvider)
	}
	vector := decoded.Data[0].Embedding
	if len(vector) != dim {
		return nil, fmt.Errorf("%w: embedding dimension mismatch: expected %d, got %d", ErrProvider, dim, len(vector))
	}
	return vector, nil
}

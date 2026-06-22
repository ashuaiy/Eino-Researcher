package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"eino-researcher/internal/config"
)

var ErrProvider = errors.New("model provider request failed")

type ChatClient interface {
	Generate(ctx context.Context, req ChatRequest) (ChatResponse, error)
}

type OpenAICompatibleClient struct {
	cfg    config.LLMConfig
	client *http.Client
}

func NewOpenAICompatibleClient(cfg config.LLMConfig) *OpenAICompatibleClient {
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 60 * time.Second
	}
	return &OpenAICompatibleClient{
		cfg:    cfg,
		client: &http.Client{Timeout: timeout},
	}
}

func (c *OpenAICompatibleClient) Generate(ctx context.Context, req ChatRequest) (ChatResponse, error) {
	if strings.TrimSpace(req.Prompt) == "" {
		return ChatResponse{}, fmt.Errorf("prompt is required")
	}

	messages := make([]Message, 0, len(req.Messages)+2)
	if req.SystemPrompt != "" {
		messages = append(messages, Message{Role: "system", Content: req.SystemPrompt})
	}
	messages = append(messages, req.Messages...)
	messages = append(messages, Message{Role: "user", Content: req.Prompt})

	payload := struct {
		Model       string    `json:"model"`
		Messages    []Message `json:"messages"`
		Temperature float32   `json:"temperature,omitempty"`
	}{
		Model:       c.cfg.Model,
		Messages:    messages,
		Temperature: req.Temperature,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return ChatResponse{}, fmt.Errorf("%w: encode chat request", ErrProvider)
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		modelEndpoint(c.cfg.BaseURL, "chat/completions"),
		bytes.NewReader(body),
	)
	if err != nil {
		return ChatResponse{}, fmt.Errorf("%w: create chat request", ErrProvider)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if c.cfg.APIKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)
	}

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return ChatResponse{}, fmt.Errorf("%w: chat request failed", ErrProvider)
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return ChatResponse{}, fmt.Errorf("%w: chat provider returned status %d", ErrProvider, resp.StatusCode)
	}

	var decoded struct {
		Model   string `json:"model"`
		Choices []struct {
			Message Message `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return ChatResponse{}, fmt.Errorf("%w: decode chat response", ErrProvider)
	}
	if len(decoded.Choices) == 0 || strings.TrimSpace(decoded.Choices[0].Message.Content) == "" {
		return ChatResponse{}, fmt.Errorf("%w: chat response contained no content", ErrProvider)
	}

	model := decoded.Model
	if model == "" {
		model = c.cfg.Model
	}
	return ChatResponse{Content: decoded.Choices[0].Message.Content, Model: model}, nil
}

func modelEndpoint(baseURL, resource string) string {
	return strings.TrimRight(baseURL, "/") + "/" + resource
}

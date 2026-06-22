package config

import "testing"

func TestValidateRequires1536EmbeddingDimensions(t *testing.T) {
	cfg := Config{Embedding: EmbeddingConfig{Dim: 768}}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected embedding dimension validation error")
	}

	cfg.Embedding.Dim = 1536
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected valid config, got %v", err)
	}
}

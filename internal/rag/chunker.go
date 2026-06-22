package rag

import "strings"

type Chunker interface {
	Split(text string) []string
}

type FixedSizeChunker struct {
	MaxRunes int
	Overlap  int
}

func NewFixedSizeChunker(maxRunes, overlap int) FixedSizeChunker {
	return FixedSizeChunker{MaxRunes: maxRunes, Overlap: overlap}
}

func (c FixedSizeChunker) Split(text string) []string {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil
	}

	maxRunes := c.MaxRunes
	if maxRunes <= 0 {
		maxRunes = 1000
	}
	overlap := c.Overlap
	if overlap < 0 || overlap >= maxRunes {
		overlap = 0
	}

	runes := []rune(text)
	if len(runes) <= maxRunes {
		return []string{text}
	}

	var chunks []string
	for start := 0; start < len(runes); {
		end := start + maxRunes
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[start:end]))
		if end == len(runes) {
			break
		}
		start = end - overlap
	}
	return chunks
}

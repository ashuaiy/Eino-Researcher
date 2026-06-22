package rag

import (
	"reflect"
	"testing"
)

func TestFixedSizeChunkerReturnsNilForWhitespace(t *testing.T) {
	chunks := NewFixedSizeChunker(4, 1).Split(" \n\t ")
	if chunks != nil {
		t.Fatalf("expected nil chunks, got %#v", chunks)
	}
}

func TestFixedSizeChunkerSplitsUnicodeByRunesWithOverlap(t *testing.T) {
	chunks := NewFixedSizeChunker(4, 1).Split("甲乙丙丁戊己")
	want := []string{"甲乙丙丁", "丁戊己"}
	if !reflect.DeepEqual(chunks, want) {
		t.Fatalf("expected %#v, got %#v", want, chunks)
	}
}

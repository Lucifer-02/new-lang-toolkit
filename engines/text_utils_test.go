package engines

import (
	"fmt"
	"testing"
	"unicode/utf8"
)

func TestSplitText1(t *testing.T) {

	text := "Hello. world!"
	limit := 10
	chunks := SplitText(text, limit)
	//print chunks
	for i, chunk := range chunks {
		fmt.Printf("Chunk %d: %s\n", i, chunk)
	}

	if len(chunks) != 2 {
		t.Errorf("Expected 2 chunks, got %d", len(chunks))
	}
}

func TestSplitText2(t *testing.T) {

	text := "Hello. world!, hoang"
	limit := 10
	chunks := SplitText(text, limit)
	//print chunks
	for i, chunk := range chunks {
		fmt.Printf("Chunk %d: %s\n", i, chunk)
	}

	if len(chunks) != 3 {
		t.Errorf("Expected 2 chunks, got %d", len(chunks))
	}
}

// A word longer than the limit with no break characters must not panic.
func TestSplitTextLongWord(t *testing.T) {
	text := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" // 30 chars, no spaces
	chunks := SplitText(text, 10)

	if got := joinLen(chunks); got != len(text) {
		t.Fatalf("chunk bytes = %d, want %d", got, len(text))
	}
	for _, c := range chunks {
		if len(c) > 10 {
			t.Errorf("chunk %q exceeds limit", c)
		}
	}
}

// Multi-byte UTF-8 text must never be cut in the middle of a rune.
func TestSplitTextUnicodeRuneBoundary(t *testing.T) {
	text := "xin chào các bạn tôi tên là Lucifer rất vui được gặp mọi người hôm nay"
	chunks := SplitText(text, 20)

	if got := joinLen(chunks); got != len(text) {
		t.Fatalf("chunk bytes = %d, want %d", got, len(text))
	}
	for i, c := range chunks {
		if !utf8.ValidString(c) {
			t.Errorf("chunk %d = %q is not valid UTF-8", i, c)
		}
	}
}

func joinLen(chunks []string) int {
	sum := 0
	for _, c := range chunks {
		sum += len(c)
	}
	return sum
}

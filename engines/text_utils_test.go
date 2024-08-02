package engines

import (
	"fmt"
	"testing"
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

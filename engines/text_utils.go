package engines

import (
	"fmt"
	"unicode/utf8"
)

// SplitText breaks text into chunks of at most `limit` bytes, preferring to
// break at sentence ends (`.`/`?`), then clause breaks (`,`/`;`), then spaces.
// If no such break exists within a window it falls back to a hard cut at a
// UTF-8 rune boundary so a multi-byte character is never split.
//
// Invariant: the summed byte lengths of the chunks equal len(text) — no bytes
// are dropped or duplicated.
func SplitText(text string, limit int) []string {
	assert(limit > 0, "limit must be positive")

	if len(text) <= limit {
		return []string{text}
	}

	chunks := []string{}
	start := 0
	for start < len(text) {
		// Whatever is left fits in a single chunk.
		if len(text)-start <= limit {
			chunks = append(chunks, text[start:])
			break
		}

		// Scan the window [start, start+limit] for the best break point.
		window := start + limit
		sentenceEnd, clauseEnd, blank := -1, -1, -1
		for i := start; i <= window; i++ {
			switch text[i] {
			case '.', '?':
				sentenceEnd = i
			case ',', ';':
				clauseEnd = i
			case ' ':
				blank = i
			}
		}

		end := sentenceEnd
		if end < start {
			end = clauseEnd
		}
		if end < start {
			end = blank
		}

		if end < start {
			// No break in the window: hard-cut, backing up to a rune
			// boundary so we never split a multi-byte character.
			cut := window
			for cut > start && !utf8.RuneStart(text[cut]) {
				cut--
			}
			if cut == start {
				cut = window // pathological input; cut anyway to make progress
			}
			chunks = append(chunks, text[start:cut])
			start = cut
			continue
		}

		// Include the break byte in the chunk.
		chunks = append(chunks, text[start:end+1])
		start = end + 1
	}

	assert(len(chunks) > 0, "len(chunks) > 0")
	// sum length of chunks must be equal to length of text
	sum := 0
	for _, chunk := range chunks {
		sum += len(chunk)
	}
	assert(sum == len(text), fmt.Sprintf("sum: %d, length: %d", sum, len(text)))

	return chunks
}

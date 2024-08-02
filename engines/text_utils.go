package engines

import "fmt"

func is_end_sentence(text string) bool {
	if text[0] == '?' || text[0] == '.' {
		return true
	}
	return false
}

func is_interrupt_sentence(text string) bool {
	if text[0] == ',' || text[0] == ';' {
		return true
	}
	return false
}

func SplitText(text string, limit int) []string {
	chunks := []string{}
	length := len(text)

	blank_pos, end_sentence_pos, interrupt_sentence_pos := 0, 0, 0

	_ = blank_pos
	count := 0
	start, end, pos := 0, 0, 0

	if length <= limit {
		return []string{text}
	}
	for {
		if text[pos] == ' ' {
			blank_pos = pos
		}

		if is_end_sentence(text[pos:]) {
			end_sentence_pos = pos
		}

		if is_interrupt_sentence(text[pos:]) {
			interrupt_sentence_pos = pos
		}

		if count == limit {

			if end_sentence_pos > end {
				end = end_sentence_pos
			} else if interrupt_sentence_pos > end {
				end = interrupt_sentence_pos
			} else {
				end = blank_pos
			}
			assert(start < end, "start < end")
			chunks = append(chunks, text[start:end+1])
			// fmt.Printf("start: %d, end: %d, pos: %d\n", start, end, pos)
			start = end + 1
			pos = start
			count = 0

		}

		if pos == length-1 {
			end = pos
			chunks = append(chunks, text[start:end+1])
			break
		}

		count++
		pos++
	}

	assert(len(chunks) > 0, "len(chunks) > 0")
	// sum length of chunks must be equal to length of text
	sum := 0
	for _, chunk := range chunks {
		sum += len(chunk)
	}
	assert(sum == length, fmt.Sprintf("sum: %d, length: %d", sum, length))

	return chunks
}

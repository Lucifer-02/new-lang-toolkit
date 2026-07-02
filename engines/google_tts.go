package engines

import (
	"fmt"
	"net/url"
	"strings"
)

type TTSParams struct {
	client string
	ie     string
	Tl     string
}

const textLimit = 200

func generateTTSUrl(params TTSParams, text string) string {
	assert(params.client != "" && params.ie != "" && params.Tl != "", "Invalid TTSParams")

	const base = "https://translate.google.com/translate_tts"

	return fmt.Sprintf("%s?client=%s&ie=%s&tl=%s&q=%s", base, params.client, params.ie, params.Tl, url.QueryEscape(text))
}

func TTSConcurrent(text string, targetLang string) []byte {
	assert(text != "", "Invalid text")

	params := TTSParams{
		client: "tw-ob",
		ie:     "UTF-8",
		Tl:     targetLang,
	}

	// Make a Request
	chunks := SplitText(text, textLimit)

	urls := make([]string, len(chunks))
	for i, chunk := range chunks {
		urls[i] = generateTTSUrl(params, strings.TrimSpace(chunk))
	}

	assert(len(urls) == len(chunks), "Invalid URLs")
	assert(len(urls) > 0, "URLs is empty")

	var combinedAudio []byte
	for _, body := range ApiRequests(urls) {
		combinedAudio = append(combinedAudio, body...)
	}

	return combinedAudio
}

package engines

import (
	"fmt"
	"io"
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

func TTS(text string, targetLang string) []byte {
	assert(text != "", "Invalid text")

	params := TTSParams{
		client: "tw-ob",
		ie:     "UTF-8",
		Tl:     targetLang,
	}

	// Make a Request
	chunks := SplitText(text, textLimit)

	var audio []byte
	for _, chunk := range chunks {

		url := generateTTSUrl(params, strings.TrimSpace(chunk))

		resp, err := ApiRequest(url)
		if err != nil {
			panic(err)
		}
		if resp.StatusCode != 200 {
			panic(fmt.Sprintf("Error: %s", resp.Status))
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		audio = append(audio, body...)
	}

	return audio
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
	// for i := range chunks {
	// 	fmt.Println("Chunk: ", i, chunks[i], len(chunks[i]))
	// }

	urls := make([]string, len(chunks))
	for i, chunk := range chunks {
		urls[i] = generateTTSUrl(params, strings.TrimSpace(chunk))
	}
	// for i := range urls {
	// 	fmt.Println("Url: ", urls[i])
	// }

	assert(len(urls) == len(chunks), "Invalid URLs")
	assert(len(urls) > 0, "URLs is empty")

	var combinedAudio []byte
	responses := ApiRequests(urls)
	for i := range len(responses) {
		if responses[i].StatusCode != 200 {
			panic(fmt.Sprintf("Error: %s", responses[i].Status))
		}

		body, err := io.ReadAll(responses[i].Body)
		if err != nil {
			panic(err)
		}

		combinedAudio = append(combinedAudio, body...)
	}

	return combinedAudio
}

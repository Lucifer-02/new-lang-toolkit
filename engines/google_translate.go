package engines

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type TransParams struct {
	Client string
	Ie     string // input encode
	Oe     string // output encode
	Dt     string // translate mode
	Sl     string // source language
	Tl     string // target language
}

func asSlice(v any, msg string) []any {
	s, ok := v.([]any)
	assert(ok, msg)
	return s
}

func extractTranslation(body []byte) []string {
	assert(len(body) > 0, "No response from server")

	//an example of the body is this: [[["CHÀO","hi",null,null,10]],null,"en",null,null,null,0.77704483,[],[["en"],null,[0.77704483],["en"]]],
	//only get CHÀO

	var jsonData []interface{}
	if err := json.Unmarshal(body, &jsonData); err != nil {
		panic(fmt.Sprintf("failed to parse translation response: %v", err))
	}
	assert(len(jsonData) > 0, "unexpected translation response shape")

	result := []string{}
	for _, data := range asSlice(jsonData[0], "unexpected translation response shape") {
		if data == nil {
			continue
		}
		segment := asSlice(data, "unexpected translation segment shape")
		assert(len(segment) > 0, "unexpected translation segment shape")
		text, ok := segment[0].(string)
		assert(ok, "unexpected translation segment shape")
		result = append(result, text)
	}

	assert(len(result) > 0, "No translation found")
	return result
}

func GoogleTranslate(text string, sourceLang string, targetLang string) string {
	const baseUrl = "https://translate.googleapis.com/translate_a/single?"

	params := TransParams{
		Sl:     sourceLang,
		Tl:     targetLang,
		Client: "gtx",
		Ie:     "UTF-8",
		Oe:     "UTF-8",
		Dt:     "t",
	}

	url := fmt.Sprintf("%sclient=%s&ie=%s&oe=%s&dt=%s&sl=%s&tl=%s&q=%s", baseUrl, params.Client, params.Ie, params.Oe, params.Dt, params.Sl, params.Tl, url.QueryEscape(text))

	body := ApiRequest(url)
	return strings.Join(extractTranslation(body), "")
}

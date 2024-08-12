package engines

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func RequestTrans(url string) (http.Response, error) {
	assert(url != "", "URL is empty")

	// Make a request
	response, err := http.Get(url)

	return *response, err
}

func readBody(response http.Response) ([]byte, error) {
	assert(response.StatusCode == 200, "Request failed")

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func extractTranslation(body []byte) []string {
	assert(len(body) > 0, "No response from server")

	var jsonData []interface{}
	json.Unmarshal(body, &jsonData)

	//an example of the body is this: [[["CHÀO","hi",null,null,10]],null,"en",null,null,null,0.77704483,[],[["en"],null,[0.77704483],["en"]]],
	//only get CHÀO

	result := []string{}
	for _, data := range jsonData[0].([]any) {
		if data == nil {
			continue
		}
		result = append(result, data.([]any)[0].(string))
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

	response := ApiRequest(url)
	if response.StatusCode != 200 {
		panic(fmt.Sprintf("Error: %s", response.Status))
	}

	body, err := readBody(response)
	if err != nil {
		panic(err)
	}

	return strings.Join(extractTranslation(body), "")
}

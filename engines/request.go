package engines

import (
	"net/http"
	"sort"
)

func ApiRequest(url string) http.Response {
	assert(url != "", "URL is empty")

	// Make a request
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	return *response
}

type resultChunk struct {
	index int
	value http.Response
}

func ApiRequests(urls []string) []http.Response {
	assert(len(urls) > 0, "URLs is empty")

	// use goroutine to make requests concurrently
	channel := make(chan resultChunk)
	for i, url := range urls {
		go func(index int, url string) {
			resp := ApiRequest(url)
			channel <- resultChunk{index, resp}
		}(i, url)
	}

	results := make([]resultChunk, len(urls))
	for i := range len(urls) {
		results[i] = <-channel
	}

  assert(len(results) == len(urls), "Invalid results")

	// sort the results
	sort.Slice(results, func(i, j int) bool {
		return results[i].index < results[j].index
	})

	responses := make([]http.Response, len(results))
	for i, result := range results {
		responses[i] = result.value
	}

	return responses
}

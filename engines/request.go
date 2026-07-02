package engines

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"
)

var httpClient = &http.Client{Timeout: 30 * time.Second}

// maxConcurrentRequests caps how many requests hit the endpoint at once. The
// endpoints are Google's unofficial ones; firing an unbounded number of
// requests (one per TTS chunk) at a long document invites HTTP 429 throttling,
// which is slower than a bounded burst.
const maxConcurrentRequests = 8

// fetch performs a single GET and returns the response body. It owns the
// http.Response lifecycle: it checks the status and closes the body itself so
// callers never touch http.Response. It returns an error rather than panicking
// so it is safe to call from a goroutine (a goroutine panic can't be recovered
// by the caller).
func fetch(url string) ([]byte, error) {
	assert(url != "", "URL is empty")

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("request failed: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

// ApiRequest performs a single GET and returns the response body, panicking on
// failure (fail-fast).
func ApiRequest(url string) []byte {
	body, err := fetch(url)
	if err != nil {
		panic(err)
	}
	return body
}

type resultChunk struct {
	index int
	body  []byte
	err   error
}

// ApiRequests performs the GETs concurrently and returns the response bodies in
// the original URL order. It panics if any request fails.
func ApiRequests(urls []string) [][]byte {
	assert(len(urls) > 0, "URLs is empty")

	// use goroutines to make requests concurrently, bounded by a semaphore so
	// no more than maxConcurrentRequests are in flight at once
	channel := make(chan resultChunk)
	sem := make(chan struct{}, maxConcurrentRequests)
	for i, url := range urls {
		go func(index int, url string) {
			sem <- struct{}{}        // acquire a slot
			defer func() { <-sem }() // release it

			body, err := fetch(url)
			channel <- resultChunk{index: index, body: body, err: err}
		}(i, url)
	}

	results := make([]resultChunk, len(urls))
	for i := range len(urls) {
		results[i] = <-channel
	}

	assert(len(results) == len(urls), "Invalid results")

	// sort the results back into the original order
	sort.Slice(results, func(i, j int) bool {
		return results[i].index < results[j].index
	})

	bodies := make([][]byte, len(results))
	for i, result := range results {
		if result.err != nil {
			panic(result.err)
		}
		bodies[i] = result.body
	}

	return bodies
}

//////////////////////////////////////////////////////////////////////
//
// Your task is to change the code to limit the crawler to at most one
// page per second, while maintaining concurrency (in other words,
// Crawl() must be called concurrently)
//
// @hint: you can achieve this by adding 3 lines
//

package main

import (
	"fmt"
	"sync"
	"time"
)

// Crawl uses `fetcher` from the `mockfetcher.go` file to imitate a
// real crawler. It crawls until the maximum depth has reached.
func Crawl(url string, depth int, wg *sync.WaitGroup, throttle <-chan time.Time) {
	defer wg.Done()

	if depth <= 0 {
		return
	}

	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("found: %s %q\n", url, body)

	wg.Add(len(urls))
	for _, u := range urls {
		<-throttle
		// Do not remove the `go` keyword, as Crawl() must be
		// called concurrently
		go Crawl(u, depth-1, wg, throttle)
	}
}

func main() {
	var wg sync.WaitGroup
	// The crawler is allowed to kick off another goroutine after every second.
	throttle := time.Tick(time.Second)

	// The WaitGroup is used to wait for all goroutines to finish.
	// So long as the count is greater than 0, the main goroutine will block with wg.Wait().
	wg.Add(1)
	Crawl("http://golang.org/", 4, &wg, throttle)
	wg.Wait()
}

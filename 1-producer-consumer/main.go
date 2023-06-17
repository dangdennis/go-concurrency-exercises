//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer scenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"time"

	"github.com/dangdennis/go-concurrency-exercises/1-producer-consumer/mockstream"
)

func producer(stream mockstream.Stream, tweetCh chan *mockstream.Tweet, doneCh chan bool) {
	for {
		tweet, err := stream.Next()
		if err == mockstream.ErrEOF {
			doneCh <- true
			return
		}

		tweetCh <- tweet
	}
}

func consumer(tweetCh chan *mockstream.Tweet) {
	for t := range tweetCh {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}
}

func main() {
	start := time.Now()
	stream := mockstream.GetMockStream()

	tweetsChan := make(chan *mockstream.Tweet)
	doneChan := make(chan bool)

	// Producer
	go producer(stream, tweetsChan, doneChan)

	// Consumer
	go consumer(tweetsChan)

	// Ensures producer stops before consumer stops
	<-doneChan

	// Close channel to wait for consumer to finish
	close(tweetsChan)

	fmt.Printf("Process took %s\n", time.Since(start))
}

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"
)

type response struct {
	code int
	body string
}

func main() {
	start := time.Now()
	var (
		requestURI  = flag.String("url", "http://localhost", "The URL you wish to hit.")
		n           = flag.Int("n", 0, "Number of requests")
		concurrency = flag.Int("conc", runtime.NumCPU(), "How many instances to run in parallel.")
		logging     = flag.Bool("log", false, "If you like logging every request.")
	)
	flag.Parse()
	requests := make(chan string)
	count := 0
	errCount := 0
	log.Printf("Concurrency: %d", *concurrency)
	go func() {
		for i := 0; i < *n; i++ {
			requests <- *requestURI
		}
		close(requests)
	}()

	var wg sync.WaitGroup
	wg.Add(*concurrency)
	results := make(chan response)

	go func() {
		wg.Wait()
		fmt.Printf("Successful requests count: %d\n", count)
		fmt.Printf("Error requests count: %d\n", errCount)
		fmt.Printf("Execution time: %s.", time.Since(start))
		close(results)
	}()

	for i := 0; i < *concurrency; i++ {
		go func() {
			defer wg.Done()
			for url := range requests {
				results <- send(&url)
			}
		}()
	}

	for res := range results {
		switch res.code {
		case 200:
			count++
		default:
			errCount++
		}
		if *logging {
			log.Printf("status: %d, body: %s", res.code, res.body)
		}
	}
}

func send(url *string) response {
	res, err := http.Get(*url)
	if err != nil {
		return response{code: http.StatusInternalServerError, body: err.Error()}
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return response{code: res.StatusCode, body: res.Status}
	}
	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	responseBody := string(responseData)
	return response{code: res.StatusCode, body: responseBody}
}

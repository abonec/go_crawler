package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

const poolSize = 5

type Result struct {
	Url   string
	Count int
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	queue := make(chan string)
	results := make(chan Result)

	var wg sync.WaitGroup
	for i := 0; i < poolSize; i++ {
		go worker(queue, results, &wg)
	}

	done := make(chan interface{})
	go printer(results, &wg, done)

	for scanner.Scan() {
		queue <- scanner.Text()
	}
	close(queue)
	wg.Wait()
	wg.Add(1)
	close(done)
	wg.Wait()
}

func worker(input chan string, result chan Result, wg *sync.WaitGroup) {
	wg.Add(1)
	for {
		url, more := <-input
		if !more {
			wg.Done()
			break
		}
		count, err := httpCount(url)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		result <- Result{url, count}
	}

}

func httpCount(url string) (int, error) {
	var client http.Client

	resp, err := client.Get(url)
	if err != nil {
		return 0, formattedError("Error while getting %s", url)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, formattedError("Error while getting %s. Got response code %d", url, resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	count := strings.Count(string(body), "go")
	return count, nil
}

func formattedError(format string, args ...interface{}) error {
	return errors.New(fmt.Sprintf(format, args...))
}

func printer(results chan Result, wg *sync.WaitGroup, done chan interface{}) {
	totalCount := 0
Loop:
	for {
		select {
		case result, more := <-results:
			if !more {
				break
			}
			fmt.Printf("Count for %s: %d\n", result.Url, result.Count)
			totalCount += result.Count
		case _, more := <-done:
			if !more {
				break Loop
			}

		}
	}
	fmt.Printf("Total: %d\n", totalCount)
	wg.Done()
}

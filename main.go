package main

import (
	"fmt"
	"os"
	"sync"
	"bufio"
)

const poolSize = 5

type Result struct {
	Url   string
	Count int
}

func main() {
	queue := make(chan string)
	results := make(chan Result)

	var wg sync.WaitGroup
	for i := 0; i < poolSize; i++ {
		go worker(queue, results, &wg)
	}

	done := make(chan interface{})
	var printerWait sync.WaitGroup
	go printer(results, &printerWait, done)

	handleInterrupt(&done, &printerWait)
	scanStdin(queue)

	close(queue)
	wg.Wait()
	close(done)
	printerWait.Wait()
}
func scanStdin(queue chan string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		queue <- scanner.Text()
	}
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


func printer(results chan Result, wg *sync.WaitGroup, done chan interface{}) {
	wg.Add(1)
	totalCount := 0

	func() {
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
					return
				}

			}
		}
	}()
	fmt.Printf("Total: %d\n", totalCount)
	wg.Done()
}

package main

import (
	"bufio"
	"os"
	"fmt"
	"net/http"
	"errors"
	"io/ioutil"
	"strings"
)

const poolSize = 5

type Result struct {
	Url   string
	Count int
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	queue := make(chan string)
	defer close(queue)
	results := make(chan Result)

	for i := 0; i < poolSize; i++ {
		go worker(queue, results)
	}
	go printer(results)

	for scanner.Scan() {
		queue <- scanner.Text()
	}

}

func worker(input chan string, result chan Result) {
	for {
		url, more := <-input
		if !more {
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

func printer(results chan Result) {
	for result := range results {
		fmt.Printf("Count for %s: %d\n", result.Url, result.Count)
	}
}

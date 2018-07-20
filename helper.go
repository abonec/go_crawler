package main

import (
	"net/http"
	"io/ioutil"
	"strings"
	"errors"
	"fmt"
	"sync"
	"os"
	"os/signal"
)

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

func handleInterrupt(done *chan interface{}, wg *sync.WaitGroup) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		close(*done)
		wg.Wait()
		os.Exit(1)
	}()
}

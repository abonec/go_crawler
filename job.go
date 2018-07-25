package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const word = "Go"

type Job interface {
	Do()
}

type Result struct {
	Url   string
	Count int
}

type job struct {
	url    string
	result chan *Result
}

func NewJob(url string, result chan *Result) *job {
	return &job{url, result}
}

func (j *job) Do() {
	count, err := httpCount(j.url)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	j.result <- &Result{j.url, count}
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
	count := strings.Count(string(body), word)
	return count, nil
}

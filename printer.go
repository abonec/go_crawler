package main

import (
	"fmt"
)

type Printer struct {
	stopped    bool
	results    <-chan *Result
	stop       chan interface{}
	totalCount int
}

func NewPrinter(results <-chan *Result) *Printer {
	return &Printer{
		stopped: false,
		results: results,
		stop:    make(chan interface{}),
	}
}

func (p *Printer) Start() {
	go func() {
		for {
			select {
			case <-p.stop:
				p.stop <- 1
				return
			case result, more := <-p.results:
				if !more {
					break
				}
				fmt.Printf("Count for %s: %d\n", result.Url, result.Count)
				p.totalCount += result.Count
			}
		}
	}()
}

func (p *Printer) Stop() {
	if p.stopped {
		return
	}
	p.stop <- 1
	<-p.stop
	close(p.stop)
	p.PrintOut()
	p.stopped = true
}

func (p *Printer) PrintOut() {
	fmt.Printf("Total: %d\n", p.totalCount)
}

package main

import (
	"fmt"
)

type Printer struct {
	results    <-chan *Result
	stop       chan interface{}
	totalCount int
}

func NewPrinter(results <-chan *Result) *Printer {
	return &Printer{
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
				break
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
	p.stop <- 1
	<-p.stop
	close(p.stop)
	p.PrintOut()
}

func (p *Printer) PrintOut() {
	fmt.Printf("Total: %d\n", p.totalCount)
}

package main

import (
	"fmt"
	"context"
)

type Printer struct {
	isStopped  bool
	results    <-chan *Result
	stopChan   chan interface{}
	wait       chan interface{}
	totalCount int
	ctx        context.Context
}

func NewPrinter(ctx context.Context, results <-chan *Result) *Printer {
	return &Printer{
		isStopped: false,
		results:   results,
		stopChan:  make(chan interface{}),
		wait:      make(chan interface{}),
		ctx:       ctx,
	}
}

func (p *Printer) stop() {
	p.PrintOut()
	p.isStopped = true
}
func (p *Printer) Start() {
	go func() {
		defer func() {
			p.stop()
			close(p.wait)
		}()
		for {
			select {
			case <-p.ctx.Done():
				return
			case <-p.stopChan:
				return
			case result, more := <-p.results:
				if !more {
					return
				}
				fmt.Printf("Count for %s: %d\n", result.Url, result.Count)
				p.totalCount += result.Count
			}
		}
	}()
}

func (p *Printer) Wait() {
	<-p.wait
}

func (p *Printer) Stop() {
	if p.isStopped {
		return
	}
	close(p.stopChan)
}

func (p *Printer) PrintOut() {
	fmt.Printf("Total: %d\n", p.totalCount)
}

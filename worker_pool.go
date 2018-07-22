package main

import (
	"sync"
	"time"
)

type WorkerPool struct {
	limit       int
	currentSize int
	wg          *sync.WaitGroup
	idleTimeout time.Duration
	jobQueue    chan Job
	done        chan interface{}
	logger      Logger
	stopped     bool

	busyWorkers int
	jobStart    chan interface{}
	jobDone     chan interface{}
	workerDone  chan interface{}
	mutex       sync.RWMutex
}

func NewWorkerPool(limit int, logger Logger) *WorkerPool {
	return &WorkerPool{
		stopped: true,
		limit:   limit, idleTimeout: time.Second, logger: logger,
		jobQueue: make(chan Job),
		jobStart: make(chan interface{}),
		jobDone:  make(chan interface{}),
		wg:       &sync.WaitGroup{},
		mutex:    sync.RWMutex{},
		done:     make(chan interface{}),
	}
}

func (p *WorkerPool) StartSupervisor() {
	for {
		select {
		case <-p.jobStart:
			p.mutex.Lock()
			p.busyWorkers++
		case <-p.jobDone:
			p.mutex.Lock()
			p.busyWorkers--
		case <-p.workerDone:
			p.mutex.Lock()
			p.currentSize--
		}
		p.mutex.Unlock()
	}
}

func (p *WorkerPool) Start() {
	p.stopped = false
	go p.StartSupervisor()
}

func (p *WorkerPool) Stop() *sync.WaitGroup {
	if !p.stopped {
		close(p.done)
		p.stopped = true
	}
	return p.wg
}

func (p *WorkerPool) Wait() {
	p.wg.Wait()
}
func (p *WorkerPool) StopAndWait() {
	p.Stop()
	p.Wait()
}

func (p *WorkerPool) SendJob(job Job) {
	p.mutex.RLock()
	if p.busyWorkers == p.currentSize && p.currentSize < p.limit {
		go p.SpawnWorker()
	}
	p.mutex.RUnlock()
	p.jobQueue <- job
}

func (p *WorkerPool) SpawnWorker() {
	p.wg.Add(1)
	defer p.wg.Done()
	p.currentSize++
	for {
		select {
		case job, ok := <-p.jobQueue:
			if !ok {
				p.workerDone <- 1
				return
			}
			p.jobStart <- 1
			job.Do()
			p.jobDone <- 1
		case _, ok := <-p.done:
			if !ok {
				return
			}
		}
	}

}

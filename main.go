package main

import (
	"bufio"
	"os"
)

const poolSize = 5

func main() {
	results := make(chan *Result)
	logger := NewLogger()

	pool := NewWorkerPool(poolSize, logger)
	pool.Start()

	printer := NewPrinter(results)
	printer.Start()

	handleInterrupt(pool, printer)

	scanStdin(pool, results)

	pool.StopAndWait()
	printer.Stop()
}
func scanStdin(pool *WorkerPool, results chan *Result) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		pool.SendJob(NewJob(scanner.Text(), results))
	}
}

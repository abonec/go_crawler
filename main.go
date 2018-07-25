package main

import (
	"bufio"
	"os"
	"context"
)

const poolSize = 5

func main() {
	results := make(chan *Result)
	logger := NewLogger()
	ctx, cancel := context.WithCancel(context.Background())

	pool := NewWorkerPool(ctx, poolSize, logger)
	pool.Start()

	printer := NewPrinter(ctx, results)
	printer.Start()

	handleInterrupt(cancel, printer)

	scanStdin(pool, results)

	pool.StopAndWait()
	close(results)
	printer.Wait()
}
func scanStdin(pool *WorkerPool, results chan *Result) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		pool.SendJob(NewJob(scanner.Text(), results))
	}
}

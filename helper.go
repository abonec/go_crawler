package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"context"
)

func formattedError(format string, args ...interface{}) error {
	return errors.New(fmt.Sprintf(format, args...))
}

func handleInterrupt(cancel context.CancelFunc, pool *WorkerPool, printer *Printer) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		cancel()
		pool.Stop()
		printer.Wait()
		os.Exit(0)
	}()
}

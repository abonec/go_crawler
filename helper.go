package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
)

func formattedError(format string, args ...interface{}) error {
	return errors.New(fmt.Sprintf(format, args...))
}

func handleInterrupt(pool *WorkerPool, printer *Printer) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		pool.Stop()
		printer.Stop()
		os.Exit(0)
	}()
}

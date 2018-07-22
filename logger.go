package main

import (
	"fmt"
	"io"
	"os"
)

type Logger interface {
	Printf(format string, v ...interface{})
}

type logger struct {
	dest io.Writer
}

func NewLogger() *logger {
	return &logger{os.Stdout}
}

func (l *logger) Printf(format string, v ...interface{}) {
	fmt.Fprintf(os.Stdout, format+"\n", v...)
}

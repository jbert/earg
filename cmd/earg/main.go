package main

import (
	"io"
	"log"
	"os"

	"github.com/jbert/earg/pkg/earg"
)

func main() {

	s, err := earg.NewArecordSource()
	if err != nil {
		log.Fatalf("Can't start arecord source: %s", err)
	}
	highFreq := 2048
	ear := earg.New(s, highFreq)
	o := earg.NewPrintObserver(os.Stdout)
	err = ear.Run(o)
	if err != nil && err != io.EOF {
		log.Fatalf("Failed to run: %w", err)
	}
}

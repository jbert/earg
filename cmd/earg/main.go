package main

import (
	"io"
	"log"
	"os"

	"github.com/jbert/earg/pkg/earg"
	"github.com/jbert/earg/pkg/earg/observer"
)

func main() {

	s, err := earg.NewArecordSource()
	if err != nil {
		log.Fatalf("Can't start arecord source: %s", err)
	}
	highFreq := 2048
	ear := earg.New(s, highFreq)
	o := observer.NewPrint(os.Stdout)
	err = ear.Run(o)
	if err != nil && err != io.EOF {
		log.Fatalf("Failed to run: %w", err)
	}
}

package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/jbert/earg/pkg/earg"
)

func main() {

	freq := 400
	dur := time.Second
	sampleRate := 16000

	source := earg.NewSineSource(sampleRate, freq, dur)
	ear := earg.New(source)
	err := ear.Run(os.Stdout)
	if err != nil && err != io.EOF {
		log.Fatalf("Failed to run: %w", err)
	}
}

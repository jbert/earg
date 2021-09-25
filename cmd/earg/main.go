package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/jbert/earg/pkg/earg"
)

func main() {

	dur := 10 * time.Second
	/*
		sampleRate := 16
		freqA := 1
		freqE := 2
	*/

	sampleRate := 16000
	freqA := 440
	freqE := 659

	sA := earg.NewSineSource(sampleRate, freqA, dur)
	sE := earg.NewSineSource(sampleRate, freqE, dur)
	mE := earg.NewScale(sE, 0.1)
	mux, err := earg.NewMux(sA, mE)
	if err != nil {
		log.Fatalf("Can't create mux: %s", err)
	}
	ear := earg.New(mux)
	err = ear.Run(os.Stdout)
	if err != nil && err != io.EOF {
		log.Fatalf("Failed to run: %w", err)
	}
}

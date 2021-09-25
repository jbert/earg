package earg

import (
	"fmt"
	"io"
	"log"
	"time"
)

type Ear struct {
	source Source
	rate   int

	wantedFullBufSize int
	readBufSize       int
	fullBuf           []float64
}

func New(s Source) *Ear {
	rate := s.SampleRate()

	// TODO: make a parameter?
	// freq range E2 (82Hz - C7 2093Hz)
	// covers operatic voice ranges

	readDur := time.Millisecond * 10
	highFreq := 2048
	e := &Ear{
		source: s,
		rate:   rate,

		wantedFullBufSize: highFreq * 2, // nyquist
		readBufSize:       rate * int(readDur) / int(time.Second),
		fullBuf:           make([]float64, 0),
	}

	if e.readBufSize < 1 {
		log.Fatalf("Read too small, can't progress (change readDur?")
	}

	//fmt.Printf("readDur %s rate %d\n", readDur, rate)
	//fmt.Printf("calc %v\n", rate*int(readDur)/int(time.Second))
	//fmt.Printf("read size %d full size %d\n", len(e.readBuf), len(e.fullBuf))

	return e
}

func appendToRingBuf(ring, buf []float64, ringBufSize int) ([]float64, bool) {
	ring = append(ring, buf...)
	if len(ring) > ringBufSize {
		ring = ring[len(ring)-ringBufSize:]
	}
	return ring, len(ring) == ringBufSize
}

func (e *Ear) Run(w io.Writer) error {
	readBuf := make([]float64, e.readBufSize)

	seenEOF := false
	for !seenEOF {
		n, err := e.source.Read(readBuf)
		if err != nil {
			if err == io.EOF {
				seenEOF = true
			} else {
				return fmt.Errorf("error readin: %w", err)
			}
		}
		fmt.Printf("Read %d samples\n", len(readBuf))

		var haveFullBuf bool
		e.fullBuf, haveFullBuf = appendToRingBuf(e.fullBuf, readBuf, e.wantedFullBufSize)

		if haveFullBuf {
			fmt.Printf("Process %d samples\n", len(e.fullBuf))
			err = e.process(w, n)
			if err != nil {
				return fmt.Errorf("Can't process: %w", err)
			}
		}
	}
	return nil
}

func (e *Ear) process(w io.Writer, numSamples int) error {
	//	fmt.Fprintf(w, "Got %d samples\n", numSamples)
	//	fmt.Fprintf(w, "%+v\n", e.readBuf[:numSamples])
	//	f := fourier.NewFFT()
	//fmt.Fprintf(w, "FFT [%v]\n", f)
	return nil
}

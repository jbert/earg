package earg

import (
	"fmt"
	"io"
	"log"
	"math/cmplx"
	"time"

	"gonum.org/v1/gonum/dsp/fourier"
)

type Ear struct {
	source Source
	rate   int

	wantedFullBufSize int
	readBufSize       int
	fullBuf           []float64

	coeffs    []complex128
	absCoeffs []float64
}

func New(s Source) *Ear {
	rate := s.SampleRate()

	// TODO: make a parameter?
	// freq range E2 (82Hz - C7 2093Hz)
	// covers operatic voice ranges

	readDur := time.Millisecond * 10
	highFreq := 2048
	wantedFullBufSize := highFreq * 2 // nyquist
	e := &Ear{
		source: s,
		rate:   rate,

		wantedFullBufSize: wantedFullBufSize,
		readBufSize:       rate * int(readDur) / int(time.Second),

		fullBuf:   make([]float64, 0),
		coeffs:    make([]complex128, wantedFullBufSize/2+1),
		absCoeffs: make([]float64, wantedFullBufSize/2+1),
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

func (e *Ear) Run(o Observer) error {
	readBuf := make([]float64, e.readBufSize)

	seenEOF := false
	for !seenEOF {
		numRead, err := e.source.Read(readBuf)
		if err != nil {
			if err == io.EOF {
				seenEOF = true
			} else {
				return fmt.Errorf("error readin: %w", err)
			}
		}

		var haveFullBuf bool
		e.fullBuf, haveFullBuf = appendToRingBuf(e.fullBuf, readBuf[:numRead], e.wantedFullBufSize)

		if haveFullBuf {
			//			fmt.Printf("Process %d samples\n", len(e.fullBuf))
			err = e.process(o)
			if err != nil {
				return fmt.Errorf("Can't process: %w", err)
			}
		}
	}
	return nil
}

func (e *Ear) process(o Observer) error {
	/*
		max := -1.0
		for _, s := range e.fullBuf {
			//		fmt.Fprintf(w, "s %9.6f - max %9.6f\n", s, max)
			if s > max {
				max = s
			}
		}
		fmt.Fprintf(w, "max is %9.6f\n", max)
	*/

	a := Analysis{
		SampleWidth: e.wantedFullBufSize,
		SampleStart: e.source.CurrentOffset(),
	}

	f := fourier.NewFFT(e.wantedFullBufSize)
	f.Coefficients(e.coeffs, e.fullBuf)
	e.setAbsCoeffs()
	//	printFFT(f, e.coeffs)
	maxFreqIndices := findMaxFreqIndices(f, e.absCoeffs)
	a.Frequencies = make([]float64, len(maxFreqIndices))
	for i, j := range maxFreqIndices {
		a.Frequencies[i] = 16000.0 / 4096.0 * float64(j)
	}
	//	fmt.Printf("Max freqs: %v\n", maxxes)

	o.Hear(a)

	return nil
}

func (e *Ear) setAbsCoeffs() {
	for i := range e.coeffs {
		e.absCoeffs[i] = cmplx.Abs(e.coeffs[i])
	}
}

func findMaxFreqIndices(f *fourier.FFT, c []float64) []int {
	n := f.Len()/2 + 1
	width := 5 // Must be odd
	maxIndexes := make([]int, 0)
ATTEMPT:
	for i := 0; i < n-width; i++ {
		possMax := c[i+width/2+1]
		for j := 0; j < width; j++ {
			if c[i+j] > possMax {
				continue ATTEMPT
			}
		}
		maxIndexes = append(maxIndexes, i+width/2+1)
	}
	return maxIndexes
}

func printFFT(f *fourier.FFT, coeffs []complex128) {
	n := f.Len()/2 + 1
	for i := 0; i < n; i++ {
		//		fftFreq := f.Freq(i)
		freq := 16000.0 / 4096.0 * float64(i)
		//		fmt.Printf("%d: %7.4f: %7.4f %7.4f\n", i, freq, cmplx.Abs(coeffs[i]), coeffs[i])
		fmt.Printf("%d: %7.4f: %7.4f\n", i, freq, cmplx.Abs(coeffs[i]))
	}
}

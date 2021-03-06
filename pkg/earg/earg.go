package earg

import (
	"fmt"
	"io"
	"log"
	"math/cmplx"
	"time"

	"github.com/jbert/earg/pkg/earg/observer"
	"gonum.org/v1/gonum/dsp/fourier"
)

type Ear struct {
	source Source
	rate   int

	sampleWindowSize int
	readBufSize      int
	fullBuf          []float64

	coeffs     []complex128
	normCoeffs []float64
}

func New(s Source) *Ear {
	rate := s.SampleRate()

	// TODO: make a parameter?
	// freq range E2 (82Hz - C7 2093Hz)
	// covers operatic voice ranges

	readDur := time.Millisecond * 10
	sampleWindowSize := 2048
	e := &Ear{
		source: s,
		rate:   rate,

		sampleWindowSize: sampleWindowSize,
		readBufSize:      rate * int(readDur) / int(time.Second),

		fullBuf:    make([]float64, 0),
		coeffs:     make([]complex128, sampleWindowSize/2+1),
		normCoeffs: make([]float64, sampleWindowSize/2+1),
	}

	if e.readBufSize < 1 {
		log.Fatalf("Read too small, can't progress (change readDur?")
	}

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
		e.fullBuf, haveFullBuf = appendToRingBuf(e.fullBuf, readBuf[:numRead], e.sampleWindowSize)

		if haveFullBuf {
			//			fmt.Printf("Process %d samples\n", len(e.fullBuf))
			err = e.process(o)
			if err != nil {
				if err == io.EOF {
					return io.EOF
				}
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

	a := observer.Analysis{
		SampleWidth: e.sampleWindowSize,
		SampleStart: e.source.CurrentOffset(),
	}

	f := fourier.NewFFT(e.sampleWindowSize)
	f.Coefficients(e.coeffs, e.fullBuf)
	e.setNormalisedCoeffs()

	a.FreqPower = make([]observer.FreqPower, len(e.normCoeffs))
	for i, p := range e.normCoeffs {
		a.FreqPower[i] = observer.FreqPower{
			Freq:  e.indexToFreq(i),
			Power: p,
		}
	}

	// Experimentally determined, meaning unclear
	// (also unclear how it scales with overall volume and number of samples)
	threshold := 2 / float64(e.sampleWindowSize)
	peakIndices := findPeaks(f, e.normCoeffs, threshold)
	a.Peaks = make([]float64, len(peakIndices))
	for i, j := range peakIndices {
		a.Peaks[i] = e.indexToFreq(j)
	}
	//	fmt.Printf("Max freqs: %v\n", maxxes)

	err := o.Hear(a)
	if err != nil {
		if err == io.EOF {
			return io.EOF
		}
		return fmt.Errorf("observer failed: %w", err)
	}

	return nil
}

func (e *Ear) indexToFreq(j int) float64 {
	// So this floats up to rate, not to max freq
	return float64(e.rate) / float64(e.sampleWindowSize) * float64(j)
}

func (e *Ear) setNormalisedCoeffs() {
	for i := range e.coeffs {
		e.normCoeffs[i] = cmplx.Abs(e.coeffs[i]) / float64(e.sampleWindowSize)
	}
}

func findPeaks(f *fourier.FFT, c []float64, threshold float64) []int {
	n := f.Len()/2 + 1
	width := 5 // Must be odd. High numbers will limit detectable low frequencies
	peakIndexes := make([]int, 0)

	//	fmt.Printf("threshold is %9.7f\n", threshold)
ATTEMPT:
	for i := 0; i < n-width; i++ {
		possPeak := c[i+width/2+1]
		if possPeak < threshold {
			continue ATTEMPT
		}
		for j := 0; j < width; j++ {
			if c[i+j] > possPeak {
				continue ATTEMPT
			}
		}
		peakIndexes = append(peakIndexes, i+width/2+1)
		//		fmt.Printf("max is %9.7f\n", possMax)
	}
	return peakIndexes
}

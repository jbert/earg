package earg

import (
	"fmt"
	"io"
	"math"
)

type Analysis struct {
	Frequencies []float64
	SampleStart int
	SampleWidth int
}

func (a Analysis) String() string {
	return fmt.Sprintf("%d (%d): %v", a.SampleStart, a.SampleWidth, a.Frequencies)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func MinCentsDiff(a float64, freqs []float64) int {
	minDiff := 1 << 32 // works on 32bit, more than enough for our porpoises
	for _, f := range freqs {
		diff := abs(CentsDiff(a, f))
		if diff < minDiff {
			minDiff = diff
		}
	}
	return minDiff
}

func CentsDiff(a, b float64) int {
	return int(math.Round(math.Log2(b/a) * 1200))
}

type Observer interface {
	Hear(a Analysis)
}

type PrintObserver struct {
	w io.Writer
}

type FuncObserver func(a Analysis)

func NewPrintObserver(w io.Writer) *PrintObserver {
	return &PrintObserver{w}
}

func (po PrintObserver) Hear(a Analysis) {
	fmt.Fprintf(po.w, "%s\n", a)
}

func NewFuncObserver(f func(a Analysis)) FuncObserver {
	return f
}

func (fo FuncObserver) Hear(a Analysis) {
	fo(a)
}

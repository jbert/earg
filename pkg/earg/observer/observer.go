package observer

import (
	"fmt"
	"io"

	"math"
)

type FreqPower struct {
	Freq  float64
	Power float64
}

type Analysis struct {
	Peaks       []float64
	FreqPower   []FreqPower
	SampleStart int
	SampleWidth int
}

func (a Analysis) String() string {
	return fmt.Sprintf("%d (%d): %v", a.SampleStart, a.SampleWidth, a.Peaks)
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

type Print struct {
	w io.Writer
}

type Func func(a Analysis) error

func NewPrint(w io.Writer) *Print {
	return &Print{w}
}

func (po Print) Hear(a Analysis) error {
	if len(a.Peaks) > 0 {
		fmt.Fprintf(po.w, "%s\n", a)
	}
	return nil
}

func NewFunc(f func(a Analysis) error) Func {
	return f
}

func (fo Func) Hear(a Analysis) error {
	return fo(a)
}

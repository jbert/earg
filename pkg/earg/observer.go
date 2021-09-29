package earg

import (
	"fmt"
	"io"
)

type Analysis struct {
	Frequencies []int
	SampleStart int
	SampleWidth int
}

func (a Analysis) String() string {
	return fmt.Sprintf("%d (%d): %v", a.SampleStart, a.SampleWidth, a.Frequencies)
}

type Observer interface {
	Hear(a Analysis)
}

type PrintObserver struct {
	w io.Writer
}

func NewPrintObserver(w io.Writer) *PrintObserver {
	return &PrintObserver{w}
}

func (po PrintObserver) Hear(a Analysis) {
	fmt.Fprintf(po.w, "%s\n", a)
}

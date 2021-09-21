package earg

import (
	"fmt"
	"io"
	"time"
)

// Sample is a single sound sample, normalised to -1/+1
type Sample float64

type SampleBuf struct {
	SampleRate int
	buf        []Sample
}

type Source interface {
	SampleRate() int
	Read(b []Sample) (int, error)
}

type Ear struct {
	source Source
	sb     SampleBuf
}

func New(s Source) *Ear {
	chunkDur := time.Millisecond * 100
	rate := s.SampleRate()

	numSamples := rate * int(chunkDur) / int(time.Second)
	buf := make([]Sample, numSamples)

	sb := SampleBuf{
		SampleRate: rate,
		buf:        buf,
	}

	return &Ear{
		source: s,
		sb:     sb,
	}
}

func (e *Ear) Run(w io.Writer) error {
	for {
		n, err := e.source.Read(e.sb.buf)
		if err != nil {
			return err
		}

		err = e.process(w, n)
		if err != nil {
			return fmt.Errorf("Can't process: %w", err)
		}
	}
	return nil
}

func (e *Ear) process(w io.Writer, numSamples int) error {
	fmt.Fprintf(w, "Got %d samples\n", numSamples)
	return nil
}

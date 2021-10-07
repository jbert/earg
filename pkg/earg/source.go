package earg

import (
	"errors"
	"fmt"
	"io"
)

type Source interface {
	SampleRate() int
	CurrentOffset() int
	Read(b []float64) (int, error)
}

type Mux struct {
	sampleRate int
	sources    []Source
}

type Scale struct {
	Source
	scale float64
}

type Clip struct {
	Source
	clip float64
}

func NewClip(s Source, clip float64) *Clip {
	return &Clip{Source: s, clip: clip}
}

func (c *Clip) Read(b []float64) (int, error) {
	n, err := c.Source.Read(b)
	// We do this on error too, might be EOF and doesn't hurt
	for i := range b {
		if b[i] > c.clip {
			b[i] = c.clip
		}
		if b[i] < -c.clip {
			b[i] = -c.clip
		}
	}
	return n, err
}

func NewScale(s Source, scale float64) *Scale {
	return &Scale{Source: s, scale: scale}
}

func (s *Scale) Read(b []float64) (int, error) {
	n, err := s.Source.Read(b)
	// Do this on error too, might be EOF and doesn't hurt
	for i := range b {
		b[i] *= s.scale
	}
	return n, err
}

func NewMux(sources ...Source) (*Mux, error) {
	if len(sources) == 0 {
		return nil, errors.New("No sources")
	}
	m := &Mux{}
	m.sampleRate = sources[0].SampleRate()
SOURCES:
	for i := range sources {
		if i == 0 {
			continue SOURCES
		}
		if sources[i].SampleRate() != m.sampleRate {
			return nil, fmt.Errorf("Source %d has incompatible sample rate: %d != %d", i, sources[i].SampleRate(), m.sampleRate)
		}
	}
	m.sources = sources
	return m, nil
}

func (m *Mux) CurrentOffset() int {
	// We read equally from all
	return m.sources[0].CurrentOffset()
}

func (m *Mux) SampleRate() int {
	return m.sampleRate
}

func (m *Mux) Read(b []float64) (int, error) {
	for j := range b {
		b[j] = 0
	}
	wantedN := len(b)
	actualN := 0
	buf := make([]float64, wantedN)
	isEOF := false

	for i := range m.sources {
		nn, err := m.sources[i].Read(buf)
		if i == 0 {
			actualN = nn
		} else {
			if nn != actualN {
				return 0, errors.New("TODO: teach Mux to handle differing partial read")
			}
		}
		if err != nil {
			if err == io.EOF {
				// if any source hits eof, all do
				isEOF = true
			} else {
				return 0, fmt.Errorf("mux: error reading from source %d: %w", i, err)
			}
		}

		for j := range b {
			b[j] += buf[j]
		}
	}
	// Normalise
	for j := range b {
		b[j] /= float64(len(m.sources))
	}
	var nilOrEof error
	if isEOF {
		nilOrEof = io.EOF
	}
	return actualN, nilOrEof
}

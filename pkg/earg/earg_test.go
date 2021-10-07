package earg

import (
	"io"
	"log"
	"testing"
	"time"

	"github.com/jbert/earg/pkg/earg/observer"
)

func TestRingBuf(t *testing.T) {
	ringBufSize := 5
	readBufSize := 2

	readBuf := make([]float64, readBufSize)
	ring := make([]float64, 0)
	var full bool
	for i := 0; i < 5; i++ {
		t.Logf("Iteration %d - adding %d to existing %d", i, len(readBuf), len(ring))
		for j := range readBuf {
			readBuf[j] = float64(i)
		}
		ring, full = appendToRingBuf(ring, readBuf, ringBufSize)
		if i >= 2 {
			if !full {
				t.Fatalf("Buffer not full after 3 iterations")
			}
			if len(ring) != ringBufSize {
				t.Fatalf("Buffer has wrong size when full")
			}
		} else {
			if full {
				t.Fatalf("buffer is reported as full too soon")
			}
			if len(ring) >= ringBufSize {
				t.Fatalf("Buffer has wrong size when not full")
			}
		}
		if ring[len(ring)-1] != float64(i) {
			t.Fatalf("Last value is incorrect [%5.2f != %5.2f]", ring[len(ring)-1], float64(i))
		}
	}
}

func TestHearSines(t *testing.T) {

	dur := 10 * time.Second

	sampleRate := 16000
	a4 := 440.0
	e4 := 659.0

	sA := NewSineSource(sampleRate, a4, dur)
	sE := NewSineSource(sampleRate, e4, dur)
	mE := NewScale(sE, 0.1)
	mux, err := NewMux(sA, mE)
	if err != nil {
		log.Fatalf("Can't create mux: %s", err)
	}
	highFreq := 4096
	ear := New(mux, highFreq)

	// Collect analyses here
	analyses := make([]observer.Analysis, 0)
	o := observer.NewFunc(func(a observer.Analysis) error {
		analyses = append(analyses, a)
		return nil
	})

	err = ear.Run(o)
	if err != nil && err != io.EOF {
		t.Fatalf("Failed to run: %s", err)
	}

	t.Logf("Collected %d analyses", len(analyses))
	if len(analyses) < 100 {
		t.Fatalf("Not enough analyses: %d", len(analyses))
	}

	allowedDiff := 2
	for _, a := range analyses {
		// Do we find our tones?
		minDiff := observer.MinCentsDiff(a4, a.Peaks)
		if minDiff > allowedDiff {
			t.Fatalf("Failed to find [%5.2f] in %s - got %d cents", a4, a, minDiff)
		}
		minDiff = observer.MinCentsDiff(e4, a.Peaks)
		if minDiff > allowedDiff {
			t.Fatalf("Failed to find [%5.2f] in %s - got %d cents", e4, a, minDiff)
		}
	}
}

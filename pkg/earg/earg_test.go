package earg

import "testing"

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

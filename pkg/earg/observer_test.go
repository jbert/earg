package earg

import "testing"

func TestCentsDiff(t *testing.T) {
	a4 := 440.0
	e4 := 659.25
	a5 := 880.0
	testCases := []struct {
		a, b          float64
		expectedCents int
	}{
		{a4, a4, 0},
		{a4, a5, 1200},
		{a5, a4, -1200},
		{a4, e4, 700},
		{e4, a4, -700},
	}

	for _, tc := range testCases {
		t.Logf("Getting cents diff for [%5.2f] : [%5.2f]", tc.a, tc.b)
		got := CentsDiff(tc.a, tc.b)
		if got != tc.expectedCents {
			t.Fatalf("Got %d expected %d", got, tc.expectedCents)
		}
	}
}

func TestMinCentsDiff(t *testing.T) {
	a4 := 440.0
	e4 := 659.25
	a5 := 880.0
	testCases := []struct {
		a             float64
		fs            []float64
		expectedCents int
	}{
		{a4, []float64{a4}, 0},
		{a4, []float64{e4}, 700},
		{a4, []float64{a4, e4}, 0},
		{a4, []float64{e4, a4}, 0},
		{a4, []float64{e4, a5}, 700},
		{a5, []float64{e4, a5}, 0},
		{a5, []float64{a4, e4}, 500},
	}

	for _, tc := range testCases {
		t.Logf("Getting cents diff for [%5.2f] : [%5.2f]", tc.a, tc.fs)
		got := MinCentsDiff(tc.a, tc.fs)
		if got != tc.expectedCents {
			t.Fatalf("Got %d expected %d", got, tc.expectedCents)
		}
	}
}

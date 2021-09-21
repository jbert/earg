package earg

import (
	"io"
	"math"
	"time"
)

type SineSource struct {
	freq       int
	sampleRate int

	maxDur  time.Duration
	now     time.Duration
	durStep time.Duration

	theta     float64
	thetaStep float64
}

func NewSineSource(sampleRate int, freq int, dur time.Duration) *SineSource {
	return &SineSource{
		freq:       freq,
		sampleRate: sampleRate,

		now:     0,
		maxDur:  dur,
		durStep: time.Second / time.Duration(sampleRate),

		theta:     0,
		thetaStep: 2 * math.Pi / float64(freq),
	}
}

func (ss *SineSource) SampleRate() int {
	return ss.sampleRate
}

func (ss *SineSource) Read(samples []Sample) (int, error) {
	numSamples := 0
	for i := range samples {
		s, err := ss.nextSample()
		if err != nil {
			return numSamples, err
		}
		samples[i] = s
		numSamples++
	}
	return numSamples, nil
}

func (ss *SineSource) nextSample() (Sample, error) {
	if ss.now > ss.maxDur {
		return 0, io.EOF
	}

	pi2 := math.Pi * 2

	sample := math.Sin(ss.theta)
	ss.theta += ss.thetaStep
	if ss.theta >= pi2 {
		ss.theta -= pi2
	}

	ss.now += ss.durStep

	return Sample(sample), nil
}

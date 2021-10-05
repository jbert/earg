package observer

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type SDL struct {
	window   *sdl.Window
	renderer *sdl.Renderer

	highFreq   int
	sampleRate int
	widthDur   time.Duration

	width  int32
	height int32

	currentX int32
}

func NewSDL(highFreq int, sampleRate int, widthDur time.Duration, width int, height int) (*SDL, error) {
	var err error
	s := &SDL{
		highFreq:   highFreq,
		sampleRate: sampleRate,
		widthDur:   widthDur,
		width:      int32(width),
		height:     int32(height),
	}

	winTitle := "Ear"
	s.window, err = sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED, s.width, s.height, sdl.WINDOW_SHOWN)
	if err != nil {
		return nil, fmt.Errorf("Failed to create sdl window: %w\n", err)
	}

	s.renderer, err = sdl.CreateRenderer(s.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return nil, fmt.Errorf("Failed to create sdl renderer: %w\n", err)
	}
	s.renderer.Clear()

	return s, nil
}

// var minFloat float64 = 1000000000000000.0
// var maxFloat float64 = -1000000000000000.0

// Piecewise linear heatmap across 5 colours:
// blue (0,0,1), cyan (0,1,1), green (0,1,0), yellow (1,1,0), red (1,0,0)
func powerToColour(p float64) (uint8, uint8, uint8) {
	/*
		if p > maxFloat {
			maxFloat = p
			fmt.Printf("%9.7f - %9.7f: %9.7f\n", minFloat, maxFloat, p)
		}
		if p < minFloat {
			minFloat = p
			fmt.Printf("%9.7f - %9.7f: %9.7f\n", minFloat, maxFloat, p)
		}
	*/
	var r, g, b uint8
	switch {
	case 0.00 <= p && p < 0.25:
		r = 0
		g = uint8((p - 0.00) * 4 * 255)
		b = 1
	case 0.25 <= p && p < 0.50:
		r = 0
		g = 1
		b = 255 - uint8((p-0.25)*4*255)
	case 0.50 <= p && p < 0.75:
		r = uint8((p - 0.50) * 4 * 255)
		g = 1
		b = 0
	case 0.75 <= p && p <= 1.00:
		r = 1
		g = uint8((p - 0.75) * 4 * 255)
		b = 0
	}
	return r, g, b
}

func (s *SDL) Hear(a Analysis) {

	/*
		s.renderer.SetDrawColor(0, 0, 0, 255)
		rect := sdl.Rect{s.currentX, 0, 5, s.height}
		s.renderer.FillRect(&rect)

		s.renderer.SetDrawColor(128, 128, 128, 255)
		for _, f := range a.Peaks {
			y := s.height - int32(float64(s.height)*f/float64(s.highFreq))
			rect = sdl.Rect{s.currentX, y, 5, 5}
			s.renderer.FillRect(&rect)
		}
	*/

	for _, fp := range a.FreqPower {

		r, g, b := powerToColour(fp.Power)
		s.renderer.SetDrawColor(r, g, b, 255)
		y := s.height - int32(float64(s.height)*fp.Freq/float64(s.highFreq))
		rect := sdl.Rect{s.currentX, y, 5, 5}
		s.renderer.FillRect(&rect)
	}

	/*
		on current x
			draw a rect for each freq in analysis
			at height: height * freq / highFreq
		current x++
	*/

	s.currentX++
	if s.currentX > s.width {
		s.currentX = 0
	}
	s.renderer.Present()
}
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

func (s *SDL) Hear(a Analysis) {

	s.renderer.SetDrawColor(0, 0, 0, 255)
	rect := sdl.Rect{s.currentX, 0, 5, s.height}
	s.renderer.FillRect(&rect)

	s.renderer.SetDrawColor(128, 128, 128, 255)
	for _, f := range a.Peaks {
		y := s.height - int32(float64(s.height)*f/float64(s.highFreq))
		rect = sdl.Rect{s.currentX, y, 5, 5}
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

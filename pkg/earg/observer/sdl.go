package observer

import (
	"fmt"
	"io"

	"github.com/veandco/go-sdl2/sdl"
)

type SDL struct {
	window   *sdl.Window
	renderer *sdl.Renderer

	highFreq int

	width  int32
	height int32

	currentX int32

	paused    bool
	stoppedCh chan struct{}
}

func (s *SDL) stopped() bool {
	stopped := false
	select {
	case _, ok := <-s.stoppedCh:
		if !ok {
			stopped = true
		}
	default:
	}
	return stopped
}

func (s *SDL) stop() {
	close(s.stoppedCh)
}

func NewSDL(highFreq int, width int, height int) (*SDL, error) {
	var err error
	s := &SDL{
		highFreq: highFreq,
		width:    int32(width),
		height:   int32(height),

		paused:    false,
		stoppedCh: make(chan struct{}),
	}

	//	sdl.Init(sdl.INIT_EVERYTHING)

	winTitle := "Ear"
	s.window, err = sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED, s.width, s.height, sdl.WINDOW_SHOWN)
	if err != nil {
		return nil, fmt.Errorf("Failed to create sdl window: %w", err)
	}

	s.renderer, err = sdl.CreateRenderer(s.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return nil, fmt.Errorf("Failed to create sdl renderer: %w", err)
	}
	err = s.renderer.Clear()
	if err != nil {
		return nil, fmt.Errorf("Can't clear renderer: %w", err)
	}

	go s.listenForEvents()

	return s, nil
}

func (s *SDL) listenForEvents() {
	// Drain any keyboard events
	for !s.stopped() {
		event := sdl.WaitEvent()
		switch ev := event.(type) {
		case *sdl.QuitEvent:
			s.stop()
		case *sdl.KeyboardEvent:
			if ev.Type == sdl.KEYDOWN {
				switch ev.Keysym.Sym {
				case sdl.K_p:
					// Data race. Don't think I care?
					s.paused = !s.paused
				case sdl.K_q:
					s.stop()
				case sdl.K_ESCAPE:
					s.stop()
				}
			}
		}
	}
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

func (s *SDL) Hear(a Analysis) error {
	var err error

	if s.stopped() {
		return io.EOF
	}
	if s.paused {
		return nil
	}

	rectWidth := int32(1)
	rectHeight := int32(s.height / int32(len(a.FreqPower)))
	if rectHeight < 1 {
		rectHeight = 1
	}
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
		err = s.renderer.SetDrawColor(r, g, b, 255)
		if err != nil {
			return fmt.Errorf("Can't setdrawcolor: %w", err)
		}
		y := s.height - int32(float64(s.height)*fp.Freq/float64(s.highFreq))
		rect := sdl.Rect{X: s.currentX, Y: y, W: rectWidth, H: rectHeight}
		err = s.renderer.FillRect(&rect)
		if err != nil {
			return fmt.Errorf("Can't fillrect: %w", err)
		}
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

	return nil
}

package observer

import (
	"fmt"
	"io"
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

func NewSDL(highFreq int, sampleRate int, widthDur time.Duration, width int, height int) (*SDL, error) {
	var err error
	s := &SDL{
		highFreq:   highFreq,
		sampleRate: sampleRate,
		widthDur:   widthDur,
		width:      int32(width),
		height:     int32(height),

		stoppedCh: make(chan struct{}),
	}

	//	sdl.Init(sdl.INIT_EVERYTHING)

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

	go s.listenForEvents()

	return s, nil
}

func (s *SDL) listenForEvents() {
	// Drain any keyboard events
	for !s.stopped() {
		event := sdl.WaitEvent()
		//		log.Printf("Got event: %v", event)
		switch ev := event.(type) {
		case *sdl.QuitEvent:
			//			log.Printf("Got quit event: %v", event)
			s.stop()
		case *sdl.KeyboardEvent:
			//			log.Printf("Got keyboard event: %v", event)
			switch ev.Keysym.Sym {
			case sdl.K_q:
				//				log.Printf("Got q keyboard event: %v", event)
				s.stop()
			case sdl.K_ESCAPE:
				//				log.Printf("Got escape keyboard event: %v", event)
				s.stop()
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

	if s.stopped() {
		return io.EOF
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

	return nil
}

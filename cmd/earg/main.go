package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/jbert/earg/pkg/earg"
	"github.com/jbert/earg/pkg/earg/observer"
)

type options struct {
	verbose uint
	gain    float64
}

func getOptions() (*options, error) {
	o := options{}
	flag.UintVar(&o.verbose, "verbose", 0, "Verbosity level")
	flag.Float64Var(&o.gain, "gain", 1.0, "Gain multiplier - may cause clipping")
	flag.Parse()

	if o.gain == 0 {
		return nil, errors.New("Zero gain is not useful")
	}
	return &o, nil
}

func main() {

	opts, err := getOptions()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n\n", err)
		flag.Usage()
		os.Exit(-1)
	}
	log.Printf("v %d\n", opts.verbose)

	if opts.verbose > 0 {
		log.Printf("Options: %+v\n", opts)
	}

	as, err := earg.NewArecordSource()
	if err != nil {
		log.Fatalf("Can't start arecord source: %s", err)
	}

	s := earg.NewClip(earg.NewScale(as, opts.gain), 1)

	highFreq := 2048
	ear := earg.New(s, highFreq)
	//	o := observer.NewPrint(os.Stdout)
	o, err := observer.NewSDL(highFreq, s.SampleRate(), 10*time.Second, 1024, 768)
	if err != nil && err != io.EOF {
		log.Fatalf("Failed to create sdl observer: %s", err)
	}
	err = ear.Run(o)
	if err != nil && err != io.EOF {
		log.Fatalf("Failed to run: %s", err)
	}
}

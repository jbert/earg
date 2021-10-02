package earg

import (
	"github.com/jbert/earg/pkg/earg/observer"
)

type Observer interface {
	Hear(a observer.Analysis)
}

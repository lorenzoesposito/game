package engine

import (
	"fmt"
	"math"
)

var (
	status string
	t0     int
)

func (e *Entity) Play(animation string, frame int) {
	if animation != status {
		t0 = frame
	}
	switch animation {
	case "Idle":
		e.idle(frame - t0)
	case "Run":
		e.run(frame - t0)
	}
	status = animation
}

func (e *Entity) idle(frame int) {
	num := int(f(float64(frame), 16))
	e.Mesh = Objects["elegant_guy"]["Idle"]["elegant"+"_"+fmt.Sprintf("%06d", num)]
}

func (e *Entity) run(frame int) {
	num := frame%16 + 1
	e.Mesh = Objects["elegant_guy"]["Run"]["elegant_"+fmt.Sprintf("%06d", num)]
}

func f(x float64, n float64) int {
	return int((math.Acos(math.Cos((1/n)*math.Pi*x)) * (n / math.Pi)) + 0.5)
}

package world

import (
	"math/rand"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const positionMargin = 100

var (
	r = rand.New(rand.NewSource(time.Now().UTC().Unix()))

	objects sync.Map // []*sdl.Rect{}
)

// World knows location of every object.
type World struct {
	// randomizer
	r *rand.Rand

	// map of all objects' positions in the world
	pos sync.Map

	// size
	H int32
	W int32

	// view port
	vp *sdl.Rect
}

// NewWorld ...
func NewWorld(width, height int32, screen *sdl.Rect) *World {
	return &World{
		r:  rand.New(rand.NewSource(time.Now().UTC().Unix())),
		pos: sync.Map{},
		W:  width,
		H:  height,
		vp: screen,
	}
}

// RandomizePos - randomizes position for the given "objID", "objW" and "objH".
func (w *World) RandomizePos(objID string, objW, objH int32) *sdl.Rect {
	var passed bool
	var pos *sdl.Rect

	for !passed {
		x := r.Int31n(w.W - objW) // decrement by width in order to fit whole object at the rightmost side
		y := r.Int31n(w.H - objH) // decrement by height in order to fit whole object at the bottom
		pos = &sdl.Rect{X: x, Y: y, W: objW, H: objH}
		passed = true

		w.pos.Range(func (key, value interface{}) bool {
			p := value.(*sdl.Rect)
			// clearenceZone has an extra margin to provide a gap in between objects
			clearenceZone := &sdl.Rect{
				X: p.X - positionMargin,
				Y: p.Y - positionMargin,
				W: p.W + positionMargin,
				H: p.H + positionMargin,
			}

			if clearenceZone.HasIntersection(pos) {
				passed = false
				return false
			}

			return true
		})

	}

	if pos != nil {
		w.pos.Store(objID, pos)
	}

	return pos
}

// Update ...
func (w *World) Update() {
	// noop
}

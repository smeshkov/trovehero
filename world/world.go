package world

import (
	"math/rand"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const posMargin = 100

var (
	r = rand.New(rand.NewSource(time.Now().UTC().Unix()))
)

// World knows location of every object.
type World struct {
	mu sync.RWMutex

	// randomizer
	Rand *rand.Rand

	// map of all objects' positions in the world
	pos map[string]*sdl.Rect

	// size
	H int32
	W int32

	// view port
	vp *sdl.Rect

	// holds player's score
	score int8

	// level
	level int8
}

// NewWorld ...
func NewWorld(width, height int32, screen *sdl.Rect, level int8) *World {
	return &World{
		Rand:  rand.New(rand.NewSource(time.Now().UTC().Unix())),
		pos:   make(map[string]*sdl.Rect),
		W:     width,
		H:     height,
		vp:    screen,
		level: level,
	}
}

// RandomizePos - randomizes position for the given "objID", "objW" and "objH".
func (w *World) RandomizePos(objID string, objW, objH int32) *sdl.Rect {
	var passed bool
	var pos *sdl.Rect

	for !passed {
		x := w.Rand.Int31n(w.W - objW) // decrement by width in order to fit whole object at the rightmost side
		y := w.Rand.Int31n(w.H - objH) // decrement by height in order to fit whole object at the bottom
		pos = &sdl.Rect{X: x, Y: y, W: objW, H: objH}
		passed = true

		for key, p := range w.pos {
			// skip itself hence it will be replaced anyway at the end
			if key == objID {
				continue
			}

			// clearenceZone has an extra margin to provide a gap in between objects
			clearenceZone := &sdl.Rect{
				X: p.X - posMargin,
				Y: p.Y - posMargin,
				W: p.W + posMargin,
				H: p.H + posMargin,
			}

			if clearenceZone.HasIntersection(pos) {
				passed = false
				break
			}
		}

	}

	if pos != nil {
		w.pos[objID] = pos
	}

	return pos
}

// IncScore increments player's score.
func (w *World) IncScore() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.score++
}

// GetScore returns player's score.
func (w *World) GetScore() int8 {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.score
}

// IncLevel increments level number.
func (w *World) IncLevel() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.level++
}

// GetLevel returns level number.
func (w *World) GetLevel() int8 {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.level
}

// Update ...
func (w *World) Update() {
	// noop
}

package world

import (
	"github.com/veandco/go-sdl2/sdl"
)

// World knows location of every object.
type World struct {
	H int32
	W int32

	screen *sdl.Rect
}

// NewWorld ...
func NewWorld(width, height int32, screen *sdl.Rect) *World {
	return &World{
		H:      height,
		W:      width,
		screen: screen,
	}
}

// Update ...
func (m *World) Update() {

}

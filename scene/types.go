package scene

import (
	"github.com/veandco/go-sdl2/sdl"
)

// SceneObject describes default API of the scene object.
type SceneObject interface {
	Update()
	Paint(r *sdl.Renderer) error
	Restart()
	Destroy()
}
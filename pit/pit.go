package pit

import (
	"sync"

	"github.com/smeshkov/trovehero/world"

	"github.com/veandco/go-sdl2/sdl"
)

// Pit represents an arbitrary pit object in the scene.
type Pit struct {
	mu sync.RWMutex

	ID string

	time int64

	X, Y int32
	W, H int32

	depth int8

	world *world.World
}

// NewPit creates new instance of the Pit.
func NewPit(id string, x, y, width, height int32, depth int8, w *world.World) *Pit {
	return &Pit{
		ID:    id,
		X:     x,
		Y:     y,
		W:     width,
		H:     height,
		depth: depth,
		world: w,
	}
}

// Depth tells how deep is the Pit.
func (p *Pit) Depth() int8 {
	// p.mu.RLock()
	// defer p.mu.RUnlock()

	return p.depth
}

// Update ...
func (p *Pit) Update() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.time++
}

// Paint ...
func (p *Pit) Paint(r *sdl.Renderer) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// fill new rectangle
	r.SetDrawColor(0, 0, 160, 255)
	r.FillRect(&sdl.Rect{X: p.X, Y: p.Y, W: p.W, H: p.H})
	r.SetDrawColor(0, 0, 0, 255)

	return nil
}

// Restart ...
func (p *Pit) Restart() {
	p.mu.Lock()
	defer p.mu.Unlock()

	pos := p.world.RandomizePos(p.ID, 150, 50)
	p = NewPit(p.ID, pos.X, pos.Y, pos.W, pos.H, -60, p.world)
}

// Destroy ...
func (p *Pit) Destroy() {}

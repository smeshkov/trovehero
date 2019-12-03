package pit

import (
	"sync"

	"github.com/veandco/go-sdl2/sdl"
)

// Pit represents an arbitrary pit object in the scene.
type Pit struct {
	mu sync.RWMutex

	time int64

	X, Y int32
	W, H int32

	depth int8
}

// NewPit creates new instance of the Pit.
func NewPit(x, y, height, width int32, depth int8) *Pit {
	return &Pit{
		X:     x,
		Y:     y,
		H:     height,
		W:     width,
		depth: depth,
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
	r.SetDrawColor(0, 0, 255, 255)
	r.FillRect(&sdl.Rect{X: p.X, Y: p.Y, W: p.W, H: p.H})
	r.SetDrawColor(0, 0, 0, 255)

	return nil
}

// Restart ...
func (p *Pit) Restart() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p = NewPit(1025/2, 250, 50, 150, -60)
}

// Destroy ...
func (p *Pit) Destroy() {}

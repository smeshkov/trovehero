package pit

import (
	// "fmt"
	"sync"

	"github.com/veandco/go-sdl2/sdl"
)

// Pit represents an arbitrary pit object in the scene.
type Pit struct {
	mu sync.RWMutex

	time int64

	// tools
	rect *sdl.Rect

	X, Y int32
	W, H int32

	depth int8
}

// NewPit creates new instance of the Pit.
func NewPit(x, y, height, width int32, depth int8) *Pit {
	return &Pit{
		X:     x - width/2,
		Y:     y - height/2,
		H:     height,
		W:     width,
		depth: depth,
	}
}

// Depth tells how deep is the Pit.
func (p *Pit) Depth() int8 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	return p.depth
}

// Update ...
func (p *Pit) Update() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.time++
}

// func (p *Pit) clearRect(r *sdl.Renderer) error {
// 	err := r.SetDrawColor(0, 0, 0, 0)
// 	if err != nil {
// 		return fmt.Errorf("could not set draw color: %w", err)
// 	}
// 	err = r.FillRect(p.rect)
// 	if err != nil {
// 		return fmt.Errorf("could not fill rectangle: %w", err)
// 	}
// 	return nil
// }

// Paint ...
func (p *Pit) Paint(r *sdl.Renderer) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// remove previous rectangle
	// err := p.clearRect(r)
	// if err != nil {
	// 	return err
	// }

	// fill new rectangle
	r.SetDrawColor(0, 0, 255, 255)
	p.rect = &sdl.Rect{X: p.X, Y: p.Y, W: p.W, H: p.H}
	r.FillRect(p.rect)
	r.SetDrawColor(0, 0, 0, 255)
	return nil
}

// Restart ...
func (p *Pit) Restart() {}

// Destroy ...
func (p *Pit) Destroy() {}

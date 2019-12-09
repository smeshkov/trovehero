package trove

import (
	"sync"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/smeshkov/trovehero/world"
)

const (
	troveW = 50
	troveH = 50
)

// Trove can be collected by the Hero.
type Trove struct {
	mu sync.RWMutex

	ID string

	time int64

	// position
	X, Y int32
	W, H int32

	world *world.World

	collected bool
}

// NewTrove creates new instance of Trove.
func NewTrove(id string, x, y int32, w *world.World) *Trove {
	return &Trove{
		ID: id,

		X: x,
		Y: y,
		W: troveW,
		H: troveH,

		world: w,
	}
}

// Update ...
func (t *Trove) Update() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.time++
}

// Paint ...
func (t *Trove) Paint(r *sdl.Renderer) error {
	t.mu.RLock()
	defer t.mu.RUnlock()

	// fill new rectangle
	r.SetDrawColor(160, 160, 0, 255)
	r.FillRect(&sdl.Rect{X: t.X, Y: t.Y, W: t.W, H: t.H})
	r.SetDrawColor(0, 0, 0, 255)

	return nil
}

// Restart ...
func (t *Trove) Restart() {
	t.mu.Lock()
	defer t.mu.Unlock()

	pos := t.world.RandomizePos(t.ID, troveW, troveH)
	t = NewTrove(t.ID, pos.X, pos.Y, t.world)
}

// Destroy ...
func (t *Trove) Destroy() {
}

// Collect collects this Trove.
func (t *Trove) Collect() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.collected = true
}

// IsCollected tells whether Trove has been collected already or not.
func (t *Trove) IsCollected() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.collected
}

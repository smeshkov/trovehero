package hero

import (
	"math"
	"sync"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/smeshkov/trovehero/pit"
	"github.com/smeshkov/trovehero/types/command"
)

const (
	gravity         = 0.1
	friction        = 0.08
	airFriction     = 0.1
	altitudeMargin  = 35
	collisionMargin = 10
)

// Hero is a playbale character.
type Hero struct {
	mu sync.RWMutex

	time int64

	// properties
	height       int32
	width        int32
	maxMoveSpeed float32
	maxJumpSpeed float32

	// coordinates
	altitude int8

	// shape
	x, y int32
	w, h int32

	// speed
	vertSpeed float32
	horSpeed  float32
	altSpeed  float32

	crashingDepth int8
	dead          bool
}

// NewHero creates new instance of Hero in given coordinates.
func NewHero(x, y int32) *Hero {
	var heroWidth int32 = 50
	var heroHeight int32 = 50

	return &Hero{
		// properties
		height:       heroHeight,
		width:        heroWidth,
		maxMoveSpeed: 4,
		maxJumpSpeed: 4,

		// shape
		x: x,
		y: y,
		h: heroHeight,
		w: heroWidth,
	}
}

// Do performes command on a Hero.
func (h *Hero) Do(t command.Type) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// can't move if not within alowed marging from the ground
	if h.altitude > altitudeMargin {
		return
	}

	switch t {
	case command.Jump:
		h.altSpeed = h.maxJumpSpeed
	case command.GoNorth:
		h.vertSpeed = -h.maxMoveSpeed
	case command.GoSouth:
		h.vertSpeed = h.maxMoveSpeed
	case command.GoWest:
		h.horSpeed = -h.maxMoveSpeed
	case command.GoEast:
		h.horSpeed = h.maxMoveSpeed
	}
}

// Update updates state of the Hero.
func (h *Hero) Update() {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.time++

	if h.horSpeed != 0 || h.vertSpeed != 0 {
		h.handleMove()
	}
	if h.crashingDepth == 0 && h.altSpeed != 0 {
		h.handleJump()
	}
	if h.crashingDepth != 0 {
		h.handleCrash()
	}
}

// Paint paints Hero to window.
func (h *Hero) Paint(r *sdl.Renderer) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// fill new rectangle
	r.SetDrawColor(255, 100, 0, 255)
	r.FillRect(h.getShape())
	r.SetDrawColor(0, 0, 0, 255)

	// i := b.time / 10 % len(b.textures)
	// if err := r.Copy(b.textures[i], nil, rect); err != nil {
	// 	return fmt.Errorf("could not copy background: %w", err)
	// }
	return nil
}

// Restart restarts state of Hero.
func (h *Hero) Restart() {
	h.mu.Lock()
	defer h.mu.Unlock()

	h = NewHero(800/2, 550)
}

// Destroy removes Hero.
func (h *Hero) Destroy() {
	// noop
}

func (h *Hero) getShape() *sdl.Rect {
	if h.altitude != 0 {
		return &sdl.Rect{
			X: h.x - int32(h.altitude/2),
			Y: h.y - int32(h.altitude/2),
			W: h.w + int32(h.altitude),
			H: h.h + int32(h.altitude),
		}
	}
	return &sdl.Rect{
		X: h.x,
		Y: h.y,
		W: h.w,
		H: h.h,
	}
}

func (h *Hero) handleCrash() {
	// crashing
	if h.altitude > h.crashingDepth {
		h.altSpeed -= gravity
		h.altitude += int8(h.altSpeed)
	} else { // crashed
		h.altSpeed = 0
		h.altitude = h.crashingDepth
		h.dead = true
	}
}

func (h *Hero) handleJump() {
	// rising
	if h.altSpeed > 0 {
		h.altitude += int8(h.altSpeed)
		h.altSpeed -= gravity
		return
	}

	// falling
	if h.altitude > 0 && h.altSpeed <= 0 {
		h.altitude = int8(math.Max(0, float64(h.altitude)+float64(h.altSpeed)))
		h.altSpeed -= gravity
		return
	}

	// landed
	if h.altitude == 0 && h.altSpeed < 0 {
		h.altSpeed = 0
		return
	}
}

func (h *Hero) handleMove() {
	var frict float64
	if h.altitude == 0 {
		frict = friction
	} else {
		frict = airFriction
	}

	if h.horSpeed != 0 {
		h.x += int32(h.horSpeed)
		if h.horSpeed > 0 {
			h.horSpeed = float32(math.Max(0, float64(h.horSpeed)-frict))
		} else {
			h.horSpeed = float32(math.Min(0, float64(h.horSpeed)+frict))
		}
	}
	if h.vertSpeed != 0 {
		h.y += int32(h.vertSpeed)
		if h.vertSpeed > 0 {
			h.vertSpeed = float32(math.Max(0, float64(h.vertSpeed)-frict))
		} else {
			h.vertSpeed = float32(math.Min(0, float64(h.vertSpeed)+frict))
		}
	}
}

// Touch checks collision with Pit.
func (h *Hero) Touch(p *pit.Pit) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.altitude > 0 { // above in the air
		return
	}
	if p.X > h.x+h.w-collisionMargin { // too far right
		return
	}
	if p.X+p.W-collisionMargin < h.x { // too far left
		return
	}
	if p.Y > h.y+h.h-collisionMargin { // too far below
		return
	}
	if p.Y+p.H-collisionMargin < h.y { // to far above
		return
	}

	h.crashingDepth = p.Depth()
}

// Location returns a location of the Hero.
func (h *Hero) Location() *sdl.Rect {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return &sdl.Rect{X: h.x, Y: h.y, W: h.w, H: h.h}
}

// IsDead ....
func (h *Hero) IsDead() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.dead
}

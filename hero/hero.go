package hero

import (
	// "log"
	"fmt"
	"math"
	"sync"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/smeshkov/trovehero/pit"
	"github.com/smeshkov/trovehero/types"
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

	// tools
	rect *sdl.Rect

	// properties
	height       int32
	width        int32
	maxMoveSpeed float32
	maxJumpSpeed float32

	// coordinates
	location *sdl.Point
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
	var coordX int32 = x - heroWidth/2
	var coordY int32 = y - heroHeight/2

	return &Hero{
		// properties
		height:       heroHeight,
		width:        heroWidth,
		maxMoveSpeed: 4,
		maxJumpSpeed: 4,

		// coordinates
		location: &sdl.Point{X: coordX, Y: coordY},

		// speed
		x: coordX,
		y: coordY,
		h: heroHeight,
		w: heroWidth,
	}
}

// Do performes command on a Hero.
func (h *Hero) Do(t types.CommandType) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// can't move if not within alowed marging from the ground
	if h.altitude > altitudeMargin {
		return
	}

	switch t {
	case types.Jump:
		h.altSpeed = h.maxJumpSpeed
	case types.GoNorth:
		h.vertSpeed = -h.maxMoveSpeed
	case types.GoSouth:
		h.vertSpeed = h.maxMoveSpeed
	case types.GoWest:
		h.horSpeed = -h.maxMoveSpeed
	case types.GoEast:
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

	// remove previous rectangle
	err := h.clearRect(r)
	if err != nil {
		return err
	}

	// fill new rectangle
	r.SetDrawColor(255, 100, 0, 255)
	h.rect = &sdl.Rect{X: h.x, Y: h.y, W: h.w, H: h.h}
	r.FillRect(h.rect)
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

func (h *Hero) resize() {
	if h.altitude != 0 {
		h.w = h.width + int32(h.altitude)
		h.h = h.height + int32(h.altitude)
		h.x = h.location.X - int32(h.altitude/2)
		h.y = h.location.Y - int32(h.altitude/2)
	} else if h.altitude == 0 {
		h.h = h.height
		h.w = h.width
		h.x = h.location.X
		h.y = h.location.Y
	}
}

func (h *Hero) clearRect(r *sdl.Renderer) error {
	err := r.SetDrawColor(0, 0, 0, 0)
	if err != nil {
		return fmt.Errorf("could not set draw color: %w", err)
	}
	err = r.FillRect(h.rect)
	if err != nil {
		return fmt.Errorf("could not fill rectangle: %w", err)
	}
	return nil
}

func (h *Hero) handleCrash() {
	// crashing
	if h.altitude > h.crashingDepth {
		h.altSpeed -= gravity
		h.altitude += int8(h.altSpeed)
		h.resize()
	} else { // crashed
		h.altSpeed = 0
		h.altitude = h.crashingDepth
		h.dead = true
	}
}

func (h *Hero) handleJump() {
	// rising
	if h.altSpeed > 0 {
		// log.Printf("rising... %.2f\n", h.altitude)
		h.altitude += int8(h.altSpeed)
		h.resize()
		h.altSpeed -= gravity
		return
	}

	// falling
	if h.altitude > 0 && h.altSpeed <= 0 {
		// log.Printf("falling...%.2f\n", h.altitude)
		h.altitude = int8(math.Max(0, float64(h.altitude)+float64(h.altSpeed)))
		h.resize()
		h.altSpeed -= gravity
		return
	}

	// landed
	if h.altitude == 0 && h.altSpeed < 0 {
		// log.Printf("landed...%.2f\n", h.altitude)
		h.altSpeed = 0
		h.h = h.height
		h.w = h.width
		h.x = h.location.X
		h.y = h.location.Y
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
		h.location.X += int32(h.horSpeed)
		h.x += int32(h.horSpeed)
		if h.horSpeed > 0 {
			h.horSpeed = float32(math.Max(0, float64(h.horSpeed)-frict))
		} else {
			h.horSpeed = float32(math.Min(0, float64(h.horSpeed)+frict))
		}
	}
	if h.vertSpeed != 0 {
		h.location.Y += int32(h.vertSpeed)
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
	// optimisation: do this expensive check only every 2nd time
	if h.time / 2 != 0 {
		return
	}

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

	if h.rect != nil {
		return h.rect
	}

	return &sdl.Rect{X: h.x, Y: h.y, W: h.w, H: h.h}
}

// IsDead ....
func (h *Hero) IsDead() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.dead
}

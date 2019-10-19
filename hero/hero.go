package hero

import (
	// "log"
	"fmt"
	"math"
	"sync"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/smeshkov/trovehero/pit"
)

const (
	gravity         = 0.1
	friction        = 0.2
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
	maxMoveSpeed float64
	maxJumpSpeed float64

	// coordinates
	coordX int32
	coordY int32

	// shape
	x, y int32
	w, h int32

	// velocity
	altitude  float64
	vertSpeed float64
	horSpeed  float64
	jumpSpeed float64

	// commands
	commands commands

	crashingDepth float64
	dead          bool
}

// NewHero creates new instance of Hero.
func NewHero(x, y int32) *Hero {
	var heroWidth int32 = 50
	var heroHeight int32 = 50
	var coordX int32 = x - heroWidth/2
	var coordY int32 = y - heroHeight/2

	return &Hero{
		// properties
		height:       heroHeight,
		width:        heroWidth,
		maxMoveSpeed: 7,
		maxJumpSpeed: 4,

		// coordinates
		coordX: coordX,
		coordY: coordY,

		// shape
		x: coordX,
		y: coordY,
		h: heroHeight,
		w: heroWidth,

		// timers
		commands: commands{},
	}
}

// Do performes command on a Hero.
func (h *Hero) Do(t CommandType) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// can't move if not within alowed marging from the ground
	if h.altitude > altitudeMargin {
		return
	}

	switch t {
	case Jump:
		h.jumpSpeed = h.maxJumpSpeed
	case Up:
		h.vertSpeed = -h.maxMoveSpeed
	case Down:
		h.vertSpeed = h.maxMoveSpeed
	case Left:
		h.horSpeed = -h.maxMoveSpeed
	case Right:
		h.horSpeed = h.maxMoveSpeed
	}
}

/* func (h *Hero) doCommand(cmd *Command) {
	switch cmd.t {
	case Jump:
		h.jumpSpeed = h.maxJumpSpeed
	case Up:
		h.vertSpeed = -h.maxMoveSpeed
	case Down:
		h.vertSpeed = h.maxMoveSpeed
	case Left:
		h.horSpeed = -h.maxMoveSpeed
	case Right:
		h.horSpeed = h.maxMoveSpeed
	}
} */

// Update updates state of the Hero.
func (h *Hero) Update() {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.time++

	// check commands
	/* for cmdType, cmd := range h.commands {
		// Time to live of the command is over,
		// remove it from register.
		if cmd.ttl <= 0 {
			delete(h.commands, cmdType)
			continue
		}

		isMoving := cmdType != Jump

		// performs command in following cases:
		// 1. Hero is standing on the ground
		// 2. Hero has jumped and the gap between Jump and Move commands is within the boundaries of the Move's TTL
		if h.altitude == 0 || isMoving {
			h.doCommand(cmd)
			delete(h.commands, cmdType)
		}

		cmd.ttl--
	} */

	if h.horSpeed != 0 || h.vertSpeed != 0 {
		h.handleMove()
	}
	if h.crashingDepth == 0 && h.jumpSpeed != 0 {
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
		h.x = h.coordX - int32(h.altitude/2)
		h.y = h.coordY - int32(h.altitude/2)
	} else if h.altitude == 0 {
		h.h = h.height
		h.w = h.width
		h.x = h.coordX
		h.y = h.coordY
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
		h.jumpSpeed -= gravity
		h.altitude += h.jumpSpeed
		h.resize()
	} else { // crashed
		h.jumpSpeed = 0
		h.altitude = h.crashingDepth
		h.dead = true
	}
}

func (h *Hero) handleJump() {
	// rising
	if h.jumpSpeed > 0 {
		// log.Printf("rising... %.2f\n", h.altitude)
		h.altitude += h.jumpSpeed
		h.resize()
		h.jumpSpeed -= gravity
		return
	}

	// falling
	if h.altitude > 0 && h.jumpSpeed <= 0 {
		// log.Printf("falling...%.2f\n", h.altitude)
		h.altitude = math.Max(0, h.altitude+h.jumpSpeed)
		h.resize()
		h.jumpSpeed -= gravity
		return
	}

	// landed
	if h.altitude == 0 && h.jumpSpeed < 0 {
		// log.Printf("landed...%.2f\n", h.altitude)
		h.jumpSpeed = 0
		h.h = h.height
		h.w = h.width
		h.x = h.coordX
		h.y = h.coordY
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
		h.coordX += int32(h.horSpeed)
		h.x += int32(h.horSpeed)
		if h.horSpeed > 0 {
			h.horSpeed = math.Max(0, h.horSpeed-frict)
		} else {
			h.horSpeed = math.Min(0, h.horSpeed+frict)
		}
	}
	if h.vertSpeed != 0 {
		h.coordY += int32(h.vertSpeed)
		h.y += int32(h.vertSpeed)
		if h.vertSpeed > 0 {
			h.vertSpeed = math.Max(0, h.vertSpeed-frict)
		} else {
			h.vertSpeed = math.Min(0, h.vertSpeed+frict)
		}
	}
}

// Touch ...
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
	// h.dead = true
}

// IsDead ....
func (h *Hero) IsDead() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.dead
}

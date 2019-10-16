package hero

import (
	"fmt"
	"math"
	"sync"
	
	"github.com/veandco/go-sdl2/sdl"

	"github.com/smeshkov/trovehero/pit"
)

const (
	gravity     = 0.1
	friction    = 1.5
	airFriction = 0.2
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

	dead bool
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
		maxMoveSpeed: 8,
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

	if _, ok := h.commands[t]; !ok {
		h.commands[t] = &Command{
			t:    t,
			ttl:  defaultTTL[t],
			time: h.time,
		}
	}
}

func (h *Hero) doCommand(cmd *Command) {
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
}

// Update updates state of the Hero.
func (h *Hero) Update() {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.time++

	// check commands
	jump := h.commands[Jump]
	for t, cmd := range h.commands {
		isMoving := t != Jump
		// performs command in following cases:
		// 1. Hero is standing on the ground
		// 2. Hero has jumped and the gap between Jump and Move commands is within the boundaries of the Move's TTL
		if h.altitude == 0 || (jump != nil && isMoving && math.Abs(float64(jump.time-cmd.time)) <= float64(cmd.ttl)) {
			h.doCommand(cmd)
			delete(h.commands, t)
		} else {
			cmd.ttl--
			if cmd.ttl == 0 {
				delete(h.commands, t)
			}
		}
	}

	h.handleJump()
	h.handleMove()
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
	if h.altitude > 0 {
		h.w = h.width + int32(h.altitude)
		h.h = h.height + int32(h.altitude)

		h.x = h.coordX - int32(h.altitude/2)
		h.y = h.coordY - int32(h.altitude/2)
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

func (h *Hero) handleJump() {
	// rising
	if h.jumpSpeed > 0 {
		// log.Println("rising...")
		h.altitude += h.jumpSpeed
		h.resize()
		h.jumpSpeed -= gravity
		return
	}

	// falling
	if h.altitude > 0 && h.jumpSpeed <= 0 {
		// log.Println("falling...")
		h.altitude = math.Max(0, h.altitude+h.jumpSpeed)
		h.resize()
		h.jumpSpeed -= gravity
		return
	}

	// landed
	if h.altitude == 0 && h.jumpSpeed < 0 {
		// log.Println("landed...")
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

	if p.X > h.x+h.w { // too far right
		return
	}
	if p.X+p.W < h.x { // too far left
		return
	}
	if p.Y > h.y+h.h { // too far below
		return
	}
	if p.Y+p.H < h.y { // to far above
		return
	}

	h.dead = true
}

// IsDead ....
func (h *Hero) IsDead() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.dead
}

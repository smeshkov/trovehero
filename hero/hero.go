package hero

import (
	"fmt"
	"math"
	"sync"

	"github.com/veandco/go-sdl2/sdl"
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
	r    *sdl.Renderer

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
}

// NewHero creates new instance of Hero.
func NewHero(r *sdl.Renderer) *Hero {
	var heroWidth int32 = 50
	var heroHeight int32 = 50
	var coordX int32 = 800/2 - heroWidth/2
	var coordY int32 = 600/2 - heroHeight/2

	return &Hero{
		// tools
		r: r,

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

	// log.Printf("performing %s...\n", c)
	if cmd, ok := h.commands[t]; ok {
		cmd.double = true
	} else {
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
		if cmd.double {
			h.vertSpeed = -h.maxMoveSpeed * 1.5
		} else {
			h.vertSpeed = -h.maxMoveSpeed
		}
	case Down:
		if cmd.double {
			h.vertSpeed = h.maxMoveSpeed * 1.5
		} else {
			h.vertSpeed = h.maxMoveSpeed
		}
	case Left:
		if cmd.double {
			h.horSpeed = -h.maxMoveSpeed * 1.5
		} else {
			h.horSpeed = -h.maxMoveSpeed
		}
	case Right:
		if cmd.double {
			h.horSpeed = h.maxMoveSpeed * 1.5
		} else {
			h.horSpeed = h.maxMoveSpeed
		}
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
func (h *Hero) Paint() error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// remove previous rectangle
	err := h.clearRect()
	if err != nil {
		return err
	}

	// fill new rectangle
	h.r.SetDrawColor(255, 100, 0, 255)
	h.rect = &sdl.Rect{X: h.x, Y: h.y, W: h.w, H: h.h}
	h.r.FillRect(h.rect)
	h.r.SetDrawColor(0, 0, 0, 255)

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

	h = NewHero(h.r)
}

// Destroy removes Hero.
func (h *Hero) Destroy() {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clearRect()
}

func (h *Hero) resize() {
	if h.altitude > 0 {
		h.w = h.width + int32(h.altitude)
		h.h = h.height + int32(h.altitude)

		h.x = h.coordX - int32(h.altitude/2)
		h.y = h.coordY - int32(h.altitude/2)
	}
}

func (h *Hero) clearRect() error {
	err := h.r.SetDrawColor(0, 0, 0, 0)
	if err != nil {
		return fmt.Errorf("could not set draw color: %w", err)
	}
	err = h.r.FillRect(h.rect)
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

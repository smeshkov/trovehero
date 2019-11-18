package enemy

import (
	"fmt"
	"log"
	"math"
	"os"
	"sync"

	"github.com/smeshkov/trovehero/world"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/smeshkov/trovehero/hero"
	"github.com/smeshkov/trovehero/types/command"
	"github.com/smeshkov/trovehero/types/direction"
	"github.com/smeshkov/trovehero/types/shape"
)

const (
	enemyMemory = 50
	enemyHeight = 50
	enemyWidth  = 50
	friction    = 0.2
	airFriction = 0.1
)

// Enemy attacks Hero.
type Enemy struct {
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
	altitude float64

	// shape
	x, y int32
	w, h int32

	// speed
	vertSpeed float32
	horSpeed  float32
	altSpeed  float32

	// AI
	sightDistnace int32
	sightWidth    int32
	direction     direction.Type
	player        *sdl.Rect
	enemyMemory   int32

	// World
	world *world.World
}

// NewEnemy creates new instance of Enemy in given coordinates.
func NewEnemy(x, y int32, w *world.World) *Enemy {
	var width int32 = enemyWidth
	var height int32 = enemyHeight
	var coordX int32 = x - width/2
	var coordY int32 = y - height/2

	return &Enemy{
		// properties
		height:       height,
		width:        width,
		maxMoveSpeed: 2,
		maxJumpSpeed: 1,

		// coordinates
		location: &sdl.Point{X: coordX, Y: coordY},

		// shape
		x: coordX,
		y: coordY,
		h: height,
		w: width,

		// AI
		sightDistnace: 50,
		sightWidth:    150,
		direction:     direction.West,

		// World
		world: w,
	}
}

func (e *Enemy) canSeeHero(hero *sdl.Rect) bool {
	var triangle [3]*sdl.Point

	switch e.direction {
	case direction.North:
		triangle = [3]*sdl.Point{
			e.location,
			&sdl.Point{X: e.location.X - e.sightWidth/2, Y: e.location.Y - e.sightDistnace},
			&sdl.Point{X: e.location.X + e.sightWidth/2, Y: e.location.Y - e.sightDistnace},
		}
	case direction.East:
		triangle = [3]*sdl.Point{
			e.location,
			&sdl.Point{X: e.location.X + e.sightDistnace, Y: e.location.Y - e.sightWidth/2},
			&sdl.Point{X: e.location.X + e.sightDistnace, Y: e.location.Y + e.sightWidth/2},
		}
	case direction.South:
		triangle = [3]*sdl.Point{
			e.location,
			&sdl.Point{X: e.location.X - e.sightWidth/2, Y: e.location.Y + e.sightDistnace},
			&sdl.Point{X: e.location.X + e.sightWidth/2, Y: e.location.Y + e.sightDistnace},
		}
	case direction.West:
		triangle = [3]*sdl.Point{
			e.location,
			&sdl.Point{X: e.location.X - e.sightDistnace, Y: e.location.Y - e.sightWidth/2},
			&sdl.Point{X: e.location.X - e.sightDistnace, Y: e.location.Y + e.sightWidth/2},
		}
	}

	// Vicinity of the enemy
	viewPort := shape.NewTriangle(triangle, nil)

	// Is hero in the vicinity of enemy
	return viewPort.OverlapsRect(hero)
}

func (e *Enemy) follow(x, y int32) direction.Type {
	if e.location.X > x {
		return direction.West
	}
	if e.location.X < x {
		return direction.East
	}
	if e.location.Y > y {
		return direction.North
	}
	return direction.South
}

// Watch checks if Enemy can see Hero.
func (e *Enemy) Watch(h *hero.Hero) {
	e.mu.Lock()
	defer e.mu.Unlock()

	heroLoc := h.Location()

	if e.canSeeHero(heroLoc) {
		log.Println("Enemy spotted Hero")
		// e.enemyMemory = enemyMemory
		e.direction = e.follow(heroLoc.X, heroLoc.Y)
	}
}

// Update updates state of the Enemy.
func (e *Enemy) Update() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.time++

	// if e.enemyMemory > 0 {
	// 	e.enemyMemory--
	// }

	e.direction = checkDirection(e.direction, e.sightDistnace, e.location.X, e.location.Y, e.world.H, e.world.W)

	if cmd, err := command.ToCommand(e.direction); err == nil {
		e.move(cmd)
	} else {
		fmt.Fprintf(os.Stderr, "enemy failed to convert direction to command: %v", err)
	}

	if e.horSpeed != 0 || e.vertSpeed != 0 {
		e.handleMove()
	}
	// if h.crashingDepth == 0 && h.altSpeed != 0 {
	// 	h.handleJump()
	// }
	// if h.crashingDepth != 0 {
	// 	h.handleCrash()
	// }
}

// move performes move command on an Enemy.
func (e *Enemy) move(t command.Type) {
	switch t {
	case command.GoNorth:
		e.vertSpeed = -e.maxMoveSpeed
	case command.GoSouth:
		e.vertSpeed = e.maxMoveSpeed
	case command.GoWest:
		e.horSpeed = -e.maxMoveSpeed
	case command.GoEast:
		e.horSpeed = e.maxMoveSpeed
	}
}

func (e *Enemy) handleMove() {
	var frict float64
	if e.altitude == 0 {
		frict = friction
	} else {
		frict = airFriction
	}

	if e.horSpeed != 0 {
		e.location.X += int32(e.horSpeed)
		e.x += int32(e.horSpeed)
		if e.horSpeed > 0 {
			e.horSpeed = float32(math.Max(0, float64(e.horSpeed)-frict))
		} else {
			e.horSpeed = float32(math.Min(0, float64(e.horSpeed)+frict))
		}
	}
	if e.vertSpeed != 0 {
		e.location.Y += int32(e.vertSpeed)
		e.y += int32(e.vertSpeed)
		if e.vertSpeed > 0 {
			e.vertSpeed = float32(math.Max(0, float64(e.vertSpeed)-frict))
		} else {
			e.vertSpeed = float32(math.Min(0, float64(e.vertSpeed)+frict))
		}
	}
}

// Paint paints Enemy to window.
func (e *Enemy) Paint(r *sdl.Renderer) error {
	e.mu.RLock()
	defer e.mu.RUnlock()

	// remove previous rectangle
	err := e.clearRect(r)
	if err != nil {
		return err
	}

	// fill new rectangle
	r.SetDrawColor(0, 128, 0, 255)
	e.rect = &sdl.Rect{X: e.x, Y: e.y, W: e.w, H: e.h}
	r.FillRect(e.rect)
	r.SetDrawColor(0, 0, 0, 255)

	return nil
}

func (e *Enemy) clearRect(r *sdl.Renderer) error {
	err := r.SetDrawColor(0, 0, 0, 0)
	if err != nil {
		return fmt.Errorf("could not set draw color: %w", err)
	}
	err = r.FillRect(e.rect)
	if err != nil {
		return fmt.Errorf("could not fill rectangle: %w", err)
	}
	return nil
}

// Restart restarts state of Enemy.
func (e *Enemy) Restart() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e = NewEnemy(30, 30, e.world)
}

// Destroy removes Enemy.
func (e *Enemy) Destroy() {
	// noop
}

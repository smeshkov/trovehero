package enemy

import (
	"fmt"
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

	// properties
	height       int32
	width        int32
	maxMoveSpeed float32
	maxJumpSpeed float32

	// coordinates
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
	// player        *sdl.Rect
	// enemyMemory   int32

	// World
	world *world.World
}

// NewEnemy creates new instance of Enemy in given coordinates.
func NewEnemy(x, y int32, w *world.World) *Enemy {
	var width int32 = enemyWidth
	var height int32 = enemyHeight

	return &Enemy{
		// properties
		height:       height,
		width:        width,
		maxMoveSpeed: 2,
		maxJumpSpeed: 1,

		// shape
		x: x,
		y: y,
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

	location := &sdl.Point{X: e.x + e.width/2, Y: e.y + e.height/2}

	// create triangle view, depending on which direction is facing
	switch e.direction {
	case direction.North:
		triangle = [3]*sdl.Point{
			location,
			&sdl.Point{X: location.X - e.sightWidth/2, Y: location.Y - e.sightDistnace},
			&sdl.Point{X: location.X + e.sightWidth/2, Y: location.Y - e.sightDistnace},
		}
	case direction.East:
		triangle = [3]*sdl.Point{
			location,
			&sdl.Point{X: location.X + e.sightDistnace, Y: location.Y - e.sightWidth/2},
			&sdl.Point{X: location.X + e.sightDistnace, Y: location.Y + e.sightWidth/2},
		}
	case direction.South:
		triangle = [3]*sdl.Point{
			location,
			&sdl.Point{X: location.X - e.sightWidth/2, Y: location.Y + e.sightDistnace},
			&sdl.Point{X: location.X + e.sightWidth/2, Y: location.Y + e.sightDistnace},
		}
	case direction.West:
		triangle = [3]*sdl.Point{
			location,
			&sdl.Point{X: location.X - e.sightDistnace, Y: location.Y - e.sightWidth/2},
			&sdl.Point{X: location.X - e.sightDistnace, Y: location.Y + e.sightWidth/2},
		}
	}

	// Vicinity of the enemy
	viewPort := shape.NewTriangle(triangle, nil)

	// Is hero in the vicinity of enemy
	return viewPort.OverlapsRect(hero)
}

func (e *Enemy) directTo(x, y int32) {
	if e.x+e.width > x {
		e.direction = direction.West
	}
	if e.x < x {
		e.direction = direction.East
	}
	if e.y+e.height > y {
		e.direction = direction.North
	}
	if e.y < y {
		e.direction = direction.South
	}
}

func (e *Enemy) directionCheck() {

	for i := 0; i < 3; i++ {

		if e.direction == direction.North && e.y-e.sightDistnace <= 0 {
			e.direction = direction.East
		} else if e.direction == direction.East && e.x+e.width+e.sightDistnace >= e.world.W {
			e.direction = direction.South
		} else if e.direction == direction.South && e.y+e.height+e.sightDistnace >= e.world.H {
			e.direction = direction.West
		} else if e.direction == direction.West && e.x-e.sightDistnace <= 0 {
			e.direction = direction.North
		}

	}
}

// Watch checks if Enemy can see Hero.
func (e *Enemy) Watch(h *hero.Hero) {
	e.mu.Lock()
	defer e.mu.Unlock()

	heroLoc := h.Location()

	if e.canSeeHero(heroLoc) {
		e.directTo(heroLoc.X, heroLoc.Y)
	}
}

// Update updates state of the Enemy.
func (e *Enemy) Update() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.time++

	e.directionCheck()

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
		e.x += int32(e.horSpeed)
		if e.horSpeed > 0 {
			e.horSpeed = float32(math.Max(0, float64(e.horSpeed)-frict))
		} else {
			e.horSpeed = float32(math.Min(0, float64(e.horSpeed)+frict))
		}
	}
	if e.vertSpeed != 0 {
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

	// fill new rectangle
	r.SetDrawColor(0, 128, 0, 255)
	r.FillRect(e.getShape())
	r.SetDrawColor(0, 0, 0, 255)

	return nil
}

func (e *Enemy) getShape() *sdl.Rect {
	if e.altitude != 0 {
		return &sdl.Rect{
			X: e.x - int32(e.altitude/2),
			Y: e.y - int32(e.altitude/2),
			W: e.w + int32(e.altitude),
			H: e.h + int32(e.altitude),
		}
	}
	return &sdl.Rect{
		X: e.x,
		Y: e.y,
		W: e.w,
		H: e.h,
	}
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

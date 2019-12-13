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

	ID string

	time int64

	// properties
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

	// World
	world *world.World
}

// NewEnemy creates new instance of Enemy in given coordinates.
func NewEnemy(id string, x, y int32, w *world.World) *Enemy {
	e := &Enemy{ID: id}
	return e.setDefaults(x, y, enemyWidth, enemyHeight, w)
}

func (e *Enemy) setDefaults(x, y, width, height int32, w *world.World) *Enemy {
	e.time = 0

	e.maxMoveSpeed = 2
	e.maxJumpSpeed = 1

	e.altitude = 0

	// shape
	e.x = x
	e.y = y
	e.h = height
	e.w = width

	// AI
	e.sightDistnace = 150
	e.sightWidth = 350
	e.direction = direction.Type(w.Rand.Int31n(3))

	// World
	e.world = w

	return e
}

func (e *Enemy) canSeeHero(hero *sdl.Rect) bool {
	var triangle [3]*sdl.Point

	location := &sdl.Point{X: e.x + e.w/2, Y: e.y + e.h/2}

	// create triangle view, depending on which direction is facing
	switch e.direction {
	case direction.North:
		triangle = [3]*sdl.Point{
			location,
			{X: location.X - e.sightWidth/2, Y: location.Y - e.sightDistnace},
			{X: location.X + e.sightWidth/2, Y: location.Y - e.sightDistnace},
		}
	case direction.East:
		triangle = [3]*sdl.Point{
			location,
			{X: location.X + e.sightDistnace, Y: location.Y - e.sightWidth/2},
			{X: location.X + e.sightDistnace, Y: location.Y + e.sightWidth/2},
		}
	case direction.South:
		triangle = [3]*sdl.Point{
			location,
			{X: location.X - e.sightWidth/2, Y: location.Y + e.sightDistnace},
			{X: location.X + e.sightWidth/2, Y: location.Y + e.sightDistnace},
		}
	case direction.West:
		triangle = [3]*sdl.Point{
			location,
			{X: location.X - e.sightDistnace, Y: location.Y - e.sightWidth/2},
			{X: location.X - e.sightDistnace, Y: location.Y + e.sightWidth/2},
		}
	}

	// Vicinity of the enemy
	viewPort := shape.NewTriangle(triangle, nil)

	// Is hero in the vicinity of enemy
	return viewPort.OverlapsRect(hero)
}

func (e *Enemy) directTo(x, y int32) {
	if e.x+e.w > x {
		e.direction = direction.West
	}
	if e.x < x {
		e.direction = direction.East
	}
	if e.y+e.h > y {
		e.direction = direction.North
	}
	if e.y < y {
		e.direction = direction.South
	}
}

func (e *Enemy) directionCheck() {

	var changed bool

	for {
		changed = false

		if e.direction == direction.North && (e.y-e.sightDistnace/4) <= 0 {
			e.direction = direction.East
			changed = true
		}
		if e.direction == direction.East && (e.x+e.w+e.sightDistnace/4) >= e.world.W {
			e.direction = direction.South
			changed = true
		}
		if e.direction == direction.South && (e.y+e.h+e.sightDistnace/4) >= e.world.H {
			e.direction = direction.West
			changed = true
		}
		if e.direction == direction.West && (e.x-e.sightDistnace/4) <= 0 {
			e.direction = direction.North
			changed = true
		}

		if !changed {
			// we are done here, hence no changes happened
			break
		}
	}

}

// Touch checks collision with Pit.
func (e *Enemy) Touch(h *hero.Hero) {
	e.mu.Lock()
	defer e.mu.Unlock()

	heroLoc := h.Location()

	if e.x > heroLoc.X+heroLoc.W { // too far right
		return
	}
	if e.x+e.w < heroLoc.X { // too far left
		return
	}
	if e.y > heroLoc.Y+heroLoc.H { // too far below
		return
	}
	if e.y+e.h < heroLoc.Y { // to far above
		return
	}

	// Hero is busted, so he dies
	h.Die()
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
	r.SetDrawColor(160, 0, 0, 255)
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
	pos := e.world.RandomizePos(e.ID, enemyWidth, enemyHeight)
	e.setDefaults(pos.X, pos.Y, pos.W, pos.H, e.world)
}

// Destroy removes Enemy.
func (e *Enemy) Destroy() {
	// noop
}

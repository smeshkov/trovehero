package enemy

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/veandco/go-sdl2/sdl"

	"github.com/smeshkov/trovehero/types/direction"
)

func Test_canSeeHero_true(t *testing.T) {
	heroLoc := &sdl.Rect{X: 100, Y: 80, W: 50, H: 50}

	e := newEnemy()

	canSee := e.canSeeHero(heroLoc)

	assert.True(t, canSee)
}

func Test_canSeeHero_false(t *testing.T) {
	heroLoc := &sdl.Rect{X: 100, Y: 200, W: 50, H: 50}

	e := newEnemy()

	canSee := e.canSeeHero(heroLoc)

	assert.False(t, canSee)
}

func newEnemy() *Enemy {
	return &Enemy{
		sightDistnace: 50,
		sightWidth:    150,
		direction:     direction.North,
		location:      &sdl.Point{X: 100, Y: 100},
	}
}

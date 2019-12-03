package enemy

import (
	"github.com/smeshkov/trovehero/world"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/veandco/go-sdl2/sdl"

	"github.com/smeshkov/trovehero/types/direction"
)

func newEnemy() *Enemy {
	return &Enemy{
		sightDistnace: 50,
		sightWidth:    150,
		direction:     direction.North,
		x:             100,
		y:             100,
	}
}

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

const testSightDistnace int32 = 50

type directionCheckTest struct {
	name               string
	x, y, areaH, areaW int32
	input              direction.Type
	expected           direction.Type
}

func Test_directionCheck(t *testing.T) {
	tests := []directionCheckTest{
		{
			name:  "no change in direction to South",
			x:     50,
			y:     50,
			areaH: 200,
			areaW: 200,

			input:    direction.South,
			expected: direction.South,
		},
		{
			name:  "change direction from South to West",
			x:     100,
			y:     50,
			areaH: 100,
			areaW: 200,

			input:    direction.South,
			expected: direction.West,
		},
		{
			name:  "change direction from West to North",
			x:     50,
			y:     100,
			areaH: 200,
			areaW: 100,

			input:    direction.West,
			expected: direction.North,
		},
		{
			name:  "change direction from North to East",
			x:     50,
			y:     50,
			areaH: 100,
			areaW: 200,

			input:    direction.North,
			expected: direction.East,
		},
		{
			name:  "change direction from East to South",
			x:     50,
			y:     50,
			areaH: 200,
			areaW: 100,

			input:    direction.East,
			expected: direction.South,
		},
		{
			name:  "left top corner: change direction from West to East",
			x:     50,
			y:     50,
			areaH: 200,
			areaW: 200,

			input:    direction.West,
			expected: direction.East,
		},
		{
			name:  "right top corner: change direction from North to South",
			x:     50,
			y:     50,
			areaH: 200,
			areaW: 100,

			input:    direction.North,
			expected: direction.South,
		},
		{
			name:  "right bottom corner: change direction from East to West",
			x:     150,
			y:     150,
			areaH: 200,
			areaW: 200,

			input:    direction.East,
			expected: direction.West,
		},
		{
			name:  "left bottom corner: change direction from South to North",
			x:     50,
			y:     150,
			areaH: 200,
			areaW: 200,

			input:    direction.South,
			expected: direction.North,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Enemy{
				sightDistnace: testSightDistnace,
				direction:     tt.input,
				x:             tt.x,
				y:             tt.y,
				world:         world.NewWorld(tt.areaW, tt.areaH, nil),
			}
			e.directionCheck()
			assert.Equal(t, tt.expected, e.direction)
		})
	}
}

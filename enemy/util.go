package enemy

import (
	"github.com/veandco/go-sdl2/sdl"

	"github.com/smeshkov/trovehero/types"
)

func followHero(enemy *sdl.Point, hero *sdl.Rect) types.Direction {
	if enemy.X > hero.X {
		return types.West
	}
	if enemy.X < hero.X {
		return types.East
	}
	if enemy.Y > hero.Y {
		return types.North
	}
	return types.South
}

func directionCheck(ed types.Direction, sight, eX, eY, wH, wW int32) types.Direction {
	if ed == types.North && eY-enemyHeight/2-sight <= 0 {
		return types.East
	}
	if ed == types.East && eX+enemyWidth/2+sight >= wW {
		return types.South
	}
	if ed == types.South && eY+enemyHeight/2+sight >= wH {
		return types.West
	}
	if ed == types.West && eX-enemyHeight/2-sight <= 0 {
		return types.North
	}
	return ed
}

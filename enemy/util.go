package enemy

import (
	"github.com/smeshkov/trovehero/types/direction"
)

func checkDirection(ed direction.Type, sight, eX, eY, wH, wW int32) direction.Type {
	d := ed

	for i := 0; i < 3; i++ {

		if d == direction.North && eY-enemyHeight/2-sight <= 0 {
			d = direction.East
		} else if d == direction.East && eX+enemyWidth/2+sight >= wW {
			d = direction.South
		} else if d == direction.South && eY+enemyHeight/2+sight >= wH {
			d = direction.West
		} else if d == direction.West && eX-enemyHeight/2-sight <= 0 {
			d = direction.North
		}

	}

	return d
}

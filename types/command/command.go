package command

import (
	"fmt"

	"github.com/smeshkov/trovehero/types/direction"
)

const (
	// GoNorth makes Hero to move North.
	GoNorth Type = iota
	// GoEast makes Hero to move East.
	GoEast
	// GoSouth makes Hero to move South.
	GoSouth
	// GoWest makes Hero to move West.
	GoWest
	// Jump direction
	Jump
	// Shoot makes Hero to shoot.
	Shoot
)

var (
	typeNames = map[Type]string{
		GoNorth: "GoNorth",
		GoEast:  "GoEast",
		GoSouth: "GoSouth",
		GoWest:  "GoWest",
		Jump:    "Jump",
		Shoot:   "Shoot",
	}
)

// Type is a type of a command.
type Type byte

// FromDirection transforms a direction type into a command type.
func FromDirection(d direction.Type) (Type, error) {
	switch d {
	case direction.North:
		return GoNorth, nil
	case direction.East:
		return GoEast, nil
	case direction.South:
		return GoSouth, nil
	case direction.West:
		return GoWest, nil
	}
	return 0, fmt.Errorf("Unknown direction %s with code %d", d, d)
}

func (t Type) String() string {
	if t < GoNorth || t > Shoot {
		return "Unknown"
	}
	return typeNames[t]
}

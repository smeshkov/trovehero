package types

import (
	"fmt"
)

const (
	// GoNorth makes Hero to move North.
	GoNorth CommandType = iota
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
	commandTypeNames = map[CommandType]string{
		GoNorth: "GoNorth",
		GoEast:  "GoEast",
		GoSouth: "GoSouth",
		GoWest:  "GoWest",
		Jump:    "Jump",
		Shoot:   "Shoot",
	}
)

// CommandType is a type of a command.
type CommandType byte

// ToCommand transforms a Direction into a CommandType.
func ToCommand(d Direction) (CommandType, error) {
	switch d {
	case North:
		return GoNorth, nil
	case East:
		return GoEast, nil
	case South:
		return GoSouth, nil
	case West:
		return GoWest, nil
	}
	return 0, fmt.Errorf("Unknown direction %s with code %d", d, d)
}

func (ct CommandType) String() string {
	if ct < GoNorth || ct > Shoot {
		return "Unknown"
	}
	return commandTypeNames[ct]
}

package hero

const (
	// Jump direction
	Jump CommandType = iota
	// Up makes Hero to move up.
	Up
	// Down makes Hero to move down.
	Down
	// Left makes Hero to move left.
	Left
	// Right makes Hero to move left.
	Right
	// Shoot makes Hero to shoot left.
	Shoot
)

var (
	commandTypeNames = map[CommandType]string{
		Jump:  "Jump",
		Up:    "Up",
		Down:  "Down",
		Left:  "Left",
		Right: "Right",
		Shoot: "Shoot",
	}

	defaultTTL = map[CommandType]byte{
		Jump:  10,
		Up:    10,
		Down:  10,
		Left:  10,
		Right: 10,
	}
)

// CommandType is a type of a command.
type CommandType byte

func (ct CommandType) String() string {
	if ct < Jump || ct > Shoot {
		return "Unknown"
	}
	return commandTypeNames[ct]
}

type commands map[CommandType]*Command

// Command is an arbitrary command.
type Command struct {
	t      CommandType
	time   int64
	ttl    byte
}

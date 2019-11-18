package direction

const (
	// North direction.
	North Type = iota
	// East direction.
	East
	// South direction.
	South
	// West direction.
	West
)

var (
	typeNames = map[Type]string{
		North: "North",
		East:  "East",
		South: "South",
		West:  "West",
	}
)

// Type is a type of direction.
type Type byte

func (t Type) String() string {
	if t < North || t > West {
		return "Unknown"
	}
	return typeNames[t]
}

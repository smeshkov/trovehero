package types

const (
	// North direction.
	North Direction = iota
	// East direction.
	East
	// South direction.
	South
	// West direction.
	West
)

var (
	directionNames = map[Direction]string{
		North: "North",
		East:  "East",
		South: "South",
		West:  "West",
	}
)

// Direction is a direction.
type Direction byte

func (d Direction) String() string {
	if d < North || d > West {
		return "Unknown"
	}
	return directionNames[d]
}

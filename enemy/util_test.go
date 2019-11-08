package enemy

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smeshkov/trovehero/types"
)

func Test_directionCheck(t *testing.T) {
	sightDistnace := int32(50)
	oldDirection := types.South
	
	newDirection := directionCheck(oldDirection, sightDistnace, 50, 50, 200, 200)

	assert.Equal(t, oldDirection, newDirection)
}

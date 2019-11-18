package shape

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/veandco/go-sdl2/sdl"
)

func Test_Triangle_ContainsPoint_1(t *testing.T) {
	point := &sdl.Point{X: 50, Y: 50}
	triangle := NewTriangle([3]*sdl.Point{{X: 50, Y: 25}, {X: 25, Y: 75}, {X: 75, Y: 75}}, nil)

	containsPoint := triangle.ContainsPoint(point)

	assert.True(t, containsPoint)
}

func Test_Triangle_ContainsPoint_2(t *testing.T) {
	point := &sdl.Point{X: 50, Y: 50}
	triangle := NewTriangle([3]*sdl.Point{{X: 25, Y: 25}, {X: 75, Y: 25}, {X: 50, Y: 75}}, nil)

	containsPoint := triangle.ContainsPoint(point)

	assert.True(t, containsPoint)
}

func Test_Triangle_ContainsPoint_3(t *testing.T) {
	point := &sdl.Point{X: 50, Y: 75}
	triangle := NewTriangle([3]*sdl.Point{{X: 25, Y: 25}, {X: 75, Y: 25}, {X: 50, Y: 75}}, nil)

	containsPoint := triangle.ContainsPoint(point)

	assert.False(t, containsPoint)
}

func Test_Triangle_OverlapsRect_1(t *testing.T) {
	rect := &sdl.Rect{X: 0, Y: 0, W: 50, H: 50}
	triangle := NewTriangle([3]*sdl.Point{{X: 25, Y: 25}, {X: 75, Y: 25}, {X: 50, Y: 75}}, nil)

	containsPoint := triangle.OverlapsRect(rect)

	assert.True(t, containsPoint)
}

func Test_Triangle_OverlapsRect_2(t *testing.T) {
	rect := &sdl.Rect{X: 295, Y: 295, W: 10, H: 10}
	triangle := NewTriangle([3]*sdl.Point{{X: 305, Y: 0}, {X: 0, Y: 600}, {X: 600, Y: 600}}, nil)

	containsPoint := triangle.OverlapsRect(rect)

	assert.True(t, containsPoint)
}

func Test_Triangle_OverlapsRect_3(t *testing.T) {
	rect := &sdl.Rect{X: 0, Y: 0, W: 600, H: 600}
	triangle := NewTriangle([3]*sdl.Point{{X: 290, Y: 310}, {X: 300, Y: 290}, {X: 310, Y: 310}}, nil)

	containsPoint := triangle.OverlapsRect(rect)

	assert.True(t, containsPoint)
}
package types

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

// Object describes default API of the scene object.
type Object interface {
	Update()
	Paint(r *sdl.Renderer) error
	Restart()
	Destroy()
}

// Triangle represent a triangle shape.
type Triangle struct {
	ps    [3]*sdl.Point
	color *sdl.Color
}

// NewTriangle creates a new triangle from given slice of Points.
func NewTriangle(points [3]*sdl.Point, color *sdl.Color) *Triangle {
	t := &Triangle{ps: points}
	return t
}

func (t *Triangle) String() string {
	return fmt.Sprintf("{[%d, %d], [%d, %d], [%d, %d]}",
		t.ps[0].X, t.ps[0].Y, t.ps[1].X, t.ps[1].Y, t.ps[2].X, t.ps[2].Y)
}

// ContainsPoint returns true of the given Point is inside the Triangle.
func (t *Triangle) ContainsPoint(point *sdl.Point) bool {
	a := t.ps[0]
	b := t.ps[1]
	c := t.ps[2]

	asX := point.X - a.X
	asY := point.Y - a.Y

	sAB := (b.X-a.X)*asY-(b.Y-a.Y)*asX > 0

	if (c.X-a.X)*asY-(c.Y-a.Y)*asX > 0 == sAB {
		return false
	}

	if (c.X-b.X)*(point.Y-b.Y)-(c.Y-b.Y)*(point.X-b.X) > 0 != sAB {
		return false
	}

	return true
}

// OverlapsRect returns true of the given Rect overlaps with the Triangle.
func (t *Triangle) OverlapsRect(rect *sdl.Rect) bool {
	a := t.ps[0]
	b := t.ps[1]
	c := t.ps[2]

	rectX := rect.X
	rectY := rect.Y
	rectW := rect.W
	rectH := rect.H

	return rect.IntersectLine(&a.X, &a.Y, &b.X, &b.Y) ||
		rect.IntersectLine(&b.X, &b.Y, &c.X, &c.Y) ||
		rect.IntersectLine(&c.X, &c.Y, &a.X, &a.Y) ||
		t.ContainsPoint(&sdl.Point{X: rectX, Y: rectY}) ||
		t.ContainsPoint(&sdl.Point{X: rectX + rectW, Y: rectY}) ||
		t.ContainsPoint(&sdl.Point{X: rectX + rectW, Y: rectY + rectH}) ||
		t.ContainsPoint(&sdl.Point{X: rectX, Y: rectY + rectH})
}

// Paint paints the Triangle to the window.
func (t *Triangle) Paint(r *sdl.Renderer) error {
	// Set color of triangle
	if t.color != nil {
		r.SetDrawColor(t.color.R, t.color.G, t.color.B, t.color.A)
	}

	// rs := ToRectangles(t.ps[:])
	//tr := t.triangle(rs)

	// r.FillRect(tr)

	// Reset colour
	r.SetDrawColor(0, 0, 0, 255)

	return nil
}

package scene

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	ttf "github.com/veandco/go-sdl2/ttf"
)

/* func getRandomPosition(previous []*sdl.Rect, h, w int32, world *world.World) (x, y int32) {
	x = rand.Int31n(world.W - w)
	y = rand.Int31n(world.H - h)

	var passed bool

	for !passed {
		passed = true
		for _, p := range previous {
			// if p.
			// if x == p.X || x == p.X + p.W || x > p.X && x < p.X + p.W {
			// 	passed = false
			// }
			// if x+w == p.X || x+w == p.X + p.W || x+w > p.X && x+w < p.X + p.W {
			// 	passed = false
			// }
			// if y == p.Y || y == p.Y + p.H || y > p.Y && y < p.Y + p.H {
			// 	passed = false
			// }
			// if x+w == p.X || x+w == p.X + p.W || x+w > p.X && x+w < p.X + p.W {
			// 	passed = false
			// }
		}
	}
} */

// DrawTitle draws a title with given "text".
func DrawTitle(r *sdl.Renderer, text string) error {
	if err := r.Clear(); err != nil {
		return fmt.Errorf("could not clear renderer: %w", err)
	}

	f, err := ttf.OpenFont("res/fonts/Flappy.ttf", 20)
	if err != nil {
		return fmt.Errorf("could not load font: %w", err)
	}
	defer f.Close()

	c := sdl.Color{R: 255, G: 100, B: 0, A: 255}
	s, err := f.RenderUTF8Solid(text, c)
	if err != nil {
		return fmt.Errorf("could not render title: %w", err)
	}
	defer s.Free()

	t, err := r.CreateTextureFromSurface(s)
	if err != nil {
		return fmt.Errorf("could not create texture: %w", err)
	}
	defer t.Destroy()

	if err := r.Copy(t, nil, nil); err != nil {
		return fmt.Errorf("could not copy texture: %w", err)
	}

	r.Present()

	return nil
}

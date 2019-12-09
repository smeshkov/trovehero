package scene

import (
	"fmt"

	"github.com/smeshkov/trovehero/trove"
	"github.com/smeshkov/trovehero/world"

	"github.com/veandco/go-sdl2/sdl"
	ttf "github.com/veandco/go-sdl2/ttf"
)

// DrawTitle draws a title with given "text".
func DrawTitle(r *sdl.Renderer, text string, color *sdl.Color) error {
	if err := r.Clear(); err != nil {
		return fmt.Errorf("could not clear renderer: %w", err)
	}

	f, err := ttf.OpenFont("res/fonts/Flappy.ttf", 20)
	if err != nil {
		return fmt.Errorf("could not load font: %w", err)
	}
	defer f.Close()

	// c := sdl.Color{R: 255, G: 100, B: 0, A: 255}
	s, err := f.RenderUTF8Solid(text, *color)
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

func createTroves(w *world.World, num int) []*trove.Trove {
	trs := make([]*trove.Trove, num)
	for i := 0; i < num; i++ {
		id := fmt.Sprintf("trove-%d", i)
		pos := w.RandomizePos(id, 50, 50)
		trs[i] = trove.NewTrove(id, pos.X, pos.Y, w)
	}
	return trs
}

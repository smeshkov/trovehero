package scene

import (
	"fmt"
	"math"

	"github.com/veandco/go-sdl2/sdl"
	ttf "github.com/veandco/go-sdl2/ttf"

	"github.com/smeshkov/trovehero/enemy"
	"github.com/smeshkov/trovehero/pit"
	"github.com/smeshkov/trovehero/trove"
	"github.com/smeshkov/trovehero/world"
)

// drawTitle draws a title with given "text".
func drawTitle(r *sdl.Renderer, text string, color *sdl.Color) error {
	if err := r.Clear(); err != nil {
		return fmt.Errorf("could not clear renderer: %w", err)
	}

	f, err := ttf.OpenFont("res/fonts/Flappy.ttf", 10)
	if err != nil {
		return fmt.Errorf("could not load font: %w", err)
	}
	defer f.Close()

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

func drawStats(w *world.World) error {
	fmt.Printf("Your score is %d, you've reached level %d\n",
		w.GetScore(), w.GetLevel())
	return nil
}

func createPits(w *world.World, num int8) []*pit.Pit {
	items := make([]*pit.Pit, num)
	var i int8
	for i = 0; i < num; i++ {
		id := fmt.Sprintf("pit-%d", i)
		width := int32(math.Max(30, float64(w.Rand.Int31n(150))))
		height := int32(math.Max(30, float64(w.Rand.Int31n(150))))
		pos := w.RandomizePos(id, width, height)
		items[i] = pit.NewPit(id, pos.X, pos.Y, width, height, int8(w.Rand.Int31n(100)), w)
	}
	return items
}

func createTroves(w *world.World, num int8) []*trove.Trove {
	items := make([]*trove.Trove, num)
	var i int8
	for i = 0; i < num; i++ {
		id := fmt.Sprintf("trove-%d", i)
		pos := w.RandomizePos(id, 50, 50)
		items[i] = trove.NewTrove(id, pos.X, pos.Y, w)
	}
	return items
}

func createEnemies(w *world.World, num int8) []*enemy.Enemy {
	items := make([]*enemy.Enemy, num)
	var i int8
	for i = 0; i < num; i++ {
		id := fmt.Sprintf("enemy-%d", i)
		pos := w.RandomizePos(id, 50, 50)
		items[i] = enemy.NewEnemy(id, pos.X, pos.Y, w)
	}
	return items
}

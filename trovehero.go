package trovehero

import (
	"fmt"
	"runtime"

	"github.com/veandco/go-sdl2/sdl"
	ttf "github.com/veandco/go-sdl2/ttf"

	"github.com/smeshkov/trovehero/scene"
)

// Run starts the game.
func Run(level int8) error {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return fmt.Errorf("could not initialize SDL: %w", err)
	}
	defer sdl.Quit()

	if err := ttf.Init(); err != nil {
		return fmt.Errorf("could not initialize TTF: %w", err)
	}
	defer ttf.Quit()

	w, r, err := sdl.CreateWindowAndRenderer(1280, 720, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("could not create window: %w", err)
	}
	defer w.Destroy()

	s, err := scene.NewScene(r, level)
	if err != nil {
		return fmt.Errorf("could not create scene: %w", err)
	}
	defer s.Destroy()

	events := make(chan sdl.Event)
	errc := s.Run(events, r)

	runtime.LockOSThread()
	for {
		select {
		case events <- sdl.WaitEvent():
		case err := <-errc:
			return err
		}
	}
}

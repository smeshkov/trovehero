package scene

import (
	"fmt"
	"log"
	"time"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/smeshkov/trovehero/enemy"
	"github.com/smeshkov/trovehero/hero"
	"github.com/smeshkov/trovehero/pit"
	"github.com/smeshkov/trovehero/types/command"
	"github.com/smeshkov/trovehero/world"
)

// Scene represent the scene of the game.
type Scene struct {
	hero  *hero.Hero
	pit   *pit.Pit
	enemy *enemy.Enemy
}

// NewScene returns new instance of the Scene.
func NewScene(r *sdl.Renderer) (*Scene, error) {
	// bg, err := img.LoadTexture(r, "res/imgs/background.png")
	// if err != nil {
	// 	return nil, fmt.Errorf("could not load background image: %w", err)
	// }

	viewPort := r.GetViewport()

	w := world.NewWorld(viewPort.W, viewPort.H, &viewPort)

	// used for storing ID of the object
	var id string

	// used for storing position of the object
	var pos *sdl.Rect

	id = "pit"
	pos = w.RandomizePos(id, 150, 50)
	p := pit.NewPit(id, pos.X, pos.Y, pos.W, pos.H, -60, w)

	id = "hero"
	pos = w.RandomizePos(id, 50, 50)
	h := hero.NewHero(id, pos.X, pos.Y, w)

	id = "enemy"
	pos = w.RandomizePos(id, 50, 50)
	e := enemy.NewEnemy(id, pos.X, pos.Y, w)

	return &Scene{hero: h, pit: p, enemy: e}, nil
}

// Run runs the Scene.
func (s *Scene) Run(events <-chan sdl.Event, r *sdl.Renderer) <-chan error {
	errc := make(chan error)

	go func() {
		defer close(errc)
		tick := time.Tick(10 * time.Millisecond)

		if err := DrawTitle(r, "Trove Hero"); err != nil {
			errc <- fmt.Errorf("could not draw title: %w", err)
		}
		time.Sleep(1 * time.Second)

		for {
			select {
			case e := <-events:
				if done := s.handleEvent(e); done {
					return
				}
			case <-tick:
				s.update()

				if s.hero.IsDead() {
					if err := DrawTitle(r, "Game Over"); err != nil {
						errc <- err
					}
					time.Sleep(1 * time.Second)
					s.restart()
				}

				if err := s.paint(r); err != nil {
					errc <- err
				}
			}
		}
	}()

	return errc
}

// handleEvent handles event and returns true if the app needs to finish execution and quite
// or false to signal to continue execution.
func (s *Scene) handleEvent(event sdl.Event) bool {
	switch event.(type) {
	case *sdl.QuitEvent:
		return true
	case *sdl.KeyboardEvent:
		return s.handleKeyboardEvent(event.(*sdl.KeyboardEvent))
	case *sdl.MouseMotionEvent, *sdl.WindowEvent, *sdl.TouchFingerEvent,
		*sdl.CommonEvent, *sdl.AudioDeviceEvent, *sdl.TextInputEvent:
	default:
		log.Printf("unknown event %T", event)
	}
	return false
}

// handleKeyboardEvent handles keyboard input event and returns true in case of exit or
// false for any other case.
func (s *Scene) handleKeyboardEvent(event *sdl.KeyboardEvent) bool {
	switch event.Keysym.Scancode {
	case sdl.SCANCODE_ESCAPE:
		return true
	case sdl.SCANCODE_SPACE:
		s.hero.Do(command.Jump)
	case sdl.SCANCODE_LEFT:
		s.hero.Do(command.GoWest)
	case sdl.SCANCODE_RIGHT:
		s.hero.Do(command.GoEast)
	case sdl.SCANCODE_UP:
		s.hero.Do(command.GoNorth)
	case sdl.SCANCODE_DOWN:
		s.hero.Do(command.GoSouth)
	}
	return false
}

func (s *Scene) update() {
	s.hero.Update()
	s.enemy.Update()
	s.pit.Update()
	s.hero.Touch(s.pit)
	s.enemy.Touch(s.hero)
	s.enemy.Watch(s.hero)
}

func (s *Scene) restart() {
	s.hero.Restart()
	s.pit.Restart()
	s.enemy.Restart()
}

func (s *Scene) paint(r *sdl.Renderer) error {
	r.Clear()

	if err := s.pit.Paint(r); err != nil {
		return err
	}
	if err := s.hero.Paint(r); err != nil {
		return err
	}
	if err := s.enemy.Paint(r); err != nil {
		return err
	}

	r.Present()
	return nil
}

// Destroy destroys the scene.
func (s *Scene) Destroy() {
	s.pit.Destroy()
	s.hero.Destroy()
	s.enemy.Destroy()
}

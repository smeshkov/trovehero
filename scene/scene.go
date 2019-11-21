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

	w := world.NewWorld(1024, 768, &sdl.Rect{W: 1024, H: 768, X: 0, Y: 0})

	p := pit.NewPit(w.W/2, 250, 50, 150, -60)
	h := hero.NewHero(w.W/2, 700)
	e := enemy.NewEnemy(30, 30, w)

	return &Scene{ hero: h, pit: p, enemy: e}, nil
}

// Run runs the Scene.
func (s *Scene) Run(events <-chan sdl.Event, r *sdl.Renderer) <-chan error {
	errc := make(chan error)

	go func() {
		defer close(errc)
		tick := time.Tick(10 * time.Millisecond)
		for {
			select {
			case e := <-events:
				if done := s.handleEvent(e); done {
					return
				}
			case <-tick:
				s.update()

				if s.hero.IsDead() {
					fmt.Println("drawing Game Over")
					if err := DrawTitle(r, "Game Over"); err != nil {
						errc <- err	
					}
					time.Sleep(1 * time.Second)
					fmt.Println("restarting")
					s.restart()
					fmt.Println("restarted")
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
	case *sdl.MouseMotionEvent, *sdl.WindowEvent, *sdl.TouchFingerEvent, *sdl.CommonEvent, *sdl.AudioDeviceEvent:
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

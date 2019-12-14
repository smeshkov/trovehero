package scene

import (
	"fmt"
	"log"
	"time"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/smeshkov/trovehero/enemy"
	"github.com/smeshkov/trovehero/hero"
	"github.com/smeshkov/trovehero/pit"
	"github.com/smeshkov/trovehero/trove"
	"github.com/smeshkov/trovehero/types/command"
	"github.com/smeshkov/trovehero/world"
)

var (
	// text colors
	orangeClr = &sdl.Color{R: 255, G: 100, B: 0, A: 255}
	redClr    = &sdl.Color{R: 210, G: 0, B: 0, A: 255}
	greenClr  = &sdl.Color{R: 0, G: 210, B: 0, A: 255}
)

// Scene represent the scene of the game.
type Scene struct {
	world   *world.World
	hero    *hero.Hero
	pits    []*pit.Pit
	trove   []*trove.Trove
	enemies []*enemy.Enemy
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

	id = "hero"
	pos = w.RandomizePos(id, 50, 50)
	h := hero.NewHero(id, pos.X, pos.Y, w)

	lvl := w.GetLevel()

	return &Scene{
		world:   w,
		hero:    h,
		pits:    createPits(w, lvl),
		trove:   createTroves(w, lvl+1),
		enemies: createEnemies(w, lvl+1),
	}, nil
}

// Run runs the Scene.
func (s *Scene) Run(events <-chan sdl.Event, r *sdl.Renderer) <-chan error {
	errc := make(chan error)

	go func() {
		defer close(errc)
		tick := time.Tick(10 * time.Millisecond)

		if err := drawTitle(r, "Trove Hero", orangeClr); err != nil {
			errc <- fmt.Errorf("could not draw title: %w", err)
		}
		time.Sleep(1 * time.Second)

		for {
			select {
			case e := <-events:
				if done := s.handleEvent(e); done {
					drawStats(s.world)
					return
				}
			case <-tick:
				s.update()

				if s.hero.IsDead() {
					if err := drawTitle(r, "Game Over", redClr); err != nil {
						errc <- err
					}
					time.Sleep(1 * time.Second)
					s.restart()
				}

				if len(s.trove) == 0 {
					if err := drawTitle(r, "You won", greenClr); err != nil {
						errc <- err
					}
					time.Sleep(1 * time.Second)
					s.world.IncLevel()
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
	for _, v := range s.pits {
		s.hero.TouchPit(v)
	}

	i := 0 // output index
	for _, t := range s.trove {
		s.hero.TouchTrove(t)
		if !t.IsCollected() {
			// copy and increment index
			s.trove[i] = t
			i++
		}
	}
	s.trove = s.trove[:i]

	for _, e := range s.enemies {
		e.Touch(s.hero)
		e.Watch(s.hero)
	}

	s.hero.Update()

	for _, v := range s.enemies {
		v.Update()
	}

	for _, v := range s.pits {
		v.Update()
	}
}

func (s *Scene) restart() {
	s.hero.Restart()

	lvl := s.world.GetLevel()
	s.pits = createPits(s.world, lvl)
	s.trove = createTroves(s.world, lvl+1)
	s.enemies = createEnemies(s.world, lvl+1)
}

func (s *Scene) paint(r *sdl.Renderer) error {
	r.Clear()

	for _, v := range s.pits {
		if err := v.Paint(r); err != nil {
			return err
		}
	}

	for _, v := range s.trove {
		if err := v.Paint(r); err != nil {
			return err
		}
	}

	if err := s.hero.Paint(r); err != nil {
		return err
	}

	for _, v := range s.enemies {
		if err := v.Paint(r); err != nil {
			return err
		}
	}

	r.Present()
	return nil
}

// Destroy destroys the scene.
func (s *Scene) Destroy() {
	for _, v := range s.pits {
		v.Destroy()
	}
	s.hero.Destroy()
	for _, v := range s.enemies {
		v.Destroy()
	}
}

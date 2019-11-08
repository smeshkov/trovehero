// Copyright 2017 Google Inc. All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to writing, software distributed
// under the License is distributed on a "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.

package scene

import (
	// "fmt"
	"log"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	// img "github.com/veandco/go-sdl2/img"

	"github.com/smeshkov/trovehero/enemy"
	"github.com/smeshkov/trovehero/hero"
	"github.com/smeshkov/trovehero/pit"
	"github.com/smeshkov/trovehero/types"
	"github.com/smeshkov/trovehero/world"
)

// Scene represent the scene of the game.
type Scene struct {
	// bg    *sdl.Texture
	// bird  *bird
	// pipes *pipes
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

	// b, err := newBird(r)
	// if err != nil {
	// 	return nil, err
	// }

	// ps, err := newPipes(r)
	// if err != nil {
	// 	return nil, err
	// }

	w := world.NewWorld(1024, 768, &sdl.Rect{W: 1024, H: 768, X: 0, Y: 0})

	p := pit.NewPit(w.W/2, 250, 50, 150, -60)
	h := hero.NewHero(w.W/2, 700)
	e := enemy.NewEnemy(30, 30, w)

	return &Scene{ /* bg: bg, bird: b, pipes: ps*/ hero: h, pit: p, enemy: e}, nil
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
					DrawTitle(r, "Game Over")
					time.Sleep(time.Second)
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
	case *sdl.MouseMotionEvent, *sdl.WindowEvent, *sdl.TouchFingerEvent, *sdl.CommonEvent:
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
		s.hero.Do(types.Jump)
	case sdl.SCANCODE_LEFT:
		s.hero.Do(types.GoWest)
	case sdl.SCANCODE_RIGHT:
		s.hero.Do(types.GoEast)
	case sdl.SCANCODE_UP:
		s.hero.Do(types.GoNorth)
	case sdl.SCANCODE_DOWN:
		s.hero.Do(types.GoSouth)
	}
	return false
}

func (s *Scene) update() {
	s.hero.Update()
	s.hero.Touch(s.pit)
	// s.pit.Update()
	s.enemy.Watch(s.hero)
	s.enemy.Update()
	// s.bird.update()
	// s.pipes.update()
	// s.pipes.touch(s.bird)
}

func (s *Scene) restart() {
	s.pit.Restart()
	s.hero.Restart()
	// s.bird.restart()
	// s.pipes.restart()
}

func (s *Scene) paint(r *sdl.Renderer) error {
	r.Clear()
	// if err := r.Copy(s.bg, nil, nil); err != nil {
	// 	return fmt.Errorf("could not copy background: %w", err)
	// }
	// if err := s.bird.paint(r); err != nil {
	// 	return err
	// }
	// if err := s.pipes.paint(r); err != nil {
	// 	return err
	// }

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

// Destroy destriys the scene.
func (s *Scene) Destroy() {
	// s.bg.Destroy()
	s.pit.Destroy()
	s.hero.Destroy()
	// s.bird.destroy()
	// s.pipes.destroy()
}

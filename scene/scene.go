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

	"github.com/smeshkov/trovehero/hero"
	"github.com/smeshkov/trovehero/pit"
)

// Scene represent the scene of the game.
type Scene struct {
	// bg    *sdl.Texture
	// bird  *bird
	// pipes *pipes
	h *hero.Hero
	p *pit.Pit
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

	h := hero.NewHero(800/2, 550)
	p := pit.NewPit(800/2, 350, 70, 150)

	return &Scene{ /* bg: bg, bird: b, pipes: ps*/ h: h, p: p}, nil
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

				if s.h.IsDead() {
					DrawTitle(r, "Game Over")
					time.Sleep(time.Second)
					s.restart()
				}
				// if s.bird.isDead() {
				// 	drawTitle(r, "Game Over")
				// 	time.Sleep(time.Second)
				// 	s.restart()
				// }

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
		s.h.Do(hero.Jump)
	case sdl.SCANCODE_LEFT:
		s.h.Do(hero.Left)
	case sdl.SCANCODE_RIGHT:
		s.h.Do(hero.Right)
	case sdl.SCANCODE_UP:
		s.h.Do(hero.Up)
	case sdl.SCANCODE_DOWN:
		s.h.Do(hero.Down)
	}
	return false
}

func (s *Scene) update() {
	s.h.Update()
	s.p.Update()
	s.h.Touch(s.p)
	// s.bird.update()
	// s.pipes.update()
	// s.pipes.touch(s.bird)
}

func (s *Scene) restart() {
	s.h.Restart()
	s.p.Restart()
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

	if err := s.h.Paint(r); err != nil {
		return err
	}
	if err := s.p.Paint(r); err != nil {
		return err
	}

	r.Present()
	return nil
}

// Destroy destriys the scene.
func (s *Scene) Destroy() {
	// s.bg.Destroy()
	s.h.Destroy()
	s.p.Destroy()
	// s.bird.destroy()
	// s.pipes.destroy()
}

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	ttf "github.com/veandco/go-sdl2/ttf"
	// "go.uber.org/zap"

	"github.com/smeshkov/trovehero/scene"
)

var (
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	// verbose    = flag.Bool("verbose", false, "enables verbose mode")
)

func main() {
	flag.Parse()

	// if err := setupLog(); err != nil {
	// 	fmt.Fprintf(os.Stderr, "%v", err)
	// 	os.Exit(2)
	// }

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(3)
	}
}

func run() error {
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error in creating cpuprofile file %s: %v", *cpuprofile, err)
			os.Exit(3)
		}
		err = pprof.StartCPUProfile(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error in profiling: %v", err)
			os.Exit(3)
		}
		defer pprof.StopCPUProfile()
	}

	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return fmt.Errorf("could not initialize SDL: %w", err)
	}
	defer sdl.Quit()

	if err := ttf.Init(); err != nil {
		return fmt.Errorf("could not initialize TTF: %w", err)
	}
	defer ttf.Quit()

	w, r, err := sdl.CreateWindowAndRenderer(1024, 768, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("could not create window: %w", err)
	}
	defer w.Destroy()

	if err := scene.DrawTitle(r, "Trove Hero"); err != nil {
		return fmt.Errorf("could not draw title: %w", err)
	}

	time.Sleep(1 * time.Second)

	s, err := scene.NewScene(r)
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

// func setupLog() error {
// 	var l *zap.Logger
// 	var err error
// 	if *verbose {
// 		l, err = zap.NewDevelopment()
// 	} else {
// 		l, err = zap.NewProduction()
// 	}

// 	if err != nil {
// 		return err
// 	}

// 	zap.ReplaceGlobals(l)
// 	zap.L().Debug("logger is ready")

// 	return nil
// }

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"

	// "go.uber.org/zap"

	"github.com/smeshkov/trovehero"
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

	if err := trovehero.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(3)
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

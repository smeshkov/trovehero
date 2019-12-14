package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/smeshkov/trovehero"
)

var (
	level = flag.Int("lvl", 0, "sets starting level, e.g. -lvl=2")
)

func main() {
	flag.Parse()
	if err := trovehero.Run(int8(*level)); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(3)
	}
}

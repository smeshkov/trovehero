package main

import (
	"fmt"
	"os"

	"github.com/smeshkov/trovehero/engine"
)

func main() {
	if err := engine.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(3)
	}
}
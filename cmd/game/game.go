package main

import (
	"fmt"
	"os"

	"github.com/smeshkov/trovehero"
)

func main() {
	if err := trovehero.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(3)
	}
}

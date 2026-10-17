package main

import (
	"fmt"
	"os"

	"github.com/technothatch/wintower/internal/ui"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("wintower %s\n", version)
		return
	}
	if err := ui.Run(version); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

package main

import (
	"os"

  "github.com/redcuckoo/bsc-checker-events/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}

package main

import (
	"os"

  "githab.com/redcuckoo/bsc-checker-events/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}

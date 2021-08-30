package main

import (
	"os"

  "github.com/Alien-Worlds/missions-api/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}

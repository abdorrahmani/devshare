package main

import (
	"github.com/abdorrahmani/devshare/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

package main

import (
	"os"

	"github.com/redjax/go-mygithub/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

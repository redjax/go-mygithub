package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/redjax/go-mygithub/cmd"
	"github.com/redjax/go-mygithub/internal/constants"
)

func main() {
	// Initialize CLI
	var cli cmd.CLI

	// Parse CLI args & create CLI context
	ctx := kong.Parse(&cli, kong.Name("mygithub"), kong.Description("CLI to interact with your Github account via the Github API & a Personal Acccess Token (PAT)."))
	// Run CLI
	err := ctx.Run()
	// Raise if CLI encounters error
	ctx.FatalIfErrorf(err)

	fmt.Printf("Github user starred repositories URL: %s\n", constants.GH_STARRED_URL)

}

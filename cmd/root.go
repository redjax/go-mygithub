package cmd

import (
	"github.com/redjax/go-mygithub/cmd/starred"
)

// Declare CLI, add subcommands
type CLI struct {
	Starred starred.GithubStarredCmd `cmd:"" help:"Get starred repositories"`
}

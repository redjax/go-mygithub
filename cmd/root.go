package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	accessToken string
)

var rootCmd = &cobra.Command{
	Use:   "mygithub",
	Short: "CLI for Github",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Persistent flag so it works on all subcommands
	rootCmd.PersistentFlags().StringVarP(&accessToken, "access-token", "t", "", "GitHub Personal Access Token (PAT)")
	viper.BindPFlag("access_token", rootCmd.PersistentFlags().Lookup("access-token"))
	viper.BindEnv("access_token", "GITHUB_TOKEN", "GH_TOKEN")
}

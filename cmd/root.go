package cmd

import (
	"github.com/hqdem/go-api-template/cmd/runserver"
	"github.com/spf13/cobra"
)

var (
	cfgFileName string
	rootCmd     = &cobra.Command{
		Use:   "go-api-template {action}",
		Short: "perfom actions for go-api-template",
		Long:  "pass the action you want to perform after name of the service (go-api-template)",
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFileName, "config", "", "path to config file")
	_ = rootCmd.MarkPersistentFlagRequired("config")

	addCommands()
}

func addCommands() {
	rootCmd.AddCommand(runserver.CreateCommand(&cfgFileName))
}

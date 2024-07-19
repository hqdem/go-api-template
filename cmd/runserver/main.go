package runserver

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var runServerCmd = &cobra.Command{
	Use:   "run api server",
	Short: "command that runs api server",
	Long:  "command that runs api server",
	//RunE: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := runServerCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	runServerCmd.Flags().String("config", "", "config file path")
}

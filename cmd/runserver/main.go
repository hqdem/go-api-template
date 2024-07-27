package runserver

import (
	"github.com/hqdem/go-api-template/pkg/commands/runserver"
	"github.com/spf13/cobra"
)

func CreateCommand(cfgPath *string) *cobra.Command {
	return &cobra.Command{
		Use:   "run api server",
		Short: "run api server",
		Long:  "run api server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runserver.RunServer(*cfgPath)
		},
	}
}

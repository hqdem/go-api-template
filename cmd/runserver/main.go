package runserver

import (
	"github.com/hqdem/go-api-template/internal/commands/runserver"
	"github.com/spf13/cobra"
)

func CreateCommand(cfgPath *string) *cobra.Command {
	return &cobra.Command{
		Use:   "run --config /path/to/cfg",
		Short: "run api server",
		Long:  "run api server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runserver.RunServer(*cfgPath)
		},
	}
}

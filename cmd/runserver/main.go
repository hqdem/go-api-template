package runserver

import (
	"github.com/hqdem/go-api-template/pkg/commands/runserver"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
)

func CreateCommand(cfgPath *string) *cobra.Command {

	_ = echo.New()
	return &cobra.Command{
		Use:   "run --config /path/to/cfg",
		Short: "run api server",
		Long:  "run api server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runserver.RunServer(*cfgPath)
		},
	}
}

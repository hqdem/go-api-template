package runserver

import (
	"github.com/hqdem/go-api-template/lib/xlog"
	"github.com/hqdem/go-api-template/pkg/config"
	"github.com/hqdem/go-api-template/pkg/connectors/postgre"
	"github.com/hqdem/go-api-template/pkg/core/facade"
	xhttp "github.com/hqdem/go-api-template/pkg/handlers/http"
)

func RunServer(cfgPath string) error {
	cfg, err := config.NewConfig(cfgPath)
	if err != nil {
		return err
	}
	err = xlog.SetDefaultLogger(cfg.Logger.Level, cfg.Logger.Development)
	if err != nil {
		return err
	}
	storage := postgre.NewConnector()
	facadeObj := facade.NewFacade(cfg, storage)

	app, err := xhttp.NewServerApp(facadeObj)
	if err != nil {
		return err
	}
	return app.Run()
}

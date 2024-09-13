package runserver

import (
	"context"
	"fmt"
	"github.com/hqdem/go-api-template/internal/config"
	"github.com/hqdem/go-api-template/internal/connectors/postgre"
	"github.com/hqdem/go-api-template/internal/core/facade"
	xhttp "github.com/hqdem/go-api-template/internal/handlers/http"
	"github.com/hqdem/go-api-template/pkg/xlog"
	"go.uber.org/zap"
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

	onPanicHook := func(ctx context.Context, panicErr error, panicStack []byte) {
		ctx = xlog.WithFields(ctx, zap.String("panic_stack", string(panicStack)))
		xlog.Error(ctx, panicErr.Error())

		// TODO: add env variable in config
		if cfg.Logger.Development {
			fmt.Println(string(panicStack)) // for debug purposes
		}
	}

	onCtxDoneHook := func(ctx context.Context) {
		ctxErr := ctx.Err()
		xlog.Error(ctx, ctxErr.Error())
	}

	app, err := xhttp.NewServerApp(facadeObj, onPanicHook, onCtxDoneHook)
	if err != nil {
		return err
	}
	return app.Run()
}

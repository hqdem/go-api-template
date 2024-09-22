package runserver

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hqdem/go-api-template/internal/config"
	"github.com/hqdem/go-api-template/internal/connectors/postgre"
	"github.com/hqdem/go-api-template/internal/core/actions"
	"github.com/hqdem/go-api-template/internal/core/facade"
	xhttp "github.com/hqdem/go-api-template/internal/handlers/http"
	"github.com/hqdem/go-api-template/pkg/xlog"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func RunServer(cfgPath string) error {
	runCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.NewConfig(cfgPath)
	if err != nil {
		return err
	}
	err = xlog.SetDefaultLogger(cfg.Logger.Level, cfg.Logger.Development)
	if err != nil {
		return err
	}
	storage := postgre.NewConnector()
	actionsImpl := actions.NewActions()
	facadeObj := facade.NewFacade(cfg, storage, actionsImpl)

	onPanicHook := func(ctx context.Context, panicErr error, panicStack []byte) {
		ctx = xlog.WithFields(ctx, zap.String("panic_stack", string(panicStack)))
		xlog.Error(ctx, panicErr.Error())

		if cfg.Env.IsDevelopment() {
			fmt.Println(string(panicStack)) // for debug purposes
		}
	}

	onCtxDoneHook := func(ctx context.Context) {
		ctxErr := ctx.Err()
		xlog.Error(ctx, ctxErr.Error())
	}

	onHandlerDoneHook := func(ctx context.Context, res any, err error) {
		if err != nil {
			xlog.Error(ctx, fmt.Sprintf("error while handle request: %v", err))
			return
		}
		jsonBytes, err := json.Marshal(res)
		if err != nil {
			xlog.Error(ctx, fmt.Sprintf("can not convert handler result to json: %v", err))
			return
		}
		xlog.Info(ctx, fmt.Sprintf("handler result: %s", string(jsonBytes)))
	}

	app, err := xhttp.NewServerApp(runCtx, facadeObj, onPanicHook, onCtxDoneHook, onHandlerDoneHook)
	if err != nil {
		return err
	}
	err = app.Run(runCtx)
	if err != nil {
		xlog.Error(runCtx, err.Error())
	}
	return nil
}

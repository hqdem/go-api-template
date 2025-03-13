package runserver

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/hqdem/go-api-template/internal/config"
	"github.com/hqdem/go-api-template/internal/connectors/postgres_repository"
	"github.com/hqdem/go-api-template/internal/core/services/ping_service"
	xhttp "github.com/hqdem/go-api-template/internal/handlers/http"
	pinghttp "github.com/hqdem/go-api-template/internal/handlers/http/ping"
	"github.com/hqdem/go-api-template/pkg/xlog"
	"go.uber.org/zap"
)

func RunServer(cfgPath string) error {
	runCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.NewConfig(cfgPath)
	if err != nil {
		return err
	}
	err = xlog.SetDefaultLogger(cfg.Logger.Level, cfg.Logger.Development)
	if err != nil {
		return err
	}

	var nopConnection *sql.DB
	pingRepo := postgres_repository.NewPingRepo(nopConnection)

	pingService := ping_service.New(pingRepo)
	pingHTTP := pinghttp.NewHTTPService(pingService)

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

	app, err := xhttp.NewServerApp(cfg, pingHTTP, onPanicHook, onCtxDoneHook, onHandlerDoneHook)
	if err != nil {
		return err
	}
	app.Run(runCtx)

	return nil
}

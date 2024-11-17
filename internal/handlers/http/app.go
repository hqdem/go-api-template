package xhttp

import (
	"context"
	"errors"
	"fmt"
	"github.com/hqdem/go-api-template/internal/core/facade"
	"github.com/hqdem/go-api-template/internal/handlers/http/middleware"
	"github.com/hqdem/go-api-template/pkg/xlog"
	"github.com/hqdem/go-api-template/pkg/xweb"
	xmiddleware "github.com/hqdem/go-api-template/pkg/xweb/middleware"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ServerApp struct {
	Facade            *facade.Facade
	onPanicHook       xweb.OnPanicFnHookT
	onCtxDoneHook     xweb.OnCtxDoneHookT
	onHandlerDoneHook xweb.OnHandlerDoneHookT
	mux               *http.ServeMux
	server            *http.Server
}

func NewServerApp(
	facade *facade.Facade,
	onPanicHook xweb.OnPanicFnHookT,
	onCtxDoneHook xweb.OnCtxDoneHookT,
	onHandlerDoneHook xweb.OnHandlerDoneHookT,
) (*ServerApp, error) {
	cfg := facade.Config
	appContainer := &ServerApp{
		Facade:            facade,
		onPanicHook:       onPanicHook,
		onCtxDoneHook:     onCtxDoneHook,
		onHandlerDoneHook: onHandlerDoneHook,
	}

	appContainer.mux = http.NewServeMux()
	routes := appContainer.GetRoutes()

	for _, route := range routes {
		appContainer.mux.HandleFunc(route.Pattern, route.Fn)
	}

	handler, err := appContainer.initMiddlewares(appContainer.mux)
	if err != nil {
		return nil, err
	}
	appContainer.server = &http.Server{
		Addr:    cfg.Server.Listen,
		Handler: handler,
	}
	appContainer.initHooks()
	return appContainer, nil
}

func (a *ServerApp) initMiddlewares(handler http.Handler) (http.Handler, error) {
	cfg := a.Facade.Config
	middlewareChain := []xmiddleware.Middleware{
		xmiddleware.TimeoutMiddleware(cfg.Handlers),
		xmiddleware.RequestIDMiddleware(),
		middleware.LogRequestIDMiddleware(),
		middleware.LogRequestMiddleware(),
	}

	// reversing middleware chain to apply mw rules
	// in top-down ordering before actual handler
	// and down-top ordering after actual handler
	for i, j := 0, len(middlewareChain)-1; i < j; i, j = i+1, j-1 {
		middlewareChain[i], middlewareChain[j] = middlewareChain[j], middlewareChain[i]
	}

	for _, mw := range middlewareChain {
		handler = mw(handler)
	}
	return handler, nil
}

func (a *ServerApp) initHooks() {
	xweb.SetPanicFnHook(a.onPanicHook)
	xweb.SetCtxDoneHook(a.onCtxDoneHook)
	xweb.SetHandlerDoneHook(a.onHandlerDoneHook)
}

func (a *ServerApp) Run(ctx context.Context) {
	logMsg := fmt.Sprintf("start listen server on addr: %s", a.server.Addr)
	xlog.Info(ctx, logMsg)

	go func() {
		if err := a.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			xlog.Fatal(ctx, fmt.Sprintf("HTTP server error: %v", err))
		}
		xlog.Info(ctx, "stopped serving new connections")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	xlog.Info(ctx, "gracefully shutting down the server")
	if err := a.server.Shutdown(shutdownCtx); err != nil {
		xlog.Fatal(ctx, fmt.Sprintf("HTTP shutdown error: %v", err))
	}
	xlog.Info(ctx, "graceful shutdown complete")
}

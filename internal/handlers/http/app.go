package xhttp

import (
	"context"
	"fmt"
	"github.com/hqdem/go-api-template/internal/core/facade"
	"github.com/hqdem/go-api-template/internal/handlers/http/middleware"
	"github.com/hqdem/go-api-template/pkg/xlog"
	"github.com/hqdem/go-api-template/pkg/xweb"
	xmiddleware "github.com/hqdem/go-api-template/pkg/xweb/middleware"
	"net/http"
)

type ServerApp struct {
	Facade        *facade.Facade
	onPanicHook   xweb.OnPanicFnHookT
	onCtxDoneHook xweb.OnCtxDoneHookT
	mux           *http.ServeMux
	server        *http.Server
}

func NewServerApp(facade *facade.Facade, onPanicHook xweb.OnPanicFnHookT, onCtxDoneHook xweb.OnCtxDoneHookT) (*ServerApp, error) {
	cfg := facade.Config
	appContainer := &ServerApp{
		Facade:        facade,
		onPanicHook:   onPanicHook,
		onCtxDoneHook: onCtxDoneHook,
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
}

func (a *ServerApp) Run() error {
	// TODO: add run context
	logMsg := fmt.Sprintf("start listen server on addr: %s", a.server.Addr)
	xlog.Info(context.Background(), logMsg)
	return a.server.ListenAndServe()
}

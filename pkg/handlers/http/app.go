package xhttp

import (
	"context"
	"fmt"
	"github.com/hqdem/go-api-template/lib/xlog"
	"github.com/hqdem/go-api-template/lib/xweb/middleware"
	"github.com/hqdem/go-api-template/pkg/core/facade"
	"net/http"
)

type ServerApp struct {
	Facade *facade.Facade
	mux    *http.ServeMux
	server *http.Server
}

func NewServerApp(facade *facade.Facade) (*ServerApp, error) {
	cfg := facade.Config
	appContainer := &ServerApp{
		Facade: facade,
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
	return appContainer, nil
}

func (a *ServerApp) initMiddlewares(handler http.Handler) (http.Handler, error) {
	cfg := a.Facade.Config
	middlewareChain := []middleware.Middleware{
		middleware.TimeoutMiddleware(cfg.Handlers),
		middleware.RequestIDMiddleware(),
	}

	for _, mw := range middlewareChain {
		handler = mw(handler)
	}
	return handler, nil
}

func (a *ServerApp) Run() error {
	// TODO: add run context
	logMsg := fmt.Sprintf("start listen server on addr: %s", a.server.Addr)
	xlog.Info(context.Background(), logMsg)
	return a.server.ListenAndServe()
}

package xhttp

import (
	"fmt"
	"github.com/hqdem/go-api-template/pkg/config"
	"github.com/hqdem/go-api-template/pkg/core/facade"
	"github.com/hqdem/go-api-template/pkg/handlers/http/middlewares"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type ServerApp struct {
	Facade         *facade.Facade
	handlersConfig *config.HandlersConfig
	mux            *http.ServeMux
	server         *http.Server
}

func NewServerApp(facade *facade.Facade) (*ServerApp, error) {
	cfg := facade.Config
	appContainer := &ServerApp{
		Facade:         facade,
		handlersConfig: cfg.Handlers,
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
	middlewareChain := []middlewares.Middleware{
		middlewares.TimeoutMiddleware(time.Second * 3),
	}

	for _, middleware := range middlewareChain {
		handler = middleware(handler)
	}
	return handler, nil
}

func (a *ServerApp) Run() error {
	logMsg := fmt.Sprintf("start listen server on addr: %s", a.server.Addr)
	zap.L().Info(logMsg)
	return a.server.ListenAndServe()
}

package xhttp

import (
	"github.com/hqdem/go-api-template/pkg/core/facade"
	"net/http"
)

type ServerApp struct {
	Facade *facade.Facade
	Listen string
	mux    *http.ServeMux
}

func NewServerApp(facade *facade.Facade) (*ServerApp, error) {
	cfg := facade.Config
	appContainer := &ServerApp{
		Facade: facade,
		Listen: cfg.Server.Listen,
	}

	appContainer.mux = http.NewServeMux()
	routes := appContainer.GetRoutes()

	for _, route := range routes {
		appContainer.mux.HandleFunc(route.Pattern, route.Fn)
	}

	err := appContainer.initMiddlewares()
	if err != nil {
		return nil, err
	}
	return appContainer, nil
}

func (a *ServerApp) initMiddlewares() error {
	// TODO: init mw here
	return nil
}

func (a *ServerApp) Run() error {
	// TODO: logging
	//logMsg := fmt.Sprintf("start listen server on addr: %s", a.Listen)
	//xlog.Info(logMsg)
	return http.ListenAndServe(a.Listen, a.mux)
}

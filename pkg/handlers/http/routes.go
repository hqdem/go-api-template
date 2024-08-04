package xhttp

import (
	"github.com/hqdem/go-api-template/lib/xweb"
	"github.com/hqdem/go-api-template/pkg/handlers/http/ping"
	"net/http"
)

type Route struct {
	Pattern string
	Fn      http.HandlerFunc
}

func (a *ServerApp) GetRoutes() []Route {
	return []Route{
		{
			Pattern: "GET /ping",
			Fn:      xweb.FacadeHandlerAdapter(a.Facade, ping.PingHandler),
		},
	}
}

package xhttp

import (
	"github.com/hqdem/go-api-template/internal/handlers/http/ping"
	"github.com/hqdem/go-api-template/pkg/xweb"
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
		{
			Pattern: "/",
			Fn:      xweb.NotFoundHandler,
		},
	}
}

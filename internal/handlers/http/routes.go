package xhttp

import (
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
			Fn:      xweb.HandlerFunc(a.pingHTTP.GetPingStatus),
		},
		{
			Pattern: "/",
			Fn:      xweb.NotFoundHandler,
		},
	}
}

package ping

import (
	"context"
	"github.com/hqdem/go-api-template/lib/xweb"
	"github.com/hqdem/go-api-template/pkg/core/facade"
	"github.com/hqdem/go-api-template/pkg/handlers/http/ping/schemas"
)

func PingHandler(ctx context.Context, _ *xweb.ResponseHeaders, facade *facade.Facade) (*schemas.PingResponseSchema, error) {
	resp := &schemas.PingResponseSchema{
		Status: "pong",
	}
	return resp, nil
}

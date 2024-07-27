package ping

import (
	"context"
	"github.com/hqdem/go-api-template/lib/xweb"
	"github.com/hqdem/go-api-template/pkg/core/actions"
	"github.com/hqdem/go-api-template/pkg/core/facade"
	"github.com/hqdem/go-api-template/pkg/handlers/http/ping/schemas"
)

func PingHandler(ctx context.Context, _ *xweb.ResponseHeaders, facade *facade.Facade) (*schemas.PingResponseSchema, error) {
	pingStatus, err := actions.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return &schemas.PingResponseSchema{
		Status: pingStatus.Status,
	}, nil
}

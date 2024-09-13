package ping

import (
	"context"
	"github.com/hqdem/go-api-template/internal/core/actions"
	"github.com/hqdem/go-api-template/internal/core/facade"
	"github.com/hqdem/go-api-template/internal/handlers/http/ping/schemas"
	"github.com/hqdem/go-api-template/pkg/xweb"
	"net/http"
)

func PingHandler(ctx context.Context, _ *xweb.ResponseHeaders, _ *http.Request, _ *facade.Facade) (*schemas.PingResponseSchema, error) {
	pingStatus, err := actions.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return &schemas.PingResponseSchema{
		Status: pingStatus.Status,
	}, nil
}
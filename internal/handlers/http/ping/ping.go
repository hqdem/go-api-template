package ping

import (
	"context"
	"github.com/hqdem/go-api-template/internal/handlers/http/ping/schemas"
	"github.com/hqdem/go-api-template/pkg/xweb"
	"net/http"
)

func (s *HTTPService) GetPingStatus(ctx context.Context, _ *xweb.ResponseHeaders, _ *http.Request) (*schemas.PingResponseSchema, error) {
	pingStatus, err := s.pingService.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return &schemas.PingResponseSchema{
		Status: pingStatus.Status,
	}, nil
}

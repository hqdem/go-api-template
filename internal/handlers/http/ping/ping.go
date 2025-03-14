package ping

import (
	"context"
	"github.com/hqdem/go-api-template/internal/handlers/http/ping/schemas"
	"github.com/hqdem/go-api-template/pkg/xweb"
	"net/http"
)

// GetPingStatus swagger info
//	@Summary		Ping service
//	@Description	Get ping status
//	@Tags			ping
//	@Produce		json
//	@Success		200	{object}	xweb.ApiOKResponse[schemas.PingResponseSchema]
//	@Failure		500	{object}	xweb.APIErrorResponse
//	@Router			/ping [get]
func (s *HTTPService) GetPingStatus(ctx context.Context, _ *xweb.ResponseHeaders, _ *http.Request) (*schemas.PingResponseSchema, error) {
	pingStatus, err := s.pingService.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return &schemas.PingResponseSchema{
		Status: pingStatus.Status,
	}, nil
}

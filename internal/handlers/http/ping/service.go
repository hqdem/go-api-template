package ping

import (
	"context"
	"github.com/hqdem/go-api-template/internal/core/entities"
)

type pingService interface {
	Ping(ctx context.Context) (entities.PingStatus, error)
}

type HTTPService struct {
	pingService pingService
}

func NewHTTPService(pingService pingService) *HTTPService {
	return &HTTPService{pingService: pingService}
}

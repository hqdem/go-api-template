package ping_service

import (
	"context"
	"github.com/hqdem/go-api-template/internal/core/entities"
	"github.com/hqdem/go-api-template/pkg/xlog"
	"go.uber.org/zap"
)

func (s *Service) Ping(ctx context.Context) (entities.PingStatus, error) {
	op := "actions.Ping"
	xlog.Info(ctx, "start operation", zap.String("operation", op))
	defer xlog.Info(ctx, "end operation", zap.String("operation", op))

	pingStatus, err := s.PingRepo.PingDB(ctx)
	return entities.PingStatus{Status: pingStatus}, err
}

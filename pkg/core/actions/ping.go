package actions

import (
	"context"
	"github.com/hqdem/go-api-template/lib/xlog"
	"github.com/hqdem/go-api-template/pkg/core/entities"
	"go.uber.org/zap"
)

func Ping(ctx context.Context) (*entities.PingStatus, error) {
	op := "actions.Ping"
	xlog.Info(ctx, "start operation", zap.String("operation", op))
	defer xlog.Info(ctx, "end operation", zap.String("operation", op))

	return entities.NewPingStatus("pong"), nil
}

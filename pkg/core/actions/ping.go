package actions

import (
	"context"
	"github.com/hqdem/go-api-template/pkg/core/entities"
	"go.uber.org/zap"
	"time"
)

func Ping(ctx context.Context) (*entities.PingStatus, error) {
	time.Sleep(time.Second * 5)
	op := "actions.Ping"
	zap.L().Info("start operation", zap.String("operation", op))
	defer zap.L().Info("end operation", zap.String("operation", op))

	return entities.NewPingStatus("pong"), nil
}

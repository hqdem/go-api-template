package core

import (
	"context"
	"github.com/hqdem/go-api-template/internal/core/entities"
)

type DBStorage interface {
	PingDB() string
}

type Actions interface {
	Ping(ctx context.Context) (*entities.PingStatus, error)
}

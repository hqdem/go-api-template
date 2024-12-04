package facade

import (
	"context"
	"github.com/hqdem/go-api-template/internal/config"
	"github.com/hqdem/go-api-template/internal/core/entities"
)

type pingService interface {
	Ping(ctx context.Context) (entities.PingStatus, error)
}

type Facade struct {
	Config      *config.Config
	PingService pingService
}

func NewFacade(cfg *config.Config, pingService pingService) *Facade {
	return &Facade{
		Config:      cfg,
		PingService: pingService,
	}
}

package facade

import (
	"github.com/hqdem/go-api-template/internal/config"
	"github.com/hqdem/go-api-template/internal/core"
)

type Facade struct {
	Config  *config.Config
	Storage core.DBStorage
	Actions core.Actions
}

func NewFacade(cfg *config.Config, storage core.DBStorage, actions core.Actions) *Facade {
	return &Facade{
		Config:  cfg,
		Storage: storage,
		Actions: actions,
	}
}

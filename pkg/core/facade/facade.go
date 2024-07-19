package facade

import (
	"github.com/hqdem/go-api-template/pkg/config"
	"github.com/hqdem/go-api-template/pkg/core"
)

type Facade struct {
	Config  *config.Config
	Storage core.DBStorage
}

package config

import (
	"github.com/hqdem/go-api-template/lib/xweb"
	"gopkg.in/yaml.v3"
	"os"
)

type ServerConf struct {
	Listen string `yaml:"listen"`
}

type LoggerConf struct {
	Level       string `yaml:"level"`
	Development bool   `yaml:"development"`
}

type Config struct {
	Server   *ServerConf          `yaml:"server"`
	Logger   *LoggerConf          `yaml:"logger"`
	Handlers *xweb.HandlersConfig `yaml:"handlers"`
}

func NewConfig(cfgPath string) (*Config, error) {
	yamlFile, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}
	cfg := Config{}

	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

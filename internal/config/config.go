package config

import (
	"fmt"
	"github.com/hqdem/go-api-template/pkg/xweb"
	"gopkg.in/yaml.v3"
	"os"
	"slices"
)

type appEnv string

func (e appEnv) IsValid() bool {
	if !slices.Contains(validEnvs, e) {
		return false
	}
	return true
}

func (e appEnv) IsDevelopment() bool {
	return e == devAppEnv
}
func (e appEnv) IsTest() bool {
	return e == testAppEnv
}
func (e appEnv) IsProd() bool {
	return e == prodAppEnv
}

var (
	devAppEnv  appEnv = "dev"
	testAppEnv appEnv = "test"
	prodAppEnv appEnv = "prod"
	validEnvs         = []appEnv{devAppEnv, testAppEnv, prodAppEnv}
)

type ServerConf struct {
	Listen string `yaml:"listen"`
}

type LoggerConf struct {
	Level       string `yaml:"level"`
	Development bool   `yaml:"development"`
}

type Config struct {
	Env      appEnv               `yaml:"env"`
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

	if !cfg.Env.IsValid() {
		return nil, fmt.Errorf("invalid environment configuration: %s, choose from %s", cfg.Env, validEnvs)
	}
	return &cfg, nil
}

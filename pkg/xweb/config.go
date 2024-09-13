package xweb

import (
	"strings"
	"time"
)

type HandlersConfig struct {
	DefaultTimeoutSecs int64                  `yaml:"default_timeout_secs"`
	HandlersTimeouts   []HandlerTimeoutConfig `yaml:"handlers_timeouts"`
}

type HandlerTimeoutConfig struct {
	Method      string `yaml:"method"`
	Endpoint    string `yaml:"endpoint"`
	TimeoutSecs int64  `yaml:"timeout_secs"`
}

func (c *HandlersConfig) GetHandlerTimeout(requestURI string, method string) time.Duration {
	for _, h := range c.HandlersTimeouts {
		requestPath := strings.Split(requestURI, "?")[0]
		if strings.EqualFold(requestPath, h.Endpoint) && strings.EqualFold(method, h.Method) {
			return time.Second * time.Duration(h.TimeoutSecs)
		}
	}
	return time.Second * time.Duration(c.DefaultTimeoutSecs)
}

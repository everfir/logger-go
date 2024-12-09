package tracer_config

import (
	"os"

	"github.com/everfir/logger-go/structs/log_level"
)

var DefaultTracerConfig = TracerConfig{
	Enable:            true,
	Compression:       No,
	ServiceName:       os.Getenv("SERVICE_NAME"),
	CollectorEndpoint: os.Getenv("OTEL_COLLECTOR_DNS"),
	ContextHandlers:   make(map[string]ContextHandler),
}

type TracerConfig struct {
	Enable bool // 开启Tracing功能

	ServiceName       string
	Level             log_level.Level
	Compression       Compression
	CollectorEndpoint string // CollectorEndpoint

	ContextHandlers map[string]ContextHandler
}

func (config *TracerConfig) FixDefault() {
	if config == nil {
		return
	}

	if config.CollectorEndpoint == "" {
		config.CollectorEndpoint = "otelcollector-service.everfir.svc.cluster.local:4317"
	}
}

func (config *TracerConfig) EnableTracing() bool {
	return config != nil && config.Enable
}

func (config *TracerConfig) Validate() bool {
	if config == nil {
		return false
	}

	if config.CollectorEndpoint == "" {
		return false
	}

	if config.Compression > Gzip {
		return false
	}

	return true
}

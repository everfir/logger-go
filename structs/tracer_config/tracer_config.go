package tracer_config

import "os"

var DefaultTracerConfig = TracerConfig{
	Enable:            true,
	Compression:       No,
	CollectorEndpoint: os.Getenv("OTEL_COLLECTOR_DNS"),
}

type TracerConfig struct {
	Enable bool // 开启Tracing功能

	Compression       Compression
	CollectorEndpoint string // CollectorEndpoint
}

func (config *TracerConfig) FixDefault() {
	if config == nil {
		return
	}

	if config.CollectorEndpoint == "" {
		config.CollectorEndpoint = "otelcollector-service.everfir.svc.cluster.local"
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

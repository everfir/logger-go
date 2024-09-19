package tracer_config

var DefaultTracerConfig = TracerConfig{
	Enable:            true,
	Compression:       No,
	CollectorEndpoint: "localhost:4318",
}

type TracerConfig struct {
	Enable bool // 开启Tracing功能

	Compression       Compression
	CollectorEndpoint string // CollectorEndpoint
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

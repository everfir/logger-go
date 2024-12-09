package tracer_config

import "context"

type ContextHandler func(context.Context) string

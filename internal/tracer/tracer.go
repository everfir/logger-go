package tracer

import (
	"context"

	"github.com/everfir/logger-go/structs/field"
	"github.com/everfir/logger-go/structs/log_level"
	"go.opentelemetry.io/otel/propagation"
)

type Tracer interface {
	Init() error
	Close() error
	FixFields(ctx context.Context, fields ...field.Field) []field.Field
	Trace(ctx context.Context, level log_level.Level, msg string, fileds ...field.Field)
	Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context
	Inject(ctx context.Context, carrier propagation.TextMapCarrier)
}

package tracer

import (
	"context"

	"github.com/everfir/logger-go/structs/field"
	"github.com/everfir/logger-go/structs/log_level"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type Tracer interface {
	Init() error
	Close() error
	FixFields(ctx context.Context, fields ...field.Field) []field.Field
	Trace(ctx context.Context, level log_level.Level, msg string, fileds ...field.Field)
	Start(ctx context.Context, name string) (context.Context, trace.Span)
	Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context
	Inject(ctx context.Context, carrier propagation.TextMapCarrier)
}

package tracer

import (
	"context"

	"github.com/everfir/logger-go/structs/field"
	"github.com/everfir/logger-go/structs/log_level"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type NoTracer struct {
}

func (tcer *NoTracer) Init() error  { return nil }
func (tcer *NoTracer) Close() error { return nil }
func (tcer *NoTracer) FixFields(ctx context.Context, fields ...field.Field) (ret []field.Field) {
	return fields
}
func (tcer *NoTracer) Trace(context.Context, log_level.Level, string, ...field.Field) {}
func (tcer *NoTracer) Start(ctx context.Context, name string) (context.Context, trace.Span) {
	return ctx, nil
}
func (tcer *NoTracer) Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	return ctx
}
func (tcer *NoTracer) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
}

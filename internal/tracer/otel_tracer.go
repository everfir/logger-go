package tracer

import (
	"context"
	"fmt"
	"time"

	"github.com/everfir/logger-go/structs/field"
	"github.com/everfir/logger-go/structs/log_level"
	"github.com/everfir/logger-go/structs/tracer_config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	trace_sdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

func NewOtelTracer(config *tracer_config.TracerConfig) *OtelTracer {
	return &OtelTracer{
		doneChan: make(chan struct{}),
		config:   config,
	}
}

type OtelTracer struct {
	doneChan   chan struct{}
	config     *tracer_config.TracerConfig
	provider   *trace_sdk.TracerProvider
	propagator propagation.TextMapPropagator
}

func (tcer *OtelTracer) Init() (err error) {
	if !tcer.config.EnableTracing() || !tcer.config.Validate() {
		return nil
	}

	// 压缩
	compression := otlptracehttp.NoCompression
	if tcer.config.Compression == tracer_config.Gzip {
		compression = otlptracehttp.GzipCompression
	}

	var exporter *otlptrace.Exporter
	exporter, err = otlptrace.New(
		context.TODO(),
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint(tcer.config.CollectorEndpoint),
			otlptracehttp.WithInsecure(),
			otlptracehttp.WithCompression(compression),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create otelExporter: %w", err)
	}

	tp := trace_sdk.NewTracerProvider(
		trace_sdk.WithSampler(trace_sdk.AlwaysSample()), // 全部采样
		trace_sdk.WithBatcher(exporter),
		trace_sdk.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("traccer"),
			),
		),
	)

	// 设置全局的provider，通过GetTracerProvider获取tracer，来开启一个流程
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagator)
	tcer.provider = tp
	tcer.propagator = otel.GetTextMapPropagator()
	return
}

func (tcer *OtelTracer) Close() (err error) {
	if tcer.provider == nil {
		return nil
	}
	defer close(tcer.doneChan)

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*20)
	defer cancel()

	if err = tcer.provider.ForceFlush(ctx); err != nil {
		err = fmt.Errorf("failed to flush tracer")
		return err
	}
	if err = tcer.provider.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

func (tcer *OtelTracer) FixFields(ctx context.Context, fields ...field.Field) []field.Field {

	var ok bool
	var span trace.Span
	if span, ok = ctx.Value("tracing").(trace.Span); !ok {
		return fields
	}
	if span == nil || span.SpanContext().IsValid() == false {
		return fields
	}

	traceID := span.SpanContext().TraceID().String()
	spanID := span.SpanContext().SpanID().String()
	fields = append(fields,
		field.String("trace_id", traceID),
		field.String("span_id", spanID),
	)
	return fields
}

func (tcer *OtelTracer) Trace(
	ctx context.Context,
	level log_level.Level,
	msg string,
	fields ...field.Field,
) {

	var ok bool
	var span trace.Span
	if span, ok = ctx.Value("tracing").(trace.Span); !ok {
		return
	}
	if span == nil || span.SpanContext().IsValid() == false {
		return
	}

	if level >= log_level.ErrorLevel {
		span.SetStatus(codes.Error, "error happended")
	} else {
		span.SetStatus(codes.Ok, "success")
	}

	var attrs []attribute.KeyValue
	for _, f := range fields {
		attrs = append(attrs, toOtelField(f))
	}

	span.AddEvent(msg, trace.WithAttributes(attrs...))
}

func toOtelField(f field.Field) attribute.KeyValue {
	switch f.Type() {
	case field.StringType:
		return attribute.String(f.Key(), f.Value().(string))
	case field.BoolType:
		return attribute.Bool(f.Key(), f.Value().(bool))
	case field.IntType:
		return attribute.Int(f.Key(), f.Value().(int))
	case field.Int8Type:
		return attribute.Int(f.Key(), int(f.Value().(int8)))
	case field.Int16Type:
		return attribute.Int(f.Key(), int(f.Value().(int16)))
	case field.Int32Type:
		return attribute.Int(f.Key(), int(f.Value().(int32)))
	case field.Int64Type:
		return attribute.Int64(f.Key(), f.Value().(int64))
	case field.UintType:
		return attribute.Int(f.Key(), int(f.Value().(uint)))
	case field.Uint8Type:
		return attribute.Int(f.Key(), int(f.Value().(uint8)))
	case field.Uint16Type:
		return attribute.Int(f.Key(), int(f.Value().(uint16)))
	case field.Uint32Type:
		return attribute.Int(f.Key(), int(f.Value().(uint32)))
	case field.Uint64Type:
		return attribute.Int64(f.Key(), int64(f.Value().(uint64)))
	case field.Float32Type:
		return attribute.Float64(f.Key(), float64(f.Value().(float32)))
	case field.Float64Type:
		return attribute.Float64(f.Key(), f.Value().(float64))
	case field.TimeType:
		return attribute.String(f.Key(), f.Value().(time.Time).Format(time.RFC3339))
	case field.DurationType:
		return attribute.Int64(f.Key(), int64(f.Value().(time.Duration)))
	default:
		return attribute.String(f.Key(), fmt.Sprintf("%v", f.Value()))
	}
}

func (tcer *OtelTracer) Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	return tcer.propagator.Extract(ctx, carrier)
}

func (tcer *OtelTracer) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	tcer.propagator.Inject(ctx, carrier)
}

func (tcer *OtelTracer) Start(ctx context.Context, name string) (context.Context, trace.Span) {
	return tcer.provider.Tracer(name).Start(ctx, name)
}

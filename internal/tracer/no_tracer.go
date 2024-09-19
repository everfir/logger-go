package tracer

import (
	"context"

	"github.com/everfir/logger-go/structs/field"
	"github.com/everfir/logger-go/structs/log_level"
)

type NoTracer struct {
}

func (tcer *NoTracer) Init() error                                                    { return nil }
func (tcer *NoTracer) Close() error                                                   { return nil }
func (tcer *NoTracer) FixFields(context.Context, ...field.Field) (ret []field.Field)  { return }
func (tcer *NoTracer) Trace(context.Context, log_level.Level, string, ...field.Field) {}

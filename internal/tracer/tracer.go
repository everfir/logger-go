package tracer

import (
	"context"

	"github.com/everfir/logger-go/structs/field"
	"github.com/everfir/logger-go/structs/log_level"
)

type Tracer interface {
	Init() error
	Close() error
	FixFields(ctx context.Context, fields ...field.Field) []field.Field
	Trace(ctx context.Context, level log_level.Level, msg string, fileds ...field.Field)
}

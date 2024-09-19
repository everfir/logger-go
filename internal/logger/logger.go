package logger

import "github.com/everfir/logger-go/structs/field"

// Logger 定义日志接口
type Logger interface {
	Debug(msg string, fields ...field.Field)
	Info(msg string, fields ...field.Field)
	Warn(msg string, fields ...field.Field)
	Error(msg string, fields ...field.Field)
	Fatal(msg string, fields ...field.Field)
}

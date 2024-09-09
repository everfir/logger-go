package logger

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/everfir/logger-go/structs/field"
	"github.com/everfir/logger-go/structs/log_config"
)

// Logger 定义日志接口
type Logger interface {
	Debug(msg string, fields ...field.Field)
	Info(msg string, fields ...field.Field)
	Warn(msg string, fields ...field.Field)
	Error(msg string, fields ...field.Field)
	Fatal(msg string, fields ...field.Field)
}

var (
	globalLogger Logger = &consoleLogger{}
)

// consoleLogger 是一个简单的控制台日志器
type consoleLogger struct{}

func (l *consoleLogger) log(level, msg string, fields ...field.Field) {
	fmt.Printf("[%s] %s", level, msg)
	for _, f := range fields {
		fmt.Printf(" %s=%v", f.Key(), f.Value())
	}
	fmt.Println()
}

func (l *consoleLogger) Debug(msg string, fields ...field.Field) { l.log("DEBUG", msg, fields...) }
func (l *consoleLogger) Info(msg string, fields ...field.Field)  { l.log("INFO", msg, fields...) }
func (l *consoleLogger) Warn(msg string, fields ...field.Field)  { l.log("WARN", msg, fields...) }
func (l *consoleLogger) Error(msg string, fields ...field.Field) { l.log("ERROR", msg, fields...) }
func (l *consoleLogger) Fatal(msg string, fields ...field.Field) {
	l.log("FATAL", msg, fields...)
	os.Exit(1)
}

// Init 初始化全局日志器
func Init(options ...Option) error {
	// 使用默认配置
	config := log_config.DefaultConfig

	// 应用所有选项
	for _, option := range options {
		option(&config)
	}

	return initWithConfig(&config)
}

// initWithConfig 使用给定的配置初始化日志器
func initWithConfig(config *log_config.LogConfig) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("[Log] Init failed: %w", err)
		}
	}()

	var dir string
	if dir, err = os.Getwd(); err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}
	if err = os.MkdirAll(dir+"/log", os.ModePerm); err != nil {
		return fmt.Errorf("failed to create log directory %s: %w", dir, err)
	}

	var logger Logger
	if logger, err = newZapLogger(config); err != nil {
		return fmt.Errorf("newZapLogger with logConfig failed: %w", err)
	}

	globalLogger = logger
	return nil
}

// 提供全局日志函数
func Debug(ctx context.Context, msg string, fields ...field.Field) {
	contextFields := fixFields(ctx)
	globalLogger.Debug(msg, append(contextFields, fields...)...)
}

func Info(ctx context.Context, msg string, fields ...field.Field) {
	contextFields := fixFields(ctx)
	globalLogger.Info(msg, append(contextFields, fields...)...)
}

func Warn(ctx context.Context, msg string, fields ...field.Field) {
	contextFields := fixFields(ctx)
	globalLogger.Warn(msg, append(contextFields, fields...)...)
}

func Error(ctx context.Context, msg string, fields ...field.Field) {
	contextFields := fixFields(ctx)
	globalLogger.Error(msg, append(contextFields, fields...)...)
}

func Fatal(ctx context.Context, msg string, fields ...field.Field) {
	contextFields := fixFields(ctx)
	globalLogger.Fatal(msg, append(contextFields, fields...)...)
}

func fixFields(ctx context.Context) (fields []field.Field) {
	fields = append(fields, field.String("Time", time.Now().String()))

	return append(fields, extractContextFields(ctx)...)
}

// extractContextFields 从 context 中提取信息并创建字段
func extractContextFields(ctx context.Context) (fields []field.Field) {

	var keyMapping = map[string]string{
		"efr_trace_id":       "TraceID",
		"efr_span_id":        "SpanID",
		"efr_parent_span_id": "ParentSpanID",
		"efr_host_name":      "Host",
		"efr_pod_name":       "PodName",
		"efr_ip":             "IP",
		"efr_client_ip":      "CIP",
	}

	var val string
	var ok bool
	for k1, k2 := range keyMapping {
		if val, ok = ctx.Value(k1).(string); ok {
			fields = append(fields, field.String(k2, val))
		}
	}

	return fields
}

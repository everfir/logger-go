package logger

import (
	"context"
	"fmt"
	"os"

	"github.com/everfir/logger-go/internal/logger"
	"github.com/everfir/logger-go/internal/tracer"
	"github.com/everfir/logger-go/structs/field"
	"github.com/everfir/logger-go/structs/log_config"
	"github.com/everfir/logger-go/structs/log_level"
)

var (
	_funcTable []func(context.Context, string, ...field.Field)
)

func init() {
	if err := Init(); err != nil {
		Error(context.TODO(), fmt.Sprintf("[Logger] Init failed:%s. use console logger", err))
		return
	}
}

type myLogger struct {
	logger.Logger
	tracer.Tracer
}

var (
	globalLogger = myLogger{
		Logger: &logger.ConsoleLogger{},
		Tracer: nil,
	}
)

// Init 初始化全局日志器
func Init(options ...Option) error {
	// 使用默认配置
	config := log_config.DefaultConfig // 应用所有选项
	for _, option := range options {
		option(&config)
	}

	return initWithConfig(&config)
}

func Close() {
	if globalLogger.Tracer != nil {
		err := globalLogger.Tracer.Close()
		if err != nil {
			Error(context.TODO(), fmt.Sprintf("Tracer close failed:%s", err))
			return
		}
	}
}

// initWithConfig 使用给定的配置初始化日志器
func initWithConfig(config *log_config.LogConfig) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("[Logger] Init failed: %w", err)
		}
	}()

	var dir string
	if dir, err = os.Getwd(); err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}
	if err = os.MkdirAll(dir+"/log", os.ModePerm); err != nil {
		return fmt.Errorf("failed to create log directory %s: %w", dir, err)
	}

	var loger logger.Logger
	if loger, err = logger.NewZapLogger(config); err != nil {
		return fmt.Errorf("failed to create zap logger with config: %w", err)
	}

	var tcer tracer.Tracer = &tracer.NoTracer{}
	if config.TracerConfig.EnableTracing() && config.TracerConfig.Validate() {
		tcer = tracer.NewOtelTracer(config.TracerConfig)
		if err = tcer.Init(); err != nil {
			return err
		}
	}

	globalLogger.Logger = loger
	globalLogger.Tracer = tcer

	return nil
}

// 提供全局日志函数
func Debug(ctx context.Context, msg string, fields ...field.Field) {
	fields = append(fields, fixFields(ctx)...)
	globalLogger.Tracer.Trace(ctx, log_level.DebugLevel, msg, fields...)
	globalLogger.Logger.Debug(msg, fields...)
}

func Info(ctx context.Context, msg string, fields ...field.Field) {
	fields = append(fields, fixFields(ctx)...)
	globalLogger.Tracer.Trace(ctx, log_level.InfoLevel, msg, fields...)
	globalLogger.Logger.Debug(msg, fields...)
}

func Warn(ctx context.Context, msg string, fields ...field.Field) {
	fields = append(fields, fixFields(ctx)...)
	globalLogger.Tracer.Trace(ctx, log_level.InfoLevel, msg, fields...)
	globalLogger.Logger.Debug(msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...field.Field) {
	fields = append(fields, fixFields(ctx)...)
	globalLogger.Tracer.Trace(ctx, log_level.ErrorLevel, msg, fields...)
	globalLogger.Logger.Debug(msg, fields...)
}

func Fatal(ctx context.Context, msg string, fields ...field.Field) {
	fields = append(fields, fixFields(ctx)...)
	globalLogger.Tracer.Trace(ctx, log_level.FatalLevel, msg, fields...)
	globalLogger.Logger.Debug(msg, fields...)
}

// TODO: 待根据环境方案更新
func fixFields(ctx context.Context) (fields []field.Field) {
	// 添加容器IP
	fields = append(fields, field.String("container.ip", "container_ip_from_env"))
	// 添加容器名
	fields = append(fields, field.String("container.name", "container_name_from_env"))
	// 添加ENV
	fields = append(fields, field.String("env", "env_from_env"))
	// 添加logger_version
	fields = append(fields, field.String("logger.version", "logger_version_from_env"))
	// 添加service_version
	fields = append(fields, field.String("service.version", "service_version_from_env"))
	// 添加service_name
	fields = append(fields, field.String("service.name", "service_name_from_env"))

	if globalLogger.Tracer != nil {
		fields = globalLogger.Tracer.FixFields(ctx, fields...)
	}
	return fields
}

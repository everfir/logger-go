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
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
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

	config *log_config.LogConfig
}

var (
	globalLogger = myLogger{
		Logger: &logger.ConsoleLogger{},
		Tracer: &tracer.NoTracer{},
		config: &log_config.DefaultConfig,
	}
)

// Init 初始化全局日志器
func Init(options ...Option) error {
	// 使用默认配置
	config := log_config.DefaultConfig // 应用所有选项
	for _, option := range options {
		option(&config)
	}
	config.FixDefault()

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

	// TODO 流程优化
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
		tcer = tracer.NewOtelTracer(config.TracerConfig, config.Level)
		if err = tcer.Init(); err != nil {
			return err
		}
	}

	globalLogger.Logger = loger
	globalLogger.Tracer = tcer
	globalLogger.config = config
	return nil
}

// 提供全局日志函数
func Debug(ctx context.Context, msg string, fields ...field.Field) {
	// env fields
	fields = append(fields, fixFields(ctx)...)

	// tracing fields
	if globalLogger.Tracer != nil {
		globalLogger.Tracer.Trace(ctx, log_level.DebugLevel, msg, fields...)
		fields = globalLogger.Tracer.FixFields(ctx, fields...)
	}
	globalLogger.Logger.Debug(msg, fields...)
}

func Info(ctx context.Context, msg string, fields ...field.Field) {
	fields = append(fields, fixFields(ctx)...)

	// tracing fields
	if globalLogger.Tracer != nil {
		globalLogger.Tracer.Trace(ctx, log_level.InfoLevel, msg, fields...)
		fields = globalLogger.Tracer.FixFields(ctx, fields...)
	}
	globalLogger.Logger.Info(msg, fields...)
}

func Warn(ctx context.Context, msg string, fields ...field.Field) {
	fields = append(fields, fixFields(ctx)...)
	if globalLogger.Tracer != nil {
		globalLogger.Tracer.Trace(ctx, log_level.InfoLevel, msg, fields...)
		fields = globalLogger.Tracer.FixFields(ctx, fields...)
	}
	globalLogger.Logger.Warn(msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...field.Field) {
	fields = append(fields, fixFields(ctx)...)
	// tracing fields
	if globalLogger.Tracer != nil {
		globalLogger.Tracer.Trace(ctx, log_level.ErrorLevel, msg, fields...)
		fields = globalLogger.Tracer.FixFields(ctx, fields...)
	}
	globalLogger.Logger.Error(msg, fields...)
}

func Fatal(ctx context.Context, msg string, fields ...field.Field) {
	fields = append(fields, fixFields(ctx)...)

	// tracing fields
	if globalLogger.Tracer != nil {
		globalLogger.Tracer.Trace(ctx, log_level.FatalLevel, msg, fields...)
		fields = globalLogger.Tracer.FixFields(ctx, fields...)
	}
	globalLogger.Logger.Fatal(msg, fields...)
}

// TODO: 待根据环境方案更新
func fixFields(ctx context.Context) (fields []field.Field) {
	// 添加容器IP
	fields = append(fields, field.String("container.ip", globalLogger.config.PodIP))
	// 添加服务名
	fields = append(fields, field.String("ServiceName", globalLogger.config.ServiceName))

	return fields
}

// ----------------------这部分功能，严格来说不算是logger的功能---------------------------------

// Extract 从上下文中提取trace信息
func Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	if globalLogger.Tracer == nil {
		return ctx
	}
	return globalLogger.Tracer.Extract(ctx, carrier)
}

// Inject 将trace信息注入到上下文中
func Inject(ctx context.Context, carrier propagation.TextMapCarrier, extra map[string]string) {
	if globalLogger.Tracer == nil {
		return
	}

	bag := baggage.FromContext(ctx)
	var members []baggage.Member = bag.Members()

	var err error
	for k, v := range extra {
		member, err := baggage.NewMember(k, v)
		if err != nil {
			Warn(ctx, "create baggage member failed", field.Any("error", err))
			continue
		}
		members = append(members, member)
	}

	bag, err = baggage.New(members...)
	if err != nil {
		Warn(ctx, "create baggage failed", field.Any("error", err))
	}
	ctx = baggage.ContextWithBaggage(ctx, bag)

	globalLogger.Tracer.Inject(ctx, carrier)
}

// Start 开始一个span
func Start(ctx context.Context, name string) (context.Context, trace.Span) {
	if globalLogger.Tracer == nil {
		return ctx, nil
	}

	ctx, span := globalLogger.Tracer.Start(ctx, name)
	return ctx, span
}

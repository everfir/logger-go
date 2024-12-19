package logger

import (
	"io"
	"os"
	"path"
	"time"

	. "github.com/everfir/logger-go/structs/field"
	"github.com/everfir/logger-go/structs/log_config"
	"github.com/everfir/logger-go/structs/log_level"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// zapLogger 实现 Logger 接口
type zapLogger struct {
	logger *zap.Logger
}

// Debug 输出调试级别的日志
func (l *zapLogger) Debug(msg string, fields ...Field) {
	l.logger.Debug(msg, toZapFields(fields)...)
}

// Info 输出信息级别的日志
func (l *zapLogger) Info(msg string, fields ...Field) {
	l.logger.Info(msg, toZapFields(fields)...)
}

// Warn 输出警告级别的日志
func (l *zapLogger) Warn(msg string, fields ...Field) {
	l.logger.Warn(msg, toZapFields(fields)...)
}

// Error 输出错误级别的日志
func (l *zapLogger) Error(msg string, fields ...Field) {
	l.logger.Error(msg, toZapFields(fields)...)
}

// Fatal 输出致命错误级别的日志
func (l *zapLogger) Fatal(msg string, fields ...Field) {
	l.logger.Fatal(msg, toZapFields(fields)...)
}

// toZapFields 将通用 Field 转换为 zap.Field
func toZapFields(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = toZapField(field)
	}
	return zapFields
}

// toZapField 根据 Field 类型创建相应的 zap.Field
func toZapField(f Field) zap.Field {
	switch f.Type() {
	case StringType:
		return zap.String(f.Key(), f.Value().(string))
	case BoolType:
		return zap.Bool(f.Key(), f.Value().(bool))
	case IntType:
		return zap.Int(f.Key(), f.Value().(int))
	case Int8Type:
		return zap.Int8(f.Key(), f.Value().(int8))
	case Int16Type:
		return zap.Int16(f.Key(), f.Value().(int16))
	case Int32Type:
		return zap.Int32(f.Key(), f.Value().(int32))
	case Int64Type:
		return zap.Int64(f.Key(), f.Value().(int64))
	case UintType:
		return zap.Uint(f.Key(), f.Value().(uint))
	case Uint8Type:
		return zap.Uint8(f.Key(), f.Value().(uint8))
	case Uint16Type:
		return zap.Uint16(f.Key(), f.Value().(uint16))
	case Uint32Type:
		return zap.Uint32(f.Key(), f.Value().(uint32))
	case Uint64Type:
		return zap.Uint64(f.Key(), f.Value().(uint64))
	case Float32Type:
		return zap.Float32(f.Key(), f.Value().(float32))
	case Float64Type:
		return zap.Float64(f.Key(), f.Value().(float64))
	case TimeType:
		return zap.Time(f.Key(), f.Value().(time.Time))
	case DurationType:
		return zap.Duration(f.Key(), f.Value().(time.Duration))
	default:
		return zap.Any(f.Key(), f.Value())
	}
}

// NewZapLogger 创建一个新的 zapLogger 实例
func NewZapLogger(config *log_config.LogConfig) (Logger, error) {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    otelSeverityEncoder,
		EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var cores []zapcore.Core

	for _, filename := range config.OutputFiles {
		var w zapcore.WriteSyncer
		if filename == "stdout" || filename == "stderr" {
			w = zapcore.AddSync(standardWriter(filename))
		} else {
			rotateLogger, err := getRotateLogger(filename, config)
			if err != nil {
				return nil, err
			}
			w = zapcore.AddSync(rotateLogger)
		}

		level := config.Level.ToZapLevel()
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			w,
			level,
		)
		cores = append(cores, core)
	}

	combinedCore := zapcore.NewTee(cores...)

	options := buildOptions(config)
	logger := zap.New(combinedCore, options...)
	return &zapLogger{logger: logger}, nil
}

// getRotateLogger 创建一个支持日志轮转的 logger
func getRotateLogger(filename string, config *log_config.LogConfig) (logger io.Writer, err error) {
	var dir string
	if dir, err = os.Getwd(); err != nil {
		return nil, err
	}

	filename = path.Join(dir+"/log", filename)
	suffix := ".%Y%m%d%H"
	logger, err = rotatelogs.New(
		filename+suffix,
		rotatelogs.WithLinkName(filename),
		// rotatelogs.WithMaxAge(time.Duration(config.MaxAge)*24*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(config.RotationTime)*time.Hour),
		rotatelogs.WithRotationCount(uint(config.MaxBackups)),
		rotatelogs.WithHandler(rotatelogs.HandlerFunc(func(e rotatelogs.Event) {
			if e.Type() != rotatelogs.FileRotatedEventType {
				return
			}
			if config.Compress {
				go compressLogFile(e.(*rotatelogs.FileRotatedEvent).PreviousFile())
			}
		})),
	)
	return
}

// compressLogFile 压缩日志文件
func compressLogFile(file string) {
	// 实现压缩逻辑，例如使用 gzip
}

// buildOptions 构建 zap logger 的选项
func buildOptions(config *log_config.LogConfig) []zap.Option {
	var opts []zap.Option
	// 添加调用者跳过级别，确保日志显示正确的调用位置
	opts = append(opts, zap.AddCallerSkip(2))

	// 如果 StackTrace 级别不是 FatalLevel，为指定级别及以上的日志添加堆栈跟踪
	opts = append(opts, zap.AddCaller())
	if config.StackTrace != log_level.FatalLevel {
		opts = append(opts, zap.AddStacktrace(config.StackTrace.ToZapLevel()))
	}
	return opts
}

// standardWriter 获取标准输出或标准错误的 writer
func standardWriter(path string) io.Writer {
	switch path {
	case "stdout":
		return os.Stdout
	case "stderr":
		return os.Stderr
	default:
		return nil
	}
}

func otelSeverityEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	switch l {
	case zapcore.DebugLevel:
		enc.AppendString("DEBUG")
	case zapcore.InfoLevel:
		enc.AppendString("INFO")
	case zapcore.WarnLevel:
		enc.AppendString("WARN")
	case zapcore.ErrorLevel:
		enc.AppendString("ERROR")
	case zapcore.DPanicLevel:
		enc.AppendString("FATAL")
	case zapcore.PanicLevel:
		enc.AppendString("FATAL")
	case zapcore.FatalLevel:
		enc.AppendString("FATAL")
	}
}

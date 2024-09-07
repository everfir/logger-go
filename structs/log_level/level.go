package log_level

import (
	"fmt"

	"go.uber.org/zap/zapcore"
)

// Level 定义日志级别
type Level int

// 定义日志级别常量
const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// 将自定义的日志级别转换为Zap的日志级别
func (level Level) ToZapLevel() zapcore.Level {
	switch level {
	case DebugLevel:
		return zapcore.DebugLevel
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	case FatalLevel:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func NewLogLevel(level string) (ret Level, err error) {
	switch level {
	case "debug":
		ret = DebugLevel
	case "info":
		ret = InfoLevel
	case "warn":
		ret = WarnLevel
	case "error":
		ret = ErrorLevel
	case "fatal":
		ret = FatalLevel
	default:
		err = fmt.Errorf("unexpect log level:%s", level)
	}
	return ret, err
}

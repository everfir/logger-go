package logger

import (
	"github.com/everfir/logger-go/structs/log_config"
	"github.com/everfir/logger-go/structs/log_level"
)

// Option 定义配置选项函数类型
type Option func(*log_config.LogConfig)

// WithLevel 设置日志级别
func WithLevel(level log_level.Level) Option {
	return func(c *log_config.LogConfig) {
		c.Level = level
	}
}

// WithStackTrace 设置堆栈跟踪级别
func WithStackTrace(level log_level.Level) Option {
	return func(c *log_config.LogConfig) {
		c.StackTrace = level
	}
}

// WithCompress 设置是否压缩旧日志文件
func WithCompress(compress bool) Option {
	return func(c *log_config.LogConfig) {
		c.Compress = compress
	}
}

// WithMaxBackups 设置旧日志文件最大保留个数
func WithMaxBackups(maxBackups int) Option {
	return func(c *log_config.LogConfig) {
		c.MaxBackups = maxBackups
	}
}

// WithRotationTime 设置日志轮转时间间隔（分钟）
func WithRotationTime(rotationTime int) Option {
	return func(c *log_config.LogConfig) {
		c.RotationTime = rotationTime
	}
}

// WithOutputFiles 设置日志输出文件名
func WithOutputFiles(outputFiles ...string) Option {
	return func(c *log_config.LogConfig) {
		c.OutputFiles = outputFiles
	}
}

// WithErrorFiles 设置错误日志文件名
func WithErrorFiles(errorFiles ...string) Option {
	return func(c *log_config.LogConfig) {
		c.ErrorFiles = errorFiles
	}
}

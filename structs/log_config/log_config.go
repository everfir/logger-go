package log_config

import (
	"os"

	"github.com/everfir/logger-go/structs/log_level"
	"github.com/everfir/logger-go/structs/tracer_config"
)

// LogConfig 定义日志配置结构
type LogConfig struct {
	// 基础信息
	ServiceName string // 服务名称
	PodIP       string // 容器IP

	Level      log_level.Level // 日志级别：定义记录哪个级别及以上的日志
	StackTrace log_level.Level // 堆栈跟踪级别：定义在哪个级别及以上的日志中包含堆栈跟踪

	Compress     bool // 旧日志文件压缩：是否压缩旧的日志文件
	MaxBackups   int  // 旧日志文件最大保留个数：超过此数量的旧文件将被删除
	RotationTime int  // 日志轮转时间间隔（分钟）：多久创建一个新的日志文件

	// 目录为当前工作目录
	OutputFiles []string // 日志输出文件名：日志文件的保存位置，可以是文件路径或 "stdout"/"stderr"

	ErrorFiles []string // 错误日志文件名：错误级别日志的额外输出位置

	// 链路追踪
	TracerConfig *tracer_config.TracerConfig
}

// 默认配置
var DefaultConfig = LogConfig{
	PodIP:       os.Getenv("PodIP"),
	ServiceName: os.Getenv("ServiceName"),

	Level:        log_level.InfoLevel,
	StackTrace:   log_level.FatalLevel,
	OutputFiles:  []string{"stdout"},
	ErrorFiles:   []string{"stderr"},
	TracerConfig: &tracer_config.DefaultTracerConfig,
}

func (config *LogConfig) FixDefault() {
	config.TracerConfig.FixDefault()
}

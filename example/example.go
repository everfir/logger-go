package main

import (
	"context"
	"time"

	"github.com/everfir/logger-go"
	"github.com/everfir/logger-go/structs/field"
	"github.com/everfir/logger-go/structs/log_level"
)

func main() {

	ctx := context.TODO()
	ctx = context.WithValue(ctx, "efr_trace_id", "xxxx_trace_id")
	ctx = context.WithValue(ctx, "efr_span_id", "xxxx_span_id")
	ctx = context.WithValue(ctx, "efr_parent_span_id", "xxxx_parent_span_id")

	ctx = context.WithValue(ctx, "efr_host_name", "train2")
	ctx = context.WithValue(ctx, "efr_pod_name", "train2_pod_1")
	ctx = context.WithValue(ctx, "efr_ip", "127.0.0.1")
	ctx = context.WithValue(ctx, "efr_client_ip", "8.8.8.8")

	if err := Init(); err != nil {
		logger.Fatal(ctx, "初始化日志器失败",
			field.String("error", err.Error()),
		)
	}

	// 测试不同级别的日志和所有类型的字段
	logger.Debug(ctx, "这是一条调试消息",
		field.String("string_field", "debug_value"),
		field.Bool("bool_field", true),
		field.Int("int_field", 42),
		field.Int8("int8_field", 8),
		field.Int16("int16_field", 16),
		field.Int32("int32_field", 32),
		field.Int64("int64_field", 64),
	)

	logger.Info(ctx, "这是一条信息消息",
		field.Uint("uint_field", 42),
		field.Uint8("uint8_field", 8),
		field.Uint16("uint16_field", 16),
		field.Uint32("uint32_field", 32),
		field.Uint64("uint64_field", 64),
		field.Float32("float32_field", 3.14),
		field.Float64("float64_field", 3.14159),
	)

	logger.Warn(ctx, "这是一条警告消息",
		field.Time("time_field", time.Now()),
		field.Duration("duration_field", 5*time.Second),
		field.Any("any_field", struct{ Name string }{"测试结构体"}),
	)

	logger.Error(ctx, "这是一条错误消息",
		field.String("error_type", "测试错误"),
		field.Int("error_code", 500),
	)

}

func Init() (err error) {
	// 初始化日志器,使用多个选项
	err = logger.Init(
		logger.WithLevel(log_level.DebugLevel),
		logger.WithStackTrace(log_level.ErrorLevel),
		logger.WithCompress(true),
		logger.WithMaxBackups(5),
		logger.WithRotationTime(60),
		logger.WithOutputFiles("test.log", "stdout"),
		logger.WithErrorFiles("error.log"),
	)
	return
}

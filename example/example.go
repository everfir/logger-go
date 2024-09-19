package main

import (
	"context"
	"fmt"
	"time"

	"github.com/everfir/logger-go"
	"github.com/everfir/logger-go/structs/field"
	"github.com/everfir/logger-go/structs/log_level"
	"go.opentelemetry.io/otel"
)

func main() {

	//--------------------脚手架负责的工作----------------------------------
	// 初始化日志库
	err := logger.Init(
		logger.WithLevel(log_level.DebugLevel),
		logger.WithOutputFiles("app.log", "stdout"),
		logger.WithTracing(true, "localhost:4318"), // 启用 tracing 并设置 Collector 地址
	)
	if err != nil {
		panic(fmt.Sprintf("初始化日志库失败: %v", err))
	}
	defer func() {
		logger.Close()
	}()

	// 创建一个根 span
	ctx, rootSpan := otel.Tracer("example").Start(context.Background(), "main")
	defer rootSpan.End()
	//--------------------脚手架负责的工作----------------------------------

	// 使用包含 span 的 context 进行日志记录
	logger.Info(ctx, "应用程序启动", field.String("version", "v1.0.0"))

	// 模拟一些操作
	for i := 0; i < 3; i++ {
		func() {
			//--------------------脚手架负责的工作----------------------------------
			ctx, span := otel.Tracer("example").Start(ctx, fmt.Sprintf("operation-%d", i))
			defer span.End()
			//--------------------脚手架负责的工作----------------------------------

			logger.Info(ctx, "执行操作", field.Int("iteration", i))

			// 模拟一些工作
			time.Sleep(100 * time.Millisecond)

			logger.Warn(ctx, "可能的问题", field.String("detail", "某些操作可能不稳定"))
		}()
	}

	// 模拟一个错误
	//--------------------脚手架负责的工作----------------------------------
	ctx, errorSpan := otel.Tracer("example").Start(ctx, "error-operation")
	logger.Error(ctx, "发生错误",
		field.String("error_type", "模拟错误"),
		field.Int("error_code", 500),
	)
	errorSpan.End()
	//--------------------脚手架负责的工作----------------------------------

	logger.Info(ctx, "应用程序结束")

}

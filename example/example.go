package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/everfir/logger-go"
	"github.com/everfir/logger-go/structs/field"
	"github.com/everfir/logger-go/structs/log_level"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// 定义全局变量
var tcer trace.Tracer
var propagator propagation.TextMapPropagator

// 初始化函数
func init() {
	// 初始化 tracer
	tcer = otel.Tracer("example")
}

// tracingMiddleware 注入 tracing 信息到 gin.Context
func tracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := logger.Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))
		ctx, span := tcer.Start(ctx, "tracingMiddleware")
		defer span.End()
		c.Set("tracing", span)

		c.Next()
	}
}

// 服务器处理函数
func serverHandler(c *gin.Context) {
	logger.Info(c, "服务端测试")

	// 响应客户端
	c.String(http.StatusOK, "Hello from server!")
}

// 客户端发送请求函数
func sendRequest(ctx context.Context) error {
	// 创建新的 span
	ctx, span := tcer.Start(ctx, "client-request")
	defer span.End()

	// 创建 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:10083", nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	// 注入 trace 信息到 header
	logger.Inject(ctx, propagation.HeaderCarrier(req.Header))
	req.Header.Set("Traceparent", "00-c1156a8801e4e6e9dd87c18071037df4-4ef6c87f0b8dc73f-01")
	ctx = context.WithValue(ctx, "tracing", span)

	// 记录发送的 header
	logger.Info(ctx, "客户端发送的 header", field.String("headers", fmt.Sprintf("%v", req.Header)))

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %v", err)
	}

	logger.Info(ctx, "收到服务器响应", field.String("body", string(body)))

	return nil
}

func main() {
	// 初始化日志库
	err := logger.Init(
		logger.WithLevel(log_level.DebugLevel),
		// logger.WithOutputFiles("app.log", "stdout"),
		// logger.WithTracing(true, "localhost:4318"),
	)
	if err != nil {
		panic(fmt.Sprintf("初始化日志库失败: %v", err))
	}
	defer logger.Close()

	// 创建一个根 span
	ctx, rootSpan := tcer.Start(context.Background(), "main")
	defer rootSpan.End()

	// 创建 Gin 引擎
	r := gin.Default()

	// 使用 tracing 中间件
	r.Use(tracingMiddleware())

	// 注册路由
	r.GET("/", serverHandler)

	// 启动 HTTP 服务器
	go func() {
		logger.Info(ctx, "启动 HTTP 服务器在 :10083")
		if err := r.Run(":10083"); err != nil {
			logger.Error(ctx, "HTTP 服务器错误", field.String("err", err.Error()))
		}
	}()

	// 等待服务器启动
	time.Sleep(time.Second)

	// 发送 HTTP 请求
	if err := sendRequest(ctx); err != nil {
		logger.Error(ctx, "发送请求失败", field.String("err", err.Error()))
	}

	logger.Info(ctx, "应用程序结束")
}

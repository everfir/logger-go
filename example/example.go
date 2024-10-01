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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// 定义全局变量
var tracer trace.Tracer
var propagator propagation.TextMapPropagator

// 初始化函数
func init() {
	// 初始化 tracer
	tracer = otel.Tracer("example")
	// 初始化 propagator
	propagator = otel.GetTextMapPropagator()
}

// 服务器处理函数
func serverHandler(w http.ResponseWriter, r *http.Request) {
	// 从请求中提取 context
	ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))
	r = r.WithContext(ctx)
	ctx = r.Context()

	logger.Info(ctx, "服务器接收到的 header", field.String("headers", fmt.Sprintf("%v", r.Header)))
	logger.Info(ctx, "服务器接收到的 traceparent", field.String("traceparent", r.Header.Get("Traceparent")))
	logger.Info(ctx, "服务器的ctx", field.Any("ctx", ctx))

	// 创建新的 span
	_, span := tracer.Start(ctx, "server-handler")
	defer span.End()

	// 响应客户端
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello from server!"))
}

// 客户端发送请求函数
func sendRequest(ctx context.Context) error {
	// 创建新的 span
	ctx, span := tracer.Start(ctx, "client-request")
	defer span.End()

	// 创建 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:10083", nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	// 注入 trace 信息到 header
	propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))
	req.Header.Set("Traceparent", "00-c1156a8801e4e6e9dd87c18071037df4-4ef6c87f0b8dc73f-01")

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
	ctx, rootSpan := tracer.Start(context.Background(), "main")
	defer rootSpan.End()

	// 启动 HTTP 服务器
	http.HandleFunc("/", serverHandler)
	go func() {
		logger.Info(ctx, "启动 HTTP 服务器在 :8080")
		if err := http.ListenAndServe(":10083", nil); err != nil {
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

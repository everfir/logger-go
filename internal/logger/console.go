package logger

import (
	"fmt"
	"os"

	"github.com/everfir/logger-go/structs/field"
)

// ConsoleLogger 是一个简单的控制台日志器
type ConsoleLogger struct{}

func (l *ConsoleLogger) log(level, msg string, fields ...field.Field) {
	fmt.Printf("[%s] %s", level, msg)
	for _, f := range fields {
		fmt.Printf(" %s=%v", f.Key(), f.Value())
	}
	fmt.Println()
}

func (l *ConsoleLogger) Debug(msg string, fields ...field.Field) { l.log("DEBUG", msg, fields...) }
func (l *ConsoleLogger) Info(msg string, fields ...field.Field)  { l.log("INFO", msg, fields...) }
func (l *ConsoleLogger) Warn(msg string, fields ...field.Field)  { l.log("WARN", msg, fields...) }
func (l *ConsoleLogger) Error(msg string, fields ...field.Field) { l.log("ERROR", msg, fields...) }
func (l *ConsoleLogger) Fatal(msg string, fields ...field.Field) {
	l.log("FATAL", msg, fields...)
	os.Exit(1)
}

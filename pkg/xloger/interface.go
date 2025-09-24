package xloger

import (
	"context"
	"fmt"
	"log"
)

type Field struct {
	Key   string
	Value interface{}
}

type Logger interface {
	Debug(ctx context.Context, msg string, fields ...Field)
	Info(ctx context.Context, msg string, fields ...Field)
	Warn(ctx context.Context, msg string, fields ...Field)
	Error(ctx context.Context, msg string, fields ...Field)
}

type StdLogger struct{}

func (l *StdLogger) Debug(ctx context.Context, msg string, fields ...Field) {
	log.Println("[DEBUG]", msg, formatFields(fields))
}

func (l *StdLogger) Info(ctx context.Context, msg string, fields ...Field) {
	log.Println("[INFO]", msg, formatFields(fields))
}

func (l *StdLogger) Warn(ctx context.Context, msg string, fields ...Field) {
	log.Println("[WARN]", msg, formatFields(fields))
}

func (l *StdLogger) Error(ctx context.Context, msg string, fields ...Field) {
	log.Println("[ERROR]", msg, formatFields(fields))
}

func formatFields(fields []Field) string {
	s := ""
	for _, f := range fields {
		s += fmt.Sprintf("%s=%v ", f.Key, f.Value)
	}
	return s
}

package logging

import (
	"time"

	"github.com/edouardparis/lntop/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Field = zapcore.Field
type ObjectEncoder = zapcore.ObjectEncoder

type Logger interface {
	Info(string, ...Field)
	Error(string, ...Field)
	Sync() error
	Debug(string, ...Field)
	With(fields ...Field) *zap.Logger
}

func String(k, v string) Field {
	return zap.String(k, v)
}

func Duration(k string, d time.Duration) Field {
	return zap.Duration(k, d)
}

func Int(k string, i int) Field {
	return zap.Int(k, i)
}

func Int64(k string, i int64) Field {
	return zap.Int64(k, i)
}

func Error(v error) Field {
	return zap.Error(v)
}

func Object(key string, val zapcore.ObjectMarshaler) Field {
	return zap.Object(key, val)
}

func New(cfg config.Logger) Logger {
	var logger Logger
	if cfg.Type == "development" {
		logger, _ = NewDevelopmentLogger()
	} else if cfg.Type == "noop" {
		logger, _ = NewNopLogger()
	} else {
		logger, _ = NewProductionLogger()
	}

	return logger
}

func NewProductionLogger() (Logger, error) {
	return zap.NewProduction()
}

func NewDevelopmentLogger() (Logger, error) {
	return zap.NewDevelopment()
}

func NewNopLogger() (Logger, error) {
	return zap.NewNop(), nil
}

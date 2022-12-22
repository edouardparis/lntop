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

func Uint64(k string, i uint64) Field {
	return zap.Uint64(k, i)
}

func Error(v error) Field {
	return zap.Error(v)
}

func Object(key string, val zapcore.ObjectMarshaler) Field {
	return zap.Object(key, val)
}

func New(cfg config.Logger) (Logger, error) {
	switch cfg.Type {
	case "development":
		return NewDevelopmentLogger(cfg.Dest)
	case "production":
		return NewProductionLogger(cfg.Dest)
	default:
		return NewDevelopmentLogger(cfg.Dest)
	}
}

func NewProductionLogger(dest string) (Logger, error) {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{dest}
	return config.Build()
}

func NewDevelopmentLogger(dest string) (Logger, error) {
	config := zap.NewDevelopmentConfig()
	config.OutputPaths = []string{dest}
	return config.Build()
}

func NewNopLogger() (Logger, error) {
	return zap.NewNop(), nil
}

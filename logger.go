package myfx

import (
	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel/log"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type LoggerFactory interface {
	Logger(pkg string) *zap.Logger
}

type loggerFactory struct {
	lp log.LoggerProvider
}

func (lf *loggerFactory) Logger(pkg string) *zap.Logger {
	return zap.New(otelzap.NewCore(pkg, otelzap.WithLoggerProvider(lf.lp)))
}

func NewLoggerFactory(lp log.LoggerProvider) LoggerFactory {
	return &loggerFactory{lp: lp}
}

func ProvidePrivateLogger(name string) fx.Option {
	return fx.Provide(
		fx.Private,
		func(lf LoggerFactory) *zap.Logger {
			return lf.Logger(name)
		},
	)
}

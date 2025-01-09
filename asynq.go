package myfx

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type asynqLogger struct {
	log *zap.SugaredLogger
}

func (log *asynqLogger) Debug(v ...interface{}) {
	log.log.Debug(v...)
}

func (log *asynqLogger) Info(v ...interface{}) {
	log.log.Info(v...)
}

func (log *asynqLogger) Warn(v ...interface{}) {
	log.log.Warn(v...)
}

func (log *asynqLogger) Error(v ...interface{}) {
	log.log.Error(v...)
}

func (log *asynqLogger) Fatal(v ...interface{}) {
	log.log.Fatal(v...)
}

func NewAsynqLogger(log *zap.Logger) asynq.Logger {
	return &asynqLogger{log: log.Sugar()}
}

func NewAsynqServer(r *redis.Client, logger asynq.Logger) *asynq.Server {
	return asynq.NewServerFromRedisClient(r, asynq.Config{
		Concurrency: 5,
		Logger:      logger,
	})
}

func NewAsynqClient(r *redis.Client, lc fx.Lifecycle) *asynq.Client {
	c := asynq.NewClientFromRedisClient(r)
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return c.Close()
		},
	})
	return c
}

package myfx

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/yankeguo/rg"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
)

type zapGormLogger struct {
	logger *zap.Logger
	level  logger.LogLevel
}

func (z *zapGormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return &zapGormLogger{
		logger: z.logger,
		level:  level,
	}
}

func (z *zapGormLogger) Info(ctx context.Context, s string, i ...interface{}) {
	if z.level < logger.Info {
		return
	}
	z.logger.Info(fmt.Sprintf(s, i...), zap.Any("ctx", ctx))
}

func (z *zapGormLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	if z.level < logger.Warn {
		return
	}
	z.logger.Warn(fmt.Sprintf(s, i...), zap.Any("ctx", ctx))
}

func (z *zapGormLogger) Error(ctx context.Context, s string, i ...interface{}) {
	if z.level < logger.Error {
		return
	}
	z.logger.Error(fmt.Sprintf(s, i...), zap.Any("ctx", ctx))
}

func (z *zapGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if z.level < logger.Silent {
		return
	}

	duration := time.Since(begin)

	switch {
	case err != nil && z.level >= logger.Error && !errors.Is(err, gorm.ErrRecordNotFound):
		sql, rows := fc()
		z.logger.Error("sql failed", zap.Any("ctx", ctx), zap.Error(err), zap.String("sql", sql), zap.Int64("rows", rows), zap.Duration("duration", duration))
	case duration > time.Millisecond*500 && z.level >= logger.Warn:
		sql, rows := fc()
		z.logger.Warn("slow sql", zap.Any("ctx", ctx), zap.String("sql", sql), zap.Int64("rows", rows), zap.Duration("duration", duration))
	case z.level >= logger.Info:
		sql, rows := fc()
		z.logger.Warn("sql executed", zap.Any("ctx", ctx), zap.String("sql", sql), zap.Int64("rows", rows), zap.Duration("duration", duration))
	}
}

func NewGormDB(vb Verbose, log *zap.Logger, tp trace.TracerProvider) (db *gorm.DB, err error) {
	defer rg.Guard(&err)

	level := logger.Warn

	if vb.Get() {
		level = logger.Info
	}

	cfg := rg.Must(mysqldriver.ParseDSN(os.Getenv("MYSQL_DSN")))

	db = rg.Must(
		gorm.Open(
			mysql.New(mysql.Config{
				DSNConfig: cfg,
			}),
			&gorm.Config{
				Logger: &zapGormLogger{
					logger: log,
					level:  level,
				},
			},
		),
	)

	rg.Must0(db.Use(
		tracing.NewPlugin(
			tracing.WithDBName(cfg.Addr+"/"+cfg.DBName),
			tracing.WithTracerProvider(tp),
		),
	))
	return
}

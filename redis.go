package myfx

import (
	"os"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/yankeguo/rg"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

func NewRedisClient(tp trace.TracerProvider, mp metric.MeterProvider) (r *redis.Client, err error) {
	defer rg.Guard(&err)
	r = redis.NewClient(rg.Must(redis.ParseURL(os.Getenv("REDIS_URL"))))
	rg.Must0(redisotel.InstrumentTracing(r, redisotel.WithTracerProvider(tp)))
	rg.Must0(redisotel.InstrumentMetrics(r, redisotel.WithMeterProvider(mp)))
	return
}

package myfx

import (
	"context"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	otellog "go.opentelemetry.io/otel/log"
	logglobal "go.opentelemetry.io/otel/log/global"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

type sampleSwitchKeyType int
type sampleSwitchValueType int

var (
	sampleSwitchKey sampleSwitchKeyType

	sampleSwitchNone sampleSwitchValueType = 0
	sampleSwitchOn   sampleSwitchValueType = 1
	sampleSwitchOff  sampleSwitchValueType = -1
)

type switchSampler struct {
	defaultValue bool
}

func (as *switchSampler) ShouldSample(p trace.SamplingParameters) trace.SamplingResult {
	var decision trace.SamplingDecision
	if as.defaultValue {
		decision = trace.RecordAndSample
		if p.ParentContext.Value(sampleSwitchKey) == sampleSwitchOff {
			decision = trace.Drop
		}
	} else {
		decision = trace.Drop
		if p.ParentContext.Value(sampleSwitchKey) == sampleSwitchOn {
			decision = trace.RecordAndSample
		}
	}
	return trace.SamplingResult{
		Decision:   decision,
		Tracestate: oteltrace.SpanContextFromContext(p.ParentContext).TraceState(),
	}
}

func (as *switchSampler) Description() string {
	return "SwitchSampler"
}

func SampleSwitch(defaultValue bool) trace.Sampler {
	return &switchSampler{
		defaultValue: defaultValue,
	}
}

func ProvideSwitchSampler(defaultValue bool) fx.Option {
	return fx.Provide(func() trace.Sampler {
		return SampleSwitch(defaultValue)
	})
}

func SetSampleSwitch(ctx context.Context, enabled bool) context.Context {
	value := sampleSwitchOff
	if enabled {
		value = sampleSwitchOn
	}
	return context.WithValue(ctx, sampleSwitchKey, value)
}

func ClearSampleSwitch(ctx context.Context) context.Context {
	return context.WithValue(ctx, sampleSwitchKey, sampleSwitchNone)
}

func NewOTELPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func NewOTELResource() (res *resource.Resource, err error) {
	return resource.New(
		context.Background(),
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithProcess(),
		resource.WithOS(),
		resource.WithContainer(),
		resource.WithHost(),
	)
}

func NewOTELTraceProvider(lc fx.Lifecycle, res *resource.Resource, s trace.Sampler) (trace.SpanExporter, oteltrace.TracerProvider, error) {
	var (
		tp  *trace.TracerProvider
		te  trace.SpanExporter
		err error
	)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) (_ error) {
			if te != nil {
				te.Shutdown(ctx)
			}
			if tp != nil {
				tp.Shutdown(ctx)
			}
			return
		},
	})

	if te, err = otlptracegrpc.New(context.Background()); err != nil {
		return nil, nil, err
	}

	tp = trace.NewTracerProvider(
		trace.WithBatcher(te, trace.WithBatchTimeout(time.Second*3)),
		trace.WithResource(res),
		trace.WithSampler(s),
	)

	return te, tp, nil
}

func NewOTELMetricProvider(lc fx.Lifecycle, res *resource.Resource) (metric.Exporter, otelmetric.MeterProvider, error) {
	var (
		me  metric.Exporter
		mp  *metric.MeterProvider
		err error
	)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) (_ error) {
			if me != nil {
				me.Shutdown(ctx)
			}
			if mp != nil {
				mp.Shutdown(ctx)
			}
			return
		},
	})

	if me, err = otlpmetricgrpc.New(context.Background()); err != nil {
		return nil, nil, err
	}

	mp = metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(me, metric.WithInterval(time.Second*3))),
		metric.WithResource(res),
	)
	return me, mp, nil
}

func NewOTELLoggerProvider(lc fx.Lifecycle, res *resource.Resource) (log.Exporter, otellog.LoggerProvider, error) {
	var (
		le  log.Exporter
		lp  *log.LoggerProvider
		err error
	)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) (_ error) {
			if le != nil {
				le.Shutdown(ctx)
			}
			if lp != nil {
				lp.Shutdown(ctx)
			}
			return
		},
	})

	if le, err = otlploggrpc.New(context.Background()); err != nil {
		return nil, nil, err
	}

	lp = log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(le)),
		log.WithResource(res),
	)
	return le, lp, nil
}

func SetupOTEL(pr propagation.TextMapPropagator, tp oteltrace.TracerProvider, mp otelmetric.MeterProvider, lp otellog.LoggerProvider) (err error) {
	otel.SetTextMapPropagator(pr)
	otel.SetTracerProvider(tp)
	otel.SetMeterProvider(mp)
	logglobal.SetLoggerProvider(lp)

	runtime.Start(runtime.WithMeterProvider(mp))
	return
}

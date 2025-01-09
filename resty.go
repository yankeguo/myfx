package myfx

import (
	"net/http"

	"github.com/dubonzi/otelresty"
	"github.com/go-resty/resty/v2"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func NewRestyClient(log *zap.Logger, tp trace.TracerProvider, vb Verbose) *resty.Client {
	r := resty.New()
	r = r.SetPreRequestHook(func(c *resty.Client, req *http.Request) (_ error) {
		if vb.Get() {
			log.Info("resty request", zap.String("method", req.Method), zap.String("url", req.URL.String()), zap.Any("ctx", req.Context()))
		}
		return
	})
	otelresty.TraceClient(r, otelresty.WithTracerProvider(tp))
	return r
}

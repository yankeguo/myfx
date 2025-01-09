package myfx

import (
	"context"

	"github.com/yankeguo/rg"
	"go.uber.org/zap"
)

func SetupRoyalGuard(vb Verbose, log *zap.Logger) {
	rg.OnGuardWithContext = func(ctx context.Context, r any) {
		if !vb.Get() {
			return
		}
		if err, ok := r.(error); ok {
			log.Error("royal guard", zap.Error(err), zap.Any("ctx", ctx))
		} else {
			log.Error("royal guard", zap.Any("err", r), zap.Any("ctx", ctx))
		}
	}
}

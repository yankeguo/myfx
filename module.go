package myfx

import "go.uber.org/fx"

var Module = fx.Module(
	"myfx",
	ProvidePrivateLogger("github.com/yankeguo/myfx"),
	fx.Provide(
		NewAsynqClient,
		NewAsynqLogger,
		NewAsynqServer,
		NewFiberServer,
		NewGormDB,
		NewLoggerFactory,
		NewOTELLoggerProvider,
		NewOTELMetricProvider,
		NewOTELPropagator,
		NewOTELResource,
		NewOTELTraceProvider,
		NewRedisClient,
		NewRedsync,
		NewRestyClient,
		NewVerbose,
		NewZhipuClient,
	),
	fx.Invoke(
		SetupRoyalGuard,
		SetupOTEL,
	),
)

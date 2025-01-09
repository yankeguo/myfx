package myfx

import "go.uber.org/fx"

var Module = fx.Module(
	"myfx",
	ProvidePrivateLogger("github.com/yankeguo/myfx"),
	fx.Provide(
		NewVerbose,
		NewLoggerFactory,
		NewOTELPropagator,
		NewOTELResource,
		NewOTELMetricProvider,
		NewOTELTraceProvider,
		NewOTELLoggerProvider,
		NewAsynqLogger,
		NewAsynqServer,
		NewAsynqClient,
		NewFiberServer,
		NewRestyClient,
		NewGormDB,
		NewRedisClient,
		NewZhipuClient,
	),
	fx.Invoke(
		SetupRoyalGuard,
		SetupOTEL,
	),
)

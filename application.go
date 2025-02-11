package enclave

import (
	"log/slog"

	"github.com/hadroncorp/geck/transport/transportfx/httpfx"
	"github.com/joho/godotenv"

	"github.com/hadroncorp/geck/application/applicationfx"
	"github.com/hadroncorp/geck/observability/observabilityfx/loggingfx"
	"github.com/hadroncorp/geck/persistence/driver/postgres/postgresfx"
	"github.com/hadroncorp/geck/persistence/persistencefx/identifierfx"
	"github.com/hadroncorp/geck/persistence/persistencefx/sqlfx"
	"github.com/hadroncorp/geck/validation/validationfx"
	"go.uber.org/fx"
)

type appOptions struct {
	fxOptions []fx.Option
}

type ApplicationOption func(o *appOptions)

func WithOptions(options ...fx.Option) ApplicationOption {
	return func(o *appOptions) {
		if o.fxOptions == nil {
			o.fxOptions = make([]fx.Option, 0, len(options))
		}
		o.fxOptions = append(o.fxOptions, options...)
	}
}

func WithPostgres() ApplicationOption {
	return func(o *appOptions) {
		opts := []fx.Option{
			postgresfx.Module,
			sqlfx.GoquModule,
			sqlfx.InterceptorModule,
			sqlfx.ObservabilityModule,
		}
		if o.fxOptions == nil {
			o.fxOptions = make([]fx.Option, 0, len(opts))
		}
		o.fxOptions = append(o.fxOptions, opts...)
	}
}

func WithServerHTTP() ApplicationOption {
	return func(o *appOptions) {
		if o.fxOptions == nil {
			o.fxOptions = make([]fx.Option, 0, 1)
		}
		o.fxOptions = append(o.fxOptions, httpfx.ServerModule)
	}
}

// NewApp allocates a runnable application object ([fx.App])
func NewApp(opts ...ApplicationOption) *fx.App {
	options := appOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	if err := godotenv.Load(); err != nil {
		slog.Warn("error loading .env file, using environment variables")
	}

	initOpts := []fx.Option{
		loggingfx.SlogModule,
		applicationfx.Module,
		validationfx.GoPlaygroundModule,
		identifierfx.KSUIDModule,
	}
	initOpts = append(initOpts, options.fxOptions...)
	app := fx.New(initOpts...)
	return app
}

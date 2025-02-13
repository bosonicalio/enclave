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

// WithOptions appends `options` ([fx.Option]) slice into an application base [fx.Option] slice.
func WithOptions(options ...fx.Option) ApplicationOption {
	return func(o *appOptions) {
		if o.fxOptions == nil {
			o.fxOptions = make([]fx.Option, 0, len(options))
		}
		o.fxOptions = append(o.fxOptions, options...)
	}
}

// WithPostgres appends [fx.Option] required to use a Postgres database in the application.
//
// Database client is configured with postgres.DBConfig, which is populated using environment variables;
// Given environment variables are specified in the same config structure by the field tag `env` whereas
// default values (if any) are specified by the field tag `envDefault`.
func WithPostgres(enableLogging bool) ApplicationOption {
	return func(o *appOptions) {
		opts := []fx.Option{
			postgresfx.Module,
			sqlfx.GoquModule,
			sqlfx.InterceptorModule,
		}
		if enableLogging {
			opts = append(opts, sqlfx.ObservabilityModule)
		}
		if o.fxOptions == nil {
			o.fxOptions = make([]fx.Option, 0, len(opts))
		}
		o.fxOptions = append(o.fxOptions, opts...)
	}
}

// WithServerHTTP appends [fx.Option] required to spin up an HTTP server.
//
// Call [httpfx.AsController] in non-`enclave` modules so `enclave` can detect, register and serve
// specified controllers through the provisioned HTTP server.
//
// Database client is configured with http.ServerConfig, which is populated using environment variables;
// Given environment variables are specified in the same config structure by the field tag `env` whereas
// default values (if any) are specified by the field tag `envDefault`.
func WithServerHTTP() ApplicationOption {
	return func(o *appOptions) {
		if o.fxOptions == nil {
			o.fxOptions = make([]fx.Option, 0, 1)
		}
		o.fxOptions = append(o.fxOptions, httpfx.ServerModule)
	}
}

// New allocates a runnable application object ([fx.App]) with a slice of [fx.Option] objects for
// basic system functionalities (logging, application-aware metadata, validation, unique-identifier generation).
//
// This routine will load environment variables from `.env` file located alongside application binary. If not existent,
// it will use operating-system (OS) environment variables instead. `enclave` recommends using OS-provisioned
// variables in cloud environments and `.env` files for local development.
//
// Use [ApplicationOption] directives to customize the final application
// (e.g. [WithServerHTTP]).
func New(opts ...ApplicationOption) *fx.App {
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

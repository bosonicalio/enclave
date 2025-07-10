package enclave

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/tesserical/geck/applicationfx"
	"github.com/tesserical/geck/observabilityfx/loggingfx"
	"github.com/tesserical/geck/persistence/postgres/postgresfx"
	"github.com/tesserical/geck/persistencefx/identifierfx"
	"github.com/tesserical/geck/persistencefx/sqlfx"
	"github.com/tesserical/geck/transportfx/httpfx"
	"github.com/tesserical/geck/validationfx"
	"go.uber.org/fx"

	"github.com/tesserical/enclave/persistencefx"
)

// RunApplication initializes the application with the provided options.
//
// It loads environment variables from a .env file if it exists, uses OS environment variables otherwise.
//
// This routine is designed to be used as the entry point for the application, setting up the necessary dependencies
// and configurations required for the application to run properly. It uses the [go.uber.org/fx] framework
// to manage the lifecycle of the application and its components. No extra steps for initialization are needed
// after calling this function, as it will automatically start the application and manage its lifecycle.
//
// In addition, this routine sets up the application with basic modules like application metadata (name, version
// and environment), logging (with stdlib [slog] package), validations ([github.com/go-playground/validator/v10]),
// identifier factory (KSUID format) and persistence (e.g. pagination) APIs.
func RunApplication(opts ...Option) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := godotenv.Load(); err != nil {
		logger.Warn("failed to load .env file, using OS environment variables", slog.String("error", err.Error()))
	}

	options := &option{
		fxOpts: []fx.Option{
			applicationfx.Module,
			loggingfx.SlogModule,
			validationfx.GoPlaygroundModule,
			identifierfx.KSUIDModule,
			persistencefx.Module,
		},
	}
	for _, opt := range opts {
		opt(options)
	}

	if options.disableDepInjectorLogs {
		options.fxOpts = append(options.fxOpts, fx.NopLogger)
	}

	fx.New(options.fxOpts...).Run()
}

// -- Options --

type option struct {
	disableDepInjectorLogs bool
	fxOpts                 []fx.Option
}

// Option is a function that modifies the enclave application options.
type Option func(*option)

// WithDisabledDepInjectorLogs disables the dependency injector logs.
func WithDisabledDepInjectorLogs() Option {
	return func(options *option) {
		options.disableDepInjectorLogs = true
	}
}

// WithFxOptions appends the provided fx options to the enclave application.
func WithFxOptions(opts ...fx.Option) Option {
	return func(options *option) {
		options.fxOpts = append(options.fxOpts, opts...)
	}
}

// WithServerHTTP adds the HTTP server module to the enclave application.
func WithServerHTTP() Option {
	return WithFxOptions(
		httpfx.ServerModule,
	)
}

// WithSQL adds the SQL module to the enclave application.
func WithSQL() Option {
	return WithFxOptions(
		sqlfx.InterceptorModule,
	)
}

// WithTransactionContextSQL adds the SQL transaction context module to the enclave application.
func WithTransactionContextSQL() Option {
	return WithFxOptions(
		sqlfx.TransactionModule,
	)
}

// WithObservabilitySQL adds the SQL observability module to the enclave application.
func WithObservabilitySQL() Option {
	return WithFxOptions(
		sqlfx.ObservabilityModule,
	)
}

// WithPostgres adds the Postgres module to the enclave application.
func WithPostgres() Option {
	return WithFxOptions(
		postgresfx.Module,
	)
}

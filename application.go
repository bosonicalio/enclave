package enclave

import (
	"log/slog"
	"testing"

	"github.com/joho/godotenv"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"github.com/bosonicalio/enclave/internal/applicationfx"
	"github.com/bosonicalio/enclave/internal/globallog"
	"github.com/bosonicalio/enclave/internal/observabilityfx/loggingfx"
	"github.com/bosonicalio/enclave/internal/persistencefx"
	"github.com/bosonicalio/enclave/internal/persistencefx/sqlfx"
	"github.com/bosonicalio/enclave/internal/transportfx/httpfx"
	"github.com/bosonicalio/enclave/internal/validationfx"
)

// NewApplication creates a new enclave application with the provided options.
//
// It allocates the application with basic modules like application metadata (name, version
// and environment) and logging (with stdlib [slog] package).
//
// The application is built using the [go.uber.org/fx] framework, which provides a powerful
// dependency injection mechanism and lifecycle management.
//
// The application can be customized with additional options, such as disabling dependency
// injector logs, adding custom fx options, or including specific modules like HTTP server,
// validation, and persistence.
//
// The returned application can be run using the [fx.App.Run] method, which will start the
// application and manage its lifecycle.
func NewApplication(opts ...Option) *fx.App {
	options := &option{
		fxOpts: []fx.Option{
			fx.RecoverFromPanics(),
			applicationfx.Module,
			loggingfx.ModuleSlog,
		},
	}
	for _, opt := range opts {
		opt(options)
	}

	if options.disableDepInjectorLogs {
		options.fxOpts = append(options.fxOpts, fx.NopLogger)
	}
	return fx.New(options.fxOpts...)
}

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
// and environment) and logging (with stdlib [slog] package).
func RunApplication(opts ...Option) {
	if err := godotenv.Load(); err != nil {
		globallog.Logger().
			Warn("failed to load .env file, using OS environment variables", slog.String("error", err.Error()))
	}
	NewApplication(opts...).Run()
}

// -- Testing --

// NewTestApplication creates a new enclave application with the provided options.
//
// It allocates the application with basic modules like application metadata (name, version
// and environment) and logging (with stdlib [slog] package).
//
// The application is built using the [go.uber.org/fx] framework, which provides a powerful
// dependency injection mechanism and lifecycle management.
//
// The application can be customized with additional options, such as disabling dependency
// injector logs, adding custom fx options, or including specific modules like HTTP server,
// validation, and persistence.
//
// The returned application can be run using the [fx.App.Run] method, which will start the
// application and manage its lifecycle.
func NewTestApplication(tb testing.TB, opts ...Option) *fxtest.App {
	options := &option{
		fxOpts: []fx.Option{
			fx.RecoverFromPanics(),
			applicationfx.Module,
			loggingfx.ModuleSlog,
		},
	}
	for _, opt := range opts {
		opt(options)
	}
	return fxtest.New(tb, options.fxOpts...)
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

// WithValidation adds the validation module to the enclave application.
//
// This module provides a validation engine based on [github.com/go-playground/validator/v10].
func WithValidation() Option {
	return WithFxOptions(
		validationfx.Module,
	)
}

// WithPersistence adds the persistence module to the enclave application.
//
// This module provides generic components for most persistence-related operations
// (e.g. pagination API, identifier factory, transaction manager).
func WithPersistence() Option {
	return WithFxOptions(
		persistencefx.Module,
	)
}

// WithSQL adds the SQL database module to the enclave application.
//
// Requires an external module to provide the database connection (sql.DB). Available drivers are in different
// go modules (e.g. enclave/postgres).
func WithSQL() Option {
	return WithFxOptions(
		sqlfx.Module,
	)
}

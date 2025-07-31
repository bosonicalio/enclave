package applicationfx

import (
	"github.com/bosonicalio/geck/environment"
	"github.com/bosonicalio/geck/version"
	"go.uber.org/fx"

	"github.com/bosonicalio/geck/application"

	"github.com/bosonicalio/enclave/internal/osenv"
)

// Module is the `uber/fx` module of the [application] package.
var Module = fx.Module("enclave/application",
	fx.Provide(
		osenv.ParseAs[config],
		newApp,
	),
	fx.Invoke(
		logAppStart,
	),
)

// -- Factory --

func newApp(cfg config) (application.Application, error) {
	app, err := application.New(
		application.WithName(cfg.Name),
		application.WithInstanceID(cfg.InstanceID),
	)
	if err != nil {
		return application.Application{}, err
	}
	app.Version, err = version.Parse(cfg.Version)
	if err != nil {
		return application.Application{}, err
	}
	app.Environment, err = environment.Parse(cfg.Environment)
	if err != nil {
		return application.Application{}, err
	}
	return app, nil
}

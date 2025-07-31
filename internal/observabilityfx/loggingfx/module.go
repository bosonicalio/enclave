package loggingfx

import (
	"go.uber.org/fx"

	"github.com/bosonicalio/geck/observability/logging"
)

// ModuleSlog is the `uber/fx` module of the [logging] package, using
// stdlib `slog` package for concrete implementations.
var ModuleSlog = fx.Module("enclave/observability/logging/slog",
	fx.Provide(
		logging.NewSlogLogger,
	),
)

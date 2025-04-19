package persistencefx

import (
	"github.com/hadroncorp/geck/persistence/paging"
	"go.uber.org/fx"
)

// Module is the [fx] module for the persistence API.
var Module = fx.Module("enclave/persistence",
	fx.Provide(
		paging.NewTokenConfig,
	),
)

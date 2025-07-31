package postgres

import (
	"github.com/bosonicalio/enclave"
)

// WithPostgres returns an enclave option that includes the Postgres module.
func WithPostgres() enclave.Option {
	return enclave.WithFxOptions(
		module,
	)
}

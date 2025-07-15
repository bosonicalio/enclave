package osenv

import (
	"context"

	"github.com/caarlos0/env/v11"
)

// ParseAs parses the given struct type containing `env` tags and loads its
// values from environment variables.
//
// Based on [env.ParseAs], extends the functionality to validate the
// parsed structure using the global validator.
func ParseAs[T any]() (T, error) {
	st, err := env.ParseAs[T]()
	if err != nil {
		var zeroVal T
		return zeroVal, err
	}
	if err = GlobalValidator().Validate(context.Background(), st); err != nil {
		var zeroVal T
		return zeroVal, err
	}
	return st, nil
}

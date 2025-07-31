package persistencefx

import (
	"errors"

	"github.com/bosonicalio/geck/persistence"
	"github.com/bosonicalio/geck/persistence/identifier"
	"github.com/bosonicalio/geck/persistence/paging"
	"github.com/samber/lo"
	"go.uber.org/fx"

	"github.com/bosonicalio/enclave/internal/globallog"
	"github.com/bosonicalio/enclave/internal/osenv"
)

// Module is the [fx] module for the persistence API.
//
// It provides a basic pagination API with [paging.TokenCipherKey] and an identifier factory
// with [identifier.FactoryKSUID] and [identifier.FactoryUUID] implementations.
//
// The basic identifier factory is KSUID-based, which is a globally unique identifier format. For scenarios where
// UUIDs are preferred, the `ID_FACTORY_DRIVER` environment variable can be set to `uuid`.
//
// For additional factories (e.g. UUID), please use concrete packages directly ([github.com/segmentio/ksuid],
// [github.com/google/uuid]).
var Module = fx.Module("enclave/persistence",
	fx.Provide(
		osenv.ParseAs[tokenConfig],
		newTokenCipherKey,
		osenv.ParseAs[identifierConfig],
		newIdentifierFactory,
		persistence.NewTxManager,
	),
)

// -- Factory --

func newTokenCipherKey(config tokenConfig) (paging.TokenCipherKey, error) {
	if config.CipherKey == "" {
		logMsg := `Using default page token cipher key. For enhanced security, please set PAGE_TOKEN_CIPHER_KEY environment variable to a 16, 24 or 32 bytes long key`
		globallog.Logger().
			Warn(logMsg)
		return paging.TokenCipherKey(lo.RandomString(32, lo.AllCharset)), nil
	}

	notValid := len(config.CipherKey) != 16 && len(config.CipherKey) != 24 && len(config.CipherKey) != 32
	if notValid {
		return nil, errors.New("invalid page token cipher key length, must be 16, 24 or 32 bytes")
	}
	return paging.TokenCipherKey(config.CipherKey), nil
}

func newIdentifierFactory(config identifierConfig) (identifier.Factory, error) {
	switch config.Driver {
	case "ksuid":
		return identifier.FactoryKSUID{}, nil
	case "uuid":
		return identifier.FactoryUUID{}, nil
	default:
		return nil, errors.New("invalid identifier factory driver, must be 'ksuid' or 'uuid'")
	}
}

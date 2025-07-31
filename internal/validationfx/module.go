package validationfx

import (
	"errors"

	"github.com/bosonicalio/geck/validation"
	"go.uber.org/fx"

	"github.com/bosonicalio/enclave/internal/osenv"
)

// Module is the `uber/fx` module of the [validation] package, offering
// implementations with the third-party `go-playground/validator` package and [validation.JSONDriver] by default.
//
// This provides a global validator for the application, which can be used to validate structures.
//
// For additional validators (using YAML, TOML and XML codecs), please instantiate [validation.Validator] directly.
var Module = fx.Module("enclave/validation",
	fx.Provide(
		osenv.ParseAs[config],
		newValidator,
	),
)

// -- Factory --

func newValidator(cfg config) (validation.Validator, error) {
	rules := make([]validation.Rule, 0, len(cfg.CustomRules))
	for _, ruleName := range cfg.CustomRules {
		switch ruleName {
		case "date":
			rules = append(rules, validation.NewDateRule())
		}
	}

	codecDriver, err := validation.ParseCodecDriver(cfg.CodecDriver)
	if err != nil {
		return nil, err
	}

	switch cfg.Driver {
	case "go-playground":
		return validation.NewGoPlaygroundValidator(
			validation.WithRules(rules...),
			validation.WithCodecDriver(codecDriver),
		), nil
	default:
		return nil, errors.New("enclave.validation: unsupported driver")
	}
}

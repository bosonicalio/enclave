package validationfx

import (
	"fmt"

	"go.uber.org/fx"

	"github.com/tesserical/geck/validation"

	"github.com/tesserical/enclave/internal/osenv"
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

// -- Config --

type config struct {
	Driver      string   `env:"VALIDATION_DRIVER" envDefault:"go-playground" validate:"required,oneof=go-playground"`
	CustomRules []string `env:"VALIDATION_CUSTOM_RULES" envDefault:"date" validate:"dive,oneof=date"`
}

// -- Factory --

func newValidator(cfg config) (validation.Validator, error) {
	driver, err := validation.FormatDriver(cfg.Driver)
	if err != nil {
		return nil, err
	}

	rules := make([]validation.Rule, 0, len(cfg.CustomRules))
	for _, ruleName := range cfg.CustomRules {
		switch ruleName {
		case "date":
			rules = append(rules, validation.NewDateRule())
		default:
			return nil, fmt.Errorf("enclave.validation: unsupported custom rule %s", ruleName)
		}
	}

	switch driver {
	case validation.GoPlaygroundDriver:
		return validation.NewGoPlaygroundValidator(
			validation.WithRules(rules...),
			validation.WithCodecDriver(validation.JSONDriver),
		), nil
	default:
		return nil, fmt.Errorf("enclave.validation: unsupported driver %s", driver)
	}
}

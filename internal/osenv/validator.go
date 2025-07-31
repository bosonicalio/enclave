package osenv

import (
	"sync"

	"github.com/bosonicalio/geck/validation"
)

var (
	_globalValidatorOnce sync.Once
	_globalValidator     validation.Validator
)

// GlobalValidator returns a global validator for environment variables.
func GlobalValidator() validation.Validator {
	_globalValidatorOnce.Do(func() {
		_globalValidator = validation.NewGoPlaygroundValidator(
			validation.WithCodecDriver(validation.EnvironmentDriver),
			validation.WithRules(
				validation.NewDateRule(),
			),
		)
	})
	return _globalValidator
}

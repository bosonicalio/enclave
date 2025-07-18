package validationfx

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tesserical/enclave/internal/osenv"
)

func TestNewValidation(t *testing.T) {
	// No valid environment variables set
	validator, err := newValidator(config{})
	assert.NotNil(t, err)
	assert.Nil(t, validator)

	// All required environment variables set - invalid codec
	validator, err = newValidator(config{
		Driver:      "go-playground",
		CodecDriver: "foo",
	})
	assert.NotNil(t, err)
	assert.Nil(t, validator)

	// All required environment variables set
	t.Setenv("VALIDATION_DRIVER", "go-playground")
	t.Setenv("VALIDATION_CODEC_DRIVER", "json")
	cfg, err := osenv.ParseAs[config]()
	assert.NoError(t, err)
	validator, err = newValidator(cfg)
	assert.Nil(t, err)
	assert.NotNil(t, validator)

	// Custom rules set
	t.Setenv("VALIDATION_CUSTOM_RULES", "date")
	cfg, err = osenv.ParseAs[config]()
	assert.NoError(t, err)
	validator, err = newValidator(cfg)
	assert.Nil(t, err)
	assert.NotNil(t, validator)
}

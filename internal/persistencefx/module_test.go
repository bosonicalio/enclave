package persistencefx

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bosonicalio/enclave/internal/osenv"
)

func TestNewCipherKey(t *testing.T) {
	// Default key not set
	key, err := newTokenCipherKey(tokenConfig{})
	assert.NoError(t, err)
	assert.NotNil(t, key)

	// Invalid key length
	t.Setenv("PAGE_TOKEN_CIPHER_KEY", strings.Repeat("a", 18))
	tokenCfg, err := osenv.ParseAs[tokenConfig]()
	assert.NoError(t, err)
	key, err = newTokenCipherKey(tokenCfg)
	assert.Error(t, err)
	assert.Nil(t, key)

	// Valid key length
	t.Setenv("PAGE_TOKEN_CIPHER_KEY", strings.Repeat("a", 16))
	tokenCfg, err = osenv.ParseAs[tokenConfig]()
	assert.NoError(t, err)
	key, err = newTokenCipherKey(tokenCfg)
	assert.NoError(t, err)
	assert.NotNil(t, key)
}

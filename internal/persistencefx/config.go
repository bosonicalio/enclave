package persistencefx

type tokenConfig struct {
	CipherKey string `env:"PAGE_TOKEN_CIPHER_KEY" validate:"omitempty,len=16|len=24|len=32"`
}

type identifierConfig struct {
	Driver string `env:"ID_FACTORY_DRIVER" envDefault:"ksuid" validate:"required,oneof=ksuid uuid"`
}

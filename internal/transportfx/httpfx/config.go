package httpfx

type serverConfig struct {
	Address          string `env:"HTTP_SERVER_ADDRESS" envDefault:":8080"`
	ErrResponseCodec string `env:"HTTP_SERVER_ERR_RESP_CODEC" envDefault:"json" validate:"omitempty,oneof=json xml text"`

	EnableTLS     bool `env:"HTTP_SERVER_ENABLE_TLS"`
	EnableAutoTLS bool `env:"HTTP_SERVER_ENABLE_AUTO_TLS"`
}

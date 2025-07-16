package applicationfx

type config struct {
	Name        string `env:"ENCLAVE_APP_NAME" validate:"required"`
	Version     string `env:"ENCLAVE_APP_VERSION" envDefault:"v0.1.0-alpha" validate:"required"`
	Environment string `env:"ENCLAVE_APP_ENVIRONMENT" envDefault:"local" validate:"required,oneof=production staging development prod stage stg dev local sandbox pilot snx"`
	InstanceID  string `env:"ENCLAVE_APP_INSTANCE_ID"`
}

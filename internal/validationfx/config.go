package validationfx

type config struct {
	Driver      string   `env:"VALIDATION_DRIVER" envDefault:"go-playground" validate:"required,oneof=go-playground"`
	CodecDriver string   `env:"VALIDATION_CODEC_DRIVER" envDefault:"json" validate:"required,oneof=json yaml toml xml"`
	CustomRules []string `env:"VALIDATION_CUSTOM_RULES" envDefault:"date" validate:"dive,oneof=date"`
}

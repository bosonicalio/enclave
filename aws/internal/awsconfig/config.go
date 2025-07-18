package awsconfig

// Config holds the configuration for the AWS SDK, including the region and endpoint URL.
type Config struct {
	Region      string `env:"AWS_REGION" envDefault:"local"`
	EndpointURL string `env:"AWS_ENDPOINT_URL" validate:"required,url"`
}

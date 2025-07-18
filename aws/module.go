package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"go.uber.org/fx"

	"github.com/tesserical/enclave/aws/internal/awsconfig"
	"github.com/tesserical/enclave/internal/osenv"
)

var module = fx.Module("enclave/aws",
	fx.Provide(
		osenv.ParseAs[awsconfig.Config],
		func(cfg awsconfig.Config) (aws.Config, error) {
			ctx := context.Background()
			if cfg.Region != "local" {
				return config.LoadDefaultConfig(ctx)
			}
			return config.LoadDefaultConfig(ctx,
				config.WithRegion("us-east-1"),
				config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("fakeAccess", "fakeSecret", "")),
			)
		},
	),
)

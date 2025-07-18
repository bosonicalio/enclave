package s3

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/samber/lo"
	"go.uber.org/fx"

	"github.com/tesserical/enclave/aws/internal/awsconfig"
)

var module = fx.Module("enclave/aws",
	fx.Provide(
		func(baseCfg awsconfig.Config, awsCfg aws.Config) *s3.Client {
			return s3.NewFromConfig(awsCfg, func(options *s3.Options) {
				if baseCfg.Region != "local" {
					return
				}
				options.UsePathStyle = true
				options.BaseEndpoint = lo.EmptyableToPtr(baseCfg.EndpointURL)
			})
		},
	),
)

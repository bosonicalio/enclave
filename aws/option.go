package aws

import (
	"github.com/tesserical/enclave"
)

// WithAmazonWebServices returns an enclave.Option that includes the AWS module.
func WithAmazonWebServices() enclave.Option {
	return enclave.WithFxOptions(
		module,
	)
}

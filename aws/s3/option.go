package s3

import (
	"github.com/tesserical/enclave"
)

// WithS3 returns an enclave.Option that includes the S3 module.
func WithS3() enclave.Option {
	return enclave.WithFxOptions(
		module,
	)
}

package kafka

import (
	"github.com/hadroncorp/geck/eventfx"

	"github.com/hadroncorp/enclave"
	"github.com/hadroncorp/enclave/kafka/kafkafx"
)

// WithKafkaEvents adds the Kafka module to the enclave application for event-driven systems.
func WithKafkaEvents() enclave.Option {
	return enclave.WithFxOptions(
		kafkafx.Module,
		eventfx.PublisherModule,
	)
}

// WithKafkaTxEvents adds the Kafka module to the enclave application for event-driven systems.
//
// It uses the Kafka Transaction API to ensure that events are published in a transactional manner.
func WithKafkaTxEvents() enclave.Option {
	return enclave.WithFxOptions(
		kafkafx.TxModule,
		eventfx.PublisherModule,
	)
}

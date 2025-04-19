package kafkafx

import (
	"go.uber.org/fx"
)

// Module is the default [go.uber.org/fx] module for Apache Kafka.
var Module = fx.Module("enclave/kafka",
	syncWriterProviders,
	readerProviders,
	fx.Invoke(
		fx.Annotate(
			registerManagerReaders,
			fx.ParamTags("", "", `group:"kafka_controllers"`),
		),
		startManagerReaders,
	),
)

// TxModule is a [go.uber.org/fx] module for Apache Kafka which offers transactional
// capabilities using Kafka Transaction API.
var TxModule = fx.Module("enclave/kafka-tx",
	txWriterProviders,
	txReaderProviders,
	fx.Invoke(
		fx.Annotate(
			registerManagerReaders,
			fx.ParamTags("", "", `group:"kafka_controllers"`),
		),
		startManagerReaders,
	),
)

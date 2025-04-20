package kafkafx

import (
	"github.com/hadroncorp/geck/configuration"
	"github.com/hadroncorp/geck/transport/stream"
	"github.com/hadroncorp/geck/transport/stream/kafka"
	"go.uber.org/fx"
)

var syncWriterProviders = fx.Provide(
	configuration.Parse[writeClientConfig],
	fx.Annotate(
		newWriterClient,
	),
	fx.Annotate(
		kafka.NewSyncWriter,
		fx.As(new(stream.Writer)),
	),
)

var readerProviders = fx.Provide(
	configuration.Parse[readerManagerConfig],
	fx.Annotate(
		newReaderManager,
		fx.As(new(kafka.ReaderManager)),
	),
)

// -- Transactional --

var txWriterProviders = fx.Provide(
	configuration.Parse[txWriteClientConfig],
	fx.Annotate(
		newTxWriteClient,
	),
	fx.Annotate(
		kafka.NewTransactionalWriter,
		fx.As(new(stream.Writer)),
	),
)

var txReaderProviders = fx.Provide(
	configuration.Parse[readerManagerConfig],
	fx.Annotate(
		newTxReaderManager,
		fx.As(new(kafka.ReaderManager)),
	),
)

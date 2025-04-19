package kafkafx

import (
	"context"
	"log/slog"
	"time"

	"github.com/hadroncorp/geck/application"
	"github.com/hadroncorp/geck/transport/stream/kafka"
	"github.com/samber/lo"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/fx"
)

// -- Writer --

type writeClientConfig struct {
	Config
	WriterConfig
}

func newWriterClient(lc fx.Lifecycle, config writeClientConfig, appConfig application.Config) (*kgo.Client, error) {
	opts := newFranzOpts(config.Config, appConfig)
	opts = append(opts, newFranzWriterOpts(config.WriterConfig)...)
	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			client.Close()
			return nil
		},
	})
	return client, nil
}

// -- Reader --

type readerManagerConfig struct {
	Config
	ConsumerConfig
	ConsumerGroupConfig
	PollBatchSize  int           `env:"KAFKA_POLL_BATCH_SIZE" envDefault:"100" validate:"min=1"`
	PollInterval   time.Duration `env:"KAFKA_POLL_INTERVAL" envDefault:"500ms" validate:"gte=0"`
	PoolSize       int           `env:"KAFKA_READER_POOL_SIZE" envDefault:"50" validate:"min=1"`
	HandlerTimeout time.Duration `env:"KAFKA_READER_HANDLER_TIMEOUT" envDefault:"30s" validate:"gte=0"`
}

func newReaderManager(logger *slog.Logger, config readerManagerConfig, appConfig application.Config) (*kafka.ChannelReaderManager, error) {
	opts := []kafka.ReaderManagerOption{
		kafka.WithReaderManagerPollBatchSize(config.PollBatchSize),
		kafka.WithReaderManagerPollInterval(config.PollInterval),
		kafka.WithReaderManagerPoolSize(config.PoolSize),
		kafka.WithReaderManagerHandlerTimeout(config.HandlerTimeout),
		kafka.WithReaderManagerErrorHandler(func(ctx context.Context, err error) {
			logger.ErrorContext(ctx, "got error from reader manager", slog.String("error", err.Error()))
		}),
	}

	config.ConsumerGroupID = lo.CoalesceOrEmpty(config.ConsumerGroupID, appConfig.Name)
	clientOpts := newFranzOpts(config.Config, appConfig)
	clientOpts = append(clientOpts, newFranzConsumerOpts(config.ConsumerConfig)...)
	clientOpts = append(clientOpts, newFranzConsumerGroupOpts(config.ConsumerGroupConfig)...)
	opts = append(opts, kafka.WithReaderManagerClientOpts(clientOpts...))
	return kafka.NewChannelReaderManager(opts...)
}

type txWriteClientConfig struct {
	Config
	WriterConfig
	TransactionalID string `env:"KAFKA_TRANSACTIONAL_ID"`
}

// -- Transactional --

func newTxWriteClient(lc fx.Lifecycle, config txWriteClientConfig, appConfig application.Config) (*kgo.Client, error) {
	txID := lo.CoalesceOrEmpty(config.TransactionalID, appConfig.InstanceName())
	opts := newFranzOpts(config.Config, appConfig)
	opts = append(opts, newFranzWriterOpts(config.WriterConfig)...)
	opts = append(opts, kgo.TransactionalID(txID))
	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			client.Close()
			return nil
		},
	})
	return client, nil
}

func newTxReaderManager(logger *slog.Logger, config readerManagerConfig, appConfig application.Config) (*kafka.ChannelReaderManager, error) {
	config.ConsumerConfig.FetchIsolationLevel = 1
	return newReaderManager(logger, config, appConfig)
}

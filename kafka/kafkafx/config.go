package kafkafx

import (
	"time"

	"github.com/hadroncorp/geck/application"
	"github.com/twmb/franz-go/pkg/kgo"
)

// -- Global --

type Config struct {
	Brokers             []string `env:"KAFKA_BROKERS" validate:"required,dive,required"`
	ClientID            string   `env:"KAFKA_CLIENT_ID"`
	InstanceID          string   `env:"KAFKA_INSTANCE_ID"`
	EnableTopicCreation bool     `env:"KAFKA_ALLOW_AUTO_TOPIC_CREATION" envDefault:"false"`
}

// newFranzOpts creates a new set of options for the Kafka client.
//
// This is generic and can be used for both readers and writers.
func newFranzOpts(config Config, appConfig application.Config) []kgo.Opt {
	if config.ClientID == "" {
		config.ClientID = appConfig.Name
	}
	if config.InstanceID == "" {
		config.InstanceID = appConfig.InstanceName()
	}

	const totalConfigs = 4
	opts := make([]kgo.Opt, 0, totalConfigs)
	opts = append(opts, kgo.SeedBrokers(config.Brokers...))
	opts = append(opts, kgo.ClientID(config.ClientID))
	opts = append(opts, kgo.InstanceID(config.InstanceID))
	if config.EnableTopicCreation {
		opts = append(opts, kgo.AllowAutoTopicCreation())
	}
	return opts
}

// -- Writer --

type WriterConfig struct {
	RequiredAcks      int16    `env:"KAFKA_REQUIRED_ACKS" envDefault:"-1" validate:"min=-1,max=1"`
	BatchCompressions []string `env:"KAFKA_BATCH_COMPRESSIONS" validate:"omitempty,dive,oneof=gzip snappy lz4 zstd none"`
}

func newFranzWriterOpts(config WriterConfig) []kgo.Opt {
	opts := make([]kgo.Opt, 0, 2)
	opts = append(opts, kgo.RequiredAcks(convertAcks(config.RequiredAcks)))
	if len(config.BatchCompressions) > 0 {
		opts = append(opts, kgo.ProducerBatchCompression(convertCompressionTypes(config.BatchCompressions)...))
	}
	return opts
}

// --- Utils ---

func convertAcks(v int16) kgo.Acks {
	switch v {
	case 0:
		return kgo.NoAck()
	case -1:
		return kgo.AllISRAcks()
	default:
		return kgo.LeaderAck()
	}
}

func convertCompressionTypes(compressionTypes []string) []kgo.CompressionCodec {
	compressions := make([]kgo.CompressionCodec, 0, len(compressionTypes))
	for _, compressionType := range compressionTypes {
		switch compressionType {
		case "gzip":
			compressions = append(compressions, kgo.GzipCompression())
		case "snappy":
			compressions = append(compressions, kgo.SnappyCompression())
		case "lz4":
			compressions = append(compressions, kgo.Lz4Compression())
		case "zstd":
			compressions = append(compressions, kgo.ZstdCompression())
		default:
			compressions = append(compressions, kgo.NoCompression())
		}
	}
	return compressions
}

// -- Consumer --

type ConsumerConfig struct {
	ResetOffset         string `env:"KAFKA_RESET_OFFSET" envDefault:"committed" validate:"required,oneof=earliest latest committed"`
	FetchIsolationLevel int8   `env:"KAFKA_FETCH_ISOLATION_LEVEL" envDefault:"0" validate:"min=0,max=1"`
}

func newFranzConsumerOpts(config ConsumerConfig) []kgo.Opt {
	opts := make([]kgo.Opt, 0, 2)
	opts = append(opts, kgo.FetchIsolationLevel(convertIsolationLevel(config.FetchIsolationLevel)))
	opts = append(opts, kgo.ConsumeResetOffset(convertResetOffset(config.ResetOffset)))
	return opts
}

// --- Utils ---

// convertResetOffset converts the reset offset string to a kgo.Offset type.
func convertResetOffset(resetOffset string) kgo.Offset {
	switch resetOffset {
	case "committed":
		return kgo.NewOffset().AtCommitted()
	case "latest":
		return kgo.NewOffset().AtEnd()
	default:
		return kgo.NewOffset().AtStart()
	}
}

func convertIsolationLevel(isolationLevel int8) kgo.IsolationLevel {
	switch isolationLevel {
	case 1:
		return kgo.ReadCommitted()
	default:
		return kgo.ReadUncommitted()
	}
}

// -- Consumer Group --

type ConsumerGroupConfig struct {
	ConsumerGroupID    string        `env:"KAFKA_CONSUMER_GROUP_ID" validate:"omitempty,lte=255"`
	DisableAutocommit  bool          `env:"KAFKA_DISABLE_AUTOCOMMIT" envDefault:"true"`
	AutocommitInterval time.Duration `env:"KAFKA_AUTOCOMMIT_INTERVAL" envDefault:"5s" validate:"gte=0"`
}

func newFranzConsumerGroupOpts(config ConsumerGroupConfig) []kgo.Opt {
	opts := make([]kgo.Opt, 0, 3)
	if config.ConsumerGroupID != "" {
		opts = append(opts, kgo.ConsumerGroup(config.ConsumerGroupID))
	}
	if config.DisableAutocommit {
		opts = append(opts, kgo.DisableAutoCommit())
	}
	opts = append(opts, kgo.AutoCommitInterval(config.AutocommitInterval))
	return opts
}

package sqlfx

import (
	"log/slog"
	"time"
)

type Config struct {
	EnableLogging      bool       `env:"SQL_ENABLE_LOGGING" envDefault:"true"`
	LogLevel           slog.Level `env:"SQL_LOG_LEVEL" envDefault:"DEBUG"`
	EnableTxContext    bool       `env:"SQL_ENABLE_TX_CONTEXT" envDefault:"false"`
	EnableTxAutoCreate bool       `env:"SQL_ENABLE_TX_AUTO_CREATE" envDefault:"false"`

	ConnectionString   string        `env:"SQL_CONNECTION_STRING" validate:"required"`
	MaxConnections     int32         `env:"SQL_MAX_CONNECTIONS" validate:"omitempty,gte=0"`
	MinConnections     int32         `env:"SQL_MIN_CONNECTIONS" validate:"omitempty,gte=0"`
	MinIdleConnections int32         `env:"SQL_MIN_IDLE_CONNECTIONS" validate:"omitempty,gte=0"`
	MaxConnLifetime    time.Duration `env:"SQL_MAX_CONN_LIFETIME" validate:"omitempty,gte=0"`
	MaxConnIdleTime    time.Duration `env:"SQL_MAX_CONN_IDLE_TIME" validate:"omitempty,gte=0"`
	HealthCheckPeriod  time.Duration `env:"SQL_HEALTH_CHECK_PERIOD" validate:"omitempty,gte=0"`
}

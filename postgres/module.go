package postgres

import (
	"context"
	"database/sql"

	"github.com/tesserical/geck/persistence/postgres"
	"go.uber.org/fx"

	"github.com/tesserical/enclave/internal/persistencefx/sqlfx"
)

var module = fx.Options(
	fx.Provide(
		newDB,
	),
	fx.Invoke(
		logDBInfo,
	),
)

// -- Factory --

func newDB(lc fx.Lifecycle, config sqlfx.Config) (*sql.DB, error) {
	opts := make([]postgres.ConnectionPoolOption, 0, 6)
	if config.MaxConnections > 0 {
		opts = append(opts, postgres.WithMaxConnections(config.MaxConnections))
	}
	if config.MinConnections > 0 {
		opts = append(opts, postgres.WithMinConnections(config.MinConnections))
	}
	if config.MaxConnLifetime > 0 {
		opts = append(opts, postgres.WithMaxConnLifetime(config.MaxConnLifetime))
	}
	if config.MaxConnIdleTime > 0 {
		opts = append(opts, postgres.WithMaxConnIdleTime(config.MaxConnIdleTime))
	}
	if config.MinIdleConnections > 0 {
		opts = append(opts, postgres.WithMinIdleConnections(config.MinIdleConnections))
	}
	if config.HealthCheckPeriod > 0 {
		opts = append(opts, postgres.WithHealthCheckPeriod(config.HealthCheckPeriod))
	}

	dbPool, err := postgres.NewConnectionPool(context.Background(), config.ConnectionString, opts...)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			return dbPool.Close()
		},
	})
	return dbPool, nil
}

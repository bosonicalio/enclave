package postgres

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"

	"github.com/bosonicalio/enclave/internal/persistencefx/sqlfx"
)

func logDBInfo(logger *slog.Logger, config sqlfx.Config) {
	if !config.EnableLogging || logger == nil {
		return
	}

	dbPool, err := pgxpool.ParseConfig(config.ConnectionString)
	if err != nil {
		return
	}

	dbConfig, err := pgx.ParseConfig(config.ConnectionString)
	if err != nil {
		return
	}

	logger.Log(context.Background(), config.LogLevel, "reading db_info",
		slog.Group("pool",
			slog.Int("max_conns", int(lo.CoalesceOrEmpty(config.MaxConnections, dbPool.MaxConns))),
			slog.Int("min_conns", int(lo.CoalesceOrEmpty(config.MinConnections, dbPool.MinConns))),
			slog.Duration("max_conn_lifetime_jitter", dbPool.MaxConnLifetimeJitter),
			slog.Duration("max_conn_lifetime", lo.CoalesceOrEmpty(config.MaxConnLifetime, dbPool.MaxConnLifetime)),
			slog.Duration("max_conn_idle_time", lo.CoalesceOrEmpty(config.MaxConnIdleTime, dbPool.MaxConnIdleTime)),
			slog.Int("min_idle_conns", int(lo.CoalesceOrEmpty(config.MinConnections, dbPool.MinIdleConns))),
			slog.Duration("health_check_period", lo.CoalesceOrEmpty(config.HealthCheckPeriod, dbPool.HealthCheckPeriod)),
		),
		slog.Group("db",
			slog.String("type", "postgres"),
			slog.String("database", dbConfig.Database),
			slog.String("host", dbConfig.Host),
			slog.String("user", dbConfig.User),
		),
	)
}

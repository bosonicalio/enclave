package sqlfx

import (
	"database/sql"
	"log/slog"

	gecksql "github.com/tesserical/geck/persistence/sql"
	"go.uber.org/fx"

	"github.com/tesserical/enclave/internal/osenv"
)

// NOTE: This might need to change in the future to support more complex configurations,
// like having both master and read database connections.

// Module provides the SQL database connection and configuration for the enclave persistence layer.
// It uses the [gecksql] package to manage SQL database interactions, with optional logging and transaction
// context propagation.
// The module is configured via environment variables, allowing for flexible deployment configurations.
var Module = fx.Module("enclave/persistence/sql",
	fx.Provide(
		osenv.ParseAs[Config],
		fx.Annotate(
			newDB,
			fx.ParamTags("", "", `optional:"true"`), // logger is optional
		),
		newTxFactory,
	),
	fx.Invoke(
		registerTxFactory,
	),
)

// -- Factory --

func newDB(cfg Config, db *sql.DB, logger *slog.Logger) gecksql.DB {
	if !cfg.EnableLogging && !cfg.EnableTxContext {
		return db
	}
	// NOTE: Order of decorators matters. Last is the first to be applied.
	var aggregateDB = gecksql.DB(db)
	if cfg.EnableTxContext {
		aggregateDB = gecksql.NewDBTxPropagator(aggregateDB)
	}
	if cfg.EnableLogging && logger != nil {
		aggregateDB = gecksql.NewDBLogger(aggregateDB, logger,
			gecksql.WithLogLevel(cfg.LogLevel),
		)
	}
	return aggregateDB
}

func newTxFactory(db gecksql.DB, cfg Config) gecksql.TxFactory {
	return gecksql.NewTxFactory(db, &sql.TxOptions{
		Isolation: cfg.TxContextIsolationLevel,
		ReadOnly:  cfg.TxContextReadOnly,
	})
}

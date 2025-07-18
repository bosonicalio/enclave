package sqlfx

import (
	"database/sql"

	"github.com/tesserical/geck/persistence"
	gecksql "github.com/tesserical/geck/persistence/sql"
)

// Registers a SQL transaction factory to the geck's Persistence API global transaction factory registry
// so it can be used by the geck's Persistence API to create transactions with the provided database connection
// and configuration options.
func registerTxFactory(db gecksql.DB, cfg Config) {
	sqlTxFactory := gecksql.NewTxFactory(db, &sql.TxOptions{
		Isolation: cfg.TxContextIsolationLevel,
		ReadOnly:  cfg.TxContextReadOnly,
	})
	persistence.RegisterTxFactory(sqlTxFactory)
}

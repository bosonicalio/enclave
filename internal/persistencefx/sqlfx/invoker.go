package sqlfx

import (
	"github.com/tesserical/geck/persistence"
	gecksql "github.com/tesserical/geck/persistence/sql"
)

func registerTxFactory(txManager *persistence.TxManager, factory gecksql.TxFactory) {
	txManager.Register(factory)
}

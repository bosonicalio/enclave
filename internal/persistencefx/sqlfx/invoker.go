package sqlfx

import (
	"github.com/bosonicalio/geck/persistence"
	gecksql "github.com/bosonicalio/geck/persistence/sql"
)

func registerTxFactory(txManager *persistence.TxManager, factory gecksql.TxFactory) {
	txManager.Register(factory)
}

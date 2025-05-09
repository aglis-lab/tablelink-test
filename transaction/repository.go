package transaction

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

//
// Gorm Transaction Repository //
//

type GormTransactionRepository interface {
	WithGormTransaction(ctx context.Context, callback func(tx GormContext) error, opts ...*sql.TxOptions) error
}

type gormTransactionRepository struct {
	GormTransaction
}

func NewGormTransactionRepository(db *gorm.DB) GormTransactionRepository {
	return &gormTransactionRepository{
		GormTransaction: GormTransaction{
			DB: db,
		},
	}
}

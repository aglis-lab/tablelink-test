package transaction

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

type GormContext struct {
	context.Context
	db *gorm.DB
}

func UnwrapContext(ctx context.Context, defaultGorm *gorm.DB) *gorm.DB {
	if trx, ok := ctx.(GormContext); ok {
		return trx.db
	}

	return defaultGorm.WithContext(ctx)
}

type GormTransaction struct {
	*gorm.DB
}

func (trx GormTransaction) WithGormTransaction(ctx context.Context, callback func(tx GormContext) error, opts ...*sql.TxOptions) error {
	if tx, ok := ctx.(GormContext); ok {
		return callback(tx)
	}

	return trx.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return callback(GormContext{Context: ctx, db: tx})
	})
}

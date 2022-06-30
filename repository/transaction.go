package repository

import (
	"context"
	"database/sql"
	"github.com/MinterTeam/explorer-sdk/helpers"
	"github.com/MinterTeam/explorer-sdk/models"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

func (rTransaction *TransactionRepository) GetByHash(hash string) (*models.Transaction, error) {
	tx := new(models.Transaction)
	err := rTransaction.db.
		NewSelect().
		Model(tx).
		Where("hash = ?", helpers.RemovePrefix(hash)).
		Scan(rTransaction.ctx)
	if err != nil {
		return nil, err
	}
	return tx, err
}

type TransactionRepository struct {
	db  *bun.DB
	ctx context.Context
}

func NewTransactionRepository(sqlDB *sql.DB, dialect *pgdialect.Dialect) *TransactionRepository {
	db := bun.NewDB(sqlDB, dialect)
	db.RegisterModel((*models.TransactionValidator)(nil))
	return &TransactionRepository{
		db:  db,
		ctx: context.Background(),
	}
}

package repository

import (
	"context"
	"github.com/MinterTeam/explorer-sdk/helpers"
	"github.com/MinterTeam/explorer-sdk/models"
	"github.com/uptrace/bun"
	"strings"
)

func (rTransaction *TransactionRepository) SaveInvalid(tx *models.InvalidTransaction) (*models.InvalidTransaction, error) {
	_, err := rTransaction.db.NewInsert().
		Model(tx).
		Exec(rTransaction.ctx)
	return tx, err
}

func (rTransaction *TransactionRepository) SaveAll(list []models.Transaction) ([]models.Transaction, error) {
	txs := list
	_, err := rTransaction.db.NewInsert().
		Model(&txs).
		Exec(rTransaction.ctx)
	return txs, err
}

func (rTransaction *TransactionRepository) Save(tx *models.Transaction) (*models.Transaction, error) {
	_, err := rTransaction.db.NewInsert().
		Model(tx).
		Exec(rTransaction.ctx)
	return tx, err
}

func (rTransaction *TransactionRepository) GetByHash(hash string) (*models.Transaction, error) {
	tx := new(models.Transaction)
	err := rTransaction.db.
		NewSelect().
		Model(tx).
		Where("hash = ?", strings.ToLower(helpers.RemovePrefix(hash))).
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

func NewTransactionRepository(db *bun.DB) *TransactionRepository {
	//db := bun.NewDB(sqlDB, dialect)
	//db.RegisterModel((*models.TransactionValidator)(nil))
	return &TransactionRepository{
		db:  db,
		ctx: context.Background(),
	}
}

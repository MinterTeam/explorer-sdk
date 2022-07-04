package repository

import (
	"context"
	"github.com/MinterTeam/explorer-sdk/models"
	"github.com/uptrace/bun"
)

func (rBalance *BalanceRepository) SaveAll(balances []models.Balance) ([]models.Balance, error) {
	list := balances

	_, err := rBalance.db.
		NewInsert().
		Model(&list).
		On("CONFLICT (address_id, coin_id) DO UPDATE").
		Exec(rBalance.ctx)

	if err != nil {
		return nil, err
	}

	return list, err
}

type BalanceRepository struct {
	db  *bun.DB
	ctx context.Context
}

func NewBalanceRepository(db *bun.DB) *BalanceRepository {
	return &BalanceRepository{
		db:  db,
		ctx: context.Background(),
	}
}

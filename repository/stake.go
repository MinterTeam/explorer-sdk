package repository

import (
	"context"
	"database/sql"
	"github.com/MinterTeam/explorer-sdk/models"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

func (rStake *StakeRepository) InsertOrUpdateStakes(stakes []models.Stake) ([]models.Stake, error) {
	list := stakes
	_, err := rStake.db.NewInsert().
		Model(&list).
		On("CONFLICT (owner_address_id, validator_id, coin_id, is_kicked) DO UPDATE").
		Set("value = EXCLUDED.value, bip_value = EXCLUDED.bip_value").
		Exec(context.Background())
	return list, err
}

type StakeRepository struct {
	db *bun.DB
}

func NewStakeRepository(sqldb *sql.DB, dialect *pgdialect.Dialect) *StakeRepository {
	return &StakeRepository{
		db: bun.NewDB(sqldb, dialect),
	}
}

package repository

import (
	"context"
	"database/sql"
	"github.com/MinterTeam/explorer-sdk/models"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

func (rStake *StakeRepository) GetDelegatorsCount() (uint64, error) {
	var count uint64
	err := rStake.db.
		NewSelect().
		Model((*models.Stake)(nil)).
		ColumnExpr("count (DISTINCT owner_address_id)").
		Scan(context.Background(), &count)
	return count, err
}

func (rStake *StakeRepository) DeleteStakesNotInListIds(idList []uint64) error {
	_, err := rStake.db.
		NewDelete().
		Model((*models.Stake)(nil)).
		Where("id not in (?)", bun.In(idList)).
		Where("is_kicked != true", bun.In(idList)).
		Exec(context.Background())
	return err
}

func (rStake *StakeRepository) DeleteStakesByValidatorIds(idList []uint64) error {
	_, err := rStake.db.
		NewDelete().
		Model((*models.Stake)(nil)).
		Where("validator_id in (?)", bun.In(idList)).
		Exec(context.Background())
	return err
}

func (rStake *StakeRepository) SaveAllStakes(stakes []*models.Stake) error {
	_, err := rStake.db.NewInsert().
		Model(&stakes).
		On("CONFLICT (owner_address_id, validator_id, coin_id) DO UPDATE").
		Exec(context.Background())
	return err
}

type StakeRepository struct {
	db *bun.DB
}

func NewStakeRepository(sqldb *sql.DB, dialect *pgdialect.Dialect) *StakeRepository {
	return &StakeRepository{
		db: bun.NewDB(sqldb, dialect),
	}
}

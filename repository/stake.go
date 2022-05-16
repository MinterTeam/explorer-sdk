package repository

import (
	"context"
	"database/sql"
	"github.com/MinterTeam/explorer-sdk/models"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"io"
)

func (rStake *StakeRepository) CopyFromReader(data io.Reader) error {
	conn, err := rStake.db.Conn(rStake.ctx)
	if err != nil {
		return err
	}
	_, err = pgdriver.CopyFrom(rStake.ctx, conn, data, "COPY stakes FROM STDIN")
	return err
}

func (rStake *StakeRepository) TruncateStakeTable() error {
	_, err := rStake.db.NewTruncateTable().
		Model((*models.Stake)(nil)).
		Restrict().
		Exec(rStake.ctx)
	return err
}

func (rStake *StakeRepository) InsertOrUpdateStakes(stakes []models.Stake) ([]models.Stake, error) {
	list := stakes
	_, err := rStake.db.NewInsert().
		Model(&list).
		On("CONFLICT (owner_address_id, validator_id, coin_id, is_kicked) DO UPDATE").
		Set("value = EXCLUDED.value, bip_value = EXCLUDED.bip_value").
		Exec(rStake.ctx)
	return list, err
}

type StakeRepository struct {
	db  *bun.DB
	ctx context.Context
}

func NewStakeRepository(sqldb *sql.DB, dialect *pgdialect.Dialect) *StakeRepository {
	return &StakeRepository{
		db:  bun.NewDB(sqldb, dialect),
		ctx: context.Background(),
	}
}

package repository

import (
	"context"
	"database/sql"
	"github.com/MinterTeam/explorer-sdk/models"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"time"
)

func (rReward *RewardRepository) GetRewardsByDay(date time.Time) ([]*models.AggregatedReward, error) {
	var list []*models.AggregatedReward
	err := rReward.db.
		NewSelect().
		Model(&list).
		Where("time_id = ?", date.Format("2006-01-02")).
		Scan(context.Background())
	return list, err
}

func (rReward *RewardRepository) SaveAggregatedRewards(rewards []*models.AggregatedReward) ([]*models.AggregatedReward, error) {
	list := rewards
	_, err := rReward.db.NewInsert().
		Model(&list).
		On("CONFLICT (time_id, address_id, validator_id, role) DO UPDATE").
		Exec(context.Background())
	return list, err
}

type RewardRepository struct {
	db *bun.DB
}

func NewRewardRepository(sqldb *sql.DB, dialect *pgdialect.Dialect) *RewardRepository {
	return &RewardRepository{
		db: bun.NewDB(sqldb, dialect),
	}
}

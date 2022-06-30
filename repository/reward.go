package repository

import (
	"context"
	"github.com/MinterTeam/explorer-sdk/models"
	"github.com/uptrace/bun"
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

func NewRewardRepository(db *bun.DB) *RewardRepository {
	return &RewardRepository{
		db: db,
	}
}

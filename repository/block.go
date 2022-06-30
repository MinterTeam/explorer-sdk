package repository

import (
	"context"
	"github.com/MinterTeam/explorer-sdk/models"
	"github.com/uptrace/bun"
)

func (rBlock *BlockRepository) Save(block *models.Block) error {
	_, err := rBlock.db.NewInsert().
		Model(block).
		Exec(rBlock.ctx)
	return err
}

func (rBlock *BlockRepository) GetLastFromDB() (*models.Block, error) {
	block := new(models.Block)
	err := rBlock.db.
		NewSelect().
		Model(block).
		Order("id desc").
		Limit(1).
		Scan(rBlock.ctx)
	if err != nil {
		return nil, err
	}
	return block, err
}

func (rBlock *BlockRepository) GetById(id uint64) (*models.Block, error) {
	block := new(models.Block)
	err := rBlock.db.
		NewSelect().
		Model(block).
		Where("id = ?", id).
		Scan(rBlock.ctx)
	if err != nil {
		return nil, err
	}
	return block, err
}

func (rBlock *BlockRepository) LinkWithValidators(links []*models.BlockValidator) error {
	list := links
	_, err := rBlock.db.NewInsert().
		Model(&list).
		Exec(rBlock.ctx)
	return err
}

type BlockRepository struct {
	db  *bun.DB
	ctx context.Context
}

func NewBlockRepository(db *bun.DB) *BlockRepository {
	//db := bun.NewDB(sqlDB, dialect)
	//db.RegisterModel((*models.BlockValidator)(nil))
	return &BlockRepository{
		db:  db,
		ctx: context.Background(),
	}
}

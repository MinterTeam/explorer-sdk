package repository

import (
	"context"
	"github.com/MinterTeam/explorer-sdk/models"
	"github.com/uptrace/bun"
	"sync"
)

func (rValidator *ValidatorRepository) FindByPk(pk string) (models.Validator, error) {

	v, ok := rValidator.pkCache.Load(pk)
	if ok {
		return v.(models.Validator), nil
	}

	var vpk models.ValidatorPublicKeys
	err := rValidator.db.
		NewSelect().
		Model(&vpk).
		Relation("Validator").
		Where("key = ?", pk).
		Scan(context.Background())

	if err != nil {
		return models.Validator{}, err
	}

	rValidator.pkCache.Store(pk, *vpk.Validator)
	return *vpk.Validator, nil
}

func (rValidator *ValidatorRepository) UpdateCache() error {
	var list []models.ValidatorPublicKeys
	err := rValidator.db.
		NewSelect().
		Model(&list).
		Relation("Validator").
		Scan(context.Background())

	if err != nil {
		return err
	}
	rValidator.pkCache = new(sync.Map)
	for _, v := range list {
		if v.Validator != nil {
			rValidator.pkCache.Store(v.Key, *v.Validator)
		}
	}
	return nil
}

func (rValidator *ValidatorRepository) Init() {
	err := rValidator.UpdateCache()
	if err != nil {
		println(err)
	}
}

type ValidatorRepository struct {
	pkCache *sync.Map
	db      *bun.DB
}

func NewValidatorRepository(db *bun.DB) *ValidatorRepository {
	return &ValidatorRepository{
		pkCache: new(sync.Map),
		db:      db,
	}
}

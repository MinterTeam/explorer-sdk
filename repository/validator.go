package repository

import (
	"context"
	"database/sql"
	"github.com/MinterTeam/explorer-sdk/models"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
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

type ValidatorRepository struct {
	pkCache *sync.Map
	db      *bun.DB
}

func NewValidatorRepository(sqldb *sql.DB, dialect *pgdialect.Dialect) *ValidatorRepository {
	return &ValidatorRepository{
		pkCache: new(sync.Map),
		db:      bun.NewDB(sqldb, dialect),
	}
}

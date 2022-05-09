package repository

import (
	"context"
	"database/sql"
	"github.com/MinterTeam/explorer-sdk/models"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"sync"
)

func (rAddress *AddressRepository) GetByAddress(addressString string) (models.Address, error) {
	a, ok := rAddress.addressCache.Load(addressString)
	if ok {
		return a.(models.Address), nil
	}

	var address models.Address
	err := rAddress.db.
		NewSelect().
		Model(&address).
		Where("address = ?", addressString).
		Scan(context.Background())

	if err != nil {
		return models.Address{}, err
	}

	rAddress.idCache.Store(address.ID, address)
	rAddress.addressCache.Store(addressString, address)
	return address, nil
}

func (rAddress *AddressRepository) GetById(id uint) (models.Address, error) {
	a, ok := rAddress.idCache.Load(id)
	if ok {
		return a.(models.Address), nil
	}

	var address models.Address
	err := rAddress.db.
		NewSelect().
		Model(&address).
		Where("id = ?", id).
		Scan(context.Background())

	if err != nil {
		return models.Address{}, err
	}

	rAddress.idCache.Store(id, address)
	rAddress.addressCache.Store(address.Address, address)
	return address, nil
}

type AddressRepository struct {
	db           *bun.DB
	idCache      *sync.Map
	addressCache *sync.Map
}

func NewAddressRepository(sqldb *sql.DB, dialect *pgdialect.Dialect) *AddressRepository {
	return &AddressRepository{
		idCache:      new(sync.Map),
		addressCache: new(sync.Map),
		db:           bun.NewDB(sqldb, dialect),
	}
}

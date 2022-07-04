package repository

import (
	"context"
	"database/sql"
	"github.com/MinterTeam/explorer-sdk/helpers"
	"github.com/MinterTeam/explorer-sdk/models"
	"github.com/uptrace/bun"
	"math"
	"strings"
	"sync"
)

func (rAddress *AddressRepository) GetWithoutBalance() ([]models.Address, error) {
	var addresses []models.Address
	err := rAddress.db.
		NewSelect().
		Model(&addresses).
		ColumnExpr("address.*").
		Join("LEFT JOIN balances as b ON address.id = b.address_id").
		Where("b.address_id IS NULL").
		Scan(rAddress.ctx)

	return addresses, err
}

func (rAddress *AddressRepository) SaveAll(addresses []string) ([]models.Address, error) {
	var list []models.Address
	for _, a := range addresses {
		_, ok := rAddress.addressCache.Load(helpers.RemovePrefix(strings.ToLower(a)))
		if ok {
			continue
		}
		list = append(list, models.Address{
			Address: helpers.RemovePrefix(strings.ToLower(a)),
		})
	}

	if len(list) == 0 {
		return nil, nil
	}

	_, err := rAddress.db.
		NewInsert().
		Model(&list).
		On("CONFLICT (address) DO NOTHING").
		Exec(rAddress.ctx)

	if err != nil {
		return nil, err
	}

	for _, a := range list {
		rAddress.idCache.Store(a.ID, a)
		rAddress.addressCache.Store(a.Address, a)
	}

	return list, err
}

func (rAddress *AddressRepository) GetByAddress(addressString string) (models.Address, error) {
	a, ok := rAddress.addressCache.Load(helpers.RemovePrefix(strings.ToLower(addressString)))
	if ok {
		return a.(models.Address), nil
	}

	var address models.Address
	err := rAddress.db.
		NewSelect().
		Model(&address).
		Where("address = ?", helpers.RemovePrefix(strings.ToLower(addressString))).
		Scan(rAddress.ctx)

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
		Scan(rAddress.ctx)

	if err != nil {
		return models.Address{}, err
	}

	rAddress.idCache.Store(id, address)
	rAddress.addressCache.Store(address.Address, address)
	return address, nil
}

func (rAddress *AddressRepository) UpdateCache() error {
	var addresses []models.Address
	err := rAddress.db.
		NewSelect().
		Model(&addresses).
		Scan(rAddress.ctx)

	if err != nil {
		return err
	}

	rAddress.idCache = new(sync.Map)
	rAddress.addressCache = new(sync.Map)

	chunkSize := 5000
	wg := new(sync.WaitGroup)

	chunksCount := int(math.Ceil(float64(len(addresses)) / float64(chunkSize)))
	wg.Add(chunksCount)

	for i := 0; i < chunksCount; i++ {
		start := chunkSize * i
		end := start + chunkSize
		if end > len(addresses) {
			end = len(addresses)
		}
		go func(list []models.Address) {
			for _, a := range list {
				rAddress.idCache.Store(a.ID, a)
				rAddress.addressCache.Store(a.Address, a)
			}
			wg.Done()
		}(addresses[start:end])
	}

	wg.Wait()

	return nil
}

func (rAddress *AddressRepository) FindIdOrCreate(addressString string) (models.Address, error) {
	a, err := rAddress.GetByAddress(helpers.RemovePrefix(strings.ToLower(addressString)))

	if err != nil && err != sql.ErrNoRows {
		return models.Address{}, err
	}

	if err != nil && err == sql.ErrNoRows {
		return rAddress.Create(addressString)
	}
	return a, err
}

func (rAddress *AddressRepository) Create(a string) (models.Address, error) {
	address := models.Address{
		Address: helpers.RemovePrefix(strings.ToLower(a)),
	}
	_, err := rAddress.db.NewInsert().
		Model(&address).
		Exec(rAddress.ctx)

	if err != nil {
		return models.Address{}, err
	}

	rAddress.idCache.Store(address.ID, address)
	rAddress.addressCache.Store(address.Address, address)
	return address, nil
}

func (rAddress *AddressRepository) Init() {
	err := rAddress.UpdateCache()
	if err != nil {
		println(err)
	}
}

type AddressRepository struct {
	db           *bun.DB
	idCache      *sync.Map
	addressCache *sync.Map
	ctx          context.Context
}

func NewAddressRepository(db *bun.DB) *AddressRepository {
	return &AddressRepository{
		idCache:      new(sync.Map),
		addressCache: new(sync.Map),
		db:           db,
		ctx:          context.Background(),
	}
}

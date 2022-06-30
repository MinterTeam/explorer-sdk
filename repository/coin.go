package repository

import (
	"context"
	"github.com/MinterTeam/explorer-sdk/models"
	"github.com/uptrace/bun"
	"math"
	"sync"
)

func (rCoin *CoinRepository) SaveTokenContracts(contracts []models.TokenContract) error {
	list := contracts
	_, err := rCoin.db.NewInsert().
		Model(&list).
		On("CONFLICT (coin_id) DO UPDATE").
		Exec(rCoin.ctx)
	return err
}

func (rCoin *CoinRepository) GetConfirmedCoinsId() ([]uint64, error) {
	var coins []models.TokenContract
	err := rCoin.db.
		NewSelect().
		Model(&coins).
		Scan(rCoin.ctx)

	if err != nil {
		return nil, err
	}
	var result []uint64
	for _, c := range coins {
		result = append(result, c.CoinId)
	}
	return result, err
}

func (rCoin *CoinRepository) UpdateOwnerBySymbol(symbol string, id uint) error {
	value := map[string]interface{}{
		"owner_address_id": id,
	}
	_, err := rCoin.db.NewUpdate().
		Model(&value).
		TableExpr("coins").
		Where("symbol = ?", symbol).
		Exec(rCoin.ctx)

	return err
}

func (rCoin *CoinRepository) UpdateAll(coins []*models.Coin) error {

	list := coins

	_, err := rCoin.db.NewUpdate().
		Model(&list).
		WherePK().
		Bulk().
		Exec(rCoin.ctx)

	for _, c := range list {
		rCoin.idCache.Store(c.ID, c)
		rCoin.symbolCache.Store(c.Symbol, c)
	}

	return err
}

func (rCoin *CoinRepository) Update(c *models.Coin) error {

	coin := c
	_, err := rCoin.db.NewUpdate().
		Model(coin).
		WherePK().
		Exec(rCoin.ctx)

	rCoin.idCache.Store(coin.ID, coin)
	rCoin.symbolCache.Store(coin.Symbol, coin)

	return err
}

func (rCoin *CoinRepository) SaveAllNewIfNotExist(coins []*models.Coin) error {
	list := coins
	_, err := rCoin.db.NewInsert().
		Model(&list).
		On("CONFLICT (symbol, version) DO UPDATE").
		Exec(rCoin.ctx)

	if err == nil {
		for _, c := range list {
			rCoin.idCache.Store(c.ID, c)
			rCoin.symbolCache.Store(c.Symbol, c)
		}
	}

	return err
}

func (rCoin *CoinRepository) Save(c *models.Coin) error {
	coin := c
	_, err := rCoin.db.NewInsert().
		Model(coin).
		On("CONFLICT (symbol, version) DO UPDATE").
		Exec(rCoin.ctx)

	if err == nil {
		rCoin.idCache.Store(coin.ID, coin)
		rCoin.symbolCache.Store(coin.Symbol, coin)
	}

	return err
}

func (rCoin *CoinRepository) Add(c *models.Coin) error {
	coin := c
	_, err := rCoin.db.NewInsert().
		Model(coin).
		Exec(rCoin.ctx)

	if err == nil {
		rCoin.idCache.Store(coin.ID, coin)
		rCoin.symbolCache.Store(coin.Symbol, coin)
	}

	return err
}

func (rCoin *CoinRepository) FindCoinIdBySymbol(symbol string) (uint, error) {
	c, err := rCoin.GetBySymbol(symbol)
	if err != nil {
		return 0, err
	}
	return c.ID, nil
}

func (rCoin *CoinRepository) FindSymbolById(id uint) (string, error) {
	c, err := rCoin.GetById(id)
	if err != nil {
		return "", err
	}
	return c.Symbol, nil
}

func (rCoin *CoinRepository) GetBySymbol(symbol string) (*models.Coin, error) {
	data, ok := rCoin.symbolCache.Load(symbol)
	if ok {
		return data.(*models.Coin), nil
	}

	coin := new(models.Coin)

	err := rCoin.db.
		NewSelect().
		Model(coin).
		Where("symbol = ?", symbol).
		Scan(rCoin.ctx)

	if err != nil {
		return nil, err
	}

	rCoin.idCache.Store(coin.ID, coin)
	rCoin.symbolCache.Store(coin.Symbol, coin)
	return coin, nil
}

func (rCoin *CoinRepository) GetById(id uint) (*models.Coin, error) {
	data, ok := rCoin.idCache.Load(id)
	if ok {
		return data.(*models.Coin), nil
	}

	coin := new(models.Coin)

	err := rCoin.db.
		NewSelect().
		Model(coin).
		Where("id = ?", id).
		Scan(rCoin.ctx)

	if err != nil {
		return nil, err
	}

	rCoin.idCache.Store(coin.ID, coin)
	rCoin.symbolCache.Store(coin.Symbol, coin)
	return coin, nil
}

func (rCoin *CoinRepository) GetCoinBySymbol(symbol string) ([]models.Coin, error) {
	var coins []models.Coin

	err := rCoin.db.
		NewSelect().
		Model(&coins).
		Where("symbol = ?", symbol).
		Scan(rCoin.ctx)

	if err != nil {
		return nil, err
	}

	return coins, err
}

func (rCoin *CoinRepository) GetAll() ([]models.Coin, error) {
	var coins []models.Coin
	err := rCoin.db.
		NewSelect().
		Model(&coins).
		Order("symbol ASC").
		Scan(rCoin.ctx)

	return coins, err
}

func (rCoin *CoinRepository) UpdateCache() error {

	list, err := rCoin.GetAll()

	if err != nil {
		return err
	}

	rCoin.idCache = new(sync.Map)
	rCoin.symbolCache = new(sync.Map)

	chunkSize := 1000
	wg := new(sync.WaitGroup)

	chunksCount := int(math.Ceil(float64(len(list)) / float64(chunkSize)))
	wg.Add(chunksCount)

	for i := 0; i < chunksCount; i++ {
		start := chunkSize * i
		end := start + chunkSize
		if end > len(list) {
			end = len(list)
		}
		go func(list []models.Coin) {
			for _, c := range list {
				rCoin.idCache.Store(c.ID, &c)
				rCoin.symbolCache.Store(c.Symbol, &c)
			}
			wg.Done()
		}(list[start:end])
	}

	wg.Wait()

	return nil
}

func (rCoin *CoinRepository) Init() {
	err := rCoin.UpdateCache()
	if err != nil {
		println(err)
	}
}

type CoinRepository struct {
	db          *bun.DB
	idCache     *sync.Map
	symbolCache *sync.Map
	ctx         context.Context
}

func NewCoinRepository(db *bun.DB) *CoinRepository {
	return &CoinRepository{
		db:          db,
		symbolCache: new(sync.Map), //TODO: добавить реализацию очистки
		idCache:     new(sync.Map), //TODO: добавить реализацию очистки
		ctx:         context.Background(),
	}
}

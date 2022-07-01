package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MinterTeam/explorer-sdk/models"
	"github.com/uptrace/bun"
	"strconv"
	"strings"
	"time"
)

func (rLp *LiquidityPoolRepository) GetPoolByPairString(pair string) (*models.LiquidityPool, error) {
	ids := strings.Split(pair, "-")
	firstCoinId, err := strconv.ParseUint(ids[0], 10, 64)
	if err != nil {
		return nil, err
	}
	secondCoinId, err := strconv.ParseUint(ids[1], 10, 64)
	if err != nil {
		return nil, err
	}
	if firstCoinId < secondCoinId {
		return rLp.GetLiquidityPoolByCoinIds(firstCoinId, secondCoinId)
	} else {
		return rLp.GetLiquidityPoolByCoinIds(secondCoinId, firstCoinId)
	}
}

func (rLp *LiquidityPoolRepository) GetPoolsByTxTags(tags map[string]string) ([]models.LiquidityPool, error) {
	pools, err := rLp.getPoolChainFromTags(tags)
	if err != nil {
		return nil, err
	}
	var idList []uint64
	for id := range pools {
		idList = append(idList, id)
	}
	return rLp.GetAllByIds(idList)
}

func (rLp *LiquidityPoolRepository) GetLastUpdateHeight() (uint64, error) {
	var lp = new(models.LiquidityPool)
	err := rLp.db.
		NewSelect().
		Model(lp).
		Order("updated_at_block_id desc").
		Limit(1).
		Scan(rLp.ctx)
	if err != nil {
		return 0, err
	}
	return lp.UpdatedAtBlockId, err
}

func (rLp *LiquidityPoolRepository) GetLastCreateHeight() (uint64, error) {
	var lp = new(models.LiquidityPool)
	err := rLp.db.
		NewSelect().
		Model(lp).
		Order("created_at_block_id desc").
		Limit(1).
		Scan(rLp.ctx)
	if err != nil {
		return 0, err
	}
	return lp.CreatedAtBlockId, err
}

func (rLp *LiquidityPoolRepository) GetLiquidityPoolByCoinIds(firstCoinId, secondCoinId uint64) (*models.LiquidityPool, error) {
	var lp = new(models.LiquidityPool)
	err := rLp.db.
		NewSelect().
		Model(lp).
		Where("first_coin_id = ? AND second_coin_id = ?", firstCoinId, secondCoinId).
		Scan(rLp.ctx)
	if err != nil {
		return nil, err
	}

	return lp, err
}

func (rLp *LiquidityPoolRepository) GetLiquidityPoolByTokenId(id uint64) (*models.LiquidityPool, error) {
	var lp = new(models.LiquidityPool)
	err := rLp.db.
		NewSelect().
		Model(lp).
		Where("token_id = ?", id).
		Scan(rLp.ctx)
	if err != nil {
		return nil, err
	}
	return lp, err
}

func (rLp *LiquidityPoolRepository) GetLiquidityPoolById(id uint64) (*models.LiquidityPool, error) {
	var lp = new(models.LiquidityPool)
	err := rLp.db.
		NewSelect().
		Model(lp).
		Where("id = ?", id).
		Scan(rLp.ctx)
	if err != nil {
		return nil, err
	}
	return lp, err
}

func (rLp *LiquidityPoolRepository) UpdateLiquidityPool(lp *models.LiquidityPool) error {
	_, err := rLp.db.NewInsert().Model(lp).On("CONFLICT (first_coin_id, second_coin_id) DO UPDATE").Exec(rLp.ctx)
	return err
}

func (rLp *LiquidityPoolRepository) UpdateLiquidityPoolById(lp *models.LiquidityPool) error {
	_, err := rLp.db.
		NewUpdate().
		Model(lp).
		Column("first_coin_volume").
		Column("second_coin_volume").
		Column("liquidity").
		Column("updated_at_block_id").
		WherePK().
		Exec(rLp.ctx)
	return err
}

func (rLp *LiquidityPoolRepository) LinkAddressLiquidityPool(addressId uint, liquidityPoolId uint64) error {
	addressLiquidityPool := &models.AddressLiquidityPool{
		LiquidityPoolId: liquidityPoolId,
		AddressId:       uint64(addressId),
	}
	_, err := rLp.db.NewInsert().
		Model(addressLiquidityPool).
		On("CONFLICT (address_id, liquidity_pool_id) DO NOTHING").
		Exec(rLp.ctx)
	return err
}

func (rLp *LiquidityPoolRepository) GetAddressLiquidityPool(addressId uint, liquidityPoolId uint64) (*models.AddressLiquidityPool, error) {
	var alp = new(models.AddressLiquidityPool)
	err := rLp.db.
		NewSelect().
		Model(alp).
		Where("address_id = ? AND liquidity_pool_id = ?", addressId, liquidityPoolId).
		Scan(rLp.ctx)
	if err != nil {
		return nil, err
	}
	return alp, err
}

func (rLp *LiquidityPoolRepository) GetAddressLiquidityPoolByCoinId(addressId uint, liquidityPoolId uint64) (*models.AddressLiquidityPool, error) {
	var alp = new(models.AddressLiquidityPool)
	err := rLp.db.
		NewSelect().
		Model(alp).
		Where("address_id = ? AND liquidity_pool_id = ?", addressId, liquidityPoolId).
		Scan(rLp.ctx)
	if err != nil {
		return nil, err
	}
	return alp, err
}

func (rLp *LiquidityPoolRepository) UpdateAddressLiquidityPool(alp *models.AddressLiquidityPool) error {
	_, err := rLp.db.NewInsert().
		Model(alp).
		On("CONFLICT (address_id, liquidity_pool_id) DO UPDATE").
		Exec(rLp.ctx)
	return err
}

func (rLp *LiquidityPoolRepository) DeleteAddressLiquidityPool(addressId uint, liquidityPoolId uint64) error {
	//_, err := rLp.db.Model().Exec(`
	//	DELETE FROM address_liquidity_pools WHERE address_id = ? and liquidity_pool_id = ?;
	//`, addressId, liquidityPoolId)

	_, err := rLp.db.NewDelete().
		Model((*models.AddressLiquidityPool)(nil)).
		Where("address_id = ?", addressId).
		Where("liquidity_pool_id = ?", liquidityPoolId).
		Exec(rLp.ctx)

	return err
}

func (rLp *LiquidityPoolRepository) UpdateAllLiquidityPool(pools []*models.AddressLiquidityPool) error {
	_, err := rLp.db.
		NewInsert().
		Model(&pools).On("CONFLICT (address_id, liquidity_pool_id) DO UPDATE").
		Exec(rLp.ctx)
	return err
}

func (rLp *LiquidityPoolRepository) GetAllByIds(ids []uint64) ([]models.LiquidityPool, error) {
	var list []models.LiquidityPool
	err := rLp.db.NewSelect().Model(&list).Where("id in (?)", bun.In(ids)).Scan(rLp.ctx)
	return list, err
}

func (rLp *LiquidityPoolRepository) SaveAllLiquidityPoolTrades(links []*models.LiquidityPoolTrade) error {
	_, err := rLp.db.NewInsert().Model(&links).Exec(rLp.ctx)
	return err
}

func (rLp *LiquidityPoolRepository) GetAll() ([]models.LiquidityPool, error) {
	var list []models.LiquidityPool
	err := rLp.db.NewSelect().Model(&list).Scan(rLp.ctx)
	return list, err
}

func (rLp *LiquidityPoolRepository) GetLastSnapshot() (*models.LiquidityPoolSnapshot, error) {
	var lps = new(models.LiquidityPoolSnapshot)
	err := rLp.db.NewSelect().Model(lps).Order("block_id desc").Limit(1).Scan(rLp.ctx)
	return lps, err
}

func (rLp *LiquidityPoolRepository) GetSnapshotsByDate(date time.Time) ([]models.LiquidityPoolSnapshot, error) {
	var list []models.LiquidityPoolSnapshot
	startDate := fmt.Sprintf("%s 00:00:00", date.Format("2006-01-02"))
	endDate := fmt.Sprintf("%s 23:59:59", date.Format("2006-01-02"))
	err := rLp.db.
		NewSelect().
		Model(&list).
		Where("created_at >= ? and created_at <= ?", startDate, endDate).
		Scan(rLp.ctx)
	return list, err
}

func (rLp *LiquidityPoolRepository) SaveLiquidityPoolSnapshots(snap []models.LiquidityPoolSnapshot) error {
	_, err := rLp.db.NewInsert().Model(&snap).Exec(rLp.ctx)
	return err
}

func (rLp *LiquidityPoolRepository) RemoveEmptyAddresses() error {
	//_, err := rLp.db.Model().Exec(`DELETE FROM address_liquidity_pools WHERE liquidity <= 0;`)
	_, err := rLp.db.NewDelete().
		Model((*models.AddressLiquidityPool)(nil)).
		Where("liquidity <= 0").
		Exec(rLp.ctx)
	return err
}

func (rLp *LiquidityPoolRepository) getPoolChainFromTags(tags map[string]string) (map[uint64][]map[string]string, error) {
	var poolsData []models.TagLiquidityPool
	err := json.Unmarshal([]byte(tags["tx.pools"]), &poolsData)
	if err != nil {
		return nil, err
	}

	data := make(map[uint64][]map[string]string)
	for _, p := range poolsData {
		firstCoinData := make(map[string]string)
		firstCoinData["coinId"] = fmt.Sprintf("%d", p.CoinIn)
		firstCoinData["volume"] = p.ValueIn

		secondCoinData := make(map[string]string)
		secondCoinData["coinId"] = fmt.Sprintf("%d", p.CoinOut)
		secondCoinData["volume"] = p.ValueIn

		data[p.PoolID] = []map[string]string{firstCoinData, secondCoinData}
	}
	return data, nil
}

type LiquidityPoolRepository struct {
	db  *bun.DB
	ctx context.Context
}

func NewLiquidityPoolRepository(db *bun.DB) *LiquidityPoolRepository {
	return &LiquidityPoolRepository{
		db:  db,
		ctx: context.Background(),
	}
}

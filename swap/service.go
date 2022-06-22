package swap

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/MinterTeam/explorer-sdk/models"
	"github.com/starwander/goraph"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"math/big"
	"sync"
)

type Service struct {
	db *bun.DB
}

func NewService(sqlDB *sql.DB, dialect *pgdialect.Dialect) *Service {
	db := bun.NewDB(sqlDB, dialect)
	return &Service{db}
}

func (s *Service) GetPoolLiquidity(pools []models.LiquidityPool, p models.LiquidityPool) *big.Float {
	if p.FirstCoinId == 0 {
		return getVolumeInBip(big.NewFloat(2), p.FirstCoinVolume)
	}

	currentVolume := p.FirstCoinVolume

	paths, err := s.FindSwapRoutePathsByGraph(pools, p.FirstCoinId, uint64(0), 4, 1500)
	if err != nil {
		paths, err = s.FindSwapRoutePathsByGraph(pools, p.SecondCoinId, uint64(0), 4, 1500)
		if err != nil {
			return big.NewFloat(0)
		}

		currentVolume = p.SecondCoinVolume
	}

	path := paths[0]

	var price *big.Float
	for i := len(path) - 2; i >= 0; i-- {
		secondCoinId := path[i+1].(uint64)
		firstCoinId := path[i].(uint64)

		for _, pool := range pools {
			if (pool.FirstCoinId == firstCoinId && pool.SecondCoinId == secondCoinId) || (pool.FirstCoinId == secondCoinId && pool.SecondCoinId == firstCoinId) {
				cprice := big.NewFloat(0)
				if pool.FirstCoinId == firstCoinId {
					cprice = computePrice(pool.SecondCoinVolume, pool.FirstCoinVolume)
				} else {
					cprice = computePrice(pool.FirstCoinVolume, pool.SecondCoinVolume)
				}

				if price == nil {
					price = cprice
				} else {
					price.Mul(price, cprice)
				}

				break
			}
		}
	}

	liquidity := getVolumeInBip(price, currentVolume)
	return liquidity.Mul(liquidity, big.NewFloat(2))
}

func (s *Service) FindSwapRoutePathsByGraph(pools []models.LiquidityPool, fromCoinId, toCoinId uint64, depth int, topN int) ([][]goraph.ID, error) {
	graph := goraph.NewGraph()
	for _, pool := range pools {
		graph.AddVertex(pool.FirstCoinId, pool.FirstCoin)
		graph.AddVertex(pool.SecondCoinId, pool.SecondCoin)
		graph.AddEdge(pool.FirstCoinId, pool.SecondCoinId, 1, nil)
		graph.AddEdge(pool.SecondCoinId, pool.FirstCoinId, 1, nil)
	}

	_, paths, err := graph.Yen(fromCoinId, toCoinId, topN)
	if err != nil {
		return nil, err
	}

	if len(paths[0]) == 0 {
		return nil, errors.New("path not found")
	}

	if depth == 0 {
		return paths, nil
	}

	var result [][]goraph.ID
	for _, path := range paths {
		if len(path) > depth || len(path) == 0 {
			break
		}

		result = append(result, path)
	}

	if len(result) == 0 {
		return nil, errors.New("path not found")
	}

	return result, nil
}

func (s *Service) MakePairsFromPaths(paths [][]goraph.ID, pools *sync.Map) ([][]Pair, error) {
	pairs := make([][]Pair, 0)
	wg := &sync.WaitGroup{}
	for _, path := range paths {
		if len(path) == 0 {
			break
		}

		wg.Add(1)
		go func(path []goraph.ID, wg *sync.WaitGroup) {
			defer wg.Done()

			currentPairs := make([]Pair, 0)
			for i := range path {
				if i == 0 {
					continue
				}

				firstCoinId, secondCoinId := path[i-1].(uint64), path[i].(uint64)
				pdata, ok := pools.Load(fmt.Sprintf("%d-%d", firstCoinId, secondCoinId))
				if !ok {
					pdata, ok = pools.Load(fmt.Sprintf("%d-%d", secondCoinId, firstCoinId))
				}

				p := pdata.(models.LiquidityPool)

				if firstCoinId == p.FirstCoinId {
					currentPairs = append(currentPairs, NewPair(
						NewTokenAmount(NewToken(p.FirstCoinId), str2bigint(p.FirstCoinVolume)),
						NewTokenAmount(NewToken(p.SecondCoinId), str2bigint(p.SecondCoinVolume)),
					))
				} else {
					currentPairs = append(currentPairs, NewPair(
						NewTokenAmount(NewToken(p.SecondCoinId), str2bigint(p.SecondCoinVolume)),
						NewTokenAmount(NewToken(p.FirstCoinId), str2bigint(p.FirstCoinVolume)),
					))
				}
			}

			pairs = append(pairs, currentPairs)
		}(path, wg)
	}
	wg.Wait()

	return pairs, nil
}

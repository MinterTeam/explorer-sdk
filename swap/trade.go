package swap

import (
	"errors"
	"github.com/MinterTeam/explorer-sdk/helpers"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
	"sync"
)

type TradeType int

const (
	TradeTypeExactInput  TradeType = 0
	TradeTypeExactOutput TradeType = 1
)

type Trade struct {
	Route          Route
	TradeType      TradeType
	InputAmount    TokenAmount
	OutputAmount   TokenAmount
	ExecutionPrice Price
	NextMidPrice   Price
	PriceImpact    *big.Int
}

func NewTrade(route Route, amount TokenAmount, tradeType TradeType) (Trade, error) {
	amounts := make([]TokenAmount, len(route.Path))
	nextPairs := make([]Pair, len(route.Pairs))

	var inputAmount, outputAmount TokenAmount
	if tradeType == TradeTypeExactInput {
		amounts[0] = amount
		for i := 0; i < len(route.Path)-1; i++ {
			tokenAmount, nextPair, err := route.Pairs[i].GetOutputAmount(amounts[i])
			if err != nil {
				return Trade{}, err
			}

			amounts[i+1], nextPairs[i] = tokenAmount, nextPair
		}

		inputAmount, outputAmount = amount, amounts[len(amounts)-1]
	} else {
		amounts[len(amounts)-1] = amount
		for i := len(route.Path) - 1; i > 0; i-- {
			tokenAmount, nextPair, err := route.Pairs[i-1].GetInputAmount(amounts[i])
			if err != nil {
				return Trade{}, err
			}

			amounts[i-1], nextPairs[i-1] = tokenAmount, nextPair
		}

		outputAmount, inputAmount = amount, amounts[0]
	}

	if inputAmount.Amount.Cmp(big.NewInt(0)) == 0 || outputAmount.Amount.Cmp(big.NewInt(0)) == 0 {
		return Trade{}, errors.New("insufficient reserve")
	}

	return Trade{
		Route:          route,
		TradeType:      tradeType,
		InputAmount:    inputAmount,
		OutputAmount:   outputAmount,
		ExecutionPrice: NewPrice(inputAmount.GetCurrency(), outputAmount.GetCurrency(), inputAmount.GetAmount(), outputAmount.GetAmount()),
		NextMidPrice:   NewPriceFromRoute(NewRoute(nextPairs, route.Input, nil)),
		PriceImpact:    computePriceImpact(route.MidPrice, inputAmount, outputAmount),
	}, nil
}

func (t *Trade) GetMaximumAmountIn(slippageTolerance float64) TokenAmount {
	if t.TradeType == TradeTypeExactInput {
		return t.InputAmount
	}

	maximumAmountIn := new(big.Int)
	inputAmount := new(big.Float).SetInt(t.InputAmount.GetAmount())
	percent := big.NewFloat(1 + slippageTolerance)
	new(big.Float).Mul(inputAmount, percent).Int(maximumAmountIn)

	return NewTokenAmount(t.InputAmount.Token, maximumAmountIn)
}

func (t *Trade) GetMinimumAmountOut(slippageTolerance float64) TokenAmount {
	if t.TradeType == TradeTypeExactOutput {
		return t.OutputAmount
	}

	minimumAmountOut := new(big.Int)
	outputAmount := new(big.Float).SetInt(t.OutputAmount.GetAmount())
	percent := big.NewFloat(1 + slippageTolerance)
	new(big.Float).Quo(outputAmount, percent).Int(minimumAmountOut)

	return NewTokenAmount(t.InputAmount.Token, minimumAmountOut)
}

type TradeOptions struct {
	MaxNumResults int
	MaxHops       int
}

func computePriceImpact(midPrice Price, inputAmount TokenAmount, outputAmount TokenAmount) *big.Int {
	mid := helpers.Pip2Bip(midPrice.Value)
	input := helpers.Pip2Bip(inputAmount.GetAmount())
	output := helpers.Pip2Bip(outputAmount.GetAmount())

	exactQuote := new(big.Float).Mul(mid, input)
	numerator := new(big.Float).Sub(exactQuote, output)
	slippage := new(big.Float).Quo(numerator, exactQuote)

	wei := new(big.Int)
	new(big.Float).Mul(slippage, big.NewFloat(params.Ether)).Int(wei)

	return wei
}

func GetBestTradeExactIn(possibleRoutes [][]Pair, originalAmountIn TokenAmount, tokenOut Token) (*Trade, error) {
	var bestTrade *Trade

	wg, mu := &sync.WaitGroup{}, sync.Mutex{}
	for _, routePairs := range possibleRoutes {
		wg.Add(1)
		go func(routePairs []Pair, wg *sync.WaitGroup) {
			defer wg.Done()

			trade, err := NewTrade(NewRoute(routePairs, originalAmountIn.GetCurrency(), &tokenOut), originalAmountIn, TradeTypeExactInput)
			if err != nil {
				return
			}

			mu.Lock()
			if bestTrade == nil {
				bestTrade = &trade
				mu.Unlock()
				return
			}

			if bestTrade.OutputAmount.Amount.Cmp(trade.OutputAmount.Amount) == -1 {
				bestTrade = &trade
			}
			mu.Unlock()
		}(routePairs, wg)
	}
	wg.Wait()

	return bestTrade, nil
}

func GetBestTradeExactOut(possibleRoutes [][]Pair, originalAmountOut TokenAmount, tokenIn Token) (*Trade, error) {
	var bestTrade *Trade

	wg, mu := &sync.WaitGroup{}, sync.Mutex{}
	for _, routePairs := range possibleRoutes {
		wg.Add(1)
		go func(routePairs []Pair, wg *sync.WaitGroup) {
			defer wg.Done()

			trade, err := NewTrade(NewRoute(routePairs, tokenIn, &originalAmountOut.Token), originalAmountOut, TradeTypeExactOutput)
			if err != nil {
				return
			}

			mu.Lock()
			if bestTrade == nil {
				bestTrade = &trade
				mu.Unlock()
				return
			}

			if bestTrade.InputAmount.Amount.Cmp(trade.InputAmount.Amount) > 0 {
				bestTrade = &trade
			}
			mu.Unlock()
		}(routePairs, wg)
	}
	wg.Wait()

	return bestTrade, nil
}
